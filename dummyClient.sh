#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Illegal number of parameters"
    echo 'USAGE: command $proxyCert $fullCsr'
    exit 1
fi

proxyCert=$1
fullCsr=$2

#step 1
#returns the URI where the info about the cert is available
step1=$(curl --cacert $proxyCert -H Content-Type: application/json -X POST -d $fullCsr https://certProxy:443/star/registration)
echo "returning value of step1: $step1"
var=$(echo "$step1" | sed -n 2p | cut -d ":" -f 2,3)


echo "URI is: $var End of uri."

#step 2
#returns status, lifetime and certificate's final URI

step2=$(curl --cacert /usr/share/ca-certificates/mozilla/server.crt $var)
echo "Step2 is: $step2"
#step 3
#returns the rea

step3=$(echo "$step2" | cut -d ' ' -f 3 | cut -d "}" -f 1)
echo "Step 3 is: $step3"

sleep 5;curl --cacert ./serverKey/cert.pem $step3

echo "end of Client"
