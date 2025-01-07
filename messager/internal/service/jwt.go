package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/models"
)

func ValidateToken(tokenString string, secretKey []byte) (*models.Claims, error) {
	logrus.Info(tokenString)
	logrus.Info(secretKey)
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Error("Invalid signing method")
			return nil, errors.New("invalid token method")
		}
		return secretKey, nil
	})

	if err != nil {
		logrus.Errorf("Token parsing error: %v", err)
		return nil, err
	}

	logrus.Info("Token parsed successfully")

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		logrus.Error("Invalid token claims")
		return nil, errors.New("invalid token claims")
	}

	logrus.Info("Token claims validated successfully")

	if claims.ExpiresAt < time.Now().Unix() {
		logrus.Error("Token expired")
		return nil, errors.New("token expired")
	}

	logrus.Info("Token is valid and not expired")
	return claims, nil
}
