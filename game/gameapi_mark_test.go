package game

import (
	"context"
	"testing"
	"time"

	"github.com/javiercbk/minesweeper/algebra"
	"github.com/javiercbk/minesweeper/models"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func TestApplyMarkOperations(t *testing.T) {
	ctx := context.Background()
	// ctxTimeout, _ := context.WithTimeout(ctx, time.Second*5)
	api, user, _ := setUp(ctx, t, username)
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
			// mark operation should have no effect
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{1, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 0,
				Op:  algebra.OpMark,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      0,
					Row:     0,
					Col:     0,
					Op:      algebra.OpMark,
					Applied: false,
				},
				Status: Status{
					Rows: 3,
					Cols: 3,
				},
			},
			expectedBoard: [][]int{
				{1, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
			},
		},
		{
			// mark operation should be applied
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{-2, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 0,
				Op:  algebra.OpMark,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      1,
					Row:     0,
					Col:     0,
					Op:      algebra.OpMark,
					Applied: true,
				},
				Status: Status{
					Rows: 3,
					Cols: 3,
				},
			},
			expectedBoard: [][]int{
				{-12, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
			},
		},
		{
			// mark operation should be applied again
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{-12, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 0,
				Op:  algebra.OpMark,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      1,
					Row:     0,
					Col:     0,
					Op:      algebra.OpMark,
					Applied: true,
				},
				Status: Status{
					Rows: 3,
					Cols: 3,
				},
			},
			expectedBoard: [][]int{
				{-22, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
			},
		},
		{
			// mark operation should be applied again thus restoring the field to the initial value
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{-22, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 0,
				Op:  algebra.OpMark,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      1,
					Row:     0,
					Col:     0,
					Op:      algebra.OpMark,
					Applied: true,
				},
				Status: Status{
					Rows: 3,
					Cols: 3,
				},
			},
			expectedBoard: [][]int{
				{-2, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
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
