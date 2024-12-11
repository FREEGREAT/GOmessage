package models

type UserModel struct {
	ID           string  `json:"user_id"`
	Nickname     string  `json:"nickname"`
	PasswordHash string  `json:"password_hash"`
	Email        string  `json:"email"`
	Age          *int    `json:"age"`
	ImageUrl     *string `json:"imge_url"`
}

type FriendList struct {
	UserID   string `json:"user_id"`
	FriendID string `json:"friend_id"`
}
