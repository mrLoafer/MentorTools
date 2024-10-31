package common

// AppError represents a structured error with code and message.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewAppError creates a new AppError with a given code and message.
func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
