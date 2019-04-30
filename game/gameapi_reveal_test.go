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
)

func TestApplyRevealOperations(t *testing.T) {
	ctx := context.Background()
	// ctxTimeout, _ := context.WithTimeout(ctx, time.Second*5)
	api, user, otherUser := setUp(ctx, t, username)
	tests := []gameTest{
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
				{-2, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
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
				{-2, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
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
				{-2, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
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
				{-2, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
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
				{-2, -10, -2},
				{-2, -3, -3},
				{-1, -2, -10},
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
			// reveal should throw an error when value is marked
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
				Op:  algebra.OpReveal,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: algebra.ErrOperationOutOfBounds.Error(),
			},
		},
		{
			// reveal should throw an error when value is marked
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
				Op:  algebra.OpReveal,
			},
			err: response.HTTPError{
				Code:    http.StatusBadRequest,
				Message: algebra.ErrOperationOutOfBounds.Error(),
			},
		},
		{
			// reveal should have no effect
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
			// should apply the reveal operation
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
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
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
			game: &models.Game{
				CreatorID: otherUser.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				// should allow other user to play non private game
				Private: false,
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
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
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
			game: &models.Game{
				CreatorID: otherUser.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				// should allow other user to play non private game
				Private: false,
			},
			initialBoard: [][]int{
				{0, -2, -10},
				{0, 1, 1},
				{0, 0, 0},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 1,
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      1,
					Row:     0,
					Col:     1,
					Op:      algebra.OpReveal,
					Applied: true,
					Result: []OperationResult{
						{
							Row:           0,
							Col:           1,
							MineProximity: 1,
							PointState:    StateRevealed,
						},
					},
				},
				Status: Status{
					Won:  true,
					Rows: 3,
					Cols: 3,
					Board: [][]int{
						{0, 1, -10},
						{0, 1, 1},
						{0, 0, 0},
					},
				},
			},
			expectedBoard: [][]int{
				{0, 1, -10},
				{0, 1, 1},
				{0, 0, 0},
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
				{0, -2, -10},
				{0, 1, 1},
				{0, 0, 0},
			},
			operation: Operation{
				ID:  1,
				Row: 0,
				Col: 2,
				Op:  algebra.OpReveal,
			},
			expectedConfirmation: OperationConfirmation{
				Operation: Operation{
					ID:      1,
					Row:     0,
					Col:     2,
					Op:      algebra.OpReveal,
					Applied: true,
					Result: []OperationResult{
						{
							Row:           0,
							Col:           2,
							MineProximity: 9,
							PointState:    StateRevealed,
						},
					},
				},
				Status: Status{
					Lost: true,
					Rows: 3,
					Cols: 3,
					Board: [][]int{
						{0, -2, 9},
						{0, 1, 1},
						{0, 0, 0},
					},
				},
			},
			expectedBoard: [][]int{
				{0, -2, 9},
				{0, 1, 1},
				{0, 0, 0},
			},
		},
	}
	assertGameTests(ctx, t, user, api, tests)
}
