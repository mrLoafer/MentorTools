package users

import (
	"MentorTools/models"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func UpdateUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		query := `UPDATE users SET email=$1, name=$2, role=$3, updated_at=NOW() WHERE id=$4`
		_, err = db.Exec(context.Background(), query, user.Email, user.Name, user.Role, userID)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User updated successfully",
		})
	}
}

func GetUserHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		var user models.User
		query := `SELECT email, name, role FROM users WHERE id=$1`
		err := db.QueryRow(context.Background(), query, userID).Scan(&user.Email, &user.Name, &user.Role)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func ListUsersHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.URL.Query().Get("role")

		var users []models.User
		query := `SELECT email, name, role FROM users`
		if role != "" {
			query += ` WHERE role=$1`
		}

		rows, err := db.Query(context.Background(), query, role)
		if err != nil {
			http.Error(w, "Error fetching users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user models.User
			err := rows.Scan(&user.Email, &user.Name, &user.Role)
			if err != nil {
				http.Error(w, "Error reading user", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}
