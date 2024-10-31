package middleware

import (
	"context"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// contextKey is a type used for storing values in context
type contextKey string

const userContextKey = contextKey("user")

// AuthMiddleware validates JWT tokens in the Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Check if the token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " from the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Fetch the public key from the environment variable
		publicKeyStr := os.Getenv("JWT_PUBLIC_KEY")
		if publicKeyStr == "" {
			http.Error(w, "Public key is missing", http.StatusInternalServerError)
			return
		}

		// Decode the PEM-encoded public key
		block, _ := pem.Decode([]byte(publicKeyStr))
		if block == nil || block.Type != "PUBLIC KEY" {
			http.Error(w, "Invalid public key format", http.StatusInternalServerError)
			return
		}

		// Parse the RSA public key
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyStr))
		if err != nil {
			http.Error(w, "Failed to parse public key", http.StatusInternalServerError)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is RS256
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("invalid signing algorithm")
			}
			return publicKey, nil
		})

		// Handle parsing or validation common
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract and validate the token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the token's expiration time
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					http.Error(w, "Token has expired", http.StatusUnauthorized)
					return
				}
			}
			// Add user information to the request context
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
	})
}

// GetUserFromContext retrieves user information from the request context
func GetUserFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	user, ok := ctx.Value(userContextKey).(jwt.MapClaims)
	return user, ok
}
