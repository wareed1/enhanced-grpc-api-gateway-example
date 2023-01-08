package main

import (
	"context"

	pbu "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/public/api/users/v1"
	usersSvcV1 "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/private/users/v1"
    "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type usersService struct {
	usersClient usersSvcV1.UsersServiceClient
}

func NewUsersService(usersClient usersSvcV1.UsersServiceClient) *usersService {
	return &usersService{
		usersClient: usersClient,
	}
}

func (u *usersService) CreateUser(ctx context.Context, request *pbu.CreateUserRequest) (*pbu.CreateUserResponse, error) {
	res, err := u.usersClient.CreateUser(ctx, &usersSvcV1.CreateUserRequest{Email: request.GetEmail()})
	if err != nil {
		return nil, err
	}

	// as you can see in this case the messages are quite similar (for now) but we have to translate
	// them between API structs and internal structs
	return &pbu.CreateUserResponse{User: &pbu.User{
		Id:    res.GetUser().GetId(),
		Email: res.GetUser().GetEmail(),
	}}, nil
}

func (u *usersService) ListUsers(ctx context.Context, request *pbu.ListUsersRequest) (*pbu.ListUsersResponse, error) {
	// This can be done async in go routines to speed things up
	allUsers, err := u.usersClient.ListUsers(ctx, &usersSvcV1.ListUsersRequest{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// map users for easy mapping
	userIdToUser := map[string]*usersSvcV1.User{}
	for _, user := range allUsers.GetUsers() {
		userIdToUser[user.Id] = user
	}

    // build result
    result := &pbu.ListUsersResponse{
        Users: &pbu.Users{
                Users: make([]*pbu.User, len(allUsers.GetUsers())),
    },
    }
    result.Users.Summary = "my users"
	for idx, user := range allUsers.GetUsers() {
		result.Users.Users[idx] = &pbu.User{
			Id:        user.GetId(),
			Email:     user.GetEmail(),
		}
	}

	return result, nil
}