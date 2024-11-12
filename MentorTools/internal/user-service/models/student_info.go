package models

// StudentInfo represents the information about a student.
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
}

type UserLinkRequest struct {
	Email string
}
