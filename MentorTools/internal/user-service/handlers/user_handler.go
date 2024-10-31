package handlers

import (
	"MentorTools/auth-service/models"
	"MentorTools/user-service/services"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SearchUsersHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем роль пользователя из токена
		userRole, err := services.GetUserRoleFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Извлекаем userID из токена
		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Определяем, кого нужно искать (учеников для учителя или учителей для ученика)
		var roleToSearch string
		if userRole == "teacher" {
			roleToSearch = "student"
		} else if userRole == "student" {
			roleToSearch = "teacher"
		} else {
			http.Error(w, "Invalid role", http.StatusForbidden)
			return
		}

		// Выполняем запрос для получения всех пользователей с противоположной ролью и проверяем наличие связи
		query := `
		SELECT u.id, u.name, u.email, 
		       CASE WHEN ts.teacher_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_linked
		FROM users u
		LEFT JOIN teacher_student ts 
		ON (u.id = ts.student_id AND ts.teacher_id = $1) OR (u.id = ts.teacher_id AND ts.student_id = $1)
		WHERE u.role = $2;
		`

		rows, err := dbpool.Query(context.Background(), query, userID, roleToSearch)
		if err != nil {
			http.Error(w, "Failed to search users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Создаём список пользователей с информацией о связи
		var users []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			IsLinked bool   `json:"is_linked"`
		}

		for rows.Next() {
			var user struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Email    string `json:"email"`
				IsLinked bool   `json:"is_linked"`
			}
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.IsLinked); err != nil {
				http.Error(w, "Error reading user data", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		// Возвращаем список пользователей с информацией о связях
		json.NewEncoder(w).Encode(users)
	}
}

func ProfileHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем идентификатор пользователя из токена
		userID, err := services.GetUserIDFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized"+err.Error(), http.StatusUnauthorized)
			return
		}

		// Выполняем запрос к базе данных для получения данных профиля пользователя
		var user models.User
		err = dbpool.QueryRow(context.Background(), "SELECT id, name, email, role FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Возвращаем данные профиля в формате JSON
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateProfileHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := services.GetUserIDFromToken(r) // Извлекаем userID из токена
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Структура для получения данных из запроса
		var updatedUser struct {
			Name   string `json:"name"`
			Email  string `json:"email"`
			Topics []int  `json:"topics"` // Ожидаем массив ID тем
		}

		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Обновляем данные пользователя в базе
		_, err = dbpool.Exec(context.Background(), "UPDATE users SET name=$1, email=$2 WHERE id=$3", updatedUser.Name, updatedUser.Email, userID)
		if err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		// Удаляем старые топики пользователя
		_, err = dbpool.Exec(context.Background(), "DELETE FROM student_topics WHERE student_id = $1", userID)
		if err != nil {
			http.Error(w, "Failed to clear existing topics", http.StatusInternalServerError)
			return
		}

		// Добавляем новые топики пользователя
		for _, topicID := range updatedUser.Topics {
			_, err = dbpool.Exec(context.Background(), "INSERT INTO student_topics (student_id, topic_id) VALUES ($1, $2)", userID, topicID)
			if err != nil {
				http.Error(w, "Failed to save topics", http.StatusInternalServerError)
				return
			}
		}

		// Возвращаем успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "Profile updated successfully"})
	}
}

func CreateTeacherStudentLink(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var linkData struct {
			TeacherID string `json:"teacher_id"`
			StudentID string `json:"student_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&linkData); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Добавляем связь в таблицу
		_, err := dbpool.Exec(context.Background(), "INSERT INTO teacher_student (teacher_id, student_id) VALUES ($1, $2)", linkData.TeacherID, linkData.StudentID)
		if err != nil {
			http.Error(w, "Failed to create link", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "Link created successfully"})
	}
}

func RemoveTeacherStudentLink(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var linkData struct {
			TeacherID string `json:"teacher_id"`
			StudentID string `json:"student_id"`
		}

		// Логируем для отладки
		log.Println("Received request to unlink:", linkData)

		if err := json.NewDecoder(r.Body).Decode(&linkData); err != nil {
			log.Println("Error decoding request body:", err)
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Удаляем связь из таблицы
		_, err := dbpool.Exec(context.Background(), "DELETE FROM teacher_student WHERE teacher_id=$1 AND student_id=$2", linkData.TeacherID, linkData.StudentID)
		if err != nil {
			log.Println("Error removing link:", err)
			http.Error(w, "Failed to remove link", http.StatusInternalServerError)
			return
		}

		log.Printf("Link removed: teacher_id=%s, student_id=%s\n", linkData.TeacherID, linkData.StudentID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "Link removed successfully"})
	}
}

func GetTeacherStudentLinks(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := services.GetUserIDFromToken(r) // Извлекаем userID из токена
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("Fetching links for user ID:", userID)

		// Получаем все связи для данного пользователя
		rows, err := dbpool.Query(context.Background(), "SELECT teacher_id, student_id FROM teacher_student WHERE teacher_id=$1 OR student_id=$1", userID)
		if err != nil {
			http.Error(w, "Failed to fetch links", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var links []struct {
			TeacherID string `json:"teacher_id"`
			StudentID string `json:"student_id"`
		}

		for rows.Next() {
			var link struct {
				TeacherID string `json:"teacher_id"`
				StudentID string `json:"student_id"`
			}
			if err := rows.Scan(&link.TeacherID, &link.StudentID); err != nil {
				http.Error(w, "Error reading data", http.StatusInternalServerError)
				return
			}
			links = append(links, link)
		}

		log.Println("Links fetched successfully:", links)
		json.NewEncoder(w).Encode(links)
	}
}

// GetUserRoleHandler - возвращает роль пользователя, основываясь на его токене
func GetUserRoleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем роль пользователя из токена
		userRole, err := services.GetUserRoleFromToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Возвращаем роль пользователя в JSON
		json.NewEncoder(w).Encode(map[string]string{"role": userRole})
	}
}

func GetAvailableContextsHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Запрос для получения тем из таблицы contexts
		rows, err := dbpool.Query(context.Background(), "SELECT id, topic_name FROM topics")
		if err != nil {
			http.Error(w, "Failed to retrieve contexts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Массив для хранения тем
		var contexts []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				http.Error(w, "Error scanning contexts", http.StatusInternalServerError)
				return
			}
			contexts = append(contexts, map[string]interface{}{
				"id":   id,
				"name": name,
			})
		}

		// Возвращаем темы в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"contexts": contexts})
	}
}
