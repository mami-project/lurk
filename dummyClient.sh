#!/bin/bash

if [ "$#" -ne 4 ]; then
    echo "Illegal number of parameters"
    echo 'USAGE: command $proxyCert $caCert $fullCsr $file'
    echo '       File refers to the cert destination'
    exit 1
fi

proxyCert=$1
caCert=$2
fullCsr=$3
saveAt=$4

#step 1
#returns the URI where the info about the cert is available
step1=$(curl --cacert $proxyCert -H Content-Type: application/json -X POST -d $fullCsr https://certProxy:443/star/registration)
echo "returning value of step1: $step1"
var=$(echo "$step1" | sed -n 2p | cut -d ":" -f 2,3)


echo "URI is: $var End of uri."

#step 2
#returns status, lifetime and certificate's final URI

step2=$(curl --cacert $proxyCert $var)
echo "Step2 is: $step2"
#step 3
#returns final URI, then retrieves the cert

step3=$(echo "$step2" | cut -d ' ' -f 3 | cut -d "}" -f 1)
echo "Step 3 is: $step3"

sleep 5
curl --cacert $caCert $step3 | sudo tee -a $saveAt


echo "end of Client"
