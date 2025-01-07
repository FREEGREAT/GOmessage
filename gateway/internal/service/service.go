package service

import (
	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	GenerateAccessToken(grpcResponse proto_user_service.LoginUserResponse) (string, error)
	ValidateAccessToken(token string) (jwt.MapClaims, error)
	GenerateRefreshToken(userID string) (string, error)
	ParseRefreshToken(tokenString string) (jwt.MapClaims, error)
}
