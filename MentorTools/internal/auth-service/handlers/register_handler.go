package handlers

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/internal/auth-service/services"
	"MentorTools/pkg/common"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

// RegisterHandler handles user registration.
func RegisterHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser models.UserRegistrationRequest

		// Decode registration data
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			response := common.NewErrorResponse("AUTH400", "Invalid request payload")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if required fields are empty
		if newUser.Email == "" || newUser.Password == "" || newUser.Role == "" || newUser.Username == "" {
			response := common.NewErrorResponse("AUTH400", "Missing required fields")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Validate email format
		if !common.IsValidEmail(newUser.Email) {
			response := common.NewErrorResponse("AUTH0002", "Invalid email format")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Register the user using the service layer
		user := models.User{
			UserBase: models.UserBase{
				Email:    newUser.Email,
				Password: newUser.Password,
				Role:     newUser.Role,
				Username: newUser.Username,
			},
		}
		appErr := services.RegisterUser(context.Background(), dbPool, user)
		if appErr != nil {
			response := common.NewErrorResponse(appErr.Code, appErr.Message)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Respond with a success message
		response := common.NewSuccessResponse("User registered successfully", nil)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
