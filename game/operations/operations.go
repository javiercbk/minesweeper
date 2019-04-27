package operations

import (
	"database/sql"
	"log"

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
	}
}

// Routes initializes all the routes with their http handlers
func (h *Handler) Routes(e *echo.Group) {
	//TODO: add operation routes
}
