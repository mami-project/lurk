package main

import (
        "bytes"
	"crypto/x509"
	"crypto/tls"
        "fmt"
	"io/ioutil"
        "os"
        "net/http"
        "net/url"

)
func main() {
        if len(os.Args) != 2 {
                panic("Illegal number of arguments: Introduce just the certificate's $uuid")
        }

        //Reads cert from file
        CACert, err := ioutil.ReadFile("./serverKey/cert.pem")
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
        //fmt.Printf("sending %v", os.Args[1])
        data := url.Values{}
        data.Set("uri",os.Args[1])
        req, err := http.NewRequest("POST", "https://CertificateAuthoritySTAR:9200/terminate", bytes.NewBufferString(data.Encode()))
        if err != nil {
                panic (err)
        }
        resp, _ := client.Do(req)
        fmt.Println(resp.Status)


}

