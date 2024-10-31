package services

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/internal/auth-service/repository"
	"MentorTools/pkg/common"
	"context"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// RegisterUser hashes the password and inserts a new user into the database.
func RegisterUser(ctx context.Context, user models.User) *common.AppError {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return common.NewAppError("AUTH500", "Error hashing password")
	}

	_, err = repository.DBPool.Exec(ctx, "CALL fnCreateUser($1, $2, $3, $4)", user.Username, hashedPassword, user.Email, user.Role)
	if err != nil {
		if strings.Contains(err.Error(), "AUTH0001") {
			return common.NewAppError("AUTH0001", "User with the same email already exists")
		}
		return common.NewAppError("AUTH500", "Database error")
	}
	return nil
}
