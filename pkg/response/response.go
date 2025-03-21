package response

import (
	"net/http"

	pkgErrors "book-store/pkg/errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// DefaultErrorMessage is the default error message.
	DefaultErrorMessage = "Something went wrong"
	// ValidationErrorCode is the validation error code.
	ValidationErrorCode = 400
	// ValidationErrorMsg is the validation error message.
	ValidationErrorMsg = "Validation error"
)

// Resp is the response format.
type Resp struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Errors    any    `json:"errors,omitempty"`
}

// NewOKResp returns a new OK response with the given data.
func NewOKResp(data any) Resp {
	return Resp{
		ErrorCode: 0,
		Message:   "Success",
		Data:      data,
	}
}

// Ok returns a new OK response with the given data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, NewOKResp(data))
}

// Unauthorized returns a new Unauthorized response with the given data.
func Unauthorized(c *gin.Context) {
	c.JSON(parseError(pkgErrors.NewUnauthorizedHTTPError()))
}

// Permission Deny returns a new Unauthorized response with the given data.
func PermissionDenied(c *gin.Context) {
	c.JSON(parseError(pkgErrors.NewPermissionDeniedHTTPError()))
}

// Unauthorized returns a new Unauthorized response with the given data.
func SystemUnderMaintenance(c *gin.Context) {
	c.JSON(parseError(pkgErrors.NewSystemUnderMaintenanceHTTPError()))
}

func parseError(err error) (int, Resp) {
	switch parsedErr := err.(type) {
	case *pkgErrors.ValidationErrorCollector:
		return http.StatusBadRequest, Resp{
			ErrorCode: ValidationErrorCode,
			Message:   ValidationErrorMsg,
			Errors:    parsedErr.Errors(),
		}
	case *pkgErrors.HTTPError:
		statusCode := parsedErr.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
		}

		return statusCode, Resp{
			ErrorCode: parsedErr.Code,
			Message:   parsedErr.Message,
		}
	default:
		return http.StatusInternalServerError, Resp{
			ErrorCode: 500,
			Message:   DefaultErrorMessage,
		}
	}
}

func adminParseError(err error) (int, Resp) {
	switch parsedErr := err.(type) {
	case *pkgErrors.ValidationErrorCollector:
		return http.StatusBadRequest, Resp{
			ErrorCode: ValidationErrorCode,
			Message:   ValidationErrorMsg,
			Errors:    parsedErr.Errors(),
		}
	case *pkgErrors.HTTPError:
		statusCode := parsedErr.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
		}

		return statusCode, Resp{
			ErrorCode: parsedErr.Code,
			Message:   parsedErr.Message,
		}
	default:
		return http.StatusInternalServerError, Resp{
			ErrorCode: 500,
			Message:   err.Error(),
		}
	}
}

// Error returns a new Error response with the given error.
func Error(c *gin.Context, err error) {
	c.JSON(parseError(err))
}

func AdminError(c *gin.Context, err error) {
	c.JSON(adminParseError(err))
}

// ErrorMapping is a map of error to HTTPError.
type ErrorMapping map[error]*pkgErrors.HTTPError

// ErrorWithMap returns a new Error response with the given error.
func ErrorWithMap(c *gin.Context, err error, eMap ErrorMapping) {
	if mongoErr, ok := err.(*mongo.CommandError); ok {
		// Handle mongo command errors separately
		// Map it to a specific HTTP error or log it, for instance
		Error(c, mongoErr)
		return
	}

	if httpErr, ok := eMap[err]; ok {
		Error(c, httpErr)
		return
	}

	AdminError(c, err)
}
