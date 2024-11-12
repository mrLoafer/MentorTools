package services

import (
	"MentorTools/pkg/common"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// UpdateUserName обновляет имя пользователя, вызывая хранимую процедуру fn_update_user_name
func UpdateUserName(ctx context.Context, dbPool *pgxpool.Pool, userID int, newName string) *common.Response {
	var code, message string

	// Вызов процедуры обновления имени пользователя
	err := dbPool.QueryRow(
		ctx,
		"SELECT code, message FROM fn_update_user_name($1, $2)",
		userID, newName,
	).Scan(&code, &message)

	if err != nil {
		return common.NewErrorResponse("USER500", "Database error during update")
	}

	// Возвращаем ответ в зависимости от кода результата
	if code != "SUCCESS" {
		return common.NewErrorResponse(code, message)
	}

	return common.NewSuccessResponse(message, nil)
}
