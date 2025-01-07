package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID   string  `json:"user_id"`
	Nickname string  `json:"nickname"`
	Age      *int    `json:"age"`
	ImageUrl *string `json:"imge_url"`
	jwt.StandardClaims
}
