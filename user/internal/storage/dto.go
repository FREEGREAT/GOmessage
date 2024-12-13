package storage

type CreateUserDTO struct {
	Nickname     string `json:"nickname"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
	Age          int    `json:"age"`
	ImageUrl     string `json:"image_url"`
}

type AddFriendsDTO struct {
	UserID   string `json:"user_id"`
	FriendID string `json:"user_id"`
}
