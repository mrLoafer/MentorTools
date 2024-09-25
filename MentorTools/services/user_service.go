package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Извлечение роли пользователя из JWT токена
func GetUserRoleFromToken(r *http.Request) (string, error) {
	// Извлекаем токен из заголовка Authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("missing token")
	}

	// Убираем префикс "Bearer " из заголовка
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil // Ваш секретный ключ
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Извлекаем данные из клеймов (claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if role, ok := claims["role"].(string); ok {
			return role, nil
		} else {
			return "", errors.New("role not found in token")
		}
	}

	return "", errors.New("invalid token claims")
}

// GetUserIDFromToken извлекает userID из JWT токена
func GetUserIDFromToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("missing token")
	}

	// Извлекаем токен, убирая "Bearer "
	tokenString = tokenString[len("Bearer "):]

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Извлекаем данные из токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["userId"].(string); ok {
			return userID, nil
		}
	}

	return "", errors.New("userID not found in token")
}
