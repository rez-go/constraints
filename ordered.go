package constraints

import (
	"fmt"

	typecons "golang.org/x/exp/constraints"
)

// Min creates a Constraint which will declare an instance is valid
// if its value is greater than or equal to refValue.
func Min[
	ValueT typecons.Ordered,
](refValue ValueT) OrderedConstraint[ValueT] {
	return Func(
		fmt.Sprintf("min %v", refValue),
		GreaterThanOrEqualTo(refValue).IsValid)
}

// Max creates a Constraint which will declare an instance is valid
// if its value is less than or equal to refValue.
func Max[
	ValueT typecons.Ordered,
](refValue ValueT) OrderedConstraint[ValueT] {
	return Func(
		fmt.Sprintf("max %v", refValue),
		LessThanOrEqualTo(refValue).IsValid)
}

// LessThan creates an Constraint which an instance will be
// declared as valid if its value is less than refValue.
func LessThan[
	ValueT typecons.Ordered,
](refValue ValueT) OrderedConstraint[ValueT] {
	return &relOpConstraint[ValueT]{ref: refValue, op: relOpLess}
}

// LessThanOrEqualTo creates an Constraint which an instance will be
// declared as valid if its value is less than or equal to refValue.
func LessThanOrEqualTo[
	ValueT typecons.Ordered,
](refValue ValueT) OrderedConstraint[ValueT] {
	return &relOpConstraint[ValueT]{ref: refValue, op: relOpLessOrEqual}
}

// GreaterThan creates an Constraint which an instance will be declared
// as valid if its value is greater than refValue.
func GreaterThan[
	ValueT typecons.Ordered,
](refValue ValueT) OrderedConstraint[ValueT] {
	return &relOpConstraint[ValueT]{ref: refValue, op: relOpGreater}
}

// GreaterThanOrEqualTo creates an Constraint which an instance will be
// declared as valid if its value is greater than or equal to refValue.
func GreaterThanOrEqualTo[
	ValueT typecons.Ordered,
](refValue ValueT) OrderedConstraint[ValueT] {
	return &relOpConstraint[ValueT]{ref: refValue, op: relOpGreaterOrEqual}
}

type OrderedConstraint[ValueT typecons.Ordered] interface {
	Constraint[ValueT]
}

var (
	_ OrderedConstraint[int] = &relOpConstraint[int]{}
	_ OrderedConstraint[int] = relOpConstraint[int]{}
)

type relOpConstraint[ValueT typecons.Ordered] struct {
	op  relOp
	ref ValueT
}

func (c relOpConstraint[ValueT]) ConstraintDescription() string {
	return fmt.Sprintf(c.op.StringFormat(), c.ref)
}

func (c relOpConstraint[ValueT]) IsValid(v ValueT) bool {
	switch c.op {
	case relOpEqual:
		return v == c.ref
	case relOpNotEqual:
		return v != c.ref
	case relOpLess:
		return v < c.ref
	case relOpLessOrEqual:
		return v <= c.ref
	case relOpGreater:
		return v > c.ref
	case relOpGreaterOrEqual:
		return v >= c.ref
	}
	return false
}

// An relOp specifies relational operator.
type relOp int

// Supported relational operators.
const (
	relOpEqual relOp = iota
	relOpNotEqual
	relOpLess
	relOpLessOrEqual
	relOpGreater
	relOpGreaterOrEqual
)

func (op relOp) String() string {
	switch op {
	case relOpEqual:
		return "equal"
	case relOpNotEqual:
		return "not equal"
	case relOpLess:
		return "less"
	case relOpLessOrEqual:
		return "less or equal"
	case relOpGreater:
		return "greater"
	case relOpGreaterOrEqual:
		return "greater or equal"
	}
	return ""
}

// Symbol returns representative symbol of the operator.
func (op relOp) Symbol() string {
	switch op {
	case relOpEqual:
		return "="
	case relOpNotEqual:
		return "≠"
	case relOpLess:
		return "<"
	case relOpLessOrEqual:
		return "≤"
	case relOpGreater:
		return ">"
	case relOpGreaterOrEqual:
		return "≥"
	}
	return "?"
}

// StringFormat returns a string which could be used to in *printf functions.
// The string format expects a value to be passed.
func (op relOp) StringFormat() string {
	switch op {
	case relOpEqual:
		return "equals %v"
	case relOpNotEqual:
		return "not equal to %v"
	case relOpLess:
		return "less than %v"
	case relOpLessOrEqual:
		return "less than or equal to %v"
	case relOpGreater:
		return "greater than %v"
	case relOpGreaterOrEqual:
		return "greater than or equal to %v"
	}
	return "(%v)"
}
