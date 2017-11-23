#!/bin/bash

if [ "$#" -ne 5 ]; then
    echo "Illegal number of parameters"
    exit 1
fi
domainName=$1 # e.g. bye.com
getDeadLineD=$2 # "day month" format, e.g. 19 06 18
getDeadLineM=$3
getDeadLineY=$4
uri=$5 # string

cert="-----BEGIN CERTIFICATE-----\nMIIEijCCA3KgAwIBAgICEk0wDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAwwgY2Fj\na2xpbmcgY3J5cHRvZ3JhcGhlciBmYWtlIFJPT1QwHhcNMTUxMDIxMjAxMTUyWhcN\nMjAxMDE5MjAxMTUyWjAfMR0wGwYDVQQDExRoYXBweSBoYWNrZXIgZmFrZSBDQTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMIKR3maBcUSsncXYzQT13D5\nNr+Z3mLxMMh3TUdt6sACmqbJ0btRlgXfMtNLM2OU1I6a3Ju+tIZSdn2v21JBwvxU\nzpZQ4zy2cimIiMQDZCQHJwzC9GZn8HaW091iz9H0Go3A7WDXwYNmsdLNRi00o14U\njoaVqaPsYrZWvRKaIRqaU0hHmS0AWwQSvN/93iMIXuyiwywmkwKbWnnxCQ/gsctK\nFUtcNrwEx9Wgj6KlhwDTyI1QWSBbxVYNyUgPFzKxrSmwMO0yNff7ho+QT9x5+Y/7\nXE59S4Mc4ZXxcXKew/gSlN9U5mvT+D2BhDtkCupdfsZNCQWp27A+b/DmrFI9NqsC\nAwEAAaOCAcIwggG+MBIGA1UdEwEB/wQIMAYBAf8CAQAwQwYDVR0eBDwwOqE4MAaC\nBC5taWwwCocIAAAAAAAAAAAwIocgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\nAAAAAAAwDgYDVR0PAQH/BAQDAgGGMH8GCCsGAQUFBwEBBHMwcTAyBggrBgEFBQcw\nAYYmaHR0cDovL2lzcmcudHJ1c3RpZC5vY3NwLmlkZW50cnVzdC5jb20wOwYIKwYB\nBQUHMAKGL2h0dHA6Ly9hcHBzLmlkZW50cnVzdC5jb20vcm9vdHMvZHN0cm9vdGNh\neDMucDdjMB8GA1UdIwQYMBaAFOmkP+6epeby1dd5YDyTpi4kjpeqMFQGA1UdIARN\nMEswCAYGZ4EMAQIBMD8GCysGAQQBgt8TAQEBMDAwLgYIKwYBBQUHAgEWImh0dHA6\nLy9jcHMucm9vdC14MS5sZXRzZW5jcnlwdC5vcmcwPAYDVR0fBDUwMzAxoC+gLYYr\naHR0cDovL2NybC5pZGVudHJ1c3QuY29tL0RTVFJPT1RDQVgzQ1JMLmNybDAdBgNV\nHQ4EFgQU+3hPEvlgFYMsnxd/NBmzLjbqQYkwDQYJKoZIhvcNAQELBQADggEBAA0Y\nAeLXOklx4hhCikUUl+BdnFfn1g0W5AiQLVNIOL6PnqXu0wjnhNyhqdwnfhYMnoy4\nidRh4lB6pz8Gf9pnlLd/DnWSV3gS+/I/mAl1dCkKby6H2V790e6IHmIK2KYm3jm+\nU++FIdGpBdsQTSdmiX/rAyuxMDM0adMkNBwTfQmZQCz6nGHw1QcSPZMvZpsC8Skv\nekzxsjF1otOrMUPNPQvtTWrVx8GlR2qfx/4xbQa1v2frNvFBCmO59goz+jnWvfTt\nj2NjwDZ7vlMBsPm16dbKYC840uvRoZjxqsdc3ChCZjqimFqlNG/xoPA8+dTicZzC\nXE9ijPIcvW6y1aa3bGw=\n-----END CERTIFICATE-----"
#echo $cert
#If today is the deadline day it deletes cronjob, else renews certificate if validity is 1/2 of its original
if date "+%d %m %y" | grep -q "$getDeadLineD $getDeadLineM $getDeadLineY"; then
       echo "Lifetime has expired"
       sudo crontab -u root -l | grep -v "$domainName $getDeadLineD $getDeadLineM $getDeadLineY $uri" | sudo crontab -u root -
       sudo rm "./starCerts/$uri/certificate.pem"
else

	endDateStuffed=$(openssl x509 -in "./starCerts/$uri/certificate.pem" -noout -enddate)

	#$a contains date in format Oct 31 21:27:59 2018 GMT
	a=$(echo $endDateStuffed | cut -d'=' -f 2)
	lengthA=${#a}
	#$b containt date in format Oct 31 21:27:59 2018 GMT
	b=${a% GMT}

	date=$(date)
	#date containts date in format Mon Nov 20 20:13:37 CET 2017
	lengthDate=${#date}
	c=$(echo $date | cut -c4-$lengthDate)

	#$d returns date in the same format: Nov 20 20:24:20 2017
	d=$(echo $c | sed 's/CET//g')

	#Does the operation using seconds to get the number of hours between
	#the cert expiration date and current date.
	date1=$(date -d "$b" +'%s')
	date2=$(date -d "$d" +'%s')
	currentValidity=$((($date1 - $date2) /3600))


	#Validity is given in days
        validityHours=$(cat "./starCerts/$uri/validity")
        validityDays=$(($validityHours / 24))
	halfValidity=$(($validityHours /2))
	echo "Halfvalidity is $halfValidity and currentValidity is $currentValidity and validityDays is $validityDays"
	validityDays=11
	#if [ $halfValidity -gt $currentValidity ]
	#then
        echo "Renewing STAR certificate"
        #sudo rm "./starCerts/$uri/certificate.pem"
        #sudo openssl x509 -req -extensions v3_req -extfile /etc/ssl/openssl.cnf -in ./starCerts/$uri/csr -CAkey ./test/test-ca.key -CA ./test/test-ca.pem -days 365 -set_serial "0x$(openssl rand -hex 18)" -out "./starCerts/$uri/certificate2.pem"

	#Creates serial
	serial=$(openssl rand -hex 18)
	#start and enddate in format YYMMDDHHMMSSZ
	startdate=$(date +%y%m%d%H%M%SZ)
	enddate=$(date -d "+$validityHours hours" +%y%m%d%H%C%SZ)
	echo $serial | sudo tee './demoCA/serial'
	yes | sudo openssl ca -policy policy_anything -out "./starCerts/$uri/certificate2.pem" -startdate $startdate -enddate $enddate -cert "./test/test-ca.pem" -keyfile "./test/test-ca.key" -infiles "./starCerts/$uri/csr"

	#Save as PEM
	sudo rm "./starCerts/$uri/certificate.pem"
	sudo openssl x509 -in "./starCerts/$uri/certificate2.pem" -out "./starCerts/$uri/certificate.pem"
	sudo rm "./starCerts/$uri/certificate2.pem"
	echo $cert | sudo tee -a "./starCerts/$uri/certificate.pem"
	#fi
fi
