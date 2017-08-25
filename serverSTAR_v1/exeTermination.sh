#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Illegal number of parameters"
    exit 1
fi
uuid=$1 # uuid4 of cert


#Terminates renewal for the cert with the uuid
       sudo crontab -u root -l | grep -v "$uuid" | sudo crontab -u root -

