package services

import (
	"fmt"
	"os"
	"time"

	"MentorTools/auth-service/models"
	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT creates a JWT token for the user with given claims
func GenerateJWT(user models.JwtData) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour) // Token validity: 24 hours
	user.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	// Retrieve private key from environment variable
	privateKey := os.Getenv("JWT_PRIVATE_KEY")
	if privateKey == "" {
		return "", fmt.Errorf("private key not found in environment variables")
	}

	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", fmt.Errorf("could not parse private key: %v", err)
	}

	// Create a new token with RS256 signing method and the provided claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, user)

	// Sign the token with the private key
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("could not sign token: %v", err)
	}

	return tokenString, nil
}
