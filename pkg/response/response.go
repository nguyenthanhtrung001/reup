package response

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"

	pkgErrors "reup/pkg/errors"

	"reup/pkg/telegram"

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

	defaultStackTraceDepth = 32
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
	c.JSON(parseError(pkgErrors.NewUnauthorizedHTTPError(), c, nil, nil))
}

// Permission Deny returns a new Unauthorized response with the given data.
func PermissionDenied(c *gin.Context) {
	c.JSON(parseError(pkgErrors.NewPermissionDeniedHTTPError(), c, nil, nil))
}

// Unauthorized returns a new Unauthorized response with the given data.
func SystemUnderMaintenance(c *gin.Context) {
	c.JSON(parseError(pkgErrors.NewSystemUnderMaintenanceHTTPError(), c, nil, nil))
}

func PanicError(c *gin.Context, err any, t telegram.Telegram, tIDs int64) {
	if err == nil {
		c.JSON(parseError(nil, c, &t, &tIDs))
	} else {
		c.JSON(parseError(err.(error), c, &t, &tIDs))
	}
}

func parseError(err error, c *gin.Context, t *telegram.Telegram, chatID *int64) (int, Resp) {
	//print error . type
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
		stackTrace := captureStackTrace()
		sendServerTelegramMessageAsync(buildInternalServerErrorDataForReportBug(err.Error(), stackTrace, c), c, *t, *chatID)

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
	c.JSON(parseError(err, c, nil, nil))
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
func captureStackTrace() []string {
	var pcs [defaultStackTraceDepth]uintptr
	n := runtime.Callers(2, pcs[:])
	if n == 0 {
		return nil
	}

	var stackTrace []string
	for _, pc := range pcs[:n] {
		f := runtime.FuncForPC(pc)
		if f != nil {
			file, line := f.FileLine(pc)
			stackTrace = append(stackTrace, fmt.Sprintf("%s:%d %s", file, line, f.Name()))
		}
	}

	return stackTrace
}

func buildInternalServerErrorDataForReportBug(errString string, backtrace []string, c *gin.Context) string {
	url := c.Request.URL.String()
	method := c.Request.Method
	params := c.Request.URL.Query().Encode()

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	body := string(bodyBytes)

	headers := c.Request.Header
	var headersBuilder strings.Builder
	for key, values := range headers {
		headersBuilder.WriteString(key + ": " + strings.Join(values, ", ") + "\n")
	}
	headersString := headersBuilder.String()

	bk := ""
	for i, line := range backtrace {
		bk += fmt.Sprintf("[%d]: %s\n", i, line)
	}

	data := "YTB LIKE SERVICE ERROR\n" +
		"Route: " + url + "\n" +
		"Method: " + method + "\n" +
		"Params: " + params + "\n" +
		"Body: " + body + "\n" +
		"Headers:\n" + headersString + "\n" +
		"Error: " + errString + "\n\n" +
		"Backtrace:\n" + bk

	return data
}
