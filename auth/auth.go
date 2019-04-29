package auth

import (
	"database/sql"
	"log"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/labstack/echo"
)

// apiFactory is a function that creates an auth API. It is stored on a var so any test can mock the API
var apiFactory = NewAPI

// Handler is a group of handlers within a route.
type Handler struct {
	logger *log.Logger
	db     *sql.DB
}

// NewHandler creates a handler for the game route
func NewHandler(logger *log.Logger, db *sql.DB) Handler {
	return Handler{
		logger: logger,
		db:     db,
	}
}

// Routes initializes all the routes with their http handlers
func (h Handler) Routes(e *echo.Group, jwtSecret string) {
	e.POST("/", h.AuthenticateFactory(jwtSecret))
}

// AuthenticateFactory creates the http handler for the login
func (h Handler) AuthenticateFactory(jwtSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		auth := Credentials{}
		err := c.Bind(&auth)
		if err != nil {
			h.logger.Printf("could not bind request data%v\n", err)
			return response.NewBadRequestResponse(c, "name and passwords are required")
		}
		if err = c.Validate(auth); err != nil {
			h.logger.Printf("validation error %v\n", err)
			return response.NewBadRequestResponse(c, err.Error())
		}
		api := apiFactory(h.logger, h.db)
		tResponse, err := api.CreateToken(ctx, jwtSecret, auth)
		if err != nil {
			return response.NewResponseFromError(c, err)
		}
		return response.NewSuccessResponse(c, tResponse)
	}
}
