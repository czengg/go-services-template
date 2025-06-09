package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) WithMessage(msg string) AppError {
	return AppError{
		Code:    e.Code,
		Message: msg,
		Status:  e.Status,
	}
}

func (e AppError) WithMessagef(format string, args ...interface{}) AppError {
	return AppError{
		Code:    e.Code,
		Message: fmt.Sprintf(format, args...),
		Status:  e.Status,
	}
}

func WriteError(w http.ResponseWriter, err error, status ...int) {
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var errorResponse AppError

	if len(status) > 0 && status[0] != 0 {
		statusCode = status[0]
	} else if appErr, ok := err.(AppError); ok {
		statusCode = appErr.Status
		errorResponse = appErr
	} else {
		statusCode = http.StatusInternalServerError
	}

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

	if len(status) > 0 && status[0] != 0 {
		errorResponse.Status = status[0]
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

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
