package main

import (
	"context"
	"log"
	"net"

	pbo "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/orders/v1"
	"google.golang.org/grpc"
)

const (
	listenAddress = "0.0.0.0:10000"
)

type ordersService struct {}

func (o *ordersService) ListOrdersWithUser(ctx context.Context, request *pbo.ListOrdersWithUserRequest) (*pbo.ListOrdersWithUserResponse, error) {
	return &pbo.ListOrdersWithUserResponse{Orders: []*pbo.Order{
		{
			Id:      "o1",
			UserId:  "1",
			Product: "product-1",
		},
		{
			Id:      "o2",
			UserId:  "1",
			Product: "product-2",
		},
	}}, nil
}

func main() {
	log.Printf("Orders service starting on %s", listenAddress)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pbo.RegisterOrdersServiceServer(s, &ordersService{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}