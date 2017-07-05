#!/bin/sh

domainName=$1
lifeTime=$2
cronTaskID=$3
day=$(date +%d) #returns current day
deadLineD=$(date -d "today + $lifeTime days" +'%d') #returns deadline in day month format, e.g. 19 06 as 2 parameters
deadLineM=$(date -d "today + $lifeTime days" +'%m')
touch myCron
crontab -l > myCron
echo "50 8 $day * * sh /root/exeAutoRenew.sh  $domainName $deadLineD $deadLineM $cronTaskID" >> myCron
#echo "0 0 $day * * sh /root/exeAutoRenew.sh  $domainName $deadLineD $deadLineM $cronTaskID" >> myCron
crontab myCron
rm myCron