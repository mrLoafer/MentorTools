package models

import "github.com/dgrijalva/jwt-go"

// UserRegistrationRequest represents the payload for a user registration request.
type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Username string `json:"username"`
}

// UserLoginRequest represents the payload for a user login request.
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a user in the database.
type User struct {
	Email    string
	Password string
	Role     string
	Username string
}

type JwtData struct {
	ID    string `json:"userId"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}
