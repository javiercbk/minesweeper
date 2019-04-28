package game

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/javiercbk/minesweeper/http/security"
	"github.com/javiercbk/minesweeper/models"
	testHelpers "github.com/javiercbk/minesweeper/testing"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	extErrors "github.com/pkg/errors"
)

const uniqueNameConstaintName = "idx_players_name"
const abcHashed = "$2y$12$Fq0ne4S2xnhZTYE7p/veuOX3X6DlF1qZYeeHhK/PY39TP7//klYkW"
const username = "benchmarkUsername"

const gameRows = 100
const gameCols = 100
const gameMines = 99

func TestMain(m *testing.M) {
	testHelpers.InitializeDB(m)
}

func setUp(ctx context.Context, t testing.TB, name string) (*Handler, security.JWTUser) {
	logger := testHelpers.NullLogger()
	db, err := testHelpers.DB()
	if err != nil {
		t.Fatalf("error connecting to database: %v\n", err)
	}
	testPlayer := &models.Player{
		Name:     name,
		Password: abcHashed,
	}
	err = testPlayer.Insert(ctx, db, boil.Infer())
	if err != nil {
		cause := extErrors.Cause(err)
		if pgerr, ok := cause.(*pq.Error); ok {
			if pgerr.Constraint == uniqueNameConstaintName {
				testPlayer, err = models.Players(qm.Where("name = ?", name)).One(ctx, db)
				if err != nil {
					t.Fatalf("error retrieving test user: %v\n", err)
				}
			}
		} else {
			t.Fatalf("error inserting test user: %v\n", err)
		}
	}
	return NewHandler(logger, db), security.JWTUser{
		ID:   testPlayer.ID,
		Name: name,
	}
}

// 10000	   3775875 ns/op	  546135 B/op	     281 allocs/op
func BenchmarkArrayStorage(b *testing.B) {
	ctx := context.Background()
	handler, user := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := handler.CreateGame(ctx, user, &pGame, handler.arrayStorageStrategy)
		if err != nil {
			b.Fatalf("error creating game %v\n", err)
		}
	}
}

// 20	1206168492 ns/op	20531937 B/op	  530139 allocs/op
func BenchmarkTableStorage(b *testing.B) {
	ctx := context.Background()
	handler, user := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := handler.CreateGame(ctx, user, &pGame, handler.tableStorageStrategy)
		if err != nil {
			b.Fatalf("error creating game %v\n", err)
		}
	}
}

// 1000	   2034873 ns/op	 1227521 B/op	   10101 allocs/op
func BenchmarkArrayRowColRetrieval(b *testing.B) {
	ctx := context.Background()
	handler, user := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	err := handler.CreateGame(ctx, user, &pGame, handler.arrayStorageStrategy)
	if err != nil {
		b.Fatalf("error creating game %v\n", err)
	}
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		randomRow := random.Intn(gameRows - 1)
		randomCol := random.Intn(gameCols - 1)
		_, err = handler.arrayRowColRetrieval(ctx, user, pGame.ID, randomRow, randomCol)
		if err != nil {
			b.Fatalf("error retrieving row col %v\n", err)
		}
	}
}

// 10000	    177058 ns/op	    3324 B/op	      68 allocs/op
func BenchmarkTableRowColRetrieval(b *testing.B) {
	ctx := context.Background()
	handler, user := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	err := handler.CreateGame(ctx, user, &pGame, handler.tableStorageStrategy)
	if err != nil {
		b.Fatalf("error creating game %v\n", err)
	}
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		randomRow := random.Intn(gameRows - 1)
		randomCol := random.Intn(gameCols - 1)
		_, err = handler.tableRowColRetrieval(ctx, user, pGame.ID, randomRow, randomCol)
		if err != nil {
			b.Fatalf("error retrieving row col %v\n", err)
		}
	}
}

// 300	   4154761 ns/op	 1226859 B/op	    9955 allocs/op
func BenchmarkArrayRowColUpdate(b *testing.B) {
	ctx := context.Background()
	handler, user := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	err := handler.CreateGame(ctx, user, &pGame, handler.arrayStorageStrategy)
	if err != nil {
		b.Fatalf("error creating game %v\n", err)
	}
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		randomRow := random.Intn(gameRows - 1)
		randomCol := random.Intn(gameCols - 1)
		err = handler.arrayUpdateRowCol(ctx, user, pGame.ID, randomRow, randomCol, 0)
		if err != nil {
			b.Fatalf("error updating row col %v\n", err)
		}
	}
}

// 1000	   1636261 ns/op	    2326 B/op	      54 allocs/op
func BenchmarkTableRowColUpdate(b *testing.B) {
	ctx := context.Background()
	handler, user := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	err := handler.CreateGame(ctx, user, &pGame, handler.tableStorageStrategy)
	if err != nil {
		b.Fatalf("error creating game %v\n", err)
	}
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		randomRow := random.Intn(gameRows - 1)
		randomCol := random.Intn(gameCols - 1)
		err = handler.tableUpdateRowCol(ctx, user, pGame.ID, randomRow, randomCol, 0)
		if err != nil {
			b.Fatalf("error updating row col %v\n", err)
		}
	}
}
