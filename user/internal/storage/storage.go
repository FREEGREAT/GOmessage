package storage

import (
	"context"

	"gomessage.com/users/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.UserModel) error
	FindAll(ctx context.Context) (u []models.UserModel, err error)
	FindOne(ctx context.Context, id string) (models.UserModel, error)
	Update(ctx context.Context, user models.UserModel) error
	Delete(ctx context.Context, id string) error
}

type Friends interface {
	AddFriends(ctx context.Context, user models.UserModel) (string, error)
	DeleteFriends(ctx context.Context, id string) error
}
