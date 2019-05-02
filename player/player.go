package player

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/labstack/echo"
)

// apiFactory is a function that creates a player API. It is stored on a var so any test can mock the API
var apiFactory = NewAPI

// Handler is a group of handlers within a route.
type Handler struct {
	logger *log.Logger
	db     *sql.DB
}

type pResponse struct {
	Player ProspectPlayer `json:"player"`
}

type uResponse struct {
	User security.JWTUser `json:"user"`
}

// NewHandler creates a handler for the game route
func NewHandler(logger *log.Logger, db *sql.DB) Handler {
	return Handler{
		logger: logger,
		db:     db,
	}
}

// Routes initializes all the routes with their http handlers
func (h Handler) Routes(e *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
	e.POST("", h.Create)
	e.GET("/current", h.RetrieveCurrent, jwtMiddleware)
}

// Create is the http handler for player creation
func (h Handler) Create(c echo.Context) error {
	pPlayer := ProspectPlayer{}
	err := c.Bind(&pPlayer)
	if err != nil {
		h.logger.Printf("could not bind request data%v\n", err)
		return response.NewBadRequestResponse(c, "name and passwords are required")
	}
	if err = c.Validate(pPlayer); err != nil {
		h.logger.Printf("validation error %v\n", err)
		return response.NewBadRequestResponse(c, err.Error())
	}
	ctx := c.Request().Context()
	api := apiFactory(h.logger, h.db)
	err = api.CreatePlayer(ctx, &pPlayer)
	if err != nil {
		return response.NewResponseFromError(c, err)
	}
	// remove the password before re sending it to the client
	pPlayer.Password = ""
	return response.NewSuccessResponse(c, pResponse{pPlayer})
}

// RetrieveCurrent is the http handler that retrieves the authenticated user
func (h Handler) RetrieveCurrent(c echo.Context) error {
	user, err := security.JWTDecode(c)
	if err == security.ErrUserNotFound {
		h.logger.Printf("error finding jwt token in context: %v\n", err)
		return response.NewErrorResponse(c, http.StatusForbidden, "authentication token was not found")
	}
	return response.NewSuccessResponse(c, uResponse{user})
}
