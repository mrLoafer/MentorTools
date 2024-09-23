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

func CreateTeacherStudentLink(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var linkData struct {
			TeacherID int `json:"teacher_id"`
			StudentID int `json:"student_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&linkData); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Добавляем связь в таблицу
		_, err := conn.Exec(context.Background(), "INSERT INTO teacher_student (teacher_id, student_id) VALUES ($1, $2)", linkData.TeacherID, linkData.StudentID)
		if err != nil {
			http.Error(w, "Failed to create link", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "Link created successfully"})
	}
}

func RemoveTeacherStudentLink(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var linkData struct {
			TeacherID int `json:"teacher_id"`
			StudentID int `json:"student_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&linkData); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Удаляем связь из таблицы
		_, err := conn.Exec(context.Background(), "DELETE FROM teacher_student WHERE teacher_id=$1 AND student_id=$2", linkData.TeacherID, linkData.StudentID)
		if err != nil {
			http.Error(w, "Failed to remove link", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "Link removed successfully"})
	}
}

func GetTeacherStudentLinks(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := services.GetUserIDFromToken(r) // Извлекаем userID из токена
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Получаем все связи для данного пользователя
		rows, err := conn.Query(context.Background(), "SELECT teacher_id, student_id FROM teacher_student WHERE teacher_id=$1 OR student_id=$1", userID)
		if err != nil {
			http.Error(w, "Failed to fetch links", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var links []struct {
			TeacherID int `json:"teacher_id"`
			StudentID int `json:"student_id"`
		}

		for rows.Next() {
			var link struct {
				TeacherID int `json:"teacher_id"`
				StudentID int `json:"student_id"`
			}
			if err := rows.Scan(&link.TeacherID, &link.StudentID); err != nil {
				http.Error(w, "Error reading data", http.StatusInternalServerError)
				return
			}
			links = append(links, link)
		}

		json.NewEncoder(w).Encode(links)
	}
}