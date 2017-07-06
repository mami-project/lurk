


Boulder/STAR Server             DNO/ STAR Proxy              Star Client




Time 0:
	Boulder is running
	Proxy is running 
	Communication between Proxy and Client is safe
	Proxy hosts a web

Time 1: 

	Client: curl --cacert /usr/share/ca-certificates/mozilla/server.crt -H \
			"Content-Type: application/json" -X POST -d @fullCSR https://cert\
			Proxy:443/star/registration

	Proxy: example.go
			function parseJsonPOST handles :443 /star/registration requests
			
			parseJsonPOST translates the block of data into a struct {csr,lifetime, validity}
			if lifetime and validity are appropiate it then proceeds to call decodeCsr.

			func decodeCsr calls helper script getCsrAsText.sh and finally returns the csr field 
			as plain text, keeping a copy in file tmpCsr

			Back in parseJsonPOST it uses the csr as plain text as a parameter for func parseFieldsOfCsr.
			This function shall retrieve the domain contained in the csr field of the already mentioned struct 
			and returns it as a string called subjectName.

			Writes the validity in a file called STARValidityCertbot and calls func addToCron.

			func addToCron adds a new cronjob by using the helper script addTask.sh. After lifetime expires,
			the script will remove itself from the cron table.

			parseJsonPOST executes function callCertbot
			
			
Time 2:


	Proxy: example.go

			function callCertbot runs cerbot application by passing the csr and the domain names as arguments.
			It uses cerbotCall.sh

	Proxy: certbot
			Normal execution of certbot with csr as a parameter but with 2 extra fields sent in the POST
			to Boulder added in acme/acme/client.py and messages.py. These 2 fields are recurrent and 
			recurrent_cert_validity.
			When certbot is called it checkes if file STARValidityCertbot exists, if so, it then reads the validity 
			value contained in it and sends it as recurrent_cert_validity in the POST as well as setting the recurrent 
			field to 'true'. 

	Boulder: wfe.go
			If Boulder function verifyPOST reads field recurrent: true in the Json then it reads the field recurrent_cert_validity
			and saves it into a temporary file called STARValidityWFE. After this operation server operations continue as normal, forwarding
			the csr until it reaches the CA.

	Boulder: ca.go

			When the time to create the certificate comes, function issueCertificate looks for the STARValidityWFE file and it finds it
			it sets the validity of the next certificate to the duration written in the file. If it doesn't find it, duration is set to 
			3 months.
			

Time 3:

	Proxy: certbot

			Proxy's certbot side receives confirmation that the certificate challenge and certificate issuance has succeeded and writes 
			the URI in a file called certId.
			Proxy will then read this file and POST its contents into /completionURL. If client tries to reach /completionURL before the
			certificate ends it returns "in progress".






