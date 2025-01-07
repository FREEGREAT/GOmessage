package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"gomessage.com/gateway/internal/service"
)

type JWTMiddleware struct {
	jwtService *service.JWTService
}

func NewJWTMiddleware(jwtService *service.JWTService) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

func (m *JWTMiddleware) Middleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Extraction токена

		tokenString := extractToken(r)
		if tokenString == "" {
			logrus.Warn("Missing authorization token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Валідація токена
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			logrus.Errorf("Token validation error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Додавання UserID в контекст
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		r = r.WithContext(ctx)

		// Виклик наступного handler
		next(w, r, params)
	}
}

// Extraction токена
func extractToken(r *http.Request) string {
	// Перевірка Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Перевірка cookie
	cookie, err := r.Cookie("jwt")
	if err == nil {
		return cookie.Value
	}

	return ""
}
