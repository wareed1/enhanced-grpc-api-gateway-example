package main

import (
	"context"
	"fmt"

	"github.com/wareed1/enhanced-grpc-api-gateway-example/insecure"
	"github.com/golang/protobuf/jsonpb"
	apiClient "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/public/api/orders/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const apiSvc = "localhost:9090"

func main() {
    conn, err := grpc.DialContext(context.TODO(), apiSvc, 
        grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")))
	if err != nil {
		panic(err)
	}

	api := apiClient.NewOrdersServiceClient(conn)

	res, err := api.ListOrdersWithUser(context.Background(), &apiClient.ListOrdersWithUserRequest{})
	if err != nil {
		panic(err)
	}

	resp, err := (&jsonpb.Marshaler{}).MarshalToString(res)
	if err != nil {
		panic(err)
	}

	fmt.Printf(resp)
}
