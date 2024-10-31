package services

import (
	"context"
	"fmt"
	"strings"

	"MentorTools/internal/auth-service/models"
	"MentorTools/pkg/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedPassword), nil
}

// RegisterUser hashes the password and calls a stored procedure to create a new user.
func RegisterUser(ctx context.Context, pool *pgxpool.Pool, user models.User) *common.AppError {

	// Hash the user's password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return common.NewAppError("AUTH500", "Internal server error")
	}

	// Call the stored procedure fn_create_user
	_, err = pool.Exec(ctx, "CALL fnCreateUser($1, $2, $3, $4)", user.Username, hashedPassword, user.Email, user.Role)
	if err != nil {
		// Check for specific known error code
		if strings.Contains(err.Error(), "AUTH0001") {
			return common.NewAppError("AUTH0001", "User with the same email already exists")
		}
		// General error response
		return common.NewAppError("AUTH500", "Internal server error")
	}

	return nil
}
