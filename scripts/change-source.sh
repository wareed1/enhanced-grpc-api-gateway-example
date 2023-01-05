#! /bin/bash

API_GW_SOURCE=api-gw/main.go
API_GATEWAY_SOURCE=gateway/gateway.go
API_CLIENT_SOURCE=apigw-client/main.go

if [[ -z "$1" ]]; then
   echo "missing directive, aborting ..."
   exit 1
fi

if [[ "$1" != "secure" && "$1" != "insecure" ]]; then
   echo "invalid directive ($1), must be 'secure' or 'insecure', aborting ..."
fi

if [[ "$1" == "secure" ]]; then
	sed -i 's/\/insecure/\/secure/' "$API_GW_SOURCE"
	sed -i 's/\&insecure\./\&secure\./' "$API_GW_SOURCE"
	sed -i 's/\/insecure/\/secure/' "$API_GATEWAY_SOURCE"
	sed -i 's/insecure\.Cert/secure\.Cert/' "$API_GATEWAY_SOURCE"
    sed -i -e '/\/insecure/ s/\"github.com/\/\/\"github.com/' "$API_CLIENT_SOURCE"
    sed -i '/^func main/a\\ 	creds, err := credentials.NewClientTLSFromFile("../certs/myCA.pem", "")\n	if err != nil {\n		panic(err)\n	}\n' "$API_CLIENT_SOURCE"
    sed -i 's/grpc.WithTransportCredentials/\/\/grpc.WithTransportCredentials/' "$API_CLIENT_SOURCE"
    sed -i '/\/\/grpc.WithTransportCredentials/a\\  	grpc.WithTransportCredentials(creds))' "$API_CLIENT_SOURCE"
else
	sed -i 's/\/secure/\/insecure/' "$API_GW_SOURCE"
	sed -i 's/\&secure\./\&insecure\./' "$API_GW_SOURCE"
	sed -i 's/\/secure/\/insecure/' "$API_GATEWAY_SOURCE"
	sed -i 's/secure\.Cert/insecure\.Cert/' "$API_GATEWAY_SOURCE"
    sed -i 's/\/\/\"github.com/\"github.com/' "$API_CLIENT_SOURCE"
    sed -i -e '/creds, err := credentials.NewClientTLSFromFile/,+4d' "$API_CLIENT_SOURCE"
    sed -i 's/\/\/grpc.WithTransportCredentials/grpc.WithTransportCredentials/' "$API_CLIENT_SOURCE"
    sed -i -e '/grpc.WithTransportCredentials(creds)/,1d' "$API_CLIENT_SOURCE"
fi