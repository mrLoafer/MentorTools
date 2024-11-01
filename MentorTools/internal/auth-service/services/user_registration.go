package services

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/pkg/common"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// RegisterUser hashes the password and inserts a new user into the database.
func RegisterUser(ctx context.Context, dbPool *pgxpool.Pool, user models.User) *common.AppError {
	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return common.NewAppError("AUTH500", "Error hashing password")
	}

	// Выполняем SQL-запрос для создания пользователя
	_, err = dbPool.Exec(ctx, "CALL fnCreateUser($1, $2, $3, $4)", user.Username, hashedPassword, user.Email, user.Role)
	if err != nil {
		if strings.Contains(err.Error(), "AUTH0001") {
			return common.NewAppError("AUTH0001", "User with the same email already exists")
		}
		return common.NewAppError("AUTH500", "Database error")
	}
	return nil
}
