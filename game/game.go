package game

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	"github.com/labstack/echo"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/types"
)

// ErrInvalidRowCols is returned when given a zero or negative row or column count
var ErrInvalidRowCols = errors.New("invalid row or column count")

// ErrTooManyMines is returned when the amount of mines is greater or equal than the board
var ErrTooManyMines = errors.New("too many mines")

// ErrNoneMines is returned when the amount of mines is negative or zero
var ErrNoneMines = errors.New("none mines")

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

// ProspectGame contains all the information needed to build a new game
type ProspectGame struct {
	ID      int64 `json:"id"`
	Rows    int   `json:"rows" validate:"required,gt=0,lte=100"`
	Cols    int   `json:"cols" validate:"required,gt=0,lte=100"`
	Mines   int   `json:"mines" validate:"required,gt=0,lte=100"`
	Private bool  `json:"private" validate:"required"`
}

// Create is the http handler that creates a game
func (h *Handler) Create(c echo.Context) error {
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
	pointsCount := pGame.Rows * pGame.Cols
	if pointsCount <= pGame.Mines {
		h.logger.Printf("validation error too many mines\n")
		return response.NewBadRequestResponse(c, "too many mines")
	}
	ctx := c.Request().Context()
	err = h.CreateGame(ctx, user, &pGame, h.ArrayStorageStrategy)
	if err != nil {
		return response.NewResponseFromError(c, err)
	}
	return response.NewSuccessResponse(c, pGame)
}

// end of http handlers

type boardPoint struct {
	row int
	col int
}

type board struct {
	rows  int
	cols  int
	mines int
	board [][]int
}

type boardStorageStrategy func(context.Context, security.JWTUser, *models.Game) error

// board storage strategy to benchmark

// ArrayStorageStrategy stores a Game board as an array inside the game
func (h *Handler) ArrayStorageStrategy(ctx context.Context, user security.JWTUser, game *models.Game) error {
	return game.Insert(ctx, h.db, boil.Infer())
}

// TableStorageStrategy stores a Game board in another table
func (h *Handler) TableStorageStrategy(ctx context.Context, user security.JWTUser, game *models.Game) error {
	flatBoard := game.Map
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		h.logger.Printf("error beggining transaction: %v\n", err)
		return err
	}
	// do not insert map
	err = game.Insert(ctx, tx, boil.Whitelist("private", "cols", "rows", "mines", "creator_id"))
	if err != nil {
		h.logger.Printf("error inserting game: %v. Rolling back game insertion\n", err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			h.logger.Printf("error rolling back game creation with error: %v\n", rollbackError)
		}
		return err
	}
	for i := range flatBoard {
		row, col := arrayToBoardPoint(i, int(game.Cols))
		gbPoint := &models.GameBoardPoint{
			GameID:        game.ID,
			Row:           int16(row),
			Col:           int16(col),
			MineProximity: int16(flatBoard[i]),
		}
		err = gbPoint.Insert(ctx, tx, boil.Infer())
		if err != nil {
			h.logger.Printf("error inserting game board point: %v. Rolling back operation\n", err)
			// just log rollback error
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				h.logger.Printf("error rolling back game board point creation with error: %v\n", rollbackError)
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		h.logger.Printf("error commiting transaction: %v. Rolling back operation\n", err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			h.logger.Printf("error rolling back operation with error: %v\n", rollbackError)
		}
		return err
	}
	return nil
}

// CreateGame creates a random board game and stores a new game in the database
func (h *Handler) CreateGame(ctx context.Context, user security.JWTUser, pGame *ProspectGame, storageStrategy boardStorageStrategy) error {
	board, err := NewBoard(pGame.Rows, pGame.Cols, pGame.Mines)
	if err != nil {
		return response.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	dbBoard := toDBBoard(pGame.Rows, pGame.Cols, board)
	game := &models.Game{
		Rows:      int16(pGame.Rows),
		Cols:      int16(pGame.Cols),
		Mines:     null.Int16From(int16(pGame.Mines)),
		CreatorID: user.ID,
		Private:   pGame.Private,
		Map:       dbBoard,
	}
	err = storageStrategy(ctx, user, game)
	if err != nil {
		return err
	}
	pGame.ID = game.ID
	return nil
}

// NewBoard creates a random minesweeper board
func NewBoard(rows, cols, mines int) ([][]int, error) {
	var initializedBoard [][]int
	if rows <= 0 || cols <= 0 {
		return initializedBoard, ErrInvalidRowCols
	}
	if mines >= (rows * cols) {
		return initializedBoard, ErrTooManyMines
	}
	if mines <= 0 {
		return initializedBoard, ErrNoneMines
	}
	b := &board{
		rows:  rows,
		cols:  cols,
		mines: mines,
		board: make([][]int, rows),
	}

	boardCartesian := make([]boardPoint, rows*cols)
	for i := range b.board {
		b.board[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			//initialize with -1
			b.board[i][j] = -1
			boardCartesian[(i*rows)+j] = boardPoint{
				row: i,
				col: j,
			}
		}
	}
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	for mines > 0 {
		mineIndex := random.Intn(len(boardCartesian) - 1)
		p := boardCartesian[mineIndex]
		b.placeMine(p.row, p.col)
		boardCartesian = append(boardCartesian[:mineIndex], boardCartesian[mineIndex+1:]...)
		mines--
	}
	return b.board, nil
}

func (b *board) placeMine(row, col int) {
	b.board[row][col] = -10
	for _, s := range b.siblingPoints(row, col) {
		if b.board[s.row][s.col] > -9 {
			b.board[s.row][s.col] = b.board[s.row][s.col] - 1
		}
	}
}

func (b *board) siblingPoints(row, col int) []boardPoint {
	rowPlaces := make([]int, 0, 3)
	colPlaces := make([]int, 0, 3)
	points := make([]boardPoint, 0, 8)
	if row > 0 {
		rowPlaces = append(rowPlaces, row-1)
	}
	rowPlaces = append(rowPlaces, row)
	if row < b.rows-1 {
		rowPlaces = append(rowPlaces, row+1)
	}
	if col > 0 {
		colPlaces = append(colPlaces, col-1)
	}
	colPlaces = append(colPlaces, col)
	if col < b.cols-1 {
		colPlaces = append(colPlaces, col+1)
	}
	for i := range rowPlaces {
		for j := range colPlaces {
			if rowPlaces[i] != row || colPlaces[j] != col {
				points = append(points, boardPoint{
					row: rowPlaces[i],
					col: colPlaces[j],
				})
			}
		}
	}
	return points
}

func toDBBoard(rows, cols int, board [][]int) types.Int64Array {
	dbBoard := make(types.Int64Array, rows*cols)
	for i := range board {
		for j := range board[i] {
			index := boardToArrayPoint(i, j, cols)
			dbBoard[index] = int64(board[i][j])
		}
	}
	return dbBoard
}

func boardToArrayPoint(row, col, colLength int) int {
	return (row * colLength) + col
}

func arrayToBoardPoint(index, colLength int) (int, int) {
	row := int(index / colLength)
	col := index - (row * colLength)
	return row, col
}
