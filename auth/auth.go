package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	"github.com/labstack/echo"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

// dummyHash is a static bcrypt
const dummyHash = "$2y$12$xJdgVt1Siwdy456cGBvY5.tlAIyorJmwSwdwMKbKtxlgmMwU2Aju2"

// ErrBadCredentials is returned when incorrect credentials are provided
var ErrBadCredentials = response.HTTPError{
	Code:    http.StatusUnauthorized,
	Message: "user name or password is incorrect",
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
func (h *Handler) Routes(e *echo.Group, jwtSecret string) {
	e.POST("/", h.AuthenticateFactory(jwtSecret))

}

// Credentials is the username and password combination
type Credentials struct {
	Name     string `json:"name" validate:"required,gt=0"`
	Password string `json:"password,omitempty" validate:"required,gt=0"`
}

// TokenResponse contains a jwt token
type TokenResponse struct {
	Token string `json:"token"`
}

// AuthenticateFactory creates the http handler for the login
func (h *Handler) AuthenticateFactory(jwtSecret string) echo.HandlerFunc {
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
		tResponse, err := h.CreateToken(ctx, jwtSecret, auth)
		if err != nil {
			return response.NewResponseFromError(c, err)
		}
		return response.NewSuccessResponse(c, tResponse)
	}
}

// end of http handlers

// CreateToken creates an authentication token that can be used to authenticate with the rest api
func (h *Handler) CreateToken(ctx context.Context, jwtSecret string, auth Credentials) (TokenResponse, error) {
	tResponse := TokenResponse{}
	player, err := models.Players(qm.Where("name = ?", auth.Name)).One(ctx, h.db)
	if err != nil && err != sql.ErrNoRows {
		h.logger.Printf("error searching for player %v\n", err)
		return tResponse, errors.New("error searching for player")
	}
	if player == nil {
		// bcrypt comparing is a slow process. If the user does not exist in the database and we reply tight away
		// an attacker might notice the request latency between a user not existing and a password being incorrect.
		// That way the attacker can brute force the API and guess user names.
		// To mitigate this risk, I perform a bcrypt comparison but I discard the result because I only want the request
		// latency to be incremented. I cannot simply time.Sleep() because the bcrypt time varies between CPUs.
		// We also need to protect the service from a DoS, but that would be a job for some other proxy server.
		bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(auth.Password))
		return tResponse, ErrBadCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(auth.Password))
	if err != nil {
		return tResponse, ErrBadCredentials
	}
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = security.JWTEncode(security.JWTUser{
		ID:   player.ID,
		Name: player.Name,
	}, time.Minute*20)

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		h.logger.Printf("error signing token %v\n", err)
		return tResponse, errors.New("error creating token")
	}
	tResponse.Token = t
	return tResponse, nil
}
