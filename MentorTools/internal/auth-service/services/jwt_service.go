package services

import (
	"fmt"
	"os"
	"time"

	"MentorTools/internal/auth-service/models"
	"MentorTools/pkg/common"
	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT creates a JWT token for the user with given claims
func GenerateJWT(user models.JwtData) (string, *common.AppError) {
	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)
	user.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	// Retrieve private key from environment variable
	privateKey := os.Getenv("JWT_PRIVATE_KEY")
	if privateKey == "" {
		return "", common.NewAppError("AUTH500", "Private key not found in environment variables")
	}

	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", common.NewAppError("AUTH500", fmt.Sprintf("Could not parse private key: %v", err))
	}

	// Create a new token with RS256 signing method and the provided claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, user)

	// Sign the token with the private key
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", common.NewAppError("AUTH500", fmt.Sprintf("Could not sign token: %v", err))
	}

	return tokenString, nil
}
