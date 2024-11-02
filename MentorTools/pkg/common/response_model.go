package common

// Response represents a standardized response model with code, message, and optional data.
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Data can hold any additional info for success responses
}

// NewErrorResponse creates a new Response for error cases.
func NewErrorResponse(code, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

// NewSuccessResponse creates a new Response for successful cases, with optional data.
func NewSuccessResponse(message string, data interface{}) *Response {
	return &Response{
		Code:    "SUCCESS",
		Message: message,
		Data:    data,
	}
}
