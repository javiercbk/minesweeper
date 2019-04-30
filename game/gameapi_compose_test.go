package game

import (
	"context"
	"testing"

	"github.com/javiercbk/minesweeper/algebra"
	"github.com/javiercbk/minesweeper/models"
)

func TestApplyComposeOperations(t *testing.T) {
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
				{1, -3, -3},
				{-1, -2, -10},
			},
			existingOperations: models.GameOperationSlice{
				{
					PlayerID:      user.ID,
					OperationID:   1,
					Operation:     "reveal",
					Row:           0,
					Col:           0,
					MineProximity: 1,
				},
				{
					PlayerID:      user.ID,
					OperationID:   2,
					Operation:     "reveal",
					Row:           1,
					Col:           0,
					MineProximity: 1,
				},
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
				DeltaOperations: []Operation{
					{
						ID:      1,
						Row:     0,
						Col:     0,
						Op:      algebra.OpReveal,
						Applied: true,
						Result: []OperationResult{
							{
								Row:           0,
								Col:           0,
								MineProximity: 1,
								PointState:    StateRevealed,
							},
						},
					},
					{
						ID:      2,
						Row:     1,
						Col:     0,
						Op:      algebra.OpReveal,
						Applied: true,
						Result: []OperationResult{
							{
								Row:           1,
								Col:           0,
								MineProximity: 1,
								PointState:    StateRevealed,
							},
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
				{1, -3, -3},
				{-1, -2, -10},
			},
		},
	}
	assertGameTests(ctx, t, user, api, tests)
}
