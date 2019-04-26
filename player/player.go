package player

import (
	"log"
	"net/http"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/labstack/echo"
)

// Handler is a group of handlers within a route.
type Handler struct {
	logger *log.Logger
}

// NewHandler creates a handler for the game route
func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Routes initializes all the routes with their http handlers
func (h *Handler) Routes(e *echo.Group) {
	e.GET("/current", h.RetrieveCurrent)
}

// RetrieveCurrent is the http handler that retrieves the authenticated user
func (h *Handler) RetrieveCurrent(c echo.Context) error {
	user, err := security.JWTDecode(c)
	if err == security.ErrUserNotFound {
		h.logger.Printf("error finding jwt token in context: %v\n", err)
		return response.NewErrorResponse(c, http.StatusForbidden, "authentication token was not found")
	}
	return response.NewSuccessResponse(c, user)
}
