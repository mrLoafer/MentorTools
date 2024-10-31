package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"MentorTools/internal/auth-service/models"
	"MentorTools/internal/auth-service/services"
	"MentorTools/pkg/common"

	"github.com/jackc/pgx/v4/pgxpool"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// isValidEmail checks if the provided email is in a valid format.
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// RegisterHandler handles user registration without generating a token.
func RegisterHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser models.UserRegistrationRequest

		// Decode registration data
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate email format
		if !isValidEmail(newUser.Email) {
			appErr := common.NewAppError("AUTH0002", "Invalid email format")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(appErr); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		// Register the user using the service layer
		user := models.User{
			Email:    newUser.Email,
			Password: newUser.Password,
			Role:     newUser.Role,
			Username: newUser.Username,
		}
		appErr := services.RegisterUser(context.Background(), dbpool, user)
		if appErr != nil {
			// Return structured error with code and message
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict) // Set HTTP status based on error type
			if err := json.NewEncoder(w).Encode(appErr); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
