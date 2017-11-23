#!/bin/sh
if [ "$#" -ne 3 ]; then
    echo "Illegal number of parameters"
    exit 1
fi
domainName=$1
lifeTime=$2
cronTaskID=$3
day=$(date +%d) #returns current day
deadLineD=$(date -d "today + $lifeTime days" +'%d') #returns deadline in day month year format, e.g. 19 06 17 as 3 parameters
deadLineM=$(date -d "today + $lifeTime days" +'%m')
deadLineY=$(date -d "today + $lifeTime days" +'%y')
validityHours=$(cat "../boulder/starCerts/$uri/validity")

sudo touch myCron
sudo chmod +wrx myCron
sudo crontab -l | sudo tee myCron
#Decide best time of the day for the renewal following these:
echo "0 */2 * * * sh /home/acme-server2/gopath/src/github.com/letsencrypt/boulder/exeAutoRenew.sh  $domainName $deadLineD $deadLineM $deadLineY $cronTaskID" | sudo tee -a myCron
sudo crontab myCron
sudo rm myCron
