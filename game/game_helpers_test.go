package game

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

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
	createPlayer(ctx, db, t, testPlayer)
	anotherTestPlayer := &models.Player{
		Name:     anotherUsername,
		Password: abcHashed,
	}
	createPlayer(ctx, db, t, anotherTestPlayer)
	return NewAPI(logger, db), security.JWTUser{
			ID:   testPlayer.ID,
			Name: name,
		}, security.JWTUser{
			ID:   anotherTestPlayer.ID,
			Name: anotherUsername,
		}
}

func createPlayer(ctx context.Context, db *sql.DB, t testing.TB, player *models.Player) {
	err := player.Insert(ctx, db, boil.Infer())
	if err != nil {
		cause := extErrors.Cause(err)
		if pgerr, ok := cause.(*pq.Error); ok {
			if pgerr.Constraint == uniqueNameConstaintName {
				player, err = models.Players(qm.Where("name = ?", player.Name)).One(ctx, db)
				if err != nil {
					t.Fatalf("error retrieving test user: %v\n", err)
				}
			}
		} else {
			t.Fatalf("error inserting test user: %v\n", err)
		}
	}
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
