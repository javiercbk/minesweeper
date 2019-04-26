package auth

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/javiercbk/minesweeper/http/security"
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
func (h *Handler) Routes(e *echo.Group, jwtSecret string) {
	e.POST("/", h.AuthenticateFactory(jwtSecret))

}

// AuthenticateFactory creates the http handler for the login
func (h *Handler) AuthenticateFactory(jwtSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		// Throws unauthorized error
		if username == "jon" && password == "shhh!" {
			return echo.ErrUnauthorized
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		token.Claims = security.JWTEncode(security.JWTUser{
			ID:   1,
			Name: "sample",
		}, time.Minute*20)

		t, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}
