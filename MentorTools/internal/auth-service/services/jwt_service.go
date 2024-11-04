package services

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/pkg/common"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT creates a JWT token for the user with given claims
func GenerateJWT(user models.JwtData) *common.Response {
	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)
	user.RegisteredClaims = jwt.RegisteredClaims{
		Subject:   fmt.Sprint(user.ID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	fmt.Printf("user.RegisteredClaims: %v\n", user.RegisteredClaims)
	// Путь к файлу с приватным ключом
	privateKeyPath := "/app/private_key.pem"

	fmt.Printf("Private key privateKeyPath: %v\n", privateKeyPath)

	// Чтение приватного ключа из файла
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("Error reading private key: %v\n", err)
		return common.NewErrorResponse("AUTH500", "Failed to read private key: "+err.Error())
	}

	// Логирование длины ключа для проверки чтения
	fmt.Printf("Private key length: %d\n", len(privateKey))

	// Парсинг приватного ключа
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return common.NewErrorResponse("AUTH500", fmt.Sprintf("Could not parse private key: %v", err))
	}

	// Создание нового токена с методом подписи RS256 и указанными claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, user)

	// Подпись токена с использованием приватного ключа
	tokenString, err := token.SignedString(key)
	if err != nil {
		return common.NewErrorResponse("AUTH500", fmt.Sprintf("Could not sign token: %v", err))
	}

	// Возврат успешного ответа с токеном
	return common.NewSuccessResponse("Token generated successfully", tokenString)
}
