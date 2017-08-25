#!/bin/sh
domainName=$1


cd /root/certbot/certbot

        a=" --csr ../../tmpCsr --agree-tos -m tutatis@gmail.com --renew-by-default -d $domainName"
        b=" --server http://172.17.0.4:4000/directory -q --webroot -w /var/www/$domainName/html"

#Executes Certbot
python main.py certonly $a $b
