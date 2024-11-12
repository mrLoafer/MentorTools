package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
}

type UserLinkRequest struct {
	Email string `json:"email"`
}
