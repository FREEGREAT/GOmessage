package service

import (
	"context"

	"gomessage.com/users/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.UserModel) error
	GetAllUsers(ctx context.Context) ([]models.UserModel, error)
	UpdateUser(ctx context.Context, user models.UserModel) error
	DeleteUser(ctx context.Context, id string) error
}
