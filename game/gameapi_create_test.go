package game

import (
	"context"
	"net/http"
	"testing"

	"github.com/javiercbk/minesweeper/http/response"
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
