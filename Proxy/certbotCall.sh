#!/bin/sh
domainName=$1
cd /root/certbot/certbot
a=" --csr ../../tmpCsr --agree-tos --renew-by-default -d $domainName"
b=" --server http://172.17.0.4:4000/directory --webroot -w /var/www/$domainName/html"
c=" --cert-path /etc/letsencrypt/live/$domainName/cert.pem --chain-path /etc/letsencrypt/live/$domainName/chain.pem --fullchain-path /etc/letsencrypt/live/$domainName/fullchain.pem "

#./certbot-auto certonly $a $b

python main.py certonly $a $b $c
