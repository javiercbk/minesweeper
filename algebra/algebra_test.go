package algebra

import (
	"fmt"
	"testing"
)

type algebraTest struct {
	oper      operationExecution
	proximity MineProximity
	expected  MineProximity
	err       error
}

func TestReveal(t *testing.T) {
	testTable := []algebraTest{
		{
			oper:      reveal,
			proximity: MineProximity(0),
			expected:  MineProximity(0),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(1),
			expected:  MineProximity(1),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(2),
			expected:  MineProximity(2),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(3),
			expected:  MineProximity(3),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(4),
			expected:  MineProximity(4),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(5),
			expected:  MineProximity(5),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(6),
			expected:  MineProximity(6),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(7),
			expected:  MineProximity(7),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(8),
			expected:  MineProximity(0),
			err:       ErrOperationOutOfBounds,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-1),
			expected:  MineProximity(0),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-2),
			expected:  MineProximity(1),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-3),
			expected:  MineProximity(2),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-4),
			expected:  MineProximity(3),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-5),
			expected:  MineProximity(4),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-6),
			expected:  MineProximity(5),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-7),
			expected:  MineProximity(6),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-8),
			expected:  MineProximity(7),
			err:       nil,
		},
		{
			oper:      reveal,
			proximity: MineProximity(-9),
			expected:  MineProximity(0),
			err:       ErrOperationOutOfBounds,
		},
	}
	for i, test := range testTable {
		proximity, err := test.oper(test.proximity)
		if proximity != test.expected {
			t.Fatalf("test %d failed: expected proximity %d but got %d", i, test.expected, proximity)
		}
		if err != test.err {
			t.Fatalf("test %d failed: expected err %v but got %v", i, test.err, err)
		}
	}
}

func TestMark(t *testing.T) {
	testTable := []algebraTest{
		{
			oper:      mark,
			proximity: MineProximity(0),
			expected:  MineProximity(0),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(1),
			expected:  MineProximity(1),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(2),
			expected:  MineProximity(2),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(3),
			expected:  MineProximity(3),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(4),
			expected:  MineProximity(4),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(5),
			expected:  MineProximity(5),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(6),
			expected:  MineProximity(6),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(7),
			expected:  MineProximity(7),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(8),
			expected:  MineProximity(0),
			err:       ErrOperationOutOfBounds,
		},
		{
			oper:      mark,
			proximity: MineProximity(-1),
			expected:  MineProximity(-11),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-2),
			expected:  MineProximity(-12),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-3),
			expected:  MineProximity(-13),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-4),
			expected:  MineProximity(-14),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-5),
			expected:  MineProximity(-15),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-6),
			expected:  MineProximity(-16),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-7),
			expected:  MineProximity(-17),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-8),
			expected:  MineProximity(-18),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-9),
			expected:  MineProximity(0),
			err:       ErrOperationOutOfBounds,
		},
		{
			oper:      mark,
			proximity: MineProximity(-11),
			expected:  MineProximity(-21),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-12),
			expected:  MineProximity(-22),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-13),
			expected:  MineProximity(-23),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-14),
			expected:  MineProximity(-24),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-15),
			expected:  MineProximity(-25),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-16),
			expected:  MineProximity(-26),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-17),
			expected:  MineProximity(-27),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-18),
			expected:  MineProximity(-28),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-19),
			expected:  MineProximity(0),
			err:       ErrOperationOutOfBounds,
		},
		{
			oper:      mark,
			proximity: MineProximity(-21),
			expected:  MineProximity(-1),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-22),
			expected:  MineProximity(-2),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-23),
			expected:  MineProximity(-3),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-24),
			expected:  MineProximity(-4),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-25),
			expected:  MineProximity(-5),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-26),
			expected:  MineProximity(-6),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-27),
			expected:  MineProximity(-7),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-28),
			expected:  MineProximity(-8),
			err:       nil,
		},
		{
			oper:      mark,
			proximity: MineProximity(-29),
			expected:  MineProximity(0),
			err:       ErrOperationOutOfBounds,
		},
	}
	for i, test := range testTable {
		proximity, err := test.oper(test.proximity)
		if proximity != test.expected {
			t.Fatalf("test %d failed: expected proximity %d but got %d", i, test.expected, proximity)
		}
		if err != test.err {
			t.Fatalf("test %d failed: expected err %v but got %v", i, test.err, err)
		}
	}
}

func compareOperation(expected, actual Operation) error {
	if expected.opType != actual.opType {
		return fmt.Errorf("expected operation type to be %d but was %d", expected.opType, actual.opType)
	}
	if expected.x != actual.x {
		return fmt.Errorf("expected x to be %d but was %d", expected.x, actual.x)
	}
	if expected.y != actual.y {
		return fmt.Errorf("expected x to be %d but was %d", expected.y, actual.y)
	}
	return nil
}

func TestCompose(t *testing.T) {
	reveal1, _ := NewOperation(OpReveal, 0, 0)
	reveal2, _ := NewOperation(OpReveal, 0, 1)
	mark1, _ := NewOperation(OpMark, 0, 0)
	mark2, _ := NewOperation(OpMark, 0, 1)
	testTable := []struct {
		oper1    Operation
		oper2    Operation
		expected CompositionResult
	}{
		{
			oper1: reveal1,
			oper2: reveal1,
			expected: CompositionResult{
				Apply: []Operation{reveal1},
			},
		},
		{
			oper1: mark1,
			oper2: mark1,
			expected: CompositionResult{
				Apply: []Operation{mark1},
			},
		},
		{
			oper1: mark1,
			oper2: mark2,
			expected: CompositionResult{
				Apply:  []Operation{mark1, mark2},
				Delta1: mark2,
				Delta2: mark1,
			},
		},
		{
			oper1: reveal1,
			oper2: reveal2,
			expected: CompositionResult{
				Apply:  []Operation{reveal1, reveal2},
				Delta1: reveal2,
				Delta2: reveal1,
			},
		},
		{
			oper1: reveal1,
			oper2: mark1,
			expected: CompositionResult{
				Apply:  []Operation{reveal1},
				Delta2: reveal1,
			},
		},
		{
			oper1: mark1,
			oper2: reveal1,
			expected: CompositionResult{
				Apply:  []Operation{mark1},
				Delta2: mark1,
			},
		},
		{
			oper1: mark1,
			oper2: reveal2,
			expected: CompositionResult{
				Apply:  []Operation{mark1, reveal2},
				Delta1: reveal2,
				Delta2: mark1,
			},
		},
	}
	for i, test := range testTable {
		composition := Compose(test.oper1, test.oper2)
		if len(composition.Apply) != len(test.expected.Apply) {
			t.Fatalf("test %d failed: expected apply to have len %d but got len %d", i, len(test.expected.Apply), len(composition.Apply))
		}
		for j := range composition.Apply {
			if composition.Apply[j].opType != test.expected.Apply[j].opType {
				t.Fatalf("test %d failed: expected apply operation index %d to be %d but was %d", i, j, test.expected.Apply[j].opType, composition.Apply[j].opType)
			}
			if composition.Apply[j].x != test.expected.Apply[j].x {
				t.Fatalf("test %d failed: expected apply operation index %d to be %d but was %d", i, j, test.expected.Apply[j].x, composition.Apply[j].x)
			}
			if composition.Apply[j].y != test.expected.Apply[j].y {
				t.Fatalf("test %d failed: expected apply operation index %d to be %d but was %d", i, j, test.expected.Apply[j].y, composition.Apply[j].y)
			}
		}
		err := compareOperation(test.expected.Delta1, composition.Delta1)
		if err != nil {
			t.Fatalf("test %d failed: delta 1 comparison failed with %s", i, err)
		}
		err = compareOperation(test.expected.Delta2, composition.Delta2)
		if err != nil {
			t.Fatalf("test %d failed: delta 2 comparison failed with %s", i, err)
		}
	}
}
