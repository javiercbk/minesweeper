package game

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

// 10	 198611863 ns/op	 3883188 B/op	   50122 allocs/op
func BenchmarkCreateGame(b *testing.B) {
	ctx := context.Background()
	api, user, _ := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := api.CreateGame(ctx, user, &pGame)
		if err != nil {
			b.Fatalf("error creating game %v\n", err)
		}
	}
}

// 10000	    177058 ns/op	    3324 B/op	      68 allocs/op
func BenchmarkRetrieveRowCol(b *testing.B) {
	ctx := context.Background()
	api, user, _ := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	err := api.CreateGame(ctx, user, &pGame)
	if err != nil {
		b.Fatalf("error creating game %v\n", err)
	}
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		randomRow := random.Intn(gameRows - 1)
		randomCol := random.Intn(gameCols - 1)
		_, err = api.retrieveRowCol(ctx, user, pGame.ID, randomRow, randomCol)
		if err != nil {
			b.Fatalf("error retrieving row col %v\n", err)
		}
	}
}

// 1000	   1636261 ns/op	    2326 B/op	      54 allocs/op
func BenchmarkUpdateRowCol(b *testing.B) {
	ctx := context.Background()
	api, user, _ := setUp(ctx, b, username)
	pGame := ProspectGame{
		Rows:    gameRows,
		Cols:    gameCols,
		Mines:   gameMines,
		Private: false,
	}
	err := api.CreateGame(ctx, user, &pGame)
	if err != nil {
		b.Fatalf("error creating game %v\n", err)
	}
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// do not count first insertion time
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		randomRow := random.Intn(gameRows - 1)
		randomCol := random.Intn(gameCols - 1)
		err = api.updateRowCol(ctx, pGame.ID, randomRow, randomCol, 0)
		if err != nil {
			b.Fatalf("error updating row col %v\n", err)
		}
	}
}
