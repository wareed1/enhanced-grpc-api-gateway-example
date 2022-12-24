package main

import (
	"context"

	pbo "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/api/orders/v1"
	usersSvcV1 "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/users/v1"
	orderSvcV1 "github.com/wareed1/enhanced-grpc-api-gateway-example/gen/proto/go/smpl/orders/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ordersService struct {
	usersClient  usersSvcV1.UsersServiceClient
	ordersClient orderSvcV1.OrdersServiceClient
}

func NewOrdersService(usersClient usersSvcV1.UsersServiceClient, ordersClient orderSvcV1.OrdersServiceClient) *ordersService {
	return &ordersService{
		usersClient:  usersClient,
		ordersClient: ordersClient,
	}
}

func (o *ordersService) ListOrdersWithUser(ctx context.Context, request *pbo.ListOrdersWithUserRequest) (*pbo.ListOrdersWithUserResponse, error) {
	// This can be done async in go routines to speed things up
	allUsers, err := o.usersClient.ListUsers(ctx, &usersSvcV1.ListUsersRequest{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}


	allOrders, err := o.ordersClient.ListOrdersWithUser(ctx, &orderSvcV1.ListOrdersWithUserRequest{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// map users for easy mapping
	userIdToUser := map[string]*usersSvcV1.User{}
	for _, user := range allUsers.GetUsers() {
		userIdToUser[user.Id] = user
	}


	// build result
	result := &pbo.ListOrdersWithUserResponse{
		Orders: &pbo.Orders{
			Orders: make([]*pbo.Order, len(allOrders.GetOrders())),
		},
	}
	result.Orders.Summary = "my orders"
	for idx, order := range allOrders.GetOrders() {
		result.Orders.Orders[idx] = &pbo.Order{
			Id:        order.GetId(),
			UserId:    order.GetUserId(),
			UserEmail: userIdToUser[order.GetUserId()].GetEmail(),
			Product:   order.GetProduct(),
		}
	}
	return result, nil
}