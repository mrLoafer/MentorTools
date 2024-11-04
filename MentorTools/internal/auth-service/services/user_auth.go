package services

import (
	"MentorTools/internal/auth-service/models"
	"MentorTools/pkg/common"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser authenticates the user by checking credentials and returning a JWT token if successful.
func AuthenticateUser(ctx context.Context, dbPool *pgxpool.Pool, loginRequest models.UserLoginRequest) *common.Response {
	var code, message string
	var email, passwordHash, roleName *string
	var userID *int

	fmt.Println("Health check of logINN service")
	// Execute the function fn_find_user_by_email and retrieve response fields
	err := dbPool.QueryRow(ctx, `SELECT code, message, user_id, email, password_hash, role_name
									  FROM fn_find_user_by_email($1)`, loginRequest.Email).
		Scan(&code, &message, &userID, &email, &passwordHash, &roleName)

	fmt.Printf("response frome DB: code - %v, msg - %v\n", code, message)
	if err != nil {
		return common.NewErrorResponse("AUTH500", "Database error: "+err.Error())
	}

	// Check if the stored procedure returned a user not found error
	if code == "AUTH0005" {
		return common.NewErrorResponse(code, message)
	}

	// Verify the user`s password
	if err := bcrypt.CompareHashAndPassword([]byte(*passwordHash), []byte(loginRequest.Password)); err != nil {
		return common.NewErrorResponse("AUTH0004", "Invalid password")
	}

	fmt.Println("Starting GenerateJWT function")
	// Generate JWT token
	tokenResponse := GenerateJWT(models.JwtData{
		ID:    *userID,
		Email: *email,
		Role:  *roleName,
	})

	// Check if token generation was successful or an error occurred
	if tokenResponse.Code != "SUCCESS" {
		return tokenResponse // Return the error response directly
	}

	// Return success response with the generated token
	return common.NewSuccessResponse("User authenticated successfully", tokenResponse.Data)
}
