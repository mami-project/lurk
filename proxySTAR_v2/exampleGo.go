package main

import (
	"net/http"
	//"fmt"
	//"os/exec"
	//"strings"
)
func main(){
	/*uriNotValid := "xxxx=https://CertificateAuthoritySTAR:9898/ff822835fec078588cda4c3364d4a94b0db1"
	uriValid := strings.SplitN(uriNotValid, "=", 2)
	uriValidString := uriValid[1] 
	getRenewalUri := []string{"--cacert", "./serverKey/cert.pem",uriValidString}
	fmt.Println(uriValidString)
	serialST,err := exec.Command("curl",getRenewalUri...).Output()
	if err != nil {
		panic(err)
	}

	serialSTString := (string)(serialST)
	fmt.Printf("the value is: %v", serialSTString)
	*/
	resp, err := http.Get("https://CertificateAuthoritySTAR:9898/ff822835fec078588cda4c3364d4a94b0db1")
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
}
