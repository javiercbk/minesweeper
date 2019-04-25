package algebra

import "errors"

// OperationType is a minesweeper algebra operation type
type OperationType int

// GameState is the current state of a game
type GameState int

const (
	// OpUnknown is an unknown operation used to identify a void operation
	OpUnknown = iota
	// OpReveal is the reveal operation type
	OpReveal
	// OpMark is the mark operation type
	OpMark
	// OpCompose is the compose operation type
	OpCompose
)

const (
	// StateOpen is the only state where a game is playable
	StateOpen GameState = iota
	// StateLoss is the only state where a game has concluded and it was lost
	StateLoss
	// StateWon is the only state where a game has concluded and it was won
	StateWon
)

// ErrUnknownOperation is returned when the operation does not exist
var ErrUnknownOperation = errors.New("unknown operation")

// MineProximity is the mine proximity value of a point in space
type MineProximity int

// operationExecution is the behaviour of all the operations of the minesweep algebra
type operationExecution func(state GameState, mineProximity MineProximity) (MineProximity, error)

// reveal is the operation that reveals a point in the board.
func reveal(state GameState, mineProximity MineProximity) (MineProximity, error) {
	return 0, nil
}

// mark is the operation that marks a point in the board as a possible or certain mine.
func mark(state GameState, mineProximity MineProximity) (MineProximity, error) {
	return 0, nil
}

// Operation is the behaviour of all the operations of the minesweep algebra
type Operation struct {
	x      int
	y      int
	opType OperationType
	Exec   operationExecution
}

// NewOperation creates an Operation object
func NewOperation(opType OperationType, x, y int) (Operation, error) {
	var err error
	oper := Operation{
		x:      x,
		y:      y,
		opType: opType,
	}
	if opType == OpReveal {
		oper.Exec = reveal
	} else if opType == OpMark {
		oper.Exec = mark
	} else {
		err = ErrUnknownOperation
	}
	return oper, err
}

// CompositionResult is the result of a Compose operation
type CompositionResult struct {
	Apply  []Operation
	Delta1 Operation
	Delta2 Operation
}

// Compose composes two operations
func Compose(oper1, oper2 Operation) CompositionResult {
	result := CompositionResult{}
	if oper1.x == oper2.x && oper1.y == oper2.y {
		// the operation is on the same point, apply the first operation
		result.Apply = []Operation{oper1}
		if oper1.opType != oper2.opType {
			// if is a different, then apply the operation 1 in the delta 2
			result.Delta2 = oper1
		}
	} else {
		// reveal both points
		result.Apply = []Operation{oper1, oper2}
		result.Delta1 = oper2
		result.Delta2 = oper1
	}
	return result
}
