package main

import (
    "fmt"
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "os"
    "io"
    "github.com/satori/go.uuid"
    "github.com/gorilla/mux"
    "os/exec"
    "encoding/pem"
)

var openPort = false //Just to make sure ports open only once
var completionURL_value = "in progress" //message that pops when STAR clientes asks for the uri if it s too early
var cronTaskID = 0 //counter for crontab

type cdn_post struct {      /*fields in the STAR client CSR*/
        Csr string   `json:"csr"`
        LifeTime int `json:"lifetime"`
        Validity int `json:"validity"`
}

type csr_struct struct {     //fields that need treatment in the received csr string
        subjectName string
}
//Main handler for STAR client requests
func parseJsonPOST(w http.ResponseWriter, r *http.Request) {
        var t cdn_post
        err := json.NewDecoder(r.Body).Decode(&t)
        if err != nil {
                panic(err)
        }
        block, _ := pem.Decode([]byte(t.Csr))
        if block == nil {
                panic("failed to parse PEM block containing the public key")
        }
        if(t.LifeTime > 365 || t.Validity > 200) { //Note that lifetime units are days but validity is in hours
                fmt.Fprintln(w, "Enter parameters are not valid. Maximum lifetime = 365, maximum validity 200.")
        }else {
                cmdS := decodeCsr(t.Csr)
                csr_fields := parseFieldsOfCsr(cmdS)
                fmt.Fprintln(w, "Received parameters are valid: LifeTime: ",t.LifeTime," Validity", t.Validity,
                " Domain:", csr_fields.subjectName)
                fmt.Fprintln(w, "Accepted, poll at /completionURL")
                //fmt.Printf("%q",csr_fields) //Uncomment this line for some extra-checking

                //Creates a file with cert validity for local certbot to read and deletes the previous one
                _, noFile := os.Stat("STARValidityCertbot")
                 if noFile == nil {
                        os.Remove("STARValidityCertbot") //Deletes previous file
                }

                toFileErr := ioutil.WriteFile("STARValidityCertbot", []byte(strconv.Itoa(t.Validity)), 0644)
                if toFileErr != nil {
                        panic(toFileErr)
                }

                addToCron(csr_fields.subjectName, t.LifeTime) 

                callCertbot(csr_fields.subjectName) /*Executes certbot for a certain domain*/

                /*file, err := ioutil.ReadFile("/root/certId") //This used to save the real URI provided by the CA
                                                                //May be used in the future
                if err != nil {
                        panic(err)
                }
                completionURL_value = string(file)
                */
                completionURL_value = uuid.NewV4().String()
                storeForRenewal(t.Validity)
                go post_cert()
                cronTaskID++ //Counter that serves together with the uuid to make each petition unique

                //Removes tmp files, comment this function if you want more information
                rmTmpFiles()

        }

}
/*
    Adds the renewal to cron, it deletes itself when lifetime. To change the renewal hours 
    go to exeAutoRenew.sh
*/
func addToCron(domainName string, lifeTime int) {
        addTaskCommand := []string{"addTask.sh", domainName,strconv.Itoa(lifeTime),strconv.Itoa(cronTaskID)}
        _,err := exec.Command("/bin/sh",addTaskCommand...).Output()
        if err != nil {
                panic(err)
        }

}
/*
Posts the certificates and keeps serving them at :9500/uuid

*/
func post_cert () {
    //fmt.Printf("\nGET the certificate from: %v", completionURL_value)
    http.HandleFunc("/" + completionURL_value, func(w http.ResponseWriter, r *http.Request) {
       http.ServeFile(w, r, "/root/starCerts/" + completionURL_value + "/certificate.pem")
        })
    if openPort != true {
        openPort = true
        err := http.ListenAndServe(":9500", nil)
       // err := http.ListenAndServe("9500", "server.crt", "server.key", nil) //Uncomment this and comment the previous one
        if err != nil {                                                       //to make  retrieving the cert https.
        panic(err)                                                            //It has been tested and it works! 
        }
    }

}
/*
Invokes post_completionURL when client asks for the certificate
*/
func answerAGet(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Responding to a GET request")
        fmt.Println(r.UserAgent())

        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, completionURL_value)

}
/*
    Lets the STAR client know when the certificate is ready and where it is.

*/
func post_completionURL() {
        router := mux.NewRouter().StrictSlash(true)
        router.HandleFunc("/completionURL", answerAGet).Methods("GET")
        err := http.ListenAndServe(":9999", router)
        //err := http.ListenAndServeTLS(":9999","server.crt", "server.key", router) //Uncomment this and comment the previous one
        if err != nil{                                                              //to make getting the URI https.
                panic(err)                                                          //It hasn't been tested YET

        }

}

/*
    Creates a copy of the csr in file tmpCsr and returns the csr as plain text
*/
func decodeCsr (csr string)(cmdS string) { 
    f, err := os.Create("tmpCsr")   
    if err != nil {
        panic(err)
    }

    _, err2 := f.WriteString(csr)
    defer f.Close()
    if err2 != nil {
        panic(err2)
    }
    cmd, err3:= exec.Command("/bin/sh", "getCsrAsText.sh").Output()
    if err3 != nil {
        panic(err3)
    }
    cmdS = (string)(cmd)
    return cmdS
}
/*
    Decodes the csr so that fields can be read.
    For now it is supposed that the required domain is correct.
*/
func parseFieldsOfCsr(cmd string)(csrFields csr_struct) {  /*Returns and array with each important field of the csr in tmpCsr*/
        fmt.Printf("String: %s FIN string", cmd)

        g := strings.SplitN(cmd, "CN=", 2) // Keeps the common name field in g[1]
        f := strings.FieldsFunc(g[1], func(r rune) bool { //f is an array with the fields requested in csr_struct
                switch r {
                        case ' ', '/', '\n' :
                                return true
                        }
                        return false
        })
        csrFields = csr_struct{     /*if csr_struct changes add the rest of the fields here*/
                subjectName: f[0]}
        return csrFields
}
/*
Stores csr, validity, certificate and URI. Lifetime is not required to store as it is controlled by a cronjob.
It also creates a file to link the cronTaskID to the uuid 
*/
func storeForRenewal (validity int) {
    var certDirName = completionURL_value //All the cert info is kept under a file named as its uri
    var linkFileName = "link" + strconv.Itoa(cronTaskID)
    var csrFileName = "csr" + strconv.Itoa(cronTaskID)
    var validityFileName = "validity" + strconv.Itoa(cronTaskID)
    var uriFileName = "uri" + strconv.Itoa(cronTaskID)
    var certFileName = "certificate.pem"

    //Creates storage directories
    if _, err := os.Stat("/root/starCerts"); os.IsNotExist(err) {
        err = os.Mkdir("/root/starCerts", 0644)
        if err != nil {
            panic(err)
        }
    }
    err := os.Mkdir("/root/starCerts/" + certDirName, 0644)
    if err != nil {
        panic(err)
    }
    if _, err := os.Stat("/root/starCerts/links"); os.IsNotExist(err) {
        err := os.Mkdir("/root/starCerts/links", 0644)
         if err != nil {
                panic(err)
         }
    }
    //links cronTaskID to URI's uuid
    e, err := os.Create("/root/starCerts/links/" + linkFileName)
    if err != nil {

        panic(err)
    }
    defer e.Close()
    e.WriteString(completionURL_value)


    //Saves the csr
    f, err := os.Create("/root/starCerts/" + certDirName + "/" + csrFileName)
    if err != nil {

        panic(err)
    }
    defer f.Close()
    in, err := os.Open("tmpCsr")
    if err != nil {
        panic(err)
    }
    defer in.Close()
    _, err = io.Copy(f, in)

    //Saves the validity
    g, err := os.Create("/root/starCerts/" + certDirName + "/" + validityFileName)
    if err != nil {

        panic(err)
    }
    defer g.Close()
    g.WriteString(strconv.Itoa(validity))

    //Saves the cert
    h, err := os.Create("/root/starCerts/" + certDirName + "/" + certFileName)
    if err != nil {

        panic(err)
    }
    defer h.Close()
    inCert, err := os.Open("ObtainedCERT.pem")
    if err != nil {
        panic(err)
    }
    defer inCert.Close()
    _, err = io.Copy(h, inCert)


    //Saves the uri
    i, err := os.Create("/root/starCerts/" + certDirName + "/" + uriFileName)
    if err != nil {

        panic(err)
    }
    defer i.Close()
    i.WriteString(completionURL_value)


}
/*
Executes certbot from certbot/certbot/main.py
Executing certbot using one of its auto executables
can destroy the changes done in it that are required
for STAR so execute using this main.py to keep the 
changes :)
*/
func callCertbot(domainName string){
        certbotCommand  := []string{"certbotCall.sh", domainName}
        ex,err := exec.Command("/bin/sh",certbotCommand...).Output()
        if err != nil {
                panic(err)
        }
        fmt.Printf("Ejecucion finalizada %s",ex)

}

/*
Executes last. 
starCerts isn't supposed to be deleted unless you
restart the proxy because it contains all the live
information.
If you plan to lauch the proxy but you had 
obtained a certificate with STAR before, then 
use: sudo rm -rf starCerts
*/
func rmTmpFiles () {
    err := os.Remove("certId")
    if err != nil {
        panic (err)
    }
    err := os.Remove("ObtainedCERT.pem")
    if err != nil {
        panic (err)
    }
    err := os.Remove("STARValidityCertbot")
    if err != nil {
        panic(err)
    }
    err := os.Remove("tmpCsr")
    if err != nil {
        panic(err)
    }

}


func main() {
    fmt.Println("Proxy STAR ON")
    go post_completionURL()
    http.HandleFunc("/star/registration", parseJsonPOST)
    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
    if err != nil {
        panic(err)
    }

}
