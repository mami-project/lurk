#!/bin/sh
domainName=$1
#crontaskID=$2
#operationCode=$3 #1 for first certificate, 2 for renewal

cd /root/certbot/certbot

        a=" --csr ../../tmpCsr --agree-tos -m tutatis@gmail.com --renew-by-default -d $domainName"
        b=" --server http://172.17.0.4:4000/directory -q --webroot -w /var/www/$domainName/html"

#Executes Certbot
python main.py certonly $a $b
#echo "A y B valen $a $b"
