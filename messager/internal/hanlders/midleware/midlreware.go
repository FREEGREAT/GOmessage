package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gommessage.com/messager/internal/service"
	"gommessage.com/messager/pkg"
)

func WSAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	if err := pkg.InitConfig(); err != nil {
		panic("Error init config at middleware")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Invalid authorization token format", http.StatusUnauthorized)
			return
		}
		logrus.Info("Received token for validation")
		claims, err := service.ValidateToken(token[len("Bearer "):], []byte(viper.GetString("jwt.secret")))
		if err != nil {
			logrus.Info("loshara")
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
