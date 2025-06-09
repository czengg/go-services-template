package client

type ClientError struct {
	Code    string
	Message string
	Details string
}

func (e *ClientError) Error() string {
	return e.Message
}

var (
	ErrInvalidBaseURL = &ClientError{
		Code:    "INVALID_BASE_URL",
		Message: "invalid base url",
		Details: "base url cannot be empty",
	}
)
