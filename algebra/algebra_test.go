package algebra

import (
	"testing"
)

type algebraTest struct {
	oper      operationExecution
	state     GameState
	proximity MineProximity
	expected  MineProximity
	err       error
}

func TestReveal(t *testing.T) {
	testTable := []algebraTest{
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(0),
			expected:  MineProximity(0),
			err:       nil,
		},
	}
	for i, test := range testTable {
		proximity, err := test.oper(test.state, test.proximity)
		if proximity != test.expected {
			t.Fatalf("test %d failed: expected proximity %d but got %d", i, test.expected, proximity)
		}
		if err != test.err {
			t.Fatalf("test %d failed: expected err %v but got %v", i, test.err, err)
		}
	}
}
