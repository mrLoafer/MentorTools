package handlers

import (
	"MentorTools/pkg/common"
	"encoding/json"
	"net/http"
)

// DashboardHandler handles requests to the dashboard page, providing user information from the JWT token.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user information from JWT token claims
	userInfo, ok := r.Context().Value("user").(map[string]interface{})
	if !ok {
		// If the token claims are invalid, return an unauthorized error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.NewErrorResponse("AUTH401", "Invalid token claims"))
		return
	}

	// If valid, return a success response with user information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(common.NewSuccessResponse("User information", userInfo))
}
