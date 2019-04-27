package game

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/labstack/echo"
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

// Create is the http handler that creates a game
func (h *Handler) Create(c echo.Context) error {
	return response.NewNotFoundResponse(c)
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

// CreateBoard creates a new minesweeper board
func CreateBoard(rows, cols, mines int) ([][]int, error) {
	var board [][]int

	return board, nil
}
