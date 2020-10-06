package errortools

import (
	"net/http"
)

var (
	// APIErrorInterface implementing the error interface
	APIErrorInterface = &APIError{}
)

// APIError to standardize the API errors
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewInternalServerError type
func (ae *APIError) NewInternalServerError(message string) *APIError {
	return &APIError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

// NewBadRequestError type
func (ae *APIError) NewBadRequestError(message string) *APIError {
	return &APIError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

// NewNotFoundError type
func (ae *APIError) NewNotFoundError(message string) *APIError {
	return &APIError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

// Error implementing the error interface
func (ae *APIError) Error() string {
	return "error"
}
