package response

import (
	"net/http"

	"github.com/labstack/echo"
)

const successMessage = "success"

// filled by compiler flag -X http.response.minesweeperVersion=value
var minesweeperVersion string

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

// NewNotFoundResponse sends a not found response
func NewNotFoundResponse(c echo.Context) error {
	code := http.StatusNotFound
	resp := ServiceResponse{
		Status: Status{
			Error:   false,
			Code:    code,
			Message: http.StatusText(http.StatusNotFound),
			Version: minesweeperVersion,
		},
	}
	return c.JSON(code, resp)
}

// NewBadRequestResponse sends a bad response with a reason
func NewBadRequestResponse(c echo.Context, message string) error {
	code := http.StatusBadRequest
	resp := ServiceResponse{
		Status: Status{
			Error:   false,
			Code:    code,
			Message: message,
			Version: minesweeperVersion,
		},
	}
	return c.JSON(code, resp)
}
