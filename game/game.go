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
	"strings"
	"time"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	"github.com/labstack/echo"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// ErrInvalidRowCols is returned when given a zero or negative row or column count
var ErrInvalidRowCols = errors.New("invalid row or column count")

// ErrTooManyMines is returned when the amount of mines is greater or equal than the board
var ErrTooManyMines = errors.New("too many mines")

// ErrNoneMines is returned when the amount of mines is negative or zero
var ErrNoneMines = errors.New("none mines")

// ErrGameNotExists is returned when attempting to make an operation with a game that does not exists
var ErrGameNotExists = errors.New("the game does not exists")

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
	err = h.CreateGame(ctx, user, &pGame)
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

// CreateGame creates a random board game and stores a new game in the database
func (h *Handler) CreateGame(ctx context.Context, user security.JWTUser, pGame *ProspectGame) error {
	board, err := NewBoard(pGame.Rows, pGame.Cols, pGame.Mines)
	if err != nil {
		return response.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	game := &models.Game{
		Rows:      int16(pGame.Rows),
		Cols:      int16(pGame.Cols),
		Mines:     null.Int16From(int16(pGame.Mines)),
		CreatorID: user.ID,
		Private:   pGame.Private,
	}
	err = h.storeGameBoard(ctx, user, game, board)
	if err != nil {
		return err
	}
	pGame.ID = game.ID
	return nil
}

// storeGameBoard stores a Game board in the game_board_points table
func (h *Handler) storeGameBoard(ctx context.Context, user security.JWTUser, game *models.Game, board [][]int) error {
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
	var bigInsert strings.Builder
	fmt.Fprintf(&bigInsert, "INSERT INTO game_board_points (game_id, row, col, mine_proximity, created_at) VALUES ")
	first := true
	creationDateStr := time.Now().UTC().Format(time.RFC3339)
	for row := range board {
		for col := range board[row] {
			if first {
				first = false
			} else {
				bigInsert.WriteString(",")
			}
			fmt.Fprintf(&bigInsert, "(%d, %d, %d, %d, '%s')", game.ID, row, col, board[row][col], creationDateStr)
		}
	}
	bigInsert.WriteString(";")
	_, err = queries.Raw(bigInsert.String()).ExecContext(ctx, tx)
	if err != nil {
		h.logger.Printf("error inserting all game board points for game %d: %s\n", game.ID, err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			h.logger.Printf("error rolling back game creation with error: %v\n", rollbackError)
		}
		return err
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

func (h *Handler) tableRowColRetrieval(ctx context.Context, user security.JWTUser, gameID int64, row, col int) (int, error) {
	gameBoardPoint, err := models.GameBoardPoints(
		qm.Select("mine_proximity"),
		qm.InnerJoin("games g on g.id = game_board_points.game_id"),
		qm.Where("game_board_points.game_id = ? AND game_board_points.row = ? AND game_board_points.col = ? AND (g.creator_id = ? OR g.private = false)", gameID, row, col, user.ID),
	).One(ctx, h.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidRowCols
		}
		h.logger.Printf("error retrieving game board point: %v\n", err)
		return 0, err
	}
	return int(gameBoardPoint.MineProximity), nil
}

func (h *Handler) tableUpdateRowCol(ctx context.Context, user security.JWTUser, gameID int64, row, col, mineProximity int) error {
	aff, err := models.GameBoardPoints(
		qm.Where("game_id = ? AND row = ? AND col = ?", gameID, row, col),
	).UpdateAll(ctx, h.db, models.M{
		"mine_proximity": mineProximity,
	})
	if err != nil {
		return err
	}
	if aff != 1 {
		return fmt.Errorf("invalid row count %d when updating a game mine proximity", aff)
	}
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
