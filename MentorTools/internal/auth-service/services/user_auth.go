package services

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/internal/auth-service/repository"
	"MentorTools/pkg/common"
	"context"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser authenticates the user.
func AuthenticateUser(ctx context.Context, loginRequest models.UserLoginRequest) (string, *common.AppError) {
	// Подключаемся к базе данных auth_db
	pool, err := repository.ConnectDB(ctx)
	if err != nil {
		return "", common.NewAppError("AUTH500", "Database connection error")
	}
	defer repository.CloseDB(pool)

	var user models.User

	// Call the fnFindUserByEmail function to retrieve user details by email
	err = pool.QueryRow(ctx, "SELECT * FROM fnFindUserByEmail($1)", loginRequest.Email).
		Scan(&user.Email, &user.Password, &user.Role, &user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", common.NewAppError("AUTH0003", "User not found")
		}
		return "", common.NewAppError("AUTH500", "Database error")
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return "", common.NewAppError("AUTH0004", "Invalid password")
	}

	// Generate JWT token
	token, err := GenerateJWT(models.JwtData{
		Email: user.Email,
		Role:  user.Role,
	})
	if err != nil {
		return "", common.NewAppError("AUTH500", "Failed to generate token")
	}

	return token, nil
}
