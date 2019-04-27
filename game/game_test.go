package game

import (
	"testing"
)

func TestSiblingPoints(t *testing.T) {
	testBoard := board{
		rows:  4,
		cols:  4,
		board: [][]int{{-1, -1, -1, -1}, {-1, -1, -1, -1}, {-1, -1, -1, -1}, {-1, -1, -1, -1}},
	}
	tests := []struct {
		b        board
		row      int
		col      int
		siblings []boardPoint
	}{
		{
			b:   testBoard,
			row: 0,
			col: 0,
			siblings: []boardPoint{
				{
					row: 0,
					col: 1,
				},
				{
					row: 1,
					col: 1,
				},
				{
					row: 1,
					col: 0,
				},
			},
		},
		{
			b:   testBoard,
			row: 1,
			col: 0,
			siblings: []boardPoint{
				{
					row: 0,
					col: 0,
				},
				{
					row: 0,
					col: 1,
				},
				{
					row: 1,
					col: 1,
				},
				{
					row: 2,
					col: 0,
				},
				{
					row: 2,
					col: 1,
				},
			},
		},
		{
			b:   testBoard,
			row: 0,
			col: 1,
			siblings: []boardPoint{
				{
					row: 0,
					col: 0,
				},
				{
					row: 1,
					col: 0,
				},
				{
					row: 1,
					col: 1,
				},
				{
					row: 0,
					col: 2,
				},
				{
					row: 1,
					col: 2,
				},
			},
		},
		{
			b:   testBoard,
			row: 1,
			col: 1,
			siblings: []boardPoint{
				{
					row: 0,
					col: 0,
				},
				{
					row: 0,
					col: 1,
				},
				{
					row: 0,
					col: 2,
				},
				{
					row: 1,
					col: 0,
				},
				{
					row: 1,
					col: 2,
				},
				{
					row: 2,
					col: 0,
				},
				{
					row: 2,
					col: 1,
				},
				{
					row: 2,
					col: 2,
				},
			},
		},
		{
			b:   testBoard,
			row: 3,
			col: 3,
			siblings: []boardPoint{
				{
					row: 2,
					col: 2,
				},
				{
					row: 2,
					col: 3,
				},
				{
					row: 3,
					col: 2,
				},
			},
		},
		{
			b:   testBoard,
			row: 0,
			col: 3,
			siblings: []boardPoint{
				{
					row: 0,
					col: 2,
				},
				{
					row: 1,
					col: 2,
				},
				{
					row: 1,
					col: 3,
				},
			},
		},
	}
	for i, test := range tests {
		siblings := test.b.siblingPoints(test.row, test.col)
		sLen := len(siblings)
		expectedLen := len(test.siblings)
		if sLen != expectedLen {
			t.Fatalf("test %d, failed: expected siblings to be %d, but was %d\n", i, expectedLen, sLen)
		}
		for _, s := range test.siblings {
			found := false
			for _, calculatedSibling := range siblings {
				if calculatedSibling.row == s.row && calculatedSibling.col == s.col {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("test %d, failed: expected sibling row %d col %d to be found but it wasn't\n", i, s.row, s.col)
			}
		}
	}
}

func TestPlaceMine(t *testing.T) {
	testBoardFactory := func() *board {
		return &board{
			rows: 4,
			cols: 4,
			board: [][]int{
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
		}
	}
	tests := []struct {
		b     *board
		row   int
		col   int
		state [][]int
	}{
		{
			b:   testBoardFactory(),
			row: 0,
			col: 0,
			state: [][]int{
				{-10, -2, -1, -1},
				{-2, -2, -1, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
		},
		{
			b:   testBoardFactory(),
			row: 1,
			col: 0,
			state: [][]int{
				{-2, -2, -1, -1},
				{-10, -2, -1, -1},
				{-2, -2, -1, -1},
				{-1, -1, -1, -1},
			},
		},
		{
			b:   testBoardFactory(),
			row: 2,
			col: 2,
			state: [][]int{
				{-1, -1, -1, -1},
				{-1, -2, -2, -2},
				{-1, -2, -10, -2},
				{-1, -2, -2, -2},
			},
		},
	}
	for i, test := range tests {
		test.b.placeMine(test.row, test.col)
		for r := range test.b.board {
			for c := range test.b.board[r] {
				if test.b.board[r][c] != test.state[r][c] {
					t.Fatalf("test %d, failed: expected row %d col %d to be %d but was %d\n", i, r, c, test.state[r][c], test.b.board[r][c])
				}
			}
		}
	}
}

func TestPlaceSeveralMines(t *testing.T) {
	testBoardFactory := func() *board {
		return &board{
			rows: 4,
			cols: 4,
			board: [][]int{
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
		}
	}
	tests := []struct {
		b     *board
		mines []boardPoint
		state [][]int
	}{
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}},
			state: [][]int{
				{-10, -10, -2, -1},
				{-3, -3, -2, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 3}},
			state: [][]int{
				{-10, -2, -2, -10},
				{-2, -2, -2, -2},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}},
			state: [][]int{
				{-10, -10, -2, -1},
				{-3, -3, -2, -1},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}, {0, 2}},
			state: [][]int{
				{-10, -10, -10, -2},
				{-3, -4, -3, -2},
				{-1, -1, -1, -1},
				{-1, -1, -1, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}, {0, 2}, {1, 0}},
			state: [][]int{
				{-10, -10, -10, -2},
				{-10, -5, -3, -2},
				{-2, -2, -1, -1},
				{-1, -1, -1, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}},
			state: [][]int{
				{-10, -10, -10, -3},
				{-10, -6, -10, -3},
				{-2, -3, -2, -2},
				{-1, -1, -1, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}},
			state: [][]int{
				{-10, -10, -10, -3},
				{-10, -7, -10, -3},
				{-10, -4, -2, -2},
				{-2, -2, -1, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}},
			state: [][]int{
				{-10, -10, -10, -3},
				{-10, -8, -10, -3},
				{-10, -10, -3, -2},
				{-3, -3, -2, -1},
			},
		},
		{
			b:     testBoardFactory(),
			mines: []boardPoint{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
			state: [][]int{
				{-10, -10, -10, -3},
				{-10, -9, -10, -4},
				{-10, -10, -10, -3},
				{-3, -4, -3, -2},
			},
		},
	}
	for i, test := range tests {
		for _, m := range test.mines {
			test.b.placeMine(m.row, m.col)
		}
		for r := range test.b.board {
			for c := range test.b.board[r] {
				if test.b.board[r][c] != test.state[r][c] {
					t.Fatalf("test %d, failed: expected row %d col %d to be %d but was %d\n", i, r, c, test.state[r][c], test.b.board[r][c])
				}
			}
		}
	}
}

func TestNewBoard(t *testing.T) {
	tests := []struct {
		rows  int
		cols  int
		mines int
		err   error
	}{
		{
			rows:  -1,
			cols:  -1,
			mines: -1,
			err:   ErrInvalidRowCols,
		},
		{
			rows:  0,
			cols:  -1,
			mines: -1,
			err:   ErrInvalidRowCols,
		},
		{
			rows:  4,
			cols:  -1,
			mines: -1,
			err:   ErrInvalidRowCols,
		},
		{
			rows:  4,
			cols:  0,
			mines: -1,
			err:   ErrInvalidRowCols,
		},
		{
			rows:  4,
			cols:  4,
			mines: -1,
			err:   ErrNoneMines,
		},
		{
			rows:  4,
			cols:  4,
			mines: 0,
			err:   ErrNoneMines,
		},
		{
			rows:  4,
			cols:  4,
			mines: 16,
			err:   ErrTooManyMines,
		},
		{
			rows:  4,
			cols:  4,
			mines: 17,
			err:   ErrTooManyMines,
		},
		{
			rows:  4,
			cols:  4,
			mines: 4,
		},
	}
	for i, test := range tests {
		board, err := NewBoard(test.rows, test.cols, test.mines)
		if err != test.err {
			t.Fatalf("test %d, failed: expected err to be %v but was %v\n", i, test.err, err)
		} else if err == nil {
			mineCount := test.mines
			for r := range board {
				for c := range board[r] {
					if board[r][c] == -10 {
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

func TestArrayToBoardAndBack(t *testing.T) {
	board := [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}}
	rows := len(board)
	cols := len(board[0])
	flatBoard := make([]int, rows*cols)
	for i := range board {
		for j := range board[i] {
			index := boardToArrayPoint(i, j, cols)
			flatBoard[index] = board[i][j]
		}
	}
	sameBoard := make([][]int, rows)
	for i := range sameBoard {
		sameBoard[i] = make([]int, cols)
	}
	for i := range flatBoard {
		row, col := arrayToBoardPoint(i, cols)
		sameBoard[row][col] = flatBoard[i]
	}
	for i := range board {
		for j := range board[i] {
			if board[i][j] != sameBoard[i][j] {
				t.Fatalf("expected row %d col %d to be %d but was %d", i, j, board[i][j], sameBoard[i][j])
			}
		}
	}
}
