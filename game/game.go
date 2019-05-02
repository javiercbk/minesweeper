package game

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	"github.com/labstack/echo"
)

// apiFactory is a function that creates a game API. It is stored on a var so any test can mock the API
var apiFactory = NewAPI

type gResponse struct {
	Game ProspectGame `json:"game"`
}

type cResponse struct {
	Confirmation OperationConfirmation `json:"confirmation"`
}

type gsResponse struct {
	Games models.GameSlice `json:"games"`
}

type sgResponse struct {
	Game StatefulGame `json:"game"`
}

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
func (h Handler) Routes(e *echo.Group) {
	e.GET("/", h.Find)
	e.GET("/:gameID", h.Retrieve)
	e.POST("/", h.Create)
	e.PATCH("/:gameID", h.Apply)

}

// Find is the http handler searchs for all the public and players open games
func (h Handler) Find(c echo.Context) error {
	user, err := security.JWTDecode(c)
	if err == security.ErrUserNotFound {
		h.logger.Printf("error finding jwt token in context: %v\n", err)
		return response.NewErrorResponse(c, http.StatusForbidden, "authentication token was not found")
	}
	ctx := c.Request().Context()
	api := apiFactory(h.logger, h.db)
	games, err := api.FindGames(ctx, user)
	if err != nil {
		return response.NewResponseFromError(c, err)
	}
	return response.NewSuccessResponse(c, gsResponse{games})
}

// Retrieve is the http handler searchs for a single game by ID
func (h Handler) Retrieve(c echo.Context) error {
	user, err := security.JWTDecode(c)
	if err == security.ErrUserNotFound {
		h.logger.Printf("error finding jwt token in context: %v\n", err)
		return response.NewErrorResponse(c, http.StatusForbidden, "authentication token was not found")
	}
	gameIDStr := c.Param("gameID")
	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		return response.HTTPError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("game %s does not exist", gameIDStr),
		}
	}
	ctx := c.Request().Context()
	api := apiFactory(h.logger, h.db)
	game, err := api.RetrieveGame(ctx, user, gameID)
	if err != nil {
		return response.NewResponseFromError(c, err)
	}
	return response.NewSuccessResponse(c, sgResponse{game})
}

// Create is the http handler that creates a game
func (h Handler) Create(c echo.Context) error {
	user, err := security.JWTDecode(c)
	if err == security.ErrUserNotFound {
		h.logger.Printf("error finding jwt token in context: %v\n", err)
		return response.NewErrorResponse(c, http.StatusForbidden, "authentication token was not found")
	}
	pGame := ProspectGame{}
	err = c.Bind(&pGame)
	if err != nil {
		h.logger.Printf("could not bind request data%v\n", err)
		return response.NewBadRequestResponse(c, "rows, cols, mines and private are required")
	}
	if err = c.Validate(pGame); err != nil {
		h.logger.Printf("validation error %v\n", err)
		return response.NewBadRequestResponse(c, err.Error())
	}
	ctx := c.Request().Context()
	api := apiFactory(h.logger, h.db)
	err = api.CreateGame(ctx, user, &pGame)
	if err != nil {
		return response.NewResponseFromError(c, err)
	}
	return response.NewSuccessResponse(c, gResponse{pGame})
}

// Apply is the http handler that
func (h Handler) Apply(c echo.Context) error {
	user, err := security.JWTDecode(c)
	if err == security.ErrUserNotFound {
		h.logger.Printf("error finding jwt token in context: %v\n", err)
		return response.NewErrorResponse(c, http.StatusForbidden, "authentication token was not found")
	}
	oper := Operation{}
	err = c.Bind(&oper)
	if err != nil {
		h.logger.Printf("could not bind request data%v\n", err)
		return response.NewBadRequestResponse(c, "id, gameId, op, row, col are required")
	}
	if err = c.Validate(oper); err != nil {
		h.logger.Printf("validation error %v\n", err)
		return response.NewBadRequestResponse(c, err.Error())
	}
	ctx := c.Request().Context()
	api := apiFactory(h.logger, h.db)
	confirmation, err := api.ApplyOperation(ctx, user, oper)
	if err != nil {
		return response.NewResponseFromError(c, err)
	}
	return response.NewSuccessResponse(c, cResponse{confirmation})
}
