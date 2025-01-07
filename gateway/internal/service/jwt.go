package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"gomessage.com/gateway/internal/models"
)

type JWTService struct {
	secretKey []byte
	issuer    string
}

const emptyString = ""

func NewJWTService() *JWTService {
	secretKey := []byte(viper.GetString("jwt.secret"))
	if len(secretKey) == 0 {
		panic("Pleace create key")
	}

	return &JWTService{
		secretKey: secretKey,
		issuer:    "gomessage_gateway",
	}
}

func (s *JWTService) GenerateToken(user_id, nickname, image_url string, age int) (string, error) {

	claims := models.Claims{
		UserID:   user_id,
		Nickname: nickname,
		Age:      &age,
		ImageUrl: &image_url,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token method")
			}
			return s.secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Перевірка claims
	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Перевірка терміну дії
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return emptyString, err
	}

	// Генерація нового токена
	return s.GenerateToken(claims.Id, claims.Nickname, *claims.ImageUrl, *claims.Age)
}
