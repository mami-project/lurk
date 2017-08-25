package main

import (
  "fmt"
  "time"
  "os"
  "os/exec"
  "io/ioutil"
  "net/http"
  "strings"
)


func checkStatus() {
    time.Sleep(10000 * time.Millisecond) //currently 10 seconds. 1s = 1000
  //  fmt.Println("Crontab updated")  //Uncomment to see a Message everytime it checks

     _, err := os.Stat("./renewTmp/renew1")
        if err == nil {
                      fmt.Printf("File deleted")
                      renewBytes, _ := ioutil.ReadFile("./renewTmp/renew1")
                      renewStr := string(renewBytes[0:len(renewBytes)])
                      contents := strings.SplitN(renewStr," ",3)
                      addToCron(contents[0], contents[1], contents[2])
                      err = os.Remove("./renewTmp/renew1") //Can't be deferred as this process never ends
                      if err != nil {
                        panic(err)
                      }
    } else {
    //fmt.Printf("File doesnt exist")
}



    checkStatus()
}

func addToCron (domainStr, lifeTimeStr, crtUuid string) {
  addTaskCommand := []string{"addTask.sh", domainStr, lifeTimeStr, crtUuid}
        fmt.Print(addTaskCommand)
        _,err := exec.Command("/bin/sh",addTaskCommand...).Output()
        if err != nil {
                fmt.Printf("El error es: %+v fin", err)
                panic(err)
        }
}

func processCancelations (w http.ResponseWriter, r *http.Request) {
  
	// Buffer the body
	bodyBuffer, _ := ioutil.ReadAll(r.Body)
	
	//Get uri
	uri := strings.SplitN(string(bodyBuffer), "=", 2)	
	//fmt.Printf("DNO wants to cancel UUID: %v", uri[1])
	
	//Delete cronjob
	exeTermination := []string{"exeTermination.sh", uri[1]} 
	_,err := exec.Command("/bin/sh",exeTermination...).Output()
	if err != nil {
                panic(err)
        }
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "Deleted")
  
}

func main() {
     go checkStatus()
     fmt.Println("RenewalManager status is: ACTIVE")
     http.HandleFunc("/terminate", processCancelations)
     err := http.ListenAndServeTLS(":9200", "cert.pem", "key.pem", nil)
     if err != nil {
	panic (err)
     }
    
}

