#!/bin/bash

if [ "$#" -ne 4 ]; then
    echo "Illegal number of parameters"
    exit 1
fi
domainName=$1 # e.g. bye.com
getDeadLineD=$2 # "day month" format, e.g. 19 06
getDeadLineM=$3 # "day month" format, e.g. 19 06
uri=$4 # string

cert="-----BEGIN CERTIFICATE-----\nMIIEijCCA3KgAwIBAgICEk0wDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAwwgY2Fj\na2xpbmcgY3J5cHRvZ3JhcGhlciBmYWtlIFJPT1QwHhcNMTUxMDIxMjAxMTUyWhcN\nMjAxMDE5MjAxMTUyWjAfMR0wGwYDVQQDExRoYXBweSBoYWNrZXIgZmFrZSBDQTCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMIKR3maBcUSsncXYzQT13D5\nNr+Z3mLxMMh3TUdt6sACmqbJ0btRlgXfMtNLM2OU1I6a3Ju+tIZSdn2v21JBwvxU\nzpZQ4zy2cimIiMQDZCQHJwzC9GZn8HaW091iz9H0Go3A7WDXwYNmsdLNRi00o14U\njoaVqaPsYrZWvRKaIRqaU0hHmS0AWwQSvN/93iMIXuyiwywmkwKbWnnxCQ/gsctK\nFUtcNrwEx9Wgj6KlhwDTyI1QWSBbxVYNyUgPFzKxrSmwMO0yNff7ho+QT9x5+Y/7\nXE59S4Mc4ZXxcXKew/gSlN9U5mvT+D2BhDtkCupdfsZNCQWp27A+b/DmrFI9NqsC\nAwEAAaOCAcIwggG+MBIGA1UdEwEB/wQIMAYBAf8CAQAwQwYDVR0eBDwwOqE4MAaC\nBC5taWwwCocIAAAAAAAAAAAwIocgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\nAAAAAAAwDgYDVR0PAQH/BAQDAgGGMH8GCCsGAQUFBwEBBHMwcTAyBggrBgEFBQcw\nAYYmaHR0cDovL2lzcmcudHJ1c3RpZC5vY3NwLmlkZW50cnVzdC5jb20wOwYIKwYB\nBQUHMAKGL2h0dHA6Ly9hcHBzLmlkZW50cnVzdC5jb20vcm9vdHMvZHN0cm9vdGNh\neDMucDdjMB8GA1UdIwQYMBaAFOmkP+6epeby1dd5YDyTpi4kjpeqMFQGA1UdIARN\nMEswCAYGZ4EMAQIBMD8GCysGAQQBgt8TAQEBMDAwLgYIKwYBBQUHAgEWImh0dHA6\nLy9jcHMucm9vdC14MS5sZXRzZW5jcnlwdC5vcmcwPAYDVR0fBDUwMzAxoC+gLYYr\naHR0cDovL2NybC5pZGVudHJ1c3QuY29tL0RTVFJPT1RDQVgzQ1JMLmNybDAdBgNV\nHQ4EFgQU+3hPEvlgFYMsnxd/NBmzLjbqQYkwDQYJKoZIhvcNAQELBQADggEBAA0Y\nAeLXOklx4hhCikUUl+BdnFfn1g0W5AiQLVNIOL6PnqXu0wjnhNyhqdwnfhYMnoy4\nidRh4lB6pz8Gf9pnlLd/DnWSV3gS+/I/mAl1dCkKby6H2V790e6IHmIK2KYm3jm+\nU++FIdGpBdsQTSdmiX/rAyuxMDM0adMkNBwTfQmZQCz6nGHw1QcSPZMvZpsC8Skv\nekzxsjF1otOrMUPNPQvtTWrVx8GlR2qfx/4xbQa1v2frNvFBCmO59goz+jnWvfTt\nj2NjwDZ7vlMBsPm16dbKYC840uvRoZjxqsdc3ChCZjqimFqlNG/xoPA8+dTicZzC\nXE9ijPIcvW6y1aa3bGw=\n-----END CERTIFICATE-----"

#Removes old certificate
if date "+%d %m" | grep -q "$getDeadLineD $getDeadLineM"; then
       echo "Lifetime expires"
       sudo crontab -u root -l | grep -v "$domainName $getDeadLineD $getDeadLineM $uri" | sudo crontab -u root -
else
        echo "Renews cert"
	sudo rm "./starCerts/$uri/certificate.pem"
	validityHours=$(cat "./starCerts/$uri/validity")
	validityDays=$(($validityHours / 24))
	sudo openssl x509 -req -extensions v3_req -extfile /etc/ssl/openssl.cnf -in ./starCerts/$uri/csr -CAkey ./test/test-ca.key -CA ./test/test-ca.pem -days $validityDays -set_serial "0x$(openssl rand -hex 18)" -out "./starCerts/$uri/certificate.pem"

echo $cert | sudo tee -a "./starCerts/$uri/certificate.pem"
fi

