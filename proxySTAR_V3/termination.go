package main

import (
        "bytes"
        "crypto/x509"
        "crypto/tls"
	"encoding/json"
        "fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
        "io/ioutil"
        "os"
        "net/http"
	"time"
        //"net/url"

)
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
/*Generates the full Json Web Token*/
func genToken(fullJson *terminationMessage)(stringJWT string) {
	//fmt.Println(fullJson)
	bytes, _ := ioutil.ReadFile("./terminationJWTKey/sample_key.priv")

	claims := jws.Claims{}
	claims["protected"] = fullJson.Protected
	claims["payload"] = fullJson.Payload
	claims["signature"] = fullJson.Signature

	rsaPrivate, _ := crypto.ParseRSAPrivateKeyFromPEM(bytes)
	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)

	b, _ := jwt.Serialize(rsaPrivate)
	
	s := string(b[:])
	fmt.Println("Gonna send: ")
	fmt.Println(s)
	return s


}

func main() {
        if len(os.Args) != 2 {
                fmt.Printf("Illegal number of arguments: Introduce just 2: command $uuid\n$uuid is the certificate's uri for renewal. You can check it at starCerts/\n")
                os.Exit(1)
        }

        //Reads cert from file
        CACert, err := ioutil.ReadFile("./serverKey/certRenewals.pem")
        if err != nil {
          panic (err)
        }
        //Parses the cert
        CA_certPool := x509.NewCertPool()
        booleanValue := CA_certPool.AppendCertsFromPEM(CACert)
        fmt.Printf("%v", booleanValue)



        client := &http.Client{
          Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                RootCAs:      CA_certPool,
            },
          },
        }

	//Ask for nonce
        nonceReq, _ := json.Marshal("I want a nonce")
        reqNonce, err := http.NewRequest("POST", "https://RenewalSTAR:9200/getNonce", bytes.NewBuffer(nonceReq))
        if err != nil{
                panic(err)
        }
        resp, _ := client.Do(reqNonce)
        fmt.Println(resp.Status)
	
        buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	nonce := buf.String()
	fmt.Printf("The nonce received is: %v", nonce)
	
	time.Sleep(3 * time.Second)
	

	url := os.Args[1]
	signatureAlg := "ES256"
	kid := "proxySTAR.com"
	status := "canceled"
	signature := "this-is-the-signature"
	protectedSubStruct := &cod_base64{Alg: signatureAlg, Kid: kid, Nonce: nonce, Url: url}
	payloadSubStruct := &terminationStatus{Status: status}
	fullCancelation := &terminationMessage{Protected: *protectedSubStruct, Payload: *payloadSubStruct, Signature: signature}
	
	var stringJWT string = genToken(fullCancelation)

	fullCanJson, _ := json.Marshal(stringJWT)
	




        
	req, err := http.NewRequest("POST", "https://RenewalSTAR:9200/terminate", bytes.NewBuffer(fullCanJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "acme-server/termination")
        if err != nil {
                panic (err)
        }
        resp, _ = client.Do(req)
        fmt.Println(resp.Status)

}
