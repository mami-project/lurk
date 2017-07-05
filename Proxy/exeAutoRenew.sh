#!/bin/bash

if [ "$#" -ne 4 ]; then
    echo "Illegal number of parameters"
    exit 1
fi
domainName=$1 # e.g. bye.com
getDeadLineD=$2 # "day month" format, e.g. 19 06
getDeadLineM=$3 # "day month" format, e.g. 19 06
cronTaskID=$4 # integer
if date "+%d %m" | grep -q "$getDeadLineD $getDeadLineM"; then
        echo "Lifetime expires"
        crontab -u root -l | grep -v "$domainName $getDeadLineD $getDeadLineM $crontaskID" | crontab -u root -
else
        echo "Renews cert"
        sh /root/certbot/certbot python main.py --webroot certonly -n -d $domainName
fi