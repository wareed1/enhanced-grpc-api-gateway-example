package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/wareed1/enhanced-grpc-api-gateway-example/gateway"
	"github.com/wareed1/enhanced-grpc-api-gateway-example/insecure"
	pbo "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/api/orders/v1"
	pbu "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/api/users/v1"
	ordersSvcV1 "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/orders/v1"
	usersSvcV1 "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	listenAddress = "0.0.0.0:9090"
	ordersSvc = "orders:9090"
	usersSvc = "users:9090"
)

func newOrdersSvcClient() (ordersSvcV1.OrdersServiceClient, error) {
	conn, err := grpc.DialContext(context.TODO(), ordersSvc, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("orders client: %w", err)
	}

	return ordersSvcV1.NewOrdersServiceClient(conn), nil
}

func newUsersSvcClient() (usersSvcV1.UsersServiceClient, error) {
	conn, err := grpc.DialContext(context.TODO(), usersSvc, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("users client: %w", err)
	}

	return usersSvcV1.NewUsersServiceClient(conn), nil
}

func logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("method %q called\n", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("method %q failed: %s\n", info.FullMethod, err)
	}
	return resp, err
}

func main() {
	log.Printf("Enhanced APIGW service starting on %s", listenAddress)

	// connect to orders svc
	ordersClient, err := newOrdersSvcClient()
	if err != nil {
		panic(err)
	}

	// connect to users svc
	usersClient, err := newUsersSvcClient()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//s := grpc.NewServer(grpc.UnaryInterceptor(logger))
	s := grpc.NewServer(
		// TODO: Replace with your own certificate!
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)

	pbo.RegisterOrdersServiceServer(s, NewOrdersService(usersClient, ordersClient))
	pbu.RegisterUsersServiceServer(s, NewUsersService(usersClient))

   	// Serve gRPC Server
	log.Printf("Serving gRPC on https://%s", listenAddress)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run("dns:///" + listenAddress)

}