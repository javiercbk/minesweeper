package player

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"

	extErrors "github.com/pkg/errors"
)

// uniqueNameConstaintName is the unique name db constraint
const uniqueNameConstaintName = "idx_players_name"

// BCryptCost is the ammount of iterations applied to bcrypt
const BCryptCost = 12

// HashPassword hashes a password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BCryptCost)
	return string(hash[0:]), err
}

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
func (h *Handler) Routes(e *echo.Group, jwtMiddleware echo.MiddlewareFunc) {
	e.POST("/", h.Create)
	e.GET("/current", h.RetrieveCurrent, jwtMiddleware)
}

// ProspectPlayer contain all the information needed to build a player
type ProspectPlayer struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name" validate:"required,gt=0"`
	Password string `json:"password,omitempty" validate:"required,gt=0"`
}

// Create is the http handler for player creation
func (h *Handler) Create(c echo.Context) error {
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
	err = h.CreatePlayer(ctx, &pPlayer)
	if err != nil {
		return response.NewResponseFromError(c, err)
	}
	// remove the password before re sending it to the client
	pPlayer.Password = ""
	return response.NewSuccessResponse(c, pPlayer)
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

// end of http handlers

// CreatePlayer creates a player in the database
func (h *Handler) CreatePlayer(ctx context.Context, pPlayer *ProspectPlayer) error {
	hashPassword, err := HashPassword(pPlayer.Password)
	if err != nil {
		h.logger.Printf("error hashing password: %v\n", err)
		return errors.New("error hashing password")
	}
	player := &models.Player{
		Name:     pPlayer.Name,
		Password: hashPassword,
	}
	err = player.Insert(ctx, h.db, boil.Infer())
	if err != nil {
		cause := extErrors.Cause(err)
		if pgerr, ok := cause.(*pq.Error); ok {
			if pgerr.Constraint == uniqueNameConstaintName {
				return response.HTTPError{
					Code:    http.StatusConflict,
					Message: fmt.Sprintf("player %s already exists", pPlayer.Name),
				}
			}
		} else {
			h.logger.Printf("error inserting player: %v\n", err)
			return errors.New("error inserting player")
		}
	}
	pPlayer.ID = player.ID
	return nil
}
