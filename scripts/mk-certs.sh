#! /bin/bash

if [[ "$HOSTNAME" == "" ]]; then
   echo "missing hostname, aborting ..."
   exit 1
fi

if [[ "$IPADDR" == "" ]]; then
   echo "missing IP address, aborting ..."
   exit 2
fi

if [[ "$OSTYPE" == "msys" ]]; then
   # Lightweight shell and GNU utilities compiled for Windows (part of MinGW)
   # needed to run in Windows Git Bash 
   export MSYS_NO_PATHCONV=1
   HOSTNAME=$HOSTNAME
fi

# generate the CA private key
openssl genrsa -out myCA.key 2048

# generate the CA certificate
openssl req -x509 -new -nodes -key myCA.key -sha256 -days 1825 -out myCA.pem -subj "/C=CA/ST=ON/L=Ottawa/O=Acme Corporation/OU=IT Department/CN=Acme Corporation CA Root"

# generate the server private key
openssl genrsa -out server.key 2048

# create CSR
openssl req -new -key server.key -out server.csr -subj "/C=CA/ST=ON/L=Ottawa/O=Acme Corporation/OU=IT Department/CN=$IPADDR"

# create a v3 ext file for SAN properties
cat > server_cert.v3.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = $HOSTNAME
DNS.2 = localhost
IP.1 = $IPADDR
IP.2 = 0.0.0.0
IP.3 = 127.0.0.1
EOF

# generate server certficate
openssl x509 -req -days 365 -in server.csr -CA myCA.pem -CAkey myCA.key -CAcreateserial -out server-cert.pem -extfile server_cert.v3.ext

if [[ "$OSTYPE" == "msys" ]]; then
   # remove any existing certificates
   certutil -delstore "Root" "Acme Corporation CA Root"
   # add new one
   certutil -addstore "Root" myCA.pem
fi
