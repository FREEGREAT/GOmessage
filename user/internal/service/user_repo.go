package service

import (
	"context"

	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/storage"
	"gomessage.com/users/pkg/utils"
)

type userService struct {
	userRepository storage.UserRepository
}

func CreateNewUserService(repo storage.UserRepository) UserService {
	return &userService{
		userRepository: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.UserModel) error {
	var err error
	user.PasswordHash, err = utils.GeneratePasswordHash(user.PasswordHash)
	if err != nil {
		panic("Error hashing password")
	}
	return s.userRepository.Create(ctx, user)
}

func (s *userService) ListOfUsers(ctx context.Context) ([]models.UserModel, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *userService) GetUser(ctx context.Context, id string) (models.UserModel, error) {
	return s.userRepository.FindOne(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user models.UserModel) error {
	var err error
	if user.PasswordHash != " " {
		user.PasswordHash, err = utils.GeneratePasswordHash(user.PasswordHash)
		if err != nil {
			panic("Error hashing password")
		}
	}
	return s.userRepository.Update(ctx, &user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) (string, error) {
	return s.userRepository.Delete(ctx, id)
}
