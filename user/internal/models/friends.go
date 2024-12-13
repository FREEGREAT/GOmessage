package models

type FriendListModel struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
}
