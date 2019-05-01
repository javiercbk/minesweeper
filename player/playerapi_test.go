package player

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/javiercbk/minesweeper/http/response"
	"github.com/javiercbk/minesweeper/models"
	testHelpers "github.com/javiercbk/minesweeper/testing"
	"github.com/volatiletech/sqlboiler/boil"
	"golang.org/x/crypto/bcrypt"
)

// abcHashed is the bcrypt hash of the string "abc" (without quotes)
const abcHashed = "$2y$12$Fq0ne4S2xnhZTYE7p/veuOX3X6DlF1qZYeeHhK/PY39TP7//klYkW"
const jwtSecret = "wow"
const username = "existingUserName"

func TestMain(m *testing.M) {
	testHelpers.InitializeDB(m)
}

func setUp(ctx context.Context, t *testing.T, existingUserName string) api {
	logger := testHelpers.NullLogger()
	db, err := testHelpers.DB()
	if err != nil {
		t.Fatalf("error connecting to database: %v\n", err)
	}
	testPlayer := &models.Player{
		Name:     existingUserName,
		Password: abcHashed,
	}
	err = testPlayer.Insert(ctx, db, boil.Infer())
	if err != nil {
		t.Fatalf("error inserting test user: %v\n", err)
	}
	return NewAPI(logger, db).(api)
}

func TestCreatePlayer(t *testing.T) {
	ctx := context.Background()
	api := setUp(ctx, t, username)
	tests := []struct {
		pPlayer ProspectPlayer
		err     error
	}{
		{
			pPlayer: ProspectPlayer{
				Name:     "user1",
				Password: "abc",
			},
			err: nil,
		},
		{
			pPlayer: ProspectPlayer{
				Name:     username,
				Password: "abc",
			},
			err: response.HTTPError{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("player %s already exists", username),
			},
		},
	}
	for i, test := range tests {
		err := api.CreatePlayer(ctx, &test.pPlayer)
		if test.err != err {
			t.Fatalf("failed test %d: expected error to be %v but was %v\n", i, test.err, err)
		}
		if err == nil {
			if test.pPlayer.ID == 0 {
				t.Fatalf("failed test %d: expected id to not be zero but was %d\n", i, test.pPlayer.ID)
			}
			player, err := models.FindPlayer(ctx, api.db, test.pPlayer.ID)
			if err != nil {
				t.Fatalf("failed test %d: expected error to be nil but was %v\n", i, err)
			}
			if player.Name != test.pPlayer.Name {
				t.Fatalf("failed test %d: expected name to be %s but was %s\n", i, test.pPlayer.Name, player.Name)
			}
			err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(test.pPlayer.Password))
			if err != nil {
				t.Fatalf("failed test %d: expected bcrypt error to be nil but was %v\n", i, err)
			}
		}
	}
}
