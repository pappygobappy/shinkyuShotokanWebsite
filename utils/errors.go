package utils

import "fmt"

// AppError represents a standardized application error
type AppError struct {
	Code    int         // HTTP status code
	Message string      // User-friendly message
	Err     interface{} // Underlying cause (optional)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Helper to create AppError with underlying error
func NewAppError(code int, message string, err interface{}) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Standard error variables
var (
	ErrNotFound     = &AppError{Code: 404, Message: "Resource not found"}
	ErrUnauthorized = &AppError{Code: 401, Message: "Unauthorized"}
	ErrForbidden    = &AppError{Code: 403, Message: "Access denied"}
	ErrValidation   = &AppError{Code: 422, Message: "Validation failed"}
	ErrInternal     = &AppError{Code: 500, Message: "Internal server error"}
	ErrConflict     = &AppError{Code: 409, Message: "Resource already exists"}
)

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}
