package main

import (
  "encoding/base64"
  "encoding/json"
  "fmt"
  "github.com/SermoDigital/jose/crypto"
  "github.com/SermoDigital/jose/jws"
  "github.com/LarryBattle/nonce-golang"
  "time"
  "os"
  "os/exec"
  "io/ioutil"
  "net/http"
  "strings"
  "strconv"
)

var renewStep int

type cod_base64 struct {

        Alg string
        Kid string
        Nonce string
        Url string

}
type terminationStatus struct {

        Status string
}

type terminationMessage struct {
        Protected cod_base64
        Payload terminationStatus
        Signature string

}

/*Paralel goroutine that checks if new STAR certs were issued*/
func checkStatus() {
    time.Sleep(time.Duration(renewStep) * time.Millisecond) // 1s = 1000
    //fmt.Println("Crontab updated")  //Uncomment to see a console message everytime it checks

     _, err := os.Stat("../boulder/renewTmp/renew1")
        if err == nil {
                      fmt.Printf("File deleted")
                      renewBytes, _ := ioutil.ReadFile("../boulder/renewTmp/renew1")
                      renewStr := string(renewBytes[0:len(renewBytes)])
                      contents := strings.SplitN(renewStr," ",3)
                      addToCron(contents[0], contents[1], contents[2])
                      err = os.Remove("../boulder/renewTmp/renew1") //Can't be deferred as this process never ends
                      if err != nil {
			fmt.Printf("Failed")
			panic(err)                        
                      }
    } else {
    //fmt.Printf("File doesnt exist")
}



    checkStatus()
}
/*Adds the new certs to a cron job that renews them while 
lifeTime doesn't expire. The cron job executes every validity
period*/
func addToCron (domainStr, lifeTimeStr, crtUuid string) {
  addTaskCommand := []string{"addTask.sh", domainStr, lifeTimeStr, crtUuid}
        fmt.Print(addTaskCommand)
        _,err := exec.Command("/bin/sh",addTaskCommand...).Output()
        if err != nil {
                fmt.Printf("El error es: %+v fin", err)
                panic(err)
        }
}
/*Checkes wheter or not the signature is valid or not
for the received JWT*/
func validateToken(accessToken string) {

	bytes, _ := ioutil.ReadFile("./terminationJWTKey/sample_key.pub")
        rsaPublic, _ := crypto.ParseRSAPublicKeyFromPEM(bytes)

        //accessToken, _ := ioutil.ReadFile("./tokenExample")
        jwt, err := jws.ParseJWT([]byte(accessToken))
        if err != nil {
                panic(err)
        }

        // Validate token
        if err = jwt.Validate(rsaPublic, crypto.SigningMethodRS256); err != nil {
                panic(err)
        }




}
/*Decodes the JWT body. Headers are ignored*/
func decodeJWT(tokenComponent []string)(JWTbody cod_base64){
	
//The struct to return
	JWTbody = cod_base64{}

 //Decodes the body and retrieves the uri field
        s,_ := base64.StdEncoding.DecodeString(tokenComponent[1])
        urlWithFiller := strings.SplitN(string(s[:]),",\"Url\":\"",2)
        urlSlice := strings.SplitN(urlWithFiller[1],"\"}",2)
        JWTbody.Url = urlSlice[0]

//Retrieves the alg
	algWithFiller := strings.SplitN(string(s[:]),"\"Alg\":\"",2)
        algSlice := strings.SplitN(algWithFiller[1],"\",",2)
        JWTbody.Alg = algSlice[0]
	
//Retrieves the kid
	kidWithFiller := strings.SplitN(string(s[:]),",\"Kid\":\"",2)
        kidSlice := strings.SplitN(kidWithFiller[1],"\",",2)
        JWTbody.Kid = kidSlice[0]

//Retrieves the nonce
	nonceWithFiller := strings.SplitN(string(s[:]),",\"Nonce\":\"",2)
        nonceSlice := strings.SplitN(nonceWithFiller[1],"\",",2)
        JWTbody.Nonce = nonceSlice[0]


	return JWTbody
}
/*Reads the JWT and deals with it.
To be valid both valid and signature are
checked*/
func processCancelations (w http.ResponseWriter, r *http.Request) {
	
	bodyBuffer, _ := ioutil.ReadAll(r.Body)
	var receivedJson string
	err := json.Unmarshal(bodyBuffer,&receivedJson)
	if err != nil{
		fmt.Printf("I failed")
	}
	fmt.Println(" TOKEN RECEIVED: ")
        //fmt.Println(receivedJson)
	validateToken(receivedJson)	
	fmt.Println("TOKEN VALIDATED")
	
	//Returns the JWT separated in: header, body and signature
	tokenComponent := strings.SplitN(receivedJson, ".", 3)
	
	
	JWT_parsed := decodeJWT(tokenComponent)
	
	//fmt.Printf("DNO wants to cancel UUID: %v", JWT_parsed.Url)
	fmt.Printf("Full JWT body once decoded is: %v",JWT_parsed)
	
	//Checks nonce
	err = nonce.CheckToken(JWT_parsed.Nonce)
	if err != nil {
		panic(err)	
	}
	
	//Removes nonce from valids list
	nonce.MarkToken(JWT_parsed.Nonce)
	fmt.Println("Nonce validated and marked as used :)")


        //Deletes certificate so it never gets served again
        certFile := "../boulder/starCerts/" + JWT_parsed.Url + "/certificate.pem"
	fmt.Printf("The certFile is %v", certFile)
        _, err = os.Stat(certFile)
	if err == nil{
		rmCert := []string{"rm", certFile}
        	_,err := exec.Command("sudo",rmCert...).Output()
        	if err != nil {
                	panic(err)
        	}	
		
		//Deletes cronjob only if cert existed
		exeTermination := []string{"exeTermination.sh", JWT_parsed.Url} 
		_,err = exec.Command("/bin/sh",exeTermination...).Output()
		if err != nil {
                	panic(err)
	        }
		w.WriteHeader(http.StatusAccepted)
	}else{
		//Returns 404 if cant find cert
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Already terminated")
	}	
  
}

/*Generates new tokens and adds them to a db*/
func nonceInstantiator(w http.ResponseWriter, r *http.Request){
	
	bodyBuffer, _ := ioutil.ReadAll(r.Body)

	var receivedMessage string
	err := json.Unmarshal(bodyBuffer,&receivedMessage)
        if err != nil{
                fmt.Printf("I failed")
        }

	fmt.Printf("Received: %s", receivedMessage)
	w.WriteHeader(http.StatusAccepted)
	valid_nonce := nonce.NewToken() 
        fmt.Fprint(w, valid_nonce)
	

}

/*Main listens at port 9200 in /terminate and /getNonce
To correctly make a cancelation a nonce has to 
be obtained via getNonce*/

func main() {
     if len(os.Args) != 2 {
	fmt.Printf("USAGE: command time.Milliseconds\nThis value sets the time between checks\n")
	os.Exit(1) 
     }
     renewStep,_ = strconv.Atoi(os.Args[1]) //renew every "renewStep" seconds
     fmt.Printf("RenewStep is: %v\n", renewStep)
     go checkStatus()
     fmt.Println("RenewalManager status is: ACTIVE")
     http.HandleFunc("/terminate", processCancelations)
     http.HandleFunc("/getNonce", nonceInstantiator)
     err := http.ListenAndServeTLS(":9200", "certRenewals.pem", "keyRenewals.pem", nil)
     if err != nil {
	panic (err)
     }
    
}

