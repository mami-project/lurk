## Contents of this README

- Installation Guide --> How to set up the environment

- Common problems --> Recommended to read

- Simulation Guide --> Obtain a STAR certificate yourself

- Round-Trips Guide --> In depth explanation of every step in STAR




## Parts involved in STAR

1. Boulder/STAR Server

2. DNO/ STAR Proxy

3. CDN/ Star Client

# INSTALLATION GUIDE

Fully read each step before doing anything :)
These are the steps to get the whole simulation going:

1. Prepare 3 VM using ubuntu 14.04 trusty and name them: Server, Proxy and Client

Now let's prepare the Server: //this is the first VM


2. Install GO and set environment variable PATH to /usr/local/go/bin. My version (`$go version`)  is "go1.8.1 linux/amd64". Go's official 
documentation available at: https://golang.org/doc/install.

3. In your home directory create: ~/gopath/src/github.com/letsencrypt/boulder and place all the files there.(the files that are currently 
in https://github.com/mami-project/lurk/tree/master/serverSTAR_v2. NOTE: When you finish copying, everything must be inside 
letsencrypt/boulder/. Doing an `$ls` in letsencrypt/boulder/ must return the files that are currently under serverSTAR_v2).

Using `$git clone https://github.com/mami-project/lurk` is the fastest way to download all the server files, however, because clone command
downloads the full repo you will have to manually delete all the other files that are not used for server: client and proxy.


4. Fully install Docker and Docker-Compose: https://docs.docker.com/compose/install/ just follow the steps and test that the hello-world 
image works. I'm using version 17.03.1-ce for Docker and version 1.12.0 for Docker-Compose.

5. Go to your file "boulder/test/config/va.json" and make sure your port config in va is : 80 for httpPort, 5001 for httpsPort and 5001 for 
tlsPort.(Some of these changes may be redundant, but this way it works, so keep it that way ;) )
6. Now go to Docker's configuration file: "boulder/docker-compose.yml" and check that in the list of ports you have 80:80 and 443:443, 
be careful NOT TO TAB as it is an illegal char (these 2 steps are already done if you download my repo, do them if you are using common Boulder)

7. Back in "boulder/docker-compose.yml" change the FAKE_DNS field to the IP of the VM that will act as your **proxy**.

8. Set ufw status to inactive: `$sudo ufw disable`
9. Check your iptables policy: `$sudo iptables -L` and set CHAIN FORWARD policy to accept if it is currently to DROP mode:
`$sudo iptables -P FORWARD ACCEPT`
10. Make sure you can reach the other machines(client and proxy) with ICMP by pinging them: `$ping remoteVM`. 

11. Open /etc/hosts and add these 2:

	172.17.0.4      acme-v01.api.letsencrypt.org boulder //this is Boulder's local IP, at least in my machine

	/*If using `ifconfig` command in the server returns you an interface called docker0 172.17.0.1 once you have lauched boulder's 
 docker for the first time, then 172.17.0.4 should be your boulder's IP too.*/

	X.X.X.X  bye.com //this is the web hosted for the example. Because it is not available in the Web, the server must know
	where it is. Subsitute the "Xs" with your proxy's IP.

12.Get a selfsigned certificate and place it in ./boulder:
`$openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365`
In my case I named the CN field "CertificateAuthoritySTAR", **REMEMBER** this name, because you will need to add it in the /etc/hosts file of the proxy.
Remember to not cypher the key file. If you accidentally did, use:
`$openssl rsa -in key.pem -out key.pem` //to extract the key


**_Now it is time to prepare the Proxy:_** //this is the second VM

1.The current proxy has been tested as root so use: `$sudo -i` to become root and paste all the files in https://github.com/mami-project/lurk/tree/master/proxySTAR_v2 
so that the end directory of files such as proxySTAR.go and termination.go is "/root/proxySTAR.go" "/root/termination.go"
Again, using git clone is the fastest way.

2. In the new directory go to certbot/ and type: 
`#./certbot-auto --os-packages-only`
`#./tools/venv.sh`
`#source ./venv/bin/activate` //*You always need this ON so remember to execute this last command again if you exit the VM.*

3. Declare global environment:  `#export SERVER=http://172.17.0.4:4000/directory`    (this is where Boulder is listening)

4. Install Go, my version (`go version) is : go1.8.2 linux/amd64

5. Set PATH to /root/go/bin. E.g.: `#PATH=$PATH:/root/go/bin`

6. Make the same icmp and iptables checks that we did preparing the server: chain policy set to ACCEPT and pinging between VMs.
(steps 8 and 9)

7. Edit "/etc/hosts" and add : "172.17.0.4 CertificateAuthoritySTAR". The domain name is that you gave to the Server's self-signed 
certificate.

8. Host some website. This is the site for which we are gonna request the certificates. The way I did it is using Apache (e.g. bye.com):
	8.1. Place an html file in "/var/www/bye.com/html/bye.html". An example could be:
  <!DOCTYPE html PUBLIC "-//IETF//DTD HTML 2.0//EN">
	<HTML>
	 <HEAD>
			<TITLE>
				 A Small bye placed in 192.168.122.125 PROXY
			</TITLE>
	 </HEAD>
	<BODY>
	 <H1>Bye</H1>
	 <P>This is a very minimal "bye bye cruel world" HTML document.</P>
	</BODY>
	</HTML>
  

	8.2. Go to "/etc/apache2/sites-available" and copy a file called "000-default.conf" in the same directory as "bye.com.conf":
	`#cp -a 000-default.conf bye.com.conf`
	8.3. Open this new file and make sure the field VirtualHost in the first line is set to *:80 and the rest of the fields look like this
	(although ServerAdmin isn't important for now so leave it out if you want):

	ServerAdmin info@bye.com
        ServerName bye.com
        ServerAlias www.bye.com
        DocumentRoot /var/www/bye.com/html

    8.4. In sites-available run this commands:
    	"#a2ensite bye.com.conf"
    	"#a2dissite 000-default.conf"
    	"#service apache2 restart"

	DONE ;)

9. Now it comes the most difficult step: Preparing the proxy for the http challenge. Create 2 new directories so the end result is like this:  "/var/www/bye.com/html/.well-known/acme-challenge"

 **IMPORTANT**: When done, make sure that you can access the directory acme-challenge from another VM, so place a "hello.html" file there
and try to reach it with curl from the server: `$curl http://YOUR_PROXY'S_IP/.well-known/acme-challenge/hello.html`

For the file contents, just take the previous "bye.html" as an e.g. but rename it to hello.html or keep it as bye.html
(but then change the curl domain too)

If it works, feel free to delete it. If it doesn't, change the file permissions going to "/var/www" in your proxy and typing
`#chmod -R 755 bye.com` and change the user so it isn't root:
`#chown -R user:user bye.com` *<----IMPORTANT: "user:user" is your name and group, so if your user in the VM is Josh from the 
Goonies--> `#chown -R Josh:Goonies bye.com`*

10.For the simulation to work, you also need to generate a certificate using openSSL so that proxy and client can use https:
`#openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365`

This generates a certificate and a key that have to be kept in the same directory as proxySTAR.go.
I am including a cert and key in the Github, but they may be outdated if you are reading this.

11.Give execution permisions to all bash scripts executing in /root/: `#sudo +x chmod *.sh`
12.Make a new file called "starCerts" in /root/ and create a txt file called "myDomains.txt" with all the domains you want to be able
to delegate. This file shall contain a list of all the domains you may delegate. The format is: "domain + new line".
E.g. `#cat myDomains.txt` should return:

```
bye.com
heelo.com
imRunningOutOfIdeasForNewDomains.com
```
12.Last step. Create a file in "/root/" called "serverKey" and a text file inside named: "cert.pem":
`#mkdir serverKey`
`#cd serverKey`
`#nano cert.pem`

cert.pem is the selfsigned certificate you obtained in the last step of Boulder VM installation, so go to the server VM, copy it and paste its contents.

13.Now add the name that you gave to the CN in your certificate with its VM's IP to the "/etc/hosts" file, e.g.
"172.17.0.4 CertificateAuthoritySTAR"

14. Let's make it possible to access Boulder's docker:

`#route add -net 172.17.0.0 gw XXXXXX netmask 255.255.0.0 dev eth0` //XXXXX should be your server's name and eth0 is the default interface.


-Now test if you can access Boulder's docker from the proxy:
1. In the server, go to the boulder file and launch it:
```
$docker-compose build
$docker-compose up
```
2.When it ends launching (A message saying to press ^C pops at the end) go to the proxy and try to ping boulder:
`$ping 172.17.0.4` // If it fails this is the moment to check that is your docker's IP with an `$ifconfig`

If it works, stop Boulder:
	```$^C
  $docker-compose stop
	$docker-compose rm
	$y```

If ping fails try `$traceroute 172.17.0.4` to see where the message gets lost. Now go to the VM where the last jump was made
and check the route tables: `$route`
You must then add the routing for the 172.17.0.0 network and gateway 255.255.0.0. In my VM the command looks like:

`#route add -net 172.17.0.0 gw acme-server2 netmask 255.255.0.0 dev eth0`

-> 172.17.0.0 refers to the docker network, acme-server2 is the name of the VM where the dockers are running and it must be 
referenced in file /etc/hosts


The client VM will be prepared during the simulation process, it doesn't requiere any specific sofware.
## COMMON PROBLEMS

1.Proxy "fails" when you lauch the client. First stop the proxy. Now type:  `#rm -rf "/etc/letsencrypt"`. Try again. If it keeps 
failing it probably is a problem with the routing. To make sure check the logfiles in proxy's VM: `cat /var/log/letsencrpy/letsencrypt.log`,
in the last paragraph it must say something similar to "No route to host". Fix it adding the routes (like explained in the previos section).
Remember that the proxy needs to be able to ping/traceroute to 172.17.0.4 (boulder docker in the server VM). On the opposite side, the client 
needs to be able to connect to the proxy and to the server's IP, not to the docker! (to retrieve the certificate).

2.Certificate is not issued. Check that your proxy is able to solve the challenge. Place an html file in /var/www/bye.com/html/.well-known/acme-challenge 
and try to access it as explained before.

3.I did the installation, and tried common problems 1 and 2 but nothing works! 
Check that your iptable's forward policy is set to accept, and that you have **_PK certificates generated with openssl for proxy and server._** 
Also, the client needs to have proxy's and server's certificate. 
In the example, client keeps server's cert in "serverKey/cert.pem" and proxy's cert is in "/usr/share/ca-certificates/mozilla/server.crt". 
Server keeps his certificate and private key (decoded!) as cert.pem and key.pem in "boulder/"
Proxy keeps his certificate and private key (decoded!) as server.crt and server.key in "/root"
Proxy keeps server's certificate in "serverKey/cert.pem"




## SIMULATION GUIDE

HOW TO RUN A FULL SIMULATION

0. If you have restarted the server VM and want to ask for new certificates or mainly if some error happens but it is not fatal for the CA (if it is still UP), 
then go to the proxy and execute: `#rm -rf /etc/letsencrypt`
1. In the server go to `~/gopath/src/github.com/letsencrypt/renewal_full`
2. Run the renewalManager in background:
 	`$go run renewalManager.go $time.Milliseconds &`
	//to update the crontab every 5s run: `go run renewalManager 5000 &`
	//Uncomment the line that says "Message" in function checkStatus() if you want to get notified when the renewal does a check.
	...and these 3 commands to run Boulder:
	
  ```
	$docker-compose rm
	$docker-compose build
	$docker-compose up
   ```
	And wait until it says "All servers running. Hit ^C to kill", the first time *it may take a while.* 
  The first command is just to free memory from old dockers.

3. Go to proxy as superuser (`$sudo -i`) and type in "~/certbot"  
`#source ./venv/bin/activate`
4. Then: `export SERVER=http://172.17.0.4:4000/directory`
4.5 If you just followed the installation there's not need to do 3 & 4, you just did them.

5. Now you are ready to go with proxy's main code: 
`go run proxySTAR.go $maxLifeTime $maxValidity $pathToDomainList $pathToCAsCert $pathToOwnCert $pathToOwnKey`
//The first two variables set the maximum lifetime(days)and validity(hours).
//Also note that you need your self signed certs to be in the valid list for the other's VM
You will see a message: "Proxy STAR status in middlebox is: ACTIVE"

IMPORTANT: At this point you can jump to step 13 if you prefer to use the auto-client, go to next step if you prefer to do 
it manually. If you chose the first option just take a moment to see the CSR structure in step 6 before jumpling to step 13.
Also note that because the purpose of STAR is automate the process, the manual client may not be up to date right now.

6. Go to client's VM and POST at https://certProxy:443/star/registration (Don't forget to add certProxy to your /etc/hosts as the Proxy's IP).
	For now, the full command looks like this:
		`curl --cacert /usr/share/ca-certificates/mozilla/server.crt -H "Content-Type: application/json" -X POST -d \
		@fullCSR2 https://certProxy:443/star/registration`

		*server.crt is an openssl cert I generated so proxy and client can use https, it is the same certificate that we
		obtained for the proxy with openSSL in the installation guide.

		*fullCSR2 is a textfile that contains a basic CSR with the domain name bye.com as SN field
		(subject name) plus lifetime and validity in format:

		{"csr":"-----BEGIN CERTIFICATE REQUEST-----
		MIICmDCCAYACAQAwUzELMAkGA1UEBhMCTU0xCjAIBgNVBAgMAWsxDDAKBgNVBAcM
		A2xsbDEKMAgGA1UECgwBdDEMMAoGA1UECwwDdGlkMRAwDgYDVQQDDAdieWUuY29t
		MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmm74/AotkDGzsdVsn+Vu
		Z4FHW2+lf3HrLcDlpWHvBl3WSLg2TJHXdl8F6GtI3w91Cws/8wT4g+W33GYB0WAF
		WIGvzTPajeZ3jQt4t98bpzbuaFZz8QCoQVuEOuk8CCQ5/Cezbml3loMtXTuR+R1c
		OuVB9sFXbpoGvGL2fbAmrTtmOY9ZoXaLQmN7sj+4TjKRtZvVdpiLRaYp608ct2h3
		E6R2Nzm0OHdI35y61jaw46WiXCM30W/V2/Ia0c35Jdy4vbPybH1+k4ajmrlwiFrO
		986AlAxvxDZIKtahQFqMdH3hEuzTR6OnDwMlDtkLXThE9XSmcAhdYd9RLC8hF33A
		SQIDAQABoAAwDQYJKoZIhvcNAQELBQADggEBAF7ja2QCYDPfJ3kY0f4eSYaAaQba
		bQ6TA2dS5AFz+WSdzBQTTa8uTzgrKOwe8mQoHhjsNHW6aRpYCxje2v0pzTMw27g2
		YmXdfEfmWsF4GHk2NZ3ECE7LwA0YlsGZXpmYkUT89+69cJiiqiWUpwaGQdbx2Ozs
		N7tlHLlLDufQubnMetOfNb7SbyMpCdNssAaj7gkkmeOHk9rjlkrpkBJf8lBb2xIo
		XFH3iswaojVO3pAKZytDPrx9tsAsLstx6Jv6+O5lfr9rS4+EAT19yeZgd/64qTjl
		yx1r4nkEp7Z/brWh4X3q8zUhBQCLSeHIXp9nWj69WGXFtTOqcyx+uruc/Qw=
		-----END CERTIFICATE REQUEST-----","lifetime":365,"validity":21}

		To generate this text file first generate a private key and a CSR:
		`openssl req -out CSR.csr -new -newkey rsa:2048 -nodes -keyout privateKey.key`

		See its contents:
		"cat CSR.csr"

		And copy it into a new file following the format above. Dont forget to add lifetime and validity at the end.


7. After the client executes the command in step 6 it will get a message similar to this one:
"Location: https://certProxy/star/registration/0"
 Now if the client goes to this URI: `$curl --cacert /usr/share/ca-certificates/mozilla/server.crt https://certProxy:443/star/registration/0`
 It will get you a message back. This response can be {pending} or {status, lifetime,the uuid4 that serves as URI}
 E.g. {valid 365 20b1bac1-db72-42f4-9620-add03c789e36}
 
8. Then the client can retrieve the chained cert by using:
	`curl --cacert ./serverKey/cert.pem https://CertificateAuthoritySTAR:9898/20b1bac1-db72-42f4-9620-add03c789e36`

	*20b1bac1...is this example's certificate URI
	*cert.pem is the certificate that validates that the communications are safe with the renewalManager running in the CA. It is
	the same cert that was obtained in the last step of the Simulation Guide for the Server side. You must have a copy of it in the
	client.


9. To check that the renewalManager has done it's job -or if you forgot to run it, you can do so now and it will still work- put the Boulder in background:
```
"^Z"
"bg"
```
And check the crontab: `$sudo crontab -l`,

10. Renewals will be at the same URI that the first certificate, so DNO is not needed adnymore and can be turned OFF if you want so.

11.Note that a new directory has been created in the Server VM, this directory contains NEEDED information for the renewal so deleting it will cause renewals 
to fail.

12. To test the termination enter in the proxy and type `$go run termination.go $uuid`,
with the $uuid being the uuid where the certificate to terminate is renewed.

13. Auto-client: To get the certificate in one try just go to the clientVM and place there the file dummyClient.sh available at 
https://github.com/mami-project/lurk, this time just copy-paste the file.

14. Execute the client: 
`sh dummyClient.sh $proxyCert $caCert $CSR $file ` 
//Cert refers to the proxy's certificate and csr is the textfile containing csr,lifetime and validty as seen in step 6. 
//File refers to a destination file for the certificate, **if it does exist it will rewrite it**.
Don't forget to add "@" before the file's name. After executing it the cert should prompt. 

:tada::tada::tada::tada::tada::tada::tada::tada::tada::tada::tada::tada::tada::tada:

15. Renewal: To terminate an auto renewal from DNO: 
`#go run termination.go $uuid `
This $uuid is kept under the "starCerts/certificateNX/renewal_uri"


## ROUND-TRIPS GUIDE

Time 0:
	SERVER:
		Boulder is running
		RenewalManager in the CA is running
	Proxy:
		Proxy is running
		Proxy hosts a web using apache

	Communication between Proxy and Client is safe: Proxy has a cert issued by openssl and the client acknowledges the site as safe.
	This cert is kept as "/root/server.crt" together with its key "/root/server.key" in the DNO and at 
	"/usr/share/ca-certificates/mozilla/server.crt" for the client.
	READ BOTH INSTALLATION AND SIMULATION GUIDE IF YOU ARE ALREADY LOST :)

Time 1:

	Client:
			It calls the proxy with https by using the certificate with argument --cacert:
```
			curl --cacert /usr/share/ca-certificates/mozilla/server.crt -H \
			"Content-Type: application/json" -X POST -d @fullCSR https://cert \
			Proxy:443/star/registration
```
	Proxy: proxySTAR.go
			function parseJsonPOST handles :443 /star/registration requests

			parseJsonPOST translates the block of data into a struct {csr,lifetime, validity}
			if lifetime and validity are appropiate it then proceeds to call decodeCsr.



			func decodeCsr uses a helper script called getCsrAsText.sh and finally returns the csr field
			as plain text, keeping a copy in the file tmpCsr, that will be deleted afterwards.

			Back in parseJsonPOST it uses the csr (now a string) as a parameter for func parseFieldsOfCsr.
			This function shall retrieve the domain contained in the csr field of the already mentioned struct and returns
			it as a string called subjectName.

			Now the code runs parseDomain()

			parseDomain func compares the domain in csr with the domains available for renewal in starCerts/myDomains.txt
			and returns a boolean true if the requested domain is a valid one. If it isn't, it returns an error message to
			the STAR client and awaits for a new request.

			At this point in time, the proxy has info about: domain, lifetime and validity and has validated that all the
			fields are OK so it sends back the URL where the info about the certificate will be posted:

			Location: https://certProxy/star/registration/$ID

			As time goes on and multiple requests are made the $ID will increase.

			Function storeIssuedCerts is called. This will storage all the information
			processed by the proxy so that Certbot can read it when it gets executed.

			If these files dont exist Certbot client will ignore STAR.

			parseJsonPOST executes function callCertbot


Time 2:


	Proxy: proxySTAR.go

			function callCertbot runs cerbot application by passing the csr and the domain names as arguments.
			It uses cerbotCall.sh to do the execution.

	Proxy: Certbot
			Normal execution of the Certbot client with csr as a parameter but with 4 extra fields sent in the POST to
			Boulder added in "acme/acme/client.py" and "acme/acme/messages.py".
			These fields are:
			{recurrent,
			recurrent-start-date,
			recurrent-end-date,
			recurrent-certificate-validity}


			recurrent: contains parameter true and it serves to turn STAR ON
			recurrent-start-date and recurrent-end-date : validity for STAR certificates
			recurrent-certificate-validity: contains the lifetime, so it is key for renewal

			When certbot is called it checks if file the tmp files with the STAR information exist, if so, it then reads the
			these files, deletes them and sends them to Boulder in the same certificate request.

	Boulder: wfe.go
			If Boulder function verifyPOST reads field "recurrent: true" in the Json sent by Certbot then it starts STAR. It
			reads all the recurrent fields and saves them into temporary files and finally to a directory with these
			features:
				The main directory is "./starCerts" and will be created in the same directory where Boulder is. This
				means, in my case it is in:

				"~/gopath/src/github.com/letsencrypt/boulder/starCerts"

				In ./starCerts files with the cert uuid(This uuid has just been created by the wfe.go) as its name will 				be created, inside each file the info for the
				certificate renewal will be storaged:

				certificate.pem
				csr
				validity

				Why lifetime isn't in this file will be explained later.

			Before the csr is saved, STAR function repairCsr will make sure the csr is valid for future operations. The
			reason why this is a MUST is because at this point the csr has lost its base64 format and has lost the
			noninformation bits(stuffing).


			After this operation, server operations continue as normal, forwarding the csr until it reaches the CA.

	Boulder: ca.go

			When the time to create the certificate comes, function issueCertificate looks for the temporary files created
			in wfe.go and if it finds them, it issues a short-term certificate. If it doesn't find it, duration is set to 3
			months.

			WARNING: In case the config file is set to other value it won't read it.

			STAR function postAtUuid serves the file certificate.pem at:
			172.17.0.4:9898/$URI
			however, only the first certificate is at 'that' uri.

			STAR function addSTARToCron creates a file called "./renewMTmp/renew1" and it
			saves info for the renewals: {domain, lifetime, uuid}
			The reason why it saves information for the renewal instead of just adding it
			to cron right away is because ca.go is running inside a docker and the docker
			doesn't have a crontab so another process running outside of the docker needs to
			pick up the information inside this file "./renewMTmp/renew1" and add it to the
			system's crontab.

	Boulder: ra.go

			Function NewCertificate checks if STAR is ON (again thanks to the tmp files) and then storages the
			certificate in the file certificate.pem mentioned before, the key part is that it adds a chain of
			certificates by using the CA root certificate.

Time 3:

	Proxy: certbot

		Proxy's certbot side receives confirmation that the certificate challenge and certificate issuance have succeeded.

	Proxy: proxySTAR.go
		proxy then storages all the info about the certificate it just issued incluiding:
		certificate.pem: It is an isntace of the first certificate.
		csr: The csr the client used.
		uri: The uri where the STAR certificates are posted (unique for each certificate)
		validity: validity of every STAR
		
		Proxy will also read the cert's UUID and send do a POST to the cert using that UUID as URI, from there he will retrieve the URI
		where cert and renewals will be posted. This URI is immediately sent to the client.
		Keeping these files isn't NEEDED for STAR but because the Proxy is the real owner of the domains, it is considered that
		knowing what name delegations are active is key. What's more, keeping the uri will allow the DNO to terminate the
		renewals.

Time 4:


	Client:
		If the client executes a curl to reach /star/registration/$ID will obtain:
		{
			status: valid,
			lifetime: 365, //just an example
			certificate: 609577eb-f47b-4443-8df6-ba926dbdcd6c //URI where the certificate is
																												//available at 172.17.0.4:9898
		}

Time 5:


		CA: renewalManager.go
			The renewal manager is running in the CA background.
			It tries to read file "./renewMTmp/renew1" every 10 seconds (this time can be modified in function checkStatus)
			and adds a cronjob that will execute itself every 24h and will look like this:

			50 8 * * * sh /home/acme-server2/gopath/src/github.com/letsencrypt/boulder/exeAutoRenew.sh  bye7.com 04 08
			6d53bdbb-beed-41d5-ae12-fe71bfe71b880ea26

			The first 5 parameters ensure the 24h execution at 8:50 everyday.
			The next one is a script: exeAutoRenew.sh that executes with 4 parameters:
			{domain, day and month when the lifetime expires(NOTE: it is the STAR lifetime, not a cert's lifetime), uri}

		CA: exeAutoRenew.sh
			This script first checks if its day&month arguments match the current day, if they match, it deletes the cronjob
			that executes this script every 24h.
			If they don't match, it means certificate's lifetime hasn't expired and will renew the cert.
			The cert is renew using openssl with all the parameters that had been saved so a new certificate will be
			generated.
			The most important part of this new issued certificate is that it contains the same key Boulder uses and that is
			located in boulder/test/test-ca.key

