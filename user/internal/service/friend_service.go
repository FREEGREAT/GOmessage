package service

import (
	"context"

	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/storage"
)

type friendService struct {
	friendRepository storage.FriendsRepository
}

// AddFriend implements FriendsService.
func (f *friendService) AddFriend(ctx context.Context, friend *models.FriendListModel) error {
	return f.friendRepository.Create(ctx, friend)
}

// DeleteFriend implements FriendsService.
func (f *friendService) DeleteFriend(ctx context.Context, friends *models.FriendListModel) error {
	return f.friendRepository.Delete(ctx,friends)
}

func CreateNewFriendService(repo storage.FriendsRepository) FriendsService {
	return &friendService{
		friendRepository: repo,
	}
}
