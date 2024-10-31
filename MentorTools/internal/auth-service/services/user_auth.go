package services

import (
	"context"

	"MentorTools/internal/auth-service/models"
	"MentorTools/internal/auth-service/repository"
	"MentorTools/pkg/common"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser authenticates the user by checking credentials and returning a JWT token if successful.
func AuthenticateUser(ctx context.Context, loginRequest models.UserLoginRequest) (string, *common.AppError) {
	var user models.User

	// Use global pool from the repository package for database access
	err := repository.DBPool.QueryRow(ctx, "SELECT * FROM fnfinduserbyemail($1)", loginRequest.Email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", common.NewAppError("AUTH0003", "User not found")
		}
		return "", common.NewAppError("AUTH500", "Database error")
	}

	// Verify the user's password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return "", common.NewAppError("AUTH0004", "Invalid password")
	}

	// Generate JWT token
	token, appErr := GenerateJWT(models.JwtData{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
	if appErr != nil {
		return "", appErr
	}

	return token, nil
}
