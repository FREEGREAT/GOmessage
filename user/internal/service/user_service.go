package service

import (
    "context"
    "gomessage.com/users/internal/models"
    "gomessage.com/users/internal/storage"
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
    return s.userRepository.Create(ctx, user)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]models.UserModel, error) {
    return s.userRepository.FindAll(ctx)
}

func (s *userService) GetUser(ctx context.Context, id string) (models.UserModel, error) {
    return s.userRepository.FindOne(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user models.UserModel) error {
    return s.userRepository.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
    return s.userRepository.Delete(ctx, id)
}
