package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*
Storages the csr in format:
---BEGIN CERTIFICATE REQUEST ---
64 char
64 char
...
---END CERTIFICATE REQUEST---
Also adds bit stuffing and fixes some chars that are not valid for openssl
*/

type NotCookedCsr struct{
	rawCsr string
}
func (strCsr *NotCookedCsr) repairCsr(g *os.File) {
	var header1, header2, body string
	i := 0
	rawLen := len(strCsr.rawCsr)
	header1 = "-----BEGIN CERTIFICATE REQUEST-----"
	header2 = "-----END CERTIFICATE REQUEST-----"
	g.WriteString("\n")
	g.WriteString(header1)

	//Replaces '-' and '_'
	strCsr.rawCsr = strings.Replace(strCsr.rawCsr, "-", "+", -1)
	strCsr.rawCsr = strings.Replace(strCsr.rawCsr, "_", "/", -1)

	//Groups the lines
	for i < rawLen {
		if (rawLen - 64) > i {
			body = strCsr.rawCsr[i : i+64]
		} else {
			body = strCsr.rawCsr[i:rawLen]
		}
		g.WriteString("\n")
		i += 64
		g.WriteString(body)

	}
	fileName := g.Name()
	numbBytes, _ := ioutil.ReadFile(fileName)

	bytesLen := (len(numbBytes) + len(header2)) * 8
	fmt.Printf("GONNA STUFF %v,", bytesLen)
	//Bit stuffing
	if bytesLen%6 != 0 {
		fmt.Printf("STUFFING 1")
		g.WriteString("=")
		bytesLen += 8
		if bytesLen%6 != 0 {
			fmt.Printf("stuffin 2")
			g.WriteString("=")
			bytesLen += 8
		}
		if bytesLen%6 != 0 {
			panic("Error, csr malformed, cant be storaged for star")
		}
	}

	g.WriteString("==\n")
	g.WriteString(header2)

}

func main(){

//fullcsr1:= "MIICmDCCAYACAQAwUzELMAkGA1UEBhMCTU0xCjAIBgNVBAgMAWsxDDAKBgNVBAcMA2xsbDEKMAgGA1UECgwBdDEMMAoGA1UECwwDdGlkMRAwDgYDVQQDDAdieWUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmm74/AotkDGzsdVsn+VuZ4FHW2+lf3HrLcDlpWHvBl3WSLg2TJHXdl8F6GtI3w91Cws/8wT4g+W33GYB0WAFWIGvzTPajeZ3jQt4t98bpzbuaFZz8QCoQVuEOuk8CCQ5/Cezbml3loMtXTuR+R1cOuVB9sFXbpoGvGL2fbAmrTtmOY9ZoXaLQmN7sj+4TjKRtZvVdpiLRaYp608ct2h3E6R2Nzm0OHdI35y61jaw46WiXCM30W/V2/Ia0c35Jdy4vbPybH1+k4ajmrlwiFrO986AlAxvxDZIKtahQFqMdH3hEuzTR6OnDwMlDtkLXThE9XSmcAhdYd9RLC8hF33ASQIDAQABoAAwDQYJKoZIhvcNAQELBQADggEBAF7ja2QCYDPfJ3kY0f4eSYaAaQbabQ6TA2dS5AFz+WSdzBQTTa8uTzgrKOwe8mQoHhjsNHW6aRpYCxje2v0pzTMw27g2YmXdfEfmWsF4GHk2NZ3ECE7LwA0YlsGZXpmYkUT89+69cJiiqiWUpwaGQdbx2OzsN7tlHLlLDufQubnMetOfNb7SbyMpCdNssAaj7gkkmeOHk9rjlkrpkBJf8lBb2xIoXFH3iswaojVO3pAKZytDPrx9tsAsLstx6Jv6+O5lfr9rS4+EAT19yeZgd/64qTjlyx1r4nkEp7Z/brWh4X3q8zUhBQCLSeHIXp9nWj69WGXFtTOqcyx+uruc/Qw"
fullcsr2:="MIICmzCCAYMCAQAwVjELMAkGA1UEBhMCRVMxDTALBgNVBAgMBE1hZGExCzAJBgNVBAcMAm1hMQowCAYDVQQKDAFUMQwwCgYDVQQLDANUSUQxETAPBgNVBAMMCGJ5ZTMuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmm74/AotkDGzsdVsn+VuZ4FHW2+lf3HrLcDlpWHvBl3WSLg2TJHXdl8F6GtI3w91Cws/8wT4g+W33GYB0WAFWIGvzTPajeZ3jQt4t98bpzbuaFZz8QCoQVuEOuk8CCQ5/Cezbml3loMtXTuR+R1cOuVB9sFXbpoGvGL2fbAmrTtmOY9ZoXaLQmN7sj+4TjKRtZvVdpiLRaYp608ct2h3E6R2Nzm0OHdI35y61jaw46WiXCM30W/V2/Ia0c35Jdy4vbPybH1+k4ajmrlwiFrO986AlAxvxDZIKtahQFqMdH3hEuzTR6OnDwMlDtkLXThE9XSmcAhdYd9RLC8hF33ASQIDAQABoAAwDQYJKoZIhvcNAQELBQADggEBAAYi963id0bCihyJDDlekqNDKuqZ/hiafmii/b9sTy0u28OXcM1EsEGSIVpsLeVdb0bITi96waHv5jgVzeOhFr+rdjsrj+/JSiYP/1zwZeivuBFAHl4zob7eNQEhCukJyR207KtnGFeVqu4EAy2rG1ZA5ra7VA1s0klGJF1LxrapaMZPn8vlvAETpHYPY3qoSWtsOXATVs6Inp8DX+94rNSNABRVV1fJ2+1i2I2AhSgToRmdKuskivWe1wNcmjMH12N6usmdWUvCwcCrf8n4A2HLcFg1kynecc4fh15brfy9kzBxvHBWpO6BMGZzmos84f0cTar3TyQMCnRCZh7jAOc"
lilcsr := NotCookedCsr{rawCsr : fullcsr2}

//Saves the csr
    f, err := os.Create("./starCerts/csr")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    lilcsr.repairCsr(f)

}
