package storage

import (
	"context"

	"gomessage.com/users/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.UserModel) (string, error)
	FindAll(ctx context.Context) (u []models.UserModel, err error)
	FindOne(ctx context.Context, user *models.UserModel) (models.UserModel, error)
	FindById(ctx context.Context, id string) (models.UserModel, error)
	Update(ctx context.Context, user *models.UserModel) error
	Delete(ctx context.Context, id string) (string, error)
}

type FriendsRepository interface {
	Create(ctx context.Context, friend *models.FriendListModel) (string, error)
	FindAll(ctx context.Context, user_id string) ([]models.FriendListModel, error)
	Delete(ctx context.Context, friends *models.FriendListModel) error
}
