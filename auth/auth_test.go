package auth

import (
	"context"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/javiercbk/minesweeper/models"
	testHelpers "github.com/javiercbk/minesweeper/testing"
	"github.com/volatiletech/sqlboiler/boil"
)

// abcHashed is the bcrypt hash of the string "abc" (without quotes)
const abcHashed = "$2y$12$Fq0ne4S2xnhZTYE7p/veuOX3X6DlF1qZYeeHhK/PY39TP7//klYkW"

const jwtSecret = "wow"

func TestMain(m *testing.M) {
	testHelpers.InitializeDB(m)
}

func setUp(ctx context.Context, t *testing.T) (*Handler, *models.Player) {
	logger := testHelpers.NullLogger()
	db, err := testHelpers.DB()
	if err != nil {
		t.Fatalf("error connecting to database: %v\n", err)
	}
	testPlayer := &models.Player{
		Name:     "abc",
		Password: abcHashed,
	}
	err = testPlayer.Insert(ctx, db, boil.Infer())
	if err != nil {
		t.Fatalf("error inserting test user: %v\n", err)
	}
	return NewHandler(logger, db), testPlayer
}

func TestAuth(t *testing.T) {
	ctx := context.Background()
	handler, testPlayer := setUp(ctx, t)
	tests := []struct {
		Name     string
		Password string
		err      error
	}{
		{
			Name:     "abc",
			Password: "abc",
			err:      nil,
		},
		{
			Name:     "abc1",
			Password: "abc",
			err:      ErrBadCredentials,
		},
		{
			Name:     "abc",
			Password: "abc1",
			err:      ErrBadCredentials,
		},
		{
			Name:     "abc1",
			Password: "abc1",
			err:      ErrBadCredentials,
		},
	}
	for i, test := range tests {
		tokenResponse, err := handler.CreateToken(ctx, jwtSecret, Credentials{
			Name:     test.Name,
			Password: test.Password,
		})
		if test.err != err {
			t.Fatalf("failed test %d: expected error to be %v but was %v\n", i, test.err, err)
		}
		if err == nil {
			if tokenResponse.Token == "" {
				t.Fatalf("failed test %d: expected jwt token %s to be generated but was\n", i, tokenResponse.Token)
			}
			token, err := jwt.Parse(tokenResponse.Token, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				t.Fatalf("failed test %d: token %s could not be parsed\n", i, tokenResponse.Token)
			}
			c, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				t.Fatalf("failed test %d: claims could not be casted\n", i)
			}
			if int64(c["id"].(float64)) != testPlayer.ID {
				t.Fatalf("failed test %d: expected user id to be %d but was %d\n", i, testPlayer.ID, c["id"])
			}
			if c["name"] != testPlayer.Name {
				t.Fatalf("failed test %d: expected user name to be %s but was %s\n", i, testPlayer.Name, c["name"])
			}
		}
	}
}
