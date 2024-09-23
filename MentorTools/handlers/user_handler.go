package handlers

import (
	"MentorTools/models"
	"MentorTools/services"
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4"
)

func SearchUsersHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем роль пользователя из токена или сессии
		userRole, err := services.GetUserRoleFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Определяем, кого нужно искать
		var roleToSearch string
		if userRole == "teacher" {
			roleToSearch = "student"
		} else if userRole == "student" {
			roleToSearch = "teacher"
		} else {
			http.Error(w, "Invalid role", http.StatusForbidden)
			return
		}

		// Выполняем запрос для поиска пользователей с противоположной ролью
		rows, err := conn.Query(context.Background(), "SELECT id, name, email FROM users WHERE role = $1", roleToSearch)
		if err != nil {
			http.Error(w, "Failed to search users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
				http.Error(w, "Error reading user data", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		// Возвращаем список пользователей
		json.NewEncoder(w).Encode(users)
	}
}

func ProfileHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем идентификатор пользователя из токена
		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized"+err.Error(), http.StatusUnauthorized)
			return
		}

		// Выполняем запрос к базе данных для получения данных профиля пользователя
		var user models.User
		err = conn.QueryRow(context.Background(), "SELECT id, name, email, role FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Возвращаем данные профиля в формате JSON
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateProfileHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := services.GetUserIDFromToken(r) // Извлекаем userID из токена
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var updatedUser models.User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Обновляем данные пользователя в базе
		_, err = conn.Exec(context.Background(), "UPDATE users SET name=$1, email=$2 WHERE id=$3", updatedUser.Name, updatedUser.Email, userID)
		if err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "Profile updated successfully"})
	}
}
