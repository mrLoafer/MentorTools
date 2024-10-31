package handlers

import (
	"MentorTools/auth-service/models"
	"MentorTools/auth-service/services"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

// LoginHandler handles user login, validates credentials, and returns a JWT token.
func LoginHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		// Decode login credentials from request body
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check user in the database (example query; adapt as needed)
		var user models.JwtData
		err := dbpool.QueryRow(
			context.Background(),
			"SELECT id, email, role FROM users WHERE email=$1 AND password=$2",
			credentials.Email, credentials.Password,
		).Scan(&user.ID, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := services.GenerateJWT(user)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Respond with the token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
