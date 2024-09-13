package middleware

import (
	"MentorTools/services"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка Authorization
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			http.Error(w, "Missing or invalid token format", http.StatusUnauthorized)
			return
		}

		// Извлекаем токен после "Bearer "
		tokenStr = tokenStr[7:]

		// Структура claims для работы с токеном
		claims := &services.Claims{}

		// Проверка токена с использованием ключа JwtKey
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(services.JwtKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Переходим к следующему обработчику
		next.ServeHTTP(w, r)
	})
}
