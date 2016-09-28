# Testing TLS logging with DTR

Assuming you have UCP & DTR running, you need to first generate certs for TLS.

```
# Create a directory to later mount into rsyslog container
mkdir ~/mysyslog && cd ~/mysyslog
```

Copy this openssl.cnf into `~/mysyslog/openssl.cnf`
```
[ req ]
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
countryName = Country Name (2 letter code)
countryName_default = US
stateOrProvinceName = State or Province Name (full name)
stateOrProvinceName_default = MN
localityName = Locality Name (eg, city)
localityName_default = Minneapolis
organizationalUnitName    = Organizational Unit Name (eg, section)
organizationalUnitName_default    = Domain Control Validated
commonName = Internet Widgits Ltd
commonName_max    = 64


[ v3_req ]
# Extensions to add to a certificate request
basicConstraints = CA:TRUE
keyUsage = nonRepudiation, keyCertSign, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 172.17.0.1
IP.2 = 192.168.0.106
```

```
# Create CA cert
openssl req -new -x509 -nodes -out ca.pem -keyout ca-key.pem -subj /CN=DaRoot -newkey rsa:2048 -sha512
# Create Rsyslog cert
openssl req -new -nodes -out csr.pem -keyout key.pem -subj /CN=Zintermediate -newkey rsa:2048 -sha512
# Sign Rsyslog cert 
openssl x509 -req -in csr.pem -CAkey ca-key.pem -CA ca.pem -days 20 -set_serial 123 -sha512 -out cert.pem -extensions v3_req -extfile openssl.cnf
# Create DTR cert
openssl req -new -nodes -out dtr-csr.pem -keyout dtr-key.pem -newkey rsa:2048 -sha512
# Sign DTR cert
openssl x509 -req -in dtr-csr.pem -CAkey ca-key.pem -CA ca.pem -days 10 -extensions v3_req -extfile openssl.cnf -set_serial 1234 -sha512 -out dtr-cert.pem
```

Run the rsyslog container and sanity test it using an alpine image. You can either build your own image using the Dockerfile in this directory or use mine from Docker Hub.

```
docker run --rm -it -v /home/edgar/mysyslog:/etc/pki/rsyslog -p 515:514 --name rsyslog hinshun/rsyslog
docker run --rm --log-driver syslog --log-opt syslog-address=tcp+tls://172.17.0.1:515 --log-opt syslog-tls-ca-cert=/home/edgar/mysyslog/ca.pem --log-opt syslog-tls-cert=/home/edgar/mysyslog/dtr-cert.pem --log-opt syslog-tls-key=/home/edgar/mysyslog/dtr-key.pem alpine echo "hello world"
```

If you receive a "hello world" in the rsyslog container, your TLS connection is working properly. Reconfigure DTR to use the TLS flags to change logging to use SSL.

```
docker run -it --rm dockerhubenterprise/dtr-dev reconfigure --ucp-insecure-tls --log-protocol tcp+tls --log-host 172.17.0.1:515 --log-tls-ca-cert="$(cat ca.pem)" --log-tls-cert="$(cat dtr-cert.pem)" --log-tls-key="$(cat dtr-key.pem)"
```
