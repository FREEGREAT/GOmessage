package models

type UserModel struct {
	ID           string  `json:"user_id"`
	Nickname     string  `json:"nickname" binding:"required"`
	PasswordHash string  `json:"password_hash" binding:"required"`
	Email        string  `json:"email" binding:"required"`
	Age          *int    `json:"age"`
	ImageUrl     *string `json:"imge_url"`
}
