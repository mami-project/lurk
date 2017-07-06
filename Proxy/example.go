/******************example.go************************/
package main

import (
    "fmt"
//    "bytes"
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    "os/exec"
//    "crypto/x509"
//    "crypto/rsa"
    "encoding/pem"
)


var completionURL_value = "in progress"
var cronTaskID = 0 //counter for crontab

type cdn_post struct {      /*fields in the STAR client CSR*/
        Csr string   `json:"csr"`
        LifeTime int `json:"lifetime"`
        Validity int `json:"validity"`
}

type csr_struct struct {     /*fields that need treatment in the received csr string*/
        subjectName string
}

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
        //fmt.Println(t)
        if(t.LifeTime > 365 || t.Validity > 200) {
                fmt.Fprintln(w, "Enter parameters are not valid. Maximum lifetime = 365, maximum validity 200.")
        }else {
                cmdS := decodeCsr(t.Csr)
                csr_fields := parseFieldsOfCsr(cmdS)
                fmt.Fprintln(w, "Received parameters are valid: LifeTime: ",t.LifeTime," Validity", t.Validity,
                " Domain:", csr_fields.subjectName)
                fmt.Fprintln(w, "Accepted, poll at /completionURL")
                fmt.Printf("%q",csr_fields)
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

                file, err := ioutil.ReadFile("/root/certId")
                if err != nil {
                        panic(err)
                }
                completionURL_value = string(file)
                cronTaskID++

        }

}
func addToCron(domainName string, lifeTime int) {
        addTaskCommand := []string{"addTask.sh", domainName,strconv.Itoa(lifeTime),strconv.Itoa(cronTaskID)}
        _,err := exec.Command("/bin/sh",addTaskCommand...).Output()
        if err != nil {
                panic(err)
        }

}
func post_completionURL() {
        fmt.Println("Executing")
        router := mux.NewRouter().StrictSlash(true)
        router.HandleFunc("/completionURL", answerAGet).Methods("GET")
        err := http.ListenAndServe(":9999", router)
        if err != nil{
                panic(err)

        }

}
func answerAGet(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Responding to a GET request")
        fmt.Println(r.UserAgent())

        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, completionURL_value)

}

func decodeCsr (csr string)(cmdS string) { /*Returns the csr as plain text*/
    f, err := os.Create("tmpCsr")   /*tmpCsr will contain a copy of the full csr*/
    if err != nil {
        panic(err)
    }

    _, err2 := f.WriteString(csr)
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
func parseFieldsOfCsr(cmd string)(csrFields csr_struct) {  /*Returns and array with each important field of the csr in tmpCsr*/
        //fmt.Printf("String: %s FIN string", cmd)

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
func callCertbot(domainName string){
        certbotCommand  := []string{"certbotCall.sh", domainName}
        ex,err := exec.Command("/bin/sh",certbotCommand...).Output()
        if err != nil {
                panic(err)
        }
        fmt.Printf("Execution ended %s",ex)

}


func main() {
    go post_completionURL()
    http.HandleFunc("/star/registration", parseJsonPOST)
    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
    if err != nil {
        panic(err)
    }
}
