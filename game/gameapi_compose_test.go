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
			// mark operation should have no effect because it is marking something that was revealed
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
		{
			// reveal operation should fail because the row is already marked
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{1, -20, -2},
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
				{
					PlayerID:      user.ID,
					OperationID:   3,
					Operation:     "mark",
					Row:           0,
					Col:           1,
					MineProximity: -20,
				},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 1,
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      0,
					Row:     0,
					Col:     1,
					Op:      algebra.OpReveal,
					Applied: false,
					Result: []OperationResult{
						{
							Row:           0,
							Col:           1,
							MineProximity: 0,
							PointState:    StateSuspectMine,
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
					{
						ID:      3,
						Row:     0,
						Col:     1,
						Op:      algebra.OpMark,
						Applied: true,
						Result: []OperationResult{
							{
								Row:           0,
								Col:           1,
								MineProximity: 0,
								PointState:    StateSuspectMine,
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
				{1, -20, -2},
				{1, -3, -3},
				{-1, -2, -10},
			},
		},
		{
			// reveal operation should have no effect because it is revealing somethin alredy revealed
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
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      0,
					Row:     0,
					Col:     0,
					Op:      algebra.OpReveal,
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
		{
			// reveal operation should be applied
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
				Row: 1,
				Col: 1,
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      3,
					Row:     1,
					Col:     1,
					Op:      algebra.OpReveal,
					Applied: true,
					Result: []OperationResult{
						{
							Row:           1,
							Col:           1,
							MineProximity: 2,
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
				{1, 2, -3},
				{-1, -2, -10},
			},
		},
		{
			// reveal operation should be applied and game should be won
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{1, -10, 1},
				{1, -2, 1},
				{0, 0, 0},
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
				Row: 1,
				Col: 1,
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      3,
					Row:     1,
					Col:     1,
					Op:      algebra.OpReveal,
					Applied: true,
					Result: []OperationResult{
						{
							Row:           1,
							Col:           1,
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
					Won:  true,
					Board: [][]int{
						{1, -10, 1},
						{1, 1, 1},
						{0, 0, 0},
					},
				},
			},
			expectedBoard: [][]int{
				{1, -10, 1},
				{1, 1, 1},
				{0, 0, 0},
			},
		},
		{
			// reveal operation should be applied and game should be lost
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   false,
			},
			initialBoard: [][]int{
				{1, -10, 1},
				{1, -2, 1},
				{0, 0, 0},
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
				Col: 1,
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      3,
					Row:     0,
					Col:     1,
					Op:      algebra.OpReveal,
					Applied: true,
					Result: []OperationResult{
						{
							Row:           0,
							Col:           1,
							MineProximity: 9,
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
					Lost: true,
					Board: [][]int{
						{1, 9, 1},
						{1, -2, 1},
						{0, 0, 0},
					},
				},
			},
			expectedBoard: [][]int{
				{1, 9, 1},
				{1, -2, 1},
				{0, 0, 0},
			},
		},
	}
	assertGameTests(ctx, t, user, api, tests)
}
