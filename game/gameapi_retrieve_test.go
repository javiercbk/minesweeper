package game

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/models"
	"github.com/volatiletech/sqlboiler/boil"
)

func assertStatefulGame(game *models.Game, board [][]int, statefulGame StatefulGame) error {
	if game.ID != statefulGame.ID {
		return fmt.Errorf("expected game id to be %d but was %d", game.ID, statefulGame.ID)
	}
	if game.Rows != statefulGame.Rows {
		return fmt.Errorf("expected game rows to be %d but was %d", game.Rows, statefulGame.Rows)
	}
	if game.Cols != statefulGame.Cols {
		return fmt.Errorf("expected game cols to be %d but was %d", game.Cols, statefulGame.Cols)
	}
	if game.Mines != statefulGame.Mines {
		return fmt.Errorf("expected game mines to be %d but was %d", game.Mines, statefulGame.Mines)
	}
	if game.Won != statefulGame.Won {
		return fmt.Errorf("expected game won to be %v but was %v", game.Won, statefulGame.Won)
	}
	if game.Private != statefulGame.Private {
		return fmt.Errorf("expected game private to be %v but was %v", game.Private, statefulGame.Private)
	}
	if game.CreatorID != statefulGame.Creator.ID {
		return fmt.Errorf("expected game creator id to be %d but was %d", game.CreatorID, statefulGame.Creator.ID)
	}
	if statefulGame.Creator.Name == "" {
		return fmt.Errorf("expected game creator name not be empty")
	}
	if statefulGame.LastOperationID != 1 {
		return fmt.Errorf("expected last operation id to be 1 but was %d", statefulGame.LastOperationID)
	}
	for row := range statefulGame.Board {
		for col := range statefulGame.Board[row] {
			if board[row][col] < 0 && statefulGame.Board[row][col].Valid && statefulGame.Board[row][col].Ptr() == nil {
				return fmt.Errorf("expected row %d, col %d to be null but was %d", row, col, statefulGame.Board[row][col].Ptr())
			}
			if board[row][col] >= 0 && statefulGame.Board[row][col].Valid && statefulGame.Board[row][col].Ptr() != nil && (*statefulGame.Board[row][col].Ptr()) != board[row][col] {
				return fmt.Errorf("expected row %d, col %d to be %d but was %d", row, col, board[row][col], (*statefulGame.Board[row][col].Ptr()))
			}
		}
	}
	return nil
}

func TestRetrieveGame(t *testing.T) {
	ctx := context.Background()
	// ctxTimeout, _ := context.WithTimeout(ctx, time.Second*5)
	api, user, otherUser := setUp(ctx, t, username)
	tests := []gameTest{
		{
			// mark operation should have no effect because it is marking something that was revealed
			game: &models.Game{
				CreatorID: otherUser.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   true,
			},
			initialBoard: [][]int{
				{1, -10, -2},
				{-2, -3, -3},
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
			},
			err: response.HTTPError{
				Code:    http.StatusNotFound,
				Message: "game 17 does not exist",
			},
		},
		{
			// mark operation should have no effect because it is marking something that was revealed
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   true,
			},
			failureGameID: 123,
			initialBoard: [][]int{
				{1, -10, -2},
				{-2, -3, -3},
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
			},
			err: response.HTTPError{
				Code:    http.StatusNotFound,
				Message: "game 123 does not exist",
			},
		},
		{
			// mark operation should have no effect because it is marking something that was revealed
			game: &models.Game{
				CreatorID: user.ID,
				Rows:      int16(3),
				Cols:      int16(3),
				Mines:     int16(2),
				Private:   true,
			},
			initialBoard: [][]int{
				{1, -10, -2},
				{-2, -3, -3},
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
			},
		},
	}
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
		if err == nil {
			gameID := test.game.ID
			if test.failureGameID != 0 {
				gameID = test.failureGameID
			}
			sg, err := api.RetrieveGame(ctx, user, gameID)
			if err != test.err {
				t.Fatalf("test %d failed: expected error to be %v but was %v", i, test.err, err)
			}
			if err == nil {
				err := assertStatefulGame(test.game, test.initialBoard, sg)
				if err != nil {
					t.Fatalf("test %d failed, %s\n", i, err)
				}
			}
		}
	}
}
