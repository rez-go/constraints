package ints

import (
	"fmt"

	"github.com/rez-go/constraints"
)

// Constraint defines interface for int types.
type Constraint interface {
	constraints.Constraint
	IsValid(v int64) bool
	ValidOrError(v int64) constraints.Error
}

var (
	// Positive is a constraint which will declare an instance as valid
	// if its value is positive.
	Positive = ConstraintFunc(
		"positive",
		func(v int64) bool { return v > 0 })

	// Negative is a constraint which will declare an instance as valid
	// if its value is negative.
	Negative = ConstraintFunc(
		"negative",
		func(v int64) bool { return v < 0 })

	// Even is a constraint which will declare an instance as valid
	// if its value is even.
	Even = ConstraintFunc(
		"even",
		func(v int64) bool { return (v & 1) == 0 })

	// Odd is a constraint which will declare an instance as valid
	// if its value is odd.
	Odd = ConstraintFunc(
		"odd",
		func(v int64) bool { return (v & 1) == 1 })

	// PowerOfTwo is a constraint which will declare an instance as valid
	// if its value is power of two.
	PowerOfTwo = ConstraintFunc(
		"power of two",
		func(v int64) bool {
			return v > 0 && (v&(v-1)) == 0
		})
)

// Min creates a Constraint which will declare an instance is valid
// if its value is greater than or equal to refValue.
func Min(refValue int64) Constraint {
	return ConstraintFunc(
		fmt.Sprintf("min %d", refValue),
		func(v int64) bool { return v >= refValue })
}

// Max creates a Constraint which will declare an instance is valid
// if its value is less than or equal to refValue.
func Max(refValue int64) Constraint {
	return ConstraintFunc(
		fmt.Sprintf("max %d", refValue),
		func(v int64) bool { return v <= refValue })
}

// Const creates an IntConstraint which will declare a value is valid
// if it matches refValue.
func Const(refValue int64) Constraint {
	return &intOp{constraints.EqOpEqual, refValue}
}

// Equals creates an IntConstraint which an instance will be
// declared as valid if its value is equal to refValue.
func Equals(refValue int64) Constraint {
	return &intOp{constraints.EqOpEqual, refValue}
}

// Op creates an IntConstraint which comprise of an operator and the
// reference value.
func Op(op constraints.EqOp, refValue int64) Constraint {
	return &intOp{op, refValue}
}

// LessThan creates an IntConstraint which an instance will be
// declared as valid if its value is less than refValue.
func LessThan(refValue int64) Constraint {
	return &intOp{ref: refValue, op: constraints.EqOpLess}
}

// LessThanOrEqualTo creates an IntConstraint which an instance will be
// declared as valid if its value is less than or equal to refValue.
func LessThanOrEqualTo(refValue int64) Constraint {
	return &intOp{ref: refValue, op: constraints.EqOpLessOrEqual}
}

// GreaterThan creates an IntConstraint which an instance will be declared
// as valid if its value is greater than refValue.
func GreaterThan(refValue int64) Constraint {
	return &intOp{ref: refValue, op: constraints.EqOpGreater}
}

// GreaterThanOrEqualTo creates an IntConstraint which an instance will be
// declared as valid if its value is greater than or equal to refValue.
func GreaterThanOrEqualTo(refValue int64) Constraint {
	return &intOp{ref: refValue, op: constraints.EqOpGreaterOrEqual}
}

var (
	_ Constraint = &intOp{}
)

type intOp struct {
	op  constraints.EqOp
	ref int64
}

func (c *intOp) ConstraintDescription() string {
	return fmt.Sprintf(c.op.StringFormat(), c.ref)
}

func (c *intOp) IsValid(v int64) bool {
	if c != nil {
		switch c.op {
		case constraints.EqOpEqual:
			return v == c.ref
		case constraints.EqOpNotEqual:
			return v != c.ref
		case constraints.EqOpLess:
			return v < c.ref
		case constraints.EqOpLessOrEqual:
			return v <= c.ref
		case constraints.EqOpGreater:
			return v > c.ref
		case constraints.EqOpGreaterOrEqual:
			return v >= c.ref
		}
	}
	return false
}

func (c *intOp) ValidOrError(v int64) constraints.Error {
	if c != nil && c.IsValid(v) {
		return nil
	}
	return constraints.ViolationError(c)
}

// ValidatorFunc is an adapter to allowe use of ordinary functions
// as a validator in constraints.
type ValidatorFunc func(i int64) bool

// ConstraintFunc creates a constraint from a validator function.
func ConstraintFunc(desc string, fn func(v int64) bool) Constraint {
	return &constraintFunc{desc, fn}
}

var (
	_ Constraint = &constraintFunc{}
)

type constraintFunc struct {
	desc string
	fn   ValidatorFunc
}

func (c *constraintFunc) ConstraintDescription() string {
	if c != nil {
		return c.desc
	}
	return "<unknown int constraint>"
}

func (c *constraintFunc) IsValid(v int64) bool {
	if c != nil {
		return c.fn(v)
	}
	return false
}

func (c *constraintFunc) ValidOrError(v int64) constraints.Error {
	if c != nil && c.IsValid(v) {
		return nil
	}
	return constraints.ViolationError(c)
}
