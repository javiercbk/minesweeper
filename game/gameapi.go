package game

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/javiercbk/minesweeper/algebra"
	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

const (
	// StateNotRevealed is an integer sent to the client that means that the point in space is not revealed
	StateNotRevealed = iota
	// StateSuspectMine is an integer sent to the client that means that the point in space is marked as suspicious
	StateSuspectMine
	// StateMarkedMine is an integer sent to the client that means that the point in space is marked as a mine
	StateMarkedMine
	// StateRevealed is an integer sent to the client that means that the point is revealed
	StateRevealed
)

// ErrInvalidRowCols is returned when given a zero or negative row or column count
var ErrInvalidRowCols = errors.New("invalid row or column count")

// ErrTooManyMines is returned when the amount of mines is greater or equal than the board
var ErrTooManyMines = errors.New("too many mines")

// ErrNoneMines is returned when the amount of mines is negative or zero
var ErrNoneMines = errors.New("none mines")

// ErrGameNotExists is returned when attempting to make an operation with a game that does not exists
var ErrGameNotExists = errors.New("the game does not exists")

// ProspectGame contains all the information needed to build a new game
type ProspectGame struct {
	ID      int64 `json:"id"`
	Rows    int   `json:"rows" validate:"required,gte=0,lt=100"`
	Cols    int   `json:"cols" validate:"required,gte=0,lt=100"`
	Mines   int   `json:"mines" validate:"required,gt=0"`
	Private bool  `json:"private" validate:"required"`
}

// OperationResult is the result of an minesweeper algebra operation application
type OperationResult struct {
	Row           int `json:"row"`
	Col           int `json:"col"`
	MineProximity int `json:"mineProximity,omitempty"`
	PointState    int `json:"pointState"`
}

// Operation is an minesweper algebra operation
type Operation struct {
	ID      int64                 `json:"id" validate:"required"`
	GameID  int64                 `json:"gameId" validate:"required"`
	Op      algebra.OperationType `json:"op" validate:"required,1|2"`
	Row     int                   `json:"row" validate:"required,gte=0,lt=100"`
	Col     int                   `json:"col" validate:"required,gte=0,lt=100"`
	Applied bool                  `json:"applied"`
	Result  []OperationResult     `json:"result,omitempty"`
}

// Status is the game status
type Status struct {
	Won   bool    `json:"won"`
	Lost  bool    `json:"lost"`
	Board [][]int `json:"board,omitempty"`
}

// OperationConfirmation is the confirmation of an operation application
type OperationConfirmation struct {
	Operation       Operation   `json:"operation,omitempty"`
	DeltaOperations []Operation `json:"deltaOperations,omitempty"`
	Status          Status      `json:"status"`
}

// API is the game api
type API struct {
	logger *log.Logger
	db     *sql.DB
}

// NewAPI creates a new game API
func NewAPI(logger *log.Logger, db *sql.DB) API {
	return API{
		logger: logger,
		db:     db,
	}
}

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
func (api API) CreateGame(ctx context.Context, user security.JWTUser, pGame *ProspectGame) error {
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
	err = api.storeGameBoard(ctx, user, game, board)
	if err != nil {
		return err
	}
	pGame.ID = game.ID
	return nil
}

// storeGameBoard stores a Game board in the game_board_points table
func (api API) storeGameBoard(ctx context.Context, user security.JWTUser, game *models.Game, board [][]int) error {
	tx, err := api.db.BeginTx(ctx, nil)
	if err != nil {
		api.logger.Printf("error beggining transaction: %v\n", err)
		return err
	}
	// do not insert map
	err = game.Insert(ctx, tx, boil.Whitelist("private", "cols", "rows", "mines", "creator_id"))
	if err != nil {
		api.logger.Printf("error inserting game: %v. Rolling back game insertion\n", err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			api.logger.Printf("error rolling back game creation with error: %v\n", rollbackError)
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
		api.logger.Printf("error inserting all game board points for game %d: %s\n", game.ID, err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			api.logger.Printf("error rolling back game creation with error: %v\n", rollbackError)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		api.logger.Printf("error commiting transaction: %v. Rolling back operation\n", err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			api.logger.Printf("error rolling back operation with error: %v\n", rollbackError)
		}
		return err
	}
	return nil
}

// ApplyOperation applies a given operation to a game and returns an operation confirmation
func (api API) ApplyOperation(ctx context.Context, user security.JWTUser, oper Operation) (OperationConfirmation, error) {
	confirmation := OperationConfirmation{}
	_, err := api.tableRowColRetrieval(ctx, user, oper.GameID, oper.Row, oper.Col)
	if err != nil {
		if err == ErrInvalidRowCols {
			return confirmation, response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
		}
		return confirmation, err
	}
	// first we need to check if there are operations to compose
	gameOperations, err := models.GameOperations(qm.Where("game_id = > and id >= ?", oper.GameID, oper.ID)).All(ctx, api.db)
	if err != nil && err != sql.ErrNoRows {
		return confirmation, err
	}
	clientOperation, err := algebra.NewOperation(oper.Op, oper.Row, oper.Col)
	if err != nil {
		// ErrUnknownOperation
		return confirmation, response.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("unknown operation %d", oper.Op),
		}
	}
	serverOperationsLen := len(gameOperations)
	if serverOperationsLen > 0 {
		serverOperations := make([]algebra.Operation, serverOperationsLen)
		for i, o := range gameOperations {
			// since this operations comes from the server, it will never throw an error
			// because we guarantee it is valid
			serverOperations[i], _ = toAlgebraOperation(o)
		}
		for _, o := range gameOperations {
			algebra.ComposeMulti(serverOperations, clientOperation)
		}

	}
	return confirmation, nil
}

func (api API) tableRowColRetrieval(ctx context.Context, user security.JWTUser, gameID int64, row, col int) (int, error) {
	gameBoardPoint, err := models.GameBoardPoints(
		qm.Select("mine_proximity"),
		qm.InnerJoin("games g on g.id = game_board_points.game_id"),
		qm.Where("game_board_points.game_id = ? AND game_board_points.row = ? AND game_board_points.col = ? AND (g.creator_id = ? OR g.private = false)", gameID, row, col, user.ID),
	).One(ctx, api.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidRowCols
		}
		api.logger.Printf("error retrieving game board point: %v\n", err)
		return 0, err
	}
	return int(gameBoardPoint.MineProximity), nil
}

func (api API) tableUpdateRowCol(ctx context.Context, user security.JWTUser, gameID int64, row, col, mineProximity int) error {
	aff, err := models.GameBoardPoints(
		qm.Where("game_id = ? AND row = ? AND col = ?", gameID, row, col),
	).UpdateAll(ctx, api.db, models.M{
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

func toAlgebraOperation(o *models.GameOperation) (algebra.Operation, error) {
	opType := algebra.OpMark
	if o.Operation == models.MineOperationReveal {
		opType = algebra.OpReveal
	}
	return algebra.NewOperation(algebra.OperationType(opType), int(o.Row), int(o.Col))
}
