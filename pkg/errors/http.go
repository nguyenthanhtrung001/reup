package errors

import "net/http"

type HTTPError struct {
	Code       int
	Message    string
	StatusCode int
}

// NewHTTPError returns a new HTTPError with the given code and message.
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

func NewBadRequestHTTPError() *HTTPError {
	return &HTTPError{
		Code:       400,
		Message:    "Bad Request",
		StatusCode: http.StatusBadRequest,
	}
}

// NewHTTPError returns a new HTTPError with the given code and message.
func NewUnauthorizedHTTPError() *HTTPError {
	return &HTTPError{
		Code:       401,
		Message:    "Unauthorized",
		StatusCode: http.StatusUnauthorized,
	}
}

// NewHTTPError returns a new HTTPError with the given code and message.
func NewSystemUnderMaintenanceHTTPError() *HTTPError {
	return &HTTPError{
		Code:       503,
		Message:    "System is currently undergoing maintenance. Please try again later.",
		StatusCode: http.StatusServiceUnavailable,
	}
}

// NewHTTPError returns a new HTTPError with the given code and message.
func NewPermissionDeniedHTTPError() *HTTPError {
	return &HTTPError{
		Code:       403,
		Message:    "Permission Denied",
		StatusCode: http.StatusForbidden,
	}
}

// Error returns the error message.
func (e HTTPError) Error() string {
	return e.Message
}
