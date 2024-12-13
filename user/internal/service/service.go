package service

import (
	"context"

	"gomessage.com/users/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.UserModel) error
	UpdateUser(ctx context.Context, user models.UserModel) error
	DeleteUser(ctx context.Context, id string) (string, error)
	ListOfUsers(ctx context.Context) ([]models.UserModel, error)
	GetUser(ctx context.Context, id string) (models.UserModel, error)
}

type FriendsService interface {
	AddFriend(ctx context.Context, friend *models.FriendListModel) error
	DeleteFriend(ctx context.Context, friends *models.FriendListModel) error
}
