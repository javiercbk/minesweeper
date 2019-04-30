package game

import (
	"context"
	"testing"

	"github.com/javiercbk/minesweeper/algebra"
	"github.com/javiercbk/minesweeper/models"
)

func TestApplyMarkOperations(t *testing.T) {
	ctx := context.Background()
	// ctxTimeout, _ := context.WithTimeout(ctx, time.Second*5)
	api, user, _ := setUp(ctx, t, username)
	tests := []gameTest{
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
					Result: []OperationResult{
						{
							Row:           0,
							Col:           0,
							MineProximity: 1,
							PointState:    StateRevealed,
						},
					},
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
					Result: []OperationResult{
						{
							Row:           0,
							Col:           0,
							MineProximity: 0,
							PointState:    StateSuspectMine,
						},
					},
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
					Result: []OperationResult{
						{
							Row:           0,
							Col:           0,
							MineProximity: 0,
							PointState:    StateMarkedMine,
						},
					},
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
					Result: []OperationResult{
						{
							Row:           0,
							Col:           0,
							MineProximity: 0,
							PointState:    StateNotRevealed,
						},
					},
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
	assertGameTests(ctx, t, user, api, tests)
}
