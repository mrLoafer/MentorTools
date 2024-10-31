package handlers

import (
	"encoding/json"
	"net/http"

	"MentorTools/internal/auth-service/models"
	"MentorTools/internal/auth-service/services"
)

// LoginHandler handles user login and token generation.
func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest models.UserLoginRequest

		// Decode login data
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Create a context with request's context
		ctx := r.Context()

		// Authenticate user using the service layer
		token, appErr := services.AuthenticateUser(ctx, loginRequest)
		if appErr != nil {
			// Set response content type and handle error based on AppError code
			w.Header().Set("Content-Type", "application/json")

			switch appErr.Code {
			case "AUTH0003": // User not found
				w.WriteHeader(http.StatusNotFound)
			case "AUTH0004": // Invalid password
				w.WriteHeader(http.StatusUnauthorized)
			default: // General internal error
				w.WriteHeader(http.StatusInternalServerError)
			}

			// Encode and send the error response
			if err := json.NewEncoder(w).Encode(appErr); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		// Respond with the generated token on successful login
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
