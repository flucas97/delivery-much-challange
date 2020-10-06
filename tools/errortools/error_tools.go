package errortools

import (
	"net/http"
)

// APIError struct pattern for handling API errors
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewInternalServerError label
func NewInternalServerError(message string) *APIError {
	return &APIError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

// NewBadRequestError label
func NewBadRequestError(message string) *APIError {
	return &APIError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

// NewNotFoundError label
func NewNotFoundError(message string) *APIError {
	return &APIError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}
