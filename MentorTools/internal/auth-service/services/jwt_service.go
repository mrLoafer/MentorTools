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

	// Read the private key from the Docker secret file path
	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return common.NewErrorResponse("AUTH500", "Failed to read private key from Docker secret: "+err.Error())
	}

	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return common.NewErrorResponse("AUTH500", fmt.Sprintf("Could not parse private key: %v", err))
	}

	// Create a new token with RS256 signing method and the provided claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, user)

	// Sign the token with the private key
	tokenString, err := token.SignedString(key)
	if err != nil {
		return common.NewErrorResponse("AUTH500", fmt.Sprintf("Could not sign token: %v", err))
	}

	// Return a success response with the token
	return common.NewSuccessResponse("Token generated successfully", tokenString)
}
