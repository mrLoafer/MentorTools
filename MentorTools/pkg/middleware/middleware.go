package middleware

import (
	"MentorTools/pkg/common"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// contextKey defines a type for storing values in context
type contextKey string

const userContextKey = contextKey("user")

// AuthMiddleware validates the JWT token and extracts user information
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendErrorResponse(w, http.StatusUnauthorized, "AUTH401", "Authorization header is missing")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			sendErrorResponse(w, http.StatusUnauthorized, "AUTH401", "Invalid token format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		publicKeyPath := "/app/public_key.pem"

		// Load and parse the public key
		publicKey, err := loadPublicKey(publicKeyPath)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "AUTH500", "Failed to load public key")
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("invalid signing algorithm")
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			sendErrorResponse(w, http.StatusUnauthorized, "AUTH401", "Invalid token")
			return
		}

		// Extract claims and validate them
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !validateClaims(claims) {
			sendErrorResponse(w, http.StatusUnauthorized, "AUTH401", "Invalid token claims")
			return
		}

		// Extract specific user information from claims
		userInfo := map[string]interface{}{
			"email": claims["email"],
			"role":  claims["role"],
			"name":  claims["name"],
		}

		// Add user information to the request context
		ctx := context.WithValue(r.Context(), userContextKey, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loadPublicKey reads and parses the PEM-encoded public key file
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// validateClaims checks if the token has expired
func validateClaims(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(float64)
	return ok && time.Unix(int64(exp), 0).After(time.Now())
}

// sendErrorResponse sends an error response in JSON format using the common response model
func sendErrorResponse(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := common.NewErrorResponse(code, message)
	json.NewEncoder(w).Encode(response)
}

// GetUserFromContext retrieves user information from the request context
func GetUserFromContext(ctx context.Context) (map[string]interface{}, bool) {
	user, ok := ctx.Value(userContextKey).(map[string]interface{})
	return user, ok
}
