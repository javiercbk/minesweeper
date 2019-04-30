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
	"github.com/javiercbk/minesweeper/models"
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

// API is the player API
type API interface {
	CreatePlayer(ctx context.Context, pPlayer *ProspectPlayer) error
}

type api struct {
	logger *log.Logger
	db     *sql.DB
}

// NewAPI creates a new player API
func NewAPI(logger *log.Logger, db *sql.DB) API {
	return api{
		logger: logger,
		db:     db,
	}
}

// ProspectPlayer contain all the information needed to build a player
type ProspectPlayer struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name" validate:"required,gt=0"`
	Password string `json:"password,omitempty" validate:"required,gt=0"`
}

// CreatePlayer creates a player in the database
func (api api) CreatePlayer(ctx context.Context, pPlayer *ProspectPlayer) error {
	hashPassword, err := HashPassword(pPlayer.Password)
	if err != nil {
		api.logger.Printf("error hashing password: %v\n", err)
		return errors.New("error hashing password")
	}
	player := &models.Player{
		Name:     pPlayer.Name,
		Password: hashPassword,
	}
	err = player.Insert(ctx, api.db, boil.Infer())
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
			api.logger.Printf("error inserting player: %v\n", err)
			return errors.New("error inserting player")
		}
	}
	pPlayer.ID = player.ID
	return nil
}
