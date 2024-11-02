package services

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/pkg/common"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser hashes the password and attempts to insert a new user into the database.
func RegisterUser(ctx context.Context, dbPool *pgxpool.Pool, user models.User) *common.Response {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return common.NewErrorResponse("AUTH500", "Error hashing password")
	}

	// Variables to hold the result of the function call
	var code, message string

	// Call the fn_create_user function to attempt user creation
	err = dbPool.QueryRow(
		ctx,
		"SELECT code, message FROM fn_create_user($1, $2, $3, $4)",
		user.Username, string(hashedPassword), user.Email, user.Role,
	).Scan(&code, &message)

	if err != nil {
		return common.NewErrorResponse("AUTH500", "Database error")
	}

	// Check the returned code to determine success or specific error
	if code != "SUCCESS" {
		return common.NewErrorResponse(code, message)
	}

	// Return a success response with optional data
	return common.NewSuccessResponse("User created successfully", map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	})
}
