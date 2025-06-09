// would be moved to a common package
package upwardli

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code and structured response
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"` // Don't include in JSON response
}

// Error implements the error interface
func (e AppError) Error() string {
	return e.Message
}

// WithMessage creates a new AppError with a custom message
func (e AppError) WithMessage(msg string) AppError {
	return AppError{
		Code:    e.Code,
		Message: msg,
		Status:  e.Status,
	}
}

// WithMessagef creates a new AppError with a formatted message
func (e AppError) WithMessagef(format string, args ...interface{}) AppError {
	return AppError{
		Code:    e.Code,
		Message: fmt.Sprintf(format, args...),
		Status:  e.Status,
	}
}

// WriteError writes an error response with proper HTTP status codes and JSON formatting
// If status is 0, it will use the error's status code or default to 500
func WriteError(w http.ResponseWriter, err error, status ...int) {
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var errorResponse AppError

	// Determine status code
	if len(status) > 0 && status[0] != 0 {
		// Use provided status code
		statusCode = status[0]
	} else if appErr, ok := err.(AppError); ok {
		// Use AppError's status code
		statusCode = appErr.Status
		errorResponse = appErr
	} else {
		// Default to 500 for unknown errors
		statusCode = http.StatusInternalServerError
	}

	// Create error response if not already an AppError
	if errorResponse.Code == "" {
		if appErr, ok := err.(AppError); ok {
			errorResponse = appErr
		} else {
			errorResponse = AppError{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			}
		}
	}

	// Override status if explicitly provided
	if len(status) > 0 && status[0] != 0 {
		errorResponse.Status = status[0]
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// WriteJSON writes a successful JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// ReadJSON reads and validates JSON from request body
func ReadJSON(r *http.Request, dest interface{}) error {
	if r.Body == nil {
		return AppError{
			Code:    "INVALID_INPUT",
			Message: "request body is empty",
			Status:  http.StatusBadRequest,
		}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict JSON parsing

	if err := decoder.Decode(dest); err != nil {
		return AppError{
			Code:    "INVALID_INPUT",
			Message: "invalid JSON format",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
