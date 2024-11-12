package handlers

import (
	"MentorTools/internal/user-service/services"
	"MentorTools/pkg/common"
	"MentorTools/pkg/middleware"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

// UpdateUserProfileHandler handles the request to update the user's name in the profile
func UpdateUserProfileHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user information from the context
		claims, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(common.NewErrorResponse("AUTH401", "Invalid token claims"))
			return
		}

		userID, ok := claims["user_id"].(int)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(common.NewErrorResponse("AUTH401", "Invalid user ID in token"))
			return
		}

		// Decode the request body to get the new name
		var updateRequest struct {
			NewName string `json:"new_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.NewErrorResponse("REQ400", "Invalid request payload"))
			return
		}

		// Check if the name is empty
		if updateRequest.NewName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.NewErrorResponse("REQ400", "Name cannot be empty"))
			return
		}

		// Call the service layer to update the user's name
		response := services.UpdateUserName(context.Background(), dbPool, userID, updateRequest.NewName)
		if response.Code != "SUCCESS" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Send a successful response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
