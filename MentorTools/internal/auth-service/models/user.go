package models

import "github.com/golang-jwt/jwt/v4"

// UserBase содержит общие поля, используемые в других моделях пользователя.
type UserBase struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Username string `json:"username"`
}

// UserRegistrationRequest represents the payload for a user registration request.
type UserRegistrationRequest struct {
	UserBase
}

// UserLoginRequest represents the payload for a user login request.
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a user in the database.
type User struct {
	UserBase
	ID int `json:"id"`
}

// JwtData represents the data stored in a JWT token.
type JwtData struct {
	ID    int    `json:"userId"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}
