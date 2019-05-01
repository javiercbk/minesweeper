package game

import (
	"context"
	"testing"
	"time"

	"github.com/javiercbk/minesweeper/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
)

func TestFindGame(t *testing.T) {
	ctx := context.Background()
	api, user, otherUser := setUp(ctx, t, username)
	games := []*models.Game{
		&models.Game{
			CreatorID: otherUser.ID,
			Rows:      int16(3),
			Cols:      int16(3),
			Mines:     int16(2),
			Private:   true,
		},
		&models.Game{
			CreatorID: otherUser.ID,
			Rows:      int16(3),
			Cols:      int16(3),
			Mines:     int16(2),
			Private:   false,
		},
		&models.Game{
			CreatorID: user.ID,
			Rows:      int16(3),
			Cols:      int16(3),
			Mines:     int16(2),
			Private:   true,
		},
		&models.Game{
			CreatorID: user.ID,
			Rows:      int16(3),
			Cols:      int16(3),
			Mines:     int16(2),
			Private:   false,
		},
		&models.Game{
			CreatorID:  user.ID,
			Rows:       int16(3),
			Cols:       int16(3),
			Mines:      int16(2),
			Private:    false,
			FinishedAt: null.NewTime(time.Now(), true),
		},
	}
	// FIXME: had to add this to only use the data from this test
	_, err := queries.Raw("TRUNCATE TABLE games CASCADE").ExecContext(ctx, api.db)
	if err != nil {
		t.Fatalf("error truncating games table: %v", err)
	}
	for i := range games {
		err := games[i].Insert(ctx, api.db, boil.Infer())
		if err != nil {
			t.Fatalf("error inserting game %d: %v", i, err)
		}
	}
	gamesFound, err := api.FindGames(ctx, user)
	if err != nil {
		t.Fatalf("error retrieving games %v", err)
	}
	gamesFoundLen := len(gamesFound)
	if gamesFoundLen != 3 {
		t.Fatalf("expected 3 games but found %d", gamesFoundLen)
	}
}
