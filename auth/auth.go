package auth

import (
	"database/sql"
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

type authentication struct {
	Name     string `json:"name" validate:"required,gt=0"`
	Password string `json:"password,omitempty" validate:"required,gt=0"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

// AuthenticateFactory creates the http handler for the login
func (h *Handler) AuthenticateFactory(jwtSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		authentication := authentication{}
		err := c.Bind(&authentication)
		if err != nil {
			h.logger.Printf("bad request from client %v\n", err)
			return response.NewBadRequestResponse(c, "name and passwords are required")
		}
		player, err := models.Players(qm.Where("name = ?", authentication.Name)).One(ctx, h.db)
		if err != nil {
			h.logger.Printf("error searching for player %v\n", err)
			return response.NewInternalErrorResponse(c, "error searching for player")
		}
		if player == nil {
			// bcrypt comparing is a slow process. If the user does not exist in the database and we reply tight away
			// an attacker might notice the request latency between a user not existing and a password being incorrect.
			// That way the attacker can brute force the API and guess user names.
			// To mitigate this risk, I perform a bcrypt comparison but I discard the result because I only want the request
			// latency to be incremented. I cannot simply time.Sleep() because the bcrypt time varies between CPUs.
			// We also need to protect the service from a DoS, but that would be a job for some other proxy server.
			bcrypt.CompareHashAndPassword([]byte(authentication.Password), []byte(dummyHash))
			return response.NewErrorResponse(c, http.StatusUnauthorized, "user name or password is incorrect")
		}
		bcrypt.CompareHashAndPassword([]byte(authentication.Password), []byte(player.Password))
		if err != nil {
			return response.NewErrorResponse(c, http.StatusUnauthorized, "user name or password is incorrect")
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
			return response.NewInternalErrorResponse(c, "error creating token")
		}
		return response.NewSuccessResponse(c, tokenResponse{t})
	}
}
