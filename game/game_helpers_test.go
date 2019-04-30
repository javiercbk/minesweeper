package game

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	testHelpers "github.com/javiercbk/minesweeper/testing"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	extErrors "github.com/pkg/errors"
)

const uniqueNameConstaintName = "idx_players_name"
const abcHashed = "$2y$12$Fq0ne4S2xnhZTYE7p/veuOX3X6DlF1qZYeeHhK/PY39TP7//klYkW"
const username = "benchmarkUsername"
const anotherUsername = "testUsername"

const gameRows = 100
const gameCols = 100
const gameMines = 99

type gameTest struct {
	game                 *models.Game
	finished             bool
	initialBoard         [][]int
	operation            Operation
	existingOperations   models.GameOperationSlice
	expectedBoard        [][]int
	expectedConfirmation OperationConfirmation
	err                  error
}

func TestMain(m *testing.M) {
	testHelpers.InitializeDB(m)
}

func setUp(ctx context.Context, t testing.TB, name string) (API, security.JWTUser, security.JWTUser) {
	logger := testHelpers.NullLogger()
	db, err := testHelpers.DB()
	if err != nil {
		t.Fatalf("error connecting to database: %v\n", err)
	}
	testPlayer := &models.Player{
		Name:     name,
		Password: abcHashed,
	}
	err = createPlayer(ctx, db, testPlayer)
	if err != nil {
		t.Fatalf("error creating player %v", err)
	}
	anotherTestPlayer := &models.Player{
		Name:     anotherUsername,
		Password: abcHashed,
	}
	err = createPlayer(ctx, db, anotherTestPlayer)
	if err != nil {
		t.Fatalf("error creating player %v", err)
	}
	return NewAPI(logger, db), security.JWTUser{
			ID:   testPlayer.ID,
			Name: name,
		}, security.JWTUser{
			ID:   anotherTestPlayer.ID,
			Name: anotherUsername,
		}
}

func createPlayer(ctx context.Context, db *sql.DB, player *models.Player) error {
	err := player.Insert(ctx, db, boil.Infer())
	if err != nil {
		cause := extErrors.Cause(err)
		if pgerr, ok := cause.(*pq.Error); ok {
			if pgerr.Constraint == uniqueNameConstaintName {
				dbPlayer, err := models.Players(qm.Where("name = ?", player.Name)).One(ctx, db)
				if err != nil {
					return fmt.Errorf("error retrieving test user: %v", err)
				}
				player.ID = dbPlayer.ID
				player.Name = dbPlayer.Name
			} else {
				return err
			}
		} else {
			return fmt.Errorf("error inserting test user: %v", err)
		}
	}
	return nil
}

func assertOperationConfirmation(o1, o2 OperationConfirmation) error {
	err := assertOperation(o1.Operation, o2.Operation)
	if err != nil {
		return err
	}
	err = assertDeltaOperations(o1.DeltaOperations, o2.DeltaOperations)
	if err != nil {
		return err
	}
	err = assertStatus(o1.Status, o2.Status)
	if err != nil {
		return err
	}
	if o1.Error != o2.Error {
		err = fmt.Errorf("expected operation confirmaation error to be %v but was %v", o1.Error, o2.Error)
	}
	return err
}

func assertOperation(o1, o2 Operation) error {
	if o1.ID != o2.ID {
		return fmt.Errorf("expected operation ID to be %v but was %v", o1.ID, o2.ID)
	}
	if o1.GameID != o2.GameID {
		return fmt.Errorf("expected operation GameID to be %v but was %v", o1.GameID, o2.GameID)
	}
	if o1.Op != o2.Op {
		return fmt.Errorf("expected operation Op to be %v but was %v", o1.Op, o2.Op)
	}
	if o1.Row != o2.Row {
		return fmt.Errorf("expected operation Row to be %v but was %v", o1.Row, o2.Row)
	}
	if o1.Col != o2.Col {
		return fmt.Errorf("expected operation Col to be %v but was %v", o1.Col, o2.Col)
	}
	if o1.Applied != o2.Applied {
		return fmt.Errorf("expected operation Applied to be %v but was %v", o1.Applied, o2.Applied)
	}
	err := assertOperationResults(o1.Result, o2.Result)
	if err != nil {
		return fmt.Errorf("expected operation Result to be %v but was %v", o1.Result, o2.Result)
	}
	return nil
}

func assertOperationResults(o1, o2 []OperationResult) error {
	var err error
	if len(o1) != len(o2) {
		return fmt.Errorf("expected operation result to have length of %d but was %d", len(o1), len(o2))
	}
	for i := range o1 {
		err = assertOperationResult(o1[i], o2[i])
		if err != nil {
			break
		}
	}
	return err
}

func assertOperationResult(o1, o2 OperationResult) error {
	if o1.Row != o2.Row {
		return fmt.Errorf("expected operation result Row to be %v but was %v", o1.Row, o2.Row)
	}
	if o1.Col != o2.Col {
		return fmt.Errorf("expected operation result Col to be %v but was %v", o1.Col, o2.Col)
	}
	if o1.MineProximity != o2.MineProximity {
		return fmt.Errorf("expected operation result MineProximity to be %v but was %v", o1.MineProximity, o2.MineProximity)
	}
	if o1.PointState != o2.PointState {
		return fmt.Errorf("expected operation result PointState to be %v but was %v", o1.PointState, o2.PointState)
	}
	return nil
}

func assertDeltaOperations(o1, o2 []Operation) error {
	var err error
	if len(o1) != len(o2) {
		return fmt.Errorf("expected delta operations to have length of %d but was %d", len(o1), len(o2))
	}
	for i := range o1 {
		err = assertOperation(o1[i], o2[i])
		if err != nil {
			break
		}
	}
	return err
}

func assertStatus(s1, s2 Status) error {
	if s1.Won != s2.Won {
		return fmt.Errorf("expected status Won to be %v but was %v", s1.Won, s2.Won)
	}
	if s1.Lost != s2.Lost {
		return fmt.Errorf("expected status Lost to be %v but was %v", s1.Lost, s2.Lost)
	}
	if s1.Rows != s2.Rows {
		return fmt.Errorf("expected status Rows to be %v but was %v", s1.Rows, s2.Rows)
	}
	if s1.Cols != s2.Cols {
		return fmt.Errorf("expected status Cols to be %v but was %v", s1.Cols, s2.Cols)
	}
	if len(s1.Board) != len(s1.Board) {
		return fmt.Errorf("expected status Board to be %v but was %v", s1.Board, s2.Board)
	}
	for r := range s1.Board {
		if len(s1.Board[r]) != len(s2.Board[r]) {
			return fmt.Errorf("expected status Board to be %v but was %v", s1.Board, s2.Board)
		}
		for c := range s1.Board[r] {
			if s1.Board[r][c] != s2.Board[r][c] {
				return fmt.Errorf("expected status Board to be %v but was %v", s1.Board, s2.Board)
			}
		}
	}
	return nil
}

func assertGameTests(ctx context.Context, t testing.TB, user security.JWTUser, api API, tests []gameTest) {
	for i, test := range tests {
		err := api.storeGameBoard(ctx, user, test.game, test.initialBoard)
		if err != nil {
			t.Fatalf("test %d, failed: error creating board %v\n", i, err)
		}
		if len(test.existingOperations) > 0 {
			for _, o := range test.existingOperations {
				o.GameID = test.game.ID
				err = o.Insert(ctx, api.db, boil.Infer())
				if err != nil {
					t.Fatalf("error inserting game operation: %v", err)
				}
			}
		}
		if test.finished {
			_, err = models.Games(qm.Where("id = ?", test.game.ID)).UpdateAll(ctx, api.db, models.M{"finished_at": time.Now().UTC()})
			if err != nil {
				t.Fatalf("test %d, failed: error setting gam as finished board %v\n", i, err)
			}
		}
		test.expectedConfirmation.Operation.GameID = test.game.ID
		if test.operation.GameID >= 0 {
			test.operation.GameID = test.game.ID
		}
		confirmation, err := api.ApplyOperation(ctx, user, test.operation)
		if err != test.err {
			t.Fatalf("test %d failed: expected err to be %v, but was %v\n", i, test.err, err)
		}
		if err == nil {
			for i := range test.expectedConfirmation.DeltaOperations {
				test.expectedConfirmation.DeltaOperations[i].GameID = test.game.ID
			}
			err = assertOperationConfirmation(test.expectedConfirmation, confirmation)
			if err != nil {
				t.Fatalf("test %d failed: %s\n", i, err.Error())
			}
			if test.expectedConfirmation.Status.Won || test.expectedConfirmation.Status.Lost {
				game, err := models.FindGame(ctx, api.db, test.game.ID)
				if err != nil {
					t.Fatalf("test %d, failed: error retrieving game %v\n", i, err)
				}
				if !game.FinishedAt.Valid {
					t.Fatalf("test %d, failed: error game was not marked as finished\n", i)
				}
			}
			board, err := retrieveFullBoard(ctx, api.db, test.game.ID, int(test.game.Rows), int(test.game.Cols))
			if err != nil {
				t.Fatalf("test %d failed: error retrieving game board %v\n", i, err)
			}
			for row := range test.expectedBoard {
				for col := range test.expectedBoard[row] {
					if test.expectedBoard[row][col] != board[row][col] {
						t.Fatalf("test %d failed: expected row %d, col %d to be %d but was %d", i, row, col, test.expectedBoard[row][col], board[row][col])
					}
				}
			}
		}
	}
}
