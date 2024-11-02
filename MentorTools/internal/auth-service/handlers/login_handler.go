package handlers

import (
	"encoding/json"
	"net/http"

	"MentorTools/internal/auth-service/models"
	"MentorTools/internal/auth-service/services"
	"MentorTools/pkg/common"
	"github.com/jackc/pgx/v4/pgxpool"
)

// LoginHandler handles user login and token generation.
func LoginHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest models.UserLoginRequest

		// Decode login data
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			response := common.NewErrorResponse("AUTH400", "Invalid request payload")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Check if required fields are empty
		if loginRequest.Email == "" || loginRequest.Password == "" {
			response := common.NewErrorResponse("AUTH400", "Missing required fields")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Validate email format
		if !common.IsValidEmail(loginRequest.Email) {
			response := common.NewErrorResponse("AUTH0002", "Invalid email format")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Create a context with request's context
		ctx := r.Context()

		// Authenticate user using the service layer
		response := services.AuthenticateUser(ctx, dbPool, loginRequest)

		// Set response content type
		w.Header().Set("Content-Type", "application/json")

		// Check response code and set appropriate HTTP status
		switch response.Code {
		case "AUTH0005": // User not found
			w.WriteHeader(http.StatusNotFound)
		case "AUTH0004": // Invalid password
			w.WriteHeader(http.StatusUnauthorized)
		case "SUCCESS":
			w.WriteHeader(http.StatusOK)
		default: // General internal error
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Encode and send the response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
