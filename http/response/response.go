package response

import (
	"net/http"

	"github.com/labstack/echo"
)

const successMessage = "success"

// filled by compiler flag -X http.response.minesweeperVersion=value
var minesweeperVersion string

// HTTPError is an error with an HTTP code
type HTTPError struct {
	Code    int
	Message string
}

// Error returns the error message of an HTTPError
func (e HTTPError) Error() string {
	if e.Message == "" {
		return http.StatusText(e.Code)
	}
	return e.Message
}

// NewHTTPError extracts an HTTP error code and a message from an error
func NewHTTPError(err error) (int, string) {
	if httpError, ok := err.(HTTPError); ok {
		return httpError.Code, httpError.Error()
	}
	return http.StatusInternalServerError, err.Error()
}

// Status is the status of the response
type Status struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// ServiceResponse is a generic service response
type ServiceResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// NewSuccessResponseWithCode sends a successful response with code
func NewSuccessResponseWithCode(c echo.Context, code int, data interface{}) error {
	resp := ServiceResponse{
		Status: Status{
			Error:   false,
			Code:    code,
			Message: successMessage,
			Version: minesweeperVersion,
		},
	}
	if data != nil {
		resp.Data = data
	}
	return c.JSON(code, resp)
}

// NewSuccessResponse sends a successful response
func NewSuccessResponse(c echo.Context, data interface{}) error {
	return NewSuccessResponseWithCode(c, http.StatusOK, data)
}

// NewErrorResponse sends an error response
func NewErrorResponse(c echo.Context, code int, message string) error {
	resp := ServiceResponse{
		Status: Status{
			Error:   true,
			Code:    code,
			Message: message,
			Version: minesweeperVersion,
		},
	}
	if resp.Status.Message == "" {
		resp.Status.Message = http.StatusText(code)
	}
	return c.JSON(code, resp)
}

// NewInternalErrorResponse sends an internal server error response
func NewInternalErrorResponse(c echo.Context, message string) error {
	return NewErrorResponse(c, http.StatusInternalServerError, message)
}

// NewNotFoundResponse sends a not found response
func NewNotFoundResponse(c echo.Context) error {
	return NewErrorResponse(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// NewBadRequestResponse sends a bad response with a reason
func NewBadRequestResponse(c echo.Context, message string) error {
	return NewErrorResponse(c, http.StatusBadRequest, message)
}

// NewResponseFromError sends an error response from an Error
func NewResponseFromError(c echo.Context, err error) error {
	code, message := NewHTTPError(err)
	return NewErrorResponse(c, code, message)
}
