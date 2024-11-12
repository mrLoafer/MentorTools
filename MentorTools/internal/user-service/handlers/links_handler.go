package handlers

import (
	"MentorTools/internal/user-service/models"
	"MentorTools/internal/user-service/services"
	"MentorTools/pkg/common"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

// GetStudentsHandler returns an http.HandlerFunc to get a list of students connected to the teacher.
func GetStudentsHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user claims from the request context
		claims, ok := r.Context().Value("user").(map[string]interface{})
		if !ok || claims["role"] != "teacher" {
			// If the user is not authenticated as a teacher, return a "Forbidden" status
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Convert the user ID from the token claims to an integer
		teacherID := int(claims["id"].(float64)) // Assuming the teacher ID is stored in the token claims

		// Call the service layer to fetch students associated with the teacher
		students, appErr := services.GetStudents(context.Background(), dbPool, teacherID)
		if appErr != nil {
			// If there is an application error (e.g., no students found), return a "Not Found" status with a JSON error response
			common.JSONResponse(w, http.StatusNotFound, common.NewErrorResponse("USER404", appErr.Message))
			return
		}

		// Return the list of students as a JSON response with a "Success" status
		common.JSONResponse(w, http.StatusOK, common.NewSuccessResponse("Students list", students))
	}
}

// CreateLinkHandler returns an http.HandlerFunc to create a link between a teacher and a new student.
func CreateLinkHandler(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user claims from the request context
		claims, ok := r.Context().Value("user").(map[string]interface{})
		if !ok || claims["role"] != "teacher" {
			// If the user is not authenticated as a teacher, return a "Forbidden" status
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Parse the incoming JSON request body into the UserLinkRequest model
		var linkRequest models.UserLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&linkRequest); err != nil {
			// If JSON decoding fails, return a "Bad Request" status with a JSON error response
			common.JSONResponse(w, http.StatusBadRequest, common.NewErrorResponse("REQ400", "Invalid request payload"))
			return
		}

		// Extract teacher ID from token claims
		teacherID := int(claims["id"].(float64)) // Convert teacher ID from token claims to integer

		// Call the service layer to create a link between the teacher and the specified student
		appErr := services.CreateLink(context.Background(), dbPool, teacherID, linkRequest.Email)
		if appErr != nil {
			// If an application error occurs (e.g., link already exists), return a "Conflict" status with a JSON error response
			common.JSONResponse(w, http.StatusConflict, common.NewErrorResponse("LINK409", appErr.Message))
			return
		}

		// Return a success message as a JSON response
		common.JSONResponse(w, http.StatusOK, common.NewSuccessResponse("Link created successfully", nil))
	}
}
