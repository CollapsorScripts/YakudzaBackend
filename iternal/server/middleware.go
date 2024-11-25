package server

import (
	"Yakudza/pkg/logger"
	"Yakudza/pkg/token"
	"context"
	"net/http"
	"strings"
)

// authMiddleware - промежуточное ПО для приватных запросов
func (route Router) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json, Authorization")
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			SetHTTPError(w, "В доступе отказано", http.StatusForbidden)
			return
		}

		logger.Info("Приватный запрос: %s [%s]", r.URL.String(), r.Method)

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		jwtToken, err := token.ParseToken(tokenString, route.cfg)
		if err != nil {
			logger.Error("Ошибка парса токена: %v", err)
			SetHTTPError(w, "Ошибка доступа, неверный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", jwtToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// publicMiddleware - промежуточное ПО для публичных запросов
func (route Router) publicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		logger.Info("Публичный запрос: %s [%s]", r.URL.String(), r.Method)

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}
