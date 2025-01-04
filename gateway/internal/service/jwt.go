package service

import (
	"time"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey   = "secret"
	errorString = ""
)

type jwtClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Age      string `json:"age"`
	ImageUrl string `json:"img_url"`
	jwt.RegisteredClaims
}

func ClaimToken(grpcResponse *proto_user_service.LoginUserResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		ID:       grpcResponse.Id,
		Username: grpcResponse.Username,
		Age:      grpcResponse.Age,
		ImageUrl: grpcResponse.ImageUrl,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    grpcResponse.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return errorString, err
	}
	return ss, nil
}
