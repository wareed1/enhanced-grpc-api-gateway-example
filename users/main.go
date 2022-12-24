package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pbu "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/users/v1"
)

const (
	listenAddress = "0.0.0.0:10000"
)

type usersService struct {

}

func (u *usersService) CreateUser(ctx context.Context, request *pbu.CreateUserRequest) (*pbu.CreateUserResponse, error) {
	return &pbu.CreateUserResponse{User: &pbu.User{
		Id:    "1",
		Email: "user1@email.com",
	}}, nil
}

func (u *usersService) ListUsers(ctx context.Context, request *pbu.ListUsersRequest) (*pbu.ListUsersResponse, error) {
	return &pbu.ListUsersResponse{
		Users: []*pbu.User{
			{
				Id:    "1",
				Email: "user1@email.com",
			},
			{
				Id:    "2",
				Email: "user2@email.com",
			},
		},
	}, nil
}

func main() {
	log.Printf("Users service starting on %s", listenAddress)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pbu.RegisterUsersServiceServer(s, &usersService{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}