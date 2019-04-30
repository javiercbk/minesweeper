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

	extErrors "github.com/pkg/errors"

	"github.com/javiercbk/minesweeper/algebra"
	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// uniqueGameOperationConstaintName is the constraint that ensures that operations are unique within a game
const uniqueGameOperationConstaintName = "idx_game_operation"

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

// proximityMask is a bit masking value to enumerate options find proximities
type proximityMask uint32

const (
	// pEmpty masks all unrevealed empty values
	pEmpty proximityMask = 1 << iota
	// prevealed masks all revealed mines
	pRevealed
	// pMines masks all
	pAll
)

// ErrInvalidRowCols is returned when given a zero or negative row or column count
var ErrInvalidRowCols = errors.New("invalid row or column count")

// ErrTooManyMines is returned when the amount of mines is greater or equal than the board
var ErrTooManyMines = errors.New("too many mines")

// ErrNoneMines is returned when the amount of mines is negative or zero
var ErrNoneMines = errors.New("none mines")

// ErrGameNotExists is returned when attempting to make an operation with a game that does not exists
var ErrGameNotExists = errors.New("the game does not exists")

// ErrGameFinished is returned when attempting to apply an operation on a concluded game
var ErrGameFinished = errors.New("the game has finished")

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
	ID      int                   `json:"id" validate:"required"`
	GameID  int64                 `json:"gameId" validate:"required"`
	Op      algebra.OperationType `json:"op" validate:"required,1|2"`
	Row     int                   `json:"row" validate:"required,gte=0,lt=100"`
	Col     int                   `json:"col" validate:"required,gte=0,lt=100"`
	Applied bool                  `json:"applied"`
	Result  []OperationResult     `json:"result,omitempty"`
}

// Status is the game status
type Status struct {
	Rows  int     `json:"-"`
	Cols  int     `json:"-"`
	Won   bool    `json:"won"`
	Lost  bool    `json:"lost"`
	Board [][]int `json:"board,omitempty"`
}

// OperationConfirmation is the confirmation of an operation application
type OperationConfirmation struct {
	Operation       Operation   `json:"operation,omitempty"`
	DeltaOperations []Operation `json:"deltaOperations,omitempty"`
	Status          Status      `json:"status,omitempty"`
	Error           error       `json:"error"`
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
		Mines:     int16(pGame.Mines),
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

// ApplyOperation applies a given operation to a game and returns an operation confirmation
func (api API) ApplyOperation(ctx context.Context, user security.JWTUser, oper Operation) (OperationConfirmation, error) {
	var err error
	confirmation := OperationConfirmation{
		Operation: Operation{
			GameID: oper.GameID,
			Row:    oper.Row,
			Col:    oper.Col,
		},
	}
	confirmationChan := make(chan OperationConfirmation)
	go api.applyOperation(ctx, user, oper, confirmationChan)
	// attempt to apply the operation within a timeout
	select {
	case confirmation = <-confirmationChan:
		err = confirmation.Error
	case <-ctx.Done():
		err = ctx.Err()
	}
	if err != nil {
		return confirmation, err
	}
	return confirmation, nil
}

func (api API) applyOperation(ctx context.Context, user security.JWTUser, oper Operation, confirmationChan chan OperationConfirmation) {
	defer close(confirmationChan)
	gameDoesNotExistError := response.HTTPError{
		Code:    http.StatusNotFound,
		Message: ErrGameNotExists.Error(),
	}
	confirmation := OperationConfirmation{
		Operation: oper,
	}
	game, err := models.Games(
		qm.Where("id = ? AND (private = false OR creator_id = ?)", oper.GameID, user.ID),
	).One(ctx, api.db)
	if err != nil {
		if err == sql.ErrNoRows {
			err = gameDoesNotExistError
		}
	} else if game == nil {
		err = gameDoesNotExistError
	} else if game.FinishedAt.Valid {
		err = response.HTTPError{
			Code:    http.StatusNotFound,
			Message: ErrGameFinished.Error(),
		}
	} else {
		confirmation.Status.Rows = int(game.Rows)
		confirmation.Status.Cols = int(game.Cols)
		err = api.attempApplyOperation(ctx, user, oper, &confirmation)
	}
	if err != nil {
		confirmation.Error = err
	}
	confirmationChan <- confirmation
}

func (api API) attempApplyOperation(ctx context.Context, user security.JWTUser, oper Operation, confirmation *OperationConfirmation) error {
	var gameOperations models.GameOperationSlice
	var mineProximity algebra.MineProximity
	operationID := oper.ID
	clientOperation, err := algebra.NewOperation(oper.Op, oper.Row, oper.Col)
	if err != nil {
		return err
	}
	// while context is not done attempt this operation
	// WARNING: This algorithm might lead to client starvation,
	// meaning that if several players are playing, the game might get bogged down with timeouts on operations
	// clashing within themselves
	for ctx.Err() == nil {
		// step 1 => check if there are older operations to apply
		gameOperations, err = models.GameOperations(
			qm.Where("game_id = ? and operation_id >= ?", oper.GameID, operationID),
			qm.OrderBy("operation_id ASC"),
		).All(ctx, api.db)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		// step 2 check if there are older, unapplied operations, that invalidate this operation
		opApplied := true
		serverOperationsLen := len(gameOperations)
		serverOperations := make([]algebra.Operation, serverOperationsLen)
		deltaOperations := make([]Operation, serverOperationsLen)
		newID := operationID
		if serverOperationsLen > 0 {
			newID = composeServerClient(gameOperations, oper.GameID, serverOperations, deltaOperations)
			opApplied = algebra.ShouldOperationApply(serverOperations, clientOperation)
		}
		confirmation.DeltaOperations = deltaOperations
		confirmation.Operation.GameID = oper.GameID
		// step 3 => retrieve the current mine proximity value
		mineProximity, err = api.retrieveRowCol(ctx, user, oper.GameID, oper.Row, oper.Col)
		if err != nil {
			if err == ErrInvalidRowCols {
				err = response.HTTPError{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				}
			}
			return err
		}
		if opApplied {
			var newMineProximity algebra.MineProximity
			// step 4a1 => apply the operation with to the current value
			newMineProximity, err = clientOperation.Exec(mineProximity)
			if err != nil {
				// ErrOperationOutOfBounds
				err = response.HTTPError{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				}
				return err
			}
			// step 4a2 => if the mine proximity is the same, after the operation, then don't apply the operation
			if newMineProximity == mineProximity {
				// operation had no action, mark as not applied
				markOperationNotApplied(confirmation, mineProximity, oper)
				break
			} else {
				// the mine proximity is different so the operation changes the actual value.
				// commit the operation.
				confirmation.Operation.ID = newID
				err = api.commitOperation(ctx, user, confirmation, newMineProximity)
				if err != nil {
					cause := extErrors.Cause(err)
					if pgerr, ok := cause.(*pq.Error); ok {
						if pgerr.Constraint == uniqueGameOperationConstaintName {
							// if the operation failed to be commited because the operation id is not unique
							// it means that some operation was commited while this operation was being process.
							// In such case retry the whole algorithm.
							continue
						}
					}
					return err
				}
				if err != nil {
					return err
				}
				confirmation.Operation.Applied = true
				break
			}
		} else {
			// operation should not be applied
			markOperationNotApplied(confirmation, mineProximity, oper)
			break
		}
	}
	return ctx.Err()
}

func (api API) commitOperation(ctx context.Context, user security.JWTUser, confirmation *OperationConfirmation, mineProximity algebra.MineProximity) error {
	tx, err := api.db.BeginTx(ctx, nil)
	if err != nil {
		api.logger.Printf("error beggining transaction for updating game row: %v\n", err)
		return err
	}
	err = api.updateRowCol(ctx, confirmation.Operation.GameID, confirmation.Operation.Row, confirmation.Operation.Col, mineProximity)
	if err != nil {
		api.logger.Printf("error updating game row: %v. Rolling back operation insertion\n", err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			api.logger.Printf("error rolling back operation with error: %v\n", rollbackError)
		}
		return err
	}
	newGameOperation := &models.GameOperation{
		GameID:        confirmation.Operation.GameID,
		Row:           int16(confirmation.Operation.Row),
		Col:           int16(confirmation.Operation.Col),
		PlayerID:      user.ID,
		MineProximity: int16(mineProximity),
		OperationID:   confirmation.Operation.ID,
		Operation:     operationTypeStr(confirmation.Operation.Op),
	}
	err = newGameOperation.Insert(ctx, tx, boil.Infer())
	if err != nil {
		api.logger.Printf("error inserting game operation: %v. Rolling back operation insertion\n", err)
		// just log rollback error
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			api.logger.Printf("error rolling back operation with error: %v\n", rollbackError)
		}
		return err
	}
	// check if the game status needs to be updated
	if mineProximity == 9 {
		confirmation.Status.Lost = true
	} else if mineProximity >= 0 && mineProximity < 9 {
		// if mine proximity is not a mine, then check if the game was won
		exists, err := models.GameBoardPoints(
			qm.Where("game_id = ? AND ((mine_proximity <= -1 AND mine_proximity > -10) OR mine_proximity = 9)", confirmation.Operation.GameID),
		).Exists(ctx, tx)
		if err != nil {
			api.logger.Printf("error checking if the game was won: %v. Rolling back operation insertion\n", err)
			// just log rollback error
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				api.logger.Printf("error rolling back operation with error: %v\n", rollbackError)
			}
			return err
		}
		if !exists {
			// no more mines detected, game won
			confirmation.Status.Won = true
		}
	}
	if confirmation.Status.Won || confirmation.Status.Lost {
		_, err = api.updateGameState(ctx, tx, confirmation.Operation.GameID, confirmation.Status.Won)
		if err != nil {
			api.logger.Printf("error setting the game won %v: %v. Rolling back operation insertion\n", confirmation.Status.Won, err)
			// just log rollback error
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				api.logger.Printf("error rolling back operation with error: %v\n", rollbackError)
			}
			return err
		}
		confirmation.Status.Board, err = retrieveFullBoard(ctx, tx, confirmation.Operation.GameID, confirmation.Status.Rows, confirmation.Status.Cols)
		if err != nil {
			api.logger.Printf("error getting the whole game board %v: %v. Rolling back operation insertion\n", confirmation.Status.Won, err)
			// just log rollback error
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				api.logger.Printf("error rolling back operation with error: %v\n", rollbackError)
			}
			return err
		}
	} else if mineProximity == 0 {
		// TODO: if the game was not lost nor one but the mine proximity was 0, reveal the sibling places with zero mines
	}
	return tx.Commit()
}

func (api API) updateGameState(ctx context.Context, tx *sql.Tx, gameID int64, won bool) (int64, error) {
	return models.Games(qm.Where("id = ?", gameID)).
		UpdateAll(ctx, tx, models.M{"won": won, "finished_at": time.Now().UTC().Format(time.RFC3339)})
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

func (api API) retrieveRowCol(ctx context.Context, user security.JWTUser, gameID int64, row, col int) (int, error) {
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

func (api API) updateRowCol(ctx context.Context, gameID int64, row, col, mineProximity int) error {
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
	if rows > 100 || cols > 100 {
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

			boardCartesian[(i*cols)+j] = boardPoint{
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

func composeServerClient(gameOperations models.GameOperationSlice, gameID int64, serverOperations []algebra.Operation, deltaOperations []Operation) int {
	newID := 0
	for i, o := range gameOperations {
		// since this operations comes from the server, it will never throw an error
		// because we guarantee it is valid
		opRes := OperationResult{
			Row:        int(o.Row),
			Col:        int(o.Col),
			PointState: proximityToState(int(o.MineProximity)),
		}
		if o.MineProximity >= 0 {
			opRes.MineProximity = int(o.MineProximity)
		}
		// we ignore this error because server operations are guaranteed to be valid
		serverOperations[i], _ = toAlgebraOperation(o)
		deltaOperations[i] = Operation{
			ID:      o.OperationID,
			GameID:  gameID,
			Row:     int(o.Row),
			Col:     int(o.Col),
			Op:      operationType(o.Operation),
			Applied: true,
			Result:  []OperationResult{opRes},
		}
		// if there are older operations, calculate the new id for the client operation
		newID = o.OperationID + 1
	}
	return newID
}

func retrieveFullBoard(ctx context.Context, executor boil.ContextExecutor, gameID int64, rows, cols int) ([][]int, error) {
	var board [][]int
	points, err := retrieveMines(ctx, executor, gameID, pAll)
	if err != nil {
		return board, err
	}
	board = make([][]int, rows)
	for i := range board {
		board[i] = make([]int, cols)
	}
	for _, p := range points {
		board[p.Row][p.Col] = int(p.MineProximity)
	}
	return board, nil
}

func retrieveMines(ctx context.Context, executor boil.ContextExecutor, gameID int64, mask proximityMask) (models.GameBoardPointSlice, error) {
	var where qm.QueryMod
	if mask&pEmpty != 0 {
		where = qm.Where("game_id = ? AND mine_proximity > -1 AND ", gameID)
	} else if mask&pRevealed != 0 {
		where = qm.Where("game_id = ? AND mine_proximity >= 0", gameID)
	} else if mask&pAll != 0 {
		where = qm.Where("game_id = ?", gameID)
	}
	return models.GameBoardPoints(
		qm.Select("row, col, mine_proximity"),
		where,
	).All(ctx, executor)
}

func toAlgebraOperation(o *models.GameOperation) (algebra.Operation, error) {
	opType := operationType(o.Operation)
	return algebra.NewOperation(opType, int(o.Row), int(o.Col))
}

func operationType(strOp string) algebra.OperationType {
	opType := algebra.OpMark
	if strOp == models.MineOperationReveal {
		opType = algebra.OpReveal
	}
	return opType
}

func operationTypeStr(t algebra.OperationType) string {
	opTypeStr := models.MineOperationMark
	if t == algebra.OpReveal {
		opTypeStr = models.MineOperationReveal
	}
	return opTypeStr
}

func proximityToState(p algebra.MineProximity) int {
	if p <= -21 {
		return StateMarkedMine
	} else if p > -21 && p <= -11 {
		return StateSuspectMine
	} else if p > -11 && p <= -1 {
		return StateNotRevealed
	}
	return StateRevealed
}

func markOperationNotApplied(confirmation *OperationConfirmation, newMineProximity algebra.MineProximity, oper Operation) {
	confirmation.Operation.ID = 0
	confirmation.Operation.Applied = false
	opResult := OperationResult{
		Row:        oper.Row,
		Col:        oper.Col,
		PointState: proximityToState(newMineProximity),
	}
	if newMineProximity >= 0 {
		// only show mine proximity if already revealed
		opResult.MineProximity = newMineProximity
	}
	confirmation.Operation.Applied = false
	confirmation.Operation.Result = []OperationResult{opResult}
}
