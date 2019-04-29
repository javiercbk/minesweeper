package game

import (
	"context"
	"net/http"
	"testing"

	"github.com/javiercbk/minesweeper/algebra"
	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/models"
)

func TestCreateGame(t *testing.T) {
	ctx := context.Background()
	api, user, _ := setUp(ctx, t, username)
	tests := []struct {
		game ProspectGame
		err  error
	}{
		{
			game: ProspectGame{
				Rows:    101,
				Cols:    gameCols,
				Mines:   1,
				Private: false,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrInvalidRowCols.Error(),
			},
		},
		{
			game: ProspectGame{
				Rows:    gameRows,
				Cols:    101,
				Mines:   1,
				Private: false,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrInvalidRowCols.Error(),
			},
		},
		{
			game: ProspectGame{
				Rows:    gameRows,
				Cols:    0,
				Mines:   1,
				Private: false,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrInvalidRowCols.Error(),
			},
		},
		{
			game: ProspectGame{
				Rows:    0,
				Cols:    gameCols,
				Mines:   1,
				Private: false,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrInvalidRowCols.Error(),
			},
		},
		{
			game: ProspectGame{
				Rows:    gameRows,
				Cols:    gameCols,
				Mines:   0,
				Private: false,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrNoneMines.Error(),
			},
		},
		{
			game: ProspectGame{
				Rows:    gameRows,
				Cols:    gameCols,
				Mines:   (gameRows * gameCols),
				Private: false,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrTooManyMines.Error(),
			},
		},
		{
			game: ProspectGame{
				Rows:    gameRows,
				Cols:    gameCols,
				Mines:   gameMines,
				Private: false,
			},
		},
	}
	for i, test := range tests {
		err := api.CreateGame(ctx, user, &test.game)
		if err != test.err {
			t.Fatalf("test %d failed: expected err to be %v, but was %v\n", i, test.err, err)
		}
		if err == nil {
			board, err := retrieveFullBoard(ctx, api.db, test.game.ID, test.game.Rows, test.game.Cols)
			if err != nil {
				t.Fatalf("test %d failed: error retrieving game board %v\n", i, err)
			}
			mineCount := test.game.Mines
			for row := range board {
				for col := range board[row] {
					if board[row][col] == -10 {
						mineCount--
					}
				}
			}
			if mineCount != 0 {
				t.Fatalf("test %d, failed: expected mine count to be zero but was %d\n", i, mineCount)
			}
		}
	}
}

func TestApplySuccessfulRevealOperations(t *testing.T) {
	ctx := context.Background()
	// ctxTimeout, _ := context.WithTimeout(ctx, time.Second*5)
	api, user, _ := setUp(ctx, t, username)
	tests := []struct {
		game                 *models.Game
		initialBoard         [][]int
		operation            Operation
		expectedBoard        [][]int
		expectedConfirmation OperationConfirmation
		err                  error
	}{
		{
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{-2, -9, -2},
				{-2, -3, -3},
				{-1, -2, -9},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 0,
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      1,
					Row:     0,
					Col:     0,
					Op:      algebra.OpReveal,
					Applied: true,
					GameID:  1,
				},
				Status: Status{
					Rows: 3,
					Cols: 3,
				},
			},
			expectedBoard: [][]int{
				{1, -9, -2},
				{-2, -3, -3},
				{-1, -2, -9},
			},
		},
	}
	for i, test := range tests {
		err := api.storeGameBoard(ctx, user, test.game, test.initialBoard)
		if err != nil {
			t.Fatalf("test %d, failed: error creating board %v\n", i, err)
		}
		test.operation.GameID = test.game.ID
		confirmation, err := api.ApplyOperation(ctx, user, test.operation)
		if err != test.err {
			t.Fatalf("test %d failed: expected err to be %v, but was %v\n", i, test.err, err)
		}
		if err == nil {
			err = assertOperationConfirmation(test.expectedConfirmation, confirmation)
			if err != nil {
				t.Fatalf("test %d failed: %s\n", i, err.Error())
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
