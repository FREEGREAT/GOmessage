package service

import (
	"context"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
)

type UserService interface {
	RegisterUser(ctx context.Context, in *proto_user_service.RegisterUserRequest) (*proto_user_service.RegisterUserResponse, error)
	LoginUser(ctx context.Context, in *proto_user_service.LoginUserRequest) (*proto_user_service.LoginUserResponse, error)
	UpdateUser(ctx context.Context, in *proto_user_service.UpdateUserRequest) (*proto_user_service.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *proto_user_service.DeleteUserRequest) (*proto_user_service.DeleteUserResponse, error)
	GetUsers(context.Context, *proto_user_service.GetUsersRequest) (*proto_user_service.GetUsersResponse, error)
	ListOfFriends(context.Context, *proto_user_service.ListOfFriendsRequest) (*proto_user_service.ListOfFriendsResponse, error)
	AddFriends(ctx context.Context, usr *proto_user_service.AddFriendsRequest) (*proto_user_service.AddFriendsResponse, error)
	DeleteFriend(context.Context, *proto_user_service.DeleteFriendsRequest) (*proto_user_service.DeleteFriendsResponse, error)
}
