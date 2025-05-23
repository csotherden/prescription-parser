package models

import "fmt"

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// ValidationError represents a model validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface for ValidationError
func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", v.Field, v.Message)
}
