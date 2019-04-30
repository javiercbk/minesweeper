package game

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/javiercbk/minesweeper/algebra"
	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
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

func TestApplyRevealOperations(t *testing.T) {
	ctx := context.Background()
	// ctxTimeout, _ := context.WithTimeout(ctx, time.Second*5)
	api, user, otherUser := setUp(ctx, t, username)
	tests := []struct {
		game                 *models.Game
		finished             bool
		initialBoard         [][]int
		operation            Operation
		expectedBoard        [][]int
		expectedConfirmation OperationConfirmation
		err                  error
	}{
		{
			game: &models.Game{
				// should fail when attempting to apply an operation in a private game
				// that was created by another user
				CreatorID: otherUser.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   true,
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
			err: response.HTTPError{
				Code:    http.StatusNotFound,
				Message: ErrGameNotExists.Error(),
			},
		},
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
				// should fail when game id does not exist
				GameID: -1,
			},
			err: response.HTTPError{
				Code:    http.StatusNotFound,
				Message: ErrGameNotExists.Error(),
			},
		},
		{
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			// should not allow to apply operations on finished games
			finished: true,
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
			err: response.HTTPError{
				Code:    http.StatusNotFound,
				Message: ErrGameFinished.Error(),
			},
		},
		{
			game: &models.Game{
				CreatorID:  user.ID,
				Rows:       int16(3),
				Cols:       int16(3),
				Mines:      int16(2),
				Private:    false,
				FinishedAt: null.NewTime(time.Now().UTC(), true),
			},
			initialBoard: [][]int{
				{-2, -9, -2},
				{-2, -3, -3},
				{-1, -2, -9},
			},
			operation: Operation{
				ID: 1,
				// should fail with this row overflow
				Row: 3,
				Col: 0,
				Op:  algebra.OpReveal,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrInvalidRowCols.Error(),
			},
		},
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
				// should fail with this column overflow
				Col: 3,
				Op:  algebra.OpReveal,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: ErrInvalidRowCols.Error(),
			},
		},
		{
			// should apply the reveal operation
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
		{
			game: &models.Game{
				CreatorID: otherUser.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				// should allow other user to play non private game
				Private: false,
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
