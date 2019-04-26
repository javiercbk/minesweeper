package game

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/labstack/echo"
)

// Handler is a group of handlers within a route.
type Handler struct {
	logger *log.Logger
	db     *sql.DB
}

// NewHandler creates a handler for the game route
func NewHandler(logger *log.Logger, db *sql.DB) *Handler {
	return &Handler{
		logger: logger,
		db:     db,
	}
}

// Routes initializes all the routes with their http handlers
func (h *Handler) Routes(e *echo.Group) {
	e.GET("/", h.Find)
	e.GET("/:gameID", h.Retrieve)
	e.POST("/", h.Create)

}

// Find is the http handler searchs for all the public and players open games
func (h *Handler) Find(c echo.Context) error {
	return response.NewNotFoundResponse(c)
}

// Retrieve is the http handler searchs for a single game by ID
func (h *Handler) Retrieve(c echo.Context) error {
	gameIDStr := c.Param("gameID")
	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		return response.NewBadRequestResponse(c, fmt.Sprintf("game id %s is not a valid id", gameIDStr))
	}
	h.logger.Printf("%d", gameID)
	return response.NewNotFoundResponse(c)
}

// Create is the http handler that creates a game
func (h *Handler) Create(c echo.Context) error {
	return response.NewNotFoundResponse(c)
}
