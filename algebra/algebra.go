package algebra

import (
	"errors"
	"math"
)

// OperationType is a minesweeper algebra operation type
type OperationType int

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

var (
	// ErrUnknownOperation is returned when the operation does not exist
	ErrUnknownOperation = errors.New("unknown operation")
	// ErrOperationOutOfBounds is returned when an operand is out of bounds or unknown
	ErrOperationOutOfBounds = errors.New("reveal operation out of bounds")
)

// MineProximity is the mine proximity value of a point in space
type MineProximity int

// operationExecution is the behaviour of all the operations of the minesweep algebra
type operationExecution func(mineProximity MineProximity) (MineProximity, error)

// reveal is the operation that reveals a point in the board.
func reveal(mineProximity MineProximity) (MineProximity, error) {
	if mineProximity >= 0 && mineProximity <= 8 {
		return mineProximity, nil
	}
	if mineProximity >= -10 && mineProximity <= -1 {
		return MineProximity(math.Abs(float64(mineProximity)) - 1), nil
	}
	return 0, ErrOperationOutOfBounds
}

// mark is the operation that marks a point in the board as a possible or certain mine.
func mark(mineProximity MineProximity) (MineProximity, error) {
	if mineProximity >= 0 && mineProximity <= 8 {
		return mineProximity, nil
	}
	if mineProximity >= -10 && mineProximity <= -1 {
		return MineProximity(mineProximity - 10), nil
	}
	if mineProximity >= -20 && mineProximity <= -11 {
		return MineProximity(mineProximity - 10), nil
	}
	if mineProximity >= -30 && mineProximity <= -21 {
		return MineProximity(mineProximity + 20), nil
	}
	return 0, ErrOperationOutOfBounds
}

// Operation is the behaviour of all the operations of the minesweep algebra
type Operation struct {
	x      int
	y      int
	opType OperationType
	exec   operationExecution
}

// Exec executes the operation over a mine proximity
func (o Operation) Exec(mineProximity MineProximity) (MineProximity, error) {
	return o.exec(mineProximity)
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
		oper.exec = reveal
	} else if opType == OpMark {
		oper.exec = mark
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
