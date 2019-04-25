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
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(1),
			expected:  MineProximity(1),
			err:       nil,
		},
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(2),
			expected:  MineProximity(2),
			err:       nil,
		},
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(3),
			expected:  MineProximity(3),
			err:       nil,
		},
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(4),
			expected:  MineProximity(4),
			err:       nil,
		},
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(5),
			expected:  MineProximity(5),
			err:       nil,
		},
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(6),
			expected:  MineProximity(6),
			err:       nil,
		},
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(7),
			expected:  MineProximity(7),
			err:       nil,
		},
		{
			oper:      reveal,
			state:     StateOpen,
			proximity: MineProximity(8),
			expected:  MineProximity(0),
			err:       ErrRevealOperationOutOfBounds,
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
