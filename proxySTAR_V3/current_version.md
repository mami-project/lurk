Note: This is a high-level guide that shows STAR's current version, for installation go to the README
## 1 Short-term certificate issuance 

**_Steps are in these order to ease comprenhension but some of them are executing at the same time_**
### 1.1 Client uses curl to ask for a certificate to the proxy.

`curl --cacert $proxyCert -H Content-Type: application/json -X POST -d $fullCsr https://certProxy:443/star/registration`

@proxyCert: certificate signed by 'certProxy', that is the name assigned in /etc/hosts to proxy's IP.   
@fullCsr contains: {"csr":".........", "lifetime":365, "validity":24}

### 1.2 DNO proxy's listens to the request

DNO is listening at `:443 /star/registration`, when a json is received it is decoded, then checks whether requested lifetime and
validity are OK, then it parses the csr and extracts the Domain requested. It compares it with a list of valid domains to delegate.
If these checks are successfull it sends a message with the final parameters used by the DNO(they may have changed), then another 
message with a header: StatusCreated, and a body: Location: .............
Location contains the *completion url*.

Then: tweaked certbot is executed, a coroutine is created to serve requests to *completion url* and client's request information is storaged.

### 1.3.1 Certbot is executed creating a default account using a fake e-mail and client's csr 

Certbot runs as usual but with extra fields: 
```
{
    "recurrent": true,
    "recurrent-start-date": "2016-01-01T00:00:00Z",
    "recurrent-end-date": "2017-01-01T00:00:00Z",
    "recurrent-certificate-validity": 604800
  }
 
 ```
 
 These extra fields are added to the new-certificate request that is sent to Boulder.
 
 ### 1.3.2 Client receives *completion url* from step 1.2
 
 Client extracts the url from `Location: url`  and does a new curl to this *completion url*. It returns a json struct containing:
```
{
    "status": "valid", 
    "lifetime": 365, // lifetime of the registration in days,
                     
    "certificate": "https://CertificateAuthoritySTAR:9898/09096ffd-5429-4a50-a80c-4fbe45a482b5"
}

```

Last field contains the URI(Located in the server) from where certificate is ready to download and also where it will be renewed.

A GET to that URI will return the certificate and show it in screen plus storage it in the directory passed as argument with the desired
name. If a certificate with the same properties exists in that file it will replace it, as that is the nature of STAR.

### 1.4 Boulder receives the certificate request with the extra fields 

When Boulder's Web Front End *WFE* receives a new valid request with `recurrent:true` it executes STAR.
Relevant information is storaged. These info includes the *csr*, *validity* and the *renewal uri*(which has just been randomly assigned
by the *WFE* using UUID v4).

Boulder continues executing but using the validity extracted from *recurrent-star-date* and *recurrent-end-date* as the certificate's
lifetime.
These certificate is storaged together with the relevant info from the request and send back to the DNO.
At the same time the cert is made available at renewal uri. 

### 1.5 DNO receives cert and GETs the renewal uri 

DNO receives the certificate but the renewal uri remains unknown to him.
It get's that first certificate UUID and makes a GET to server's `port :9898` that will return him the real renewal-uri.
This uri is send back to the Client in step 3.2 and also storaged in the proxy.

## 2 Renewal

### 2.1 Server creates a new STAR cert

When star cert is created, a coroutine starts serving a *cert file* at the *renewal uri*, this serving stays still until the ends of
times(TO DO stop serving at *renewal uri* after a while when termination is requested?).

When a petition is received it serves the *cert file* that contains the cert chain in pem format. It also sends 2 new headers: 
`Not-After` and `Not-Before` on top of the default ones:

```
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 3248
Content-Type: text/plain; charset=utf-8
Last-Modified: Thu, 23 Nov 2017 18:30:09 GMT
Not-After: Nov 23 17:30:00 2017 GMT
Not-Before: Nov 24 16:30:00 2017 GMT
Date: Thu, 23 Nov 2017 18:30:14 GMT

```

Because the coroutine always serves this *cert file*, renewal will consist on updating this file.

### 2.2 RenewalManager 

Boulder aside, renewalManager is the process that really does the renewals. 
When a new STAR cert is created renewalManager adds a cronjob. This cronjob executes every 2h. When executing it checks if STAR lifetime 
has expired, if it has it kills itself and deletes the *cert file*. Else, checks if current cert's validity is less than half of the
validity and renews it if so. 

Renewal is made using a CA that replicates ACME using the csr and validity stored. Then the full chain replaces *cert file*.

## 3. Termination 

### 3.1 DNO wants to terminate 'x' renewal.
DNO keeps every tuple {csr, renewal-uri, validity} it has issued so it has the *renewal-uri* of every STAR cert.

`go run termination.go 09096ffd-5429-4a50-a80c-4fbe45a482b5`

This command stops renewals for *renewal-uri* : 0909...
In order to do so, DNO requests a nonce to the Server and sends a JWT when it receives it.

### 2.RenewalManager receives termination request

Renewal manager deletes the cronjob + *cert file* but still serves at *renewal-uri* as it was pointed in **2.Renewal in TO DO?**. New 
client requests to *renewal-uri* will return `Order status: canceled`. 
RenewalManager is listening at `:9200 /getNonce` and `:9200 /terminate` and uses its own cert for TLS.
