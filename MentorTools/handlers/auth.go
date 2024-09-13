package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mentor_tools/models"
	"mentor_tools/services"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error encrypting password", http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO users (email, password, name, role, created_at, updated_at) 
                  VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = db.Exec(context.Background(), query, user.Email, string(hashedPassword), user.Name, user.Role, time.Now(), time.Now())
		if err != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully!")
	}
}

func LoginHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var storedPassword string
		var role string

		err = db.QueryRow(context.Background(), "SELECT password, role FROM users WHERE email=$1", creds.Email).Scan(&storedPassword, &role)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := services.GenerateJWT(creds.Email, role)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the protected route!")
}
