package stdtypes

import (
	"fmt"

	typcon "golang.org/x/exp/constraints"

	"github.com/rez-go/constraints"
)

type numeric interface {
	typcon.Integer | typcon.Float
}

// NumericConstraint defines interface for int types.
type NumericConstraint[ValueT numeric] interface {
	constraints.Constraint[ValueT]
}

type Int8Constraint interface {
	NumericConstraint[int8]
}

var (
	IntPositive   = Positive[int]()
	Int8Positive  = Positive[int8]()
	Int16Positive = Positive[int16]()
	Int32Positive = Positive[int32]()
	Int64Positive = Positive[int64]()
)

var _ Int8Constraint = Int8Positive

// Positive creates a constraint that will declare an instance as valid
// if its value is positive.
func Positive[ValueT numeric]() NumericConstraint[ValueT] {
	return constraints.Func("positive", PositiveCheck[ValueT])
}

func PositiveCheck[ValueT numeric](v ValueT) bool { return v > 0 }

var (
	IntNegative   = Negative[int]()
	Int8Negative  = Negative[int8]()
	Int16Negative = Negative[int16]()
	Int32Negative = Negative[int32]()
	Int64Negative = Negative[int64]()
)

// Negative creates a constraint that will declare an instance as valid
// if its value is negative.
func Negative[ValueT numeric]() NumericConstraint[ValueT] {
	return constraints.Func("negative", NegativeCheck[ValueT])
}

func NegativeCheck[ValueT numeric](v ValueT) bool { return v < 0 }

var (
	IntEven   = Even[int]()
	Int8Even  = Even[int8]()
	Int16Even = Even[int16]()
	Int32Even = Even[int32]()
	Int64Even = Even[int64]()
)

// Even creates a constraint that will declare an instance as valid
// if its value is even.
func Even[ValueT typcon.Integer]() NumericConstraint[ValueT] {
	return constraints.Func("even", EvenCheck[ValueT])
}

func EvenCheck[ValueT typcon.Integer](v ValueT) bool { return (v & 1) == 0 }

// Odd creates a constraint that will declare an instance as valid
// if its value is odd.
func Odd[ValueT typcon.Integer]() NumericConstraint[ValueT] {
	return constraints.Func("odd", OddCheck[ValueT])
}

func OddCheck[ValueT typcon.Integer](v ValueT) bool { return (v & 1) == 1 }

// PowerOfTwo creates a constraint that will declare an instance as valid
// if its value is power of two.
func PowerOfTwo[ValueT typcon.Integer]() NumericConstraint[ValueT] {
	return constraints.Func("power of two", PowerOfTwoCheck[ValueT])
}

func PowerOfTwoCheck[ValueT typcon.Integer](v ValueT) bool {
	return v > 0 && (v&(v-1)) == 0
}

// Min creates a Constraint which will declare an instance is valid
// if its value is greater than or equal to refValue.
func Min[
	ValueT numeric,
](refValue ValueT) NumericConstraint[ValueT] {
	return constraints.Func(
		fmt.Sprintf("min %v", refValue),
		func(v ValueT) bool { return v >= refValue })
}

// Max creates a Constraint which will declare an instance is valid
// if its value is less than or equal to refValue.
func Max[
	ValueT numeric,
](refValue ValueT) NumericConstraint[ValueT] {
	return constraints.Func(
		fmt.Sprintf("max %v", refValue),
		func(v ValueT) bool { return v <= refValue })
}

// // Match creates an Constraint which will declare a value is valid
// // if it matches refValue.
// func Match[
// 	ValueT numeric,
// ](refValue ValueT) NumericConstraint[ValueT] {
// 	return &orderedOp[ValueT]{constraints.EqOpEqual, refValue}
// }

// Equals creates an Constraint which an instance will be
// declared as valid if its value is equal to refValue.
func Equals[
	ValueT numeric,
](refValue ValueT) NumericConstraint[ValueT] {
	return &orderedOp[ValueT]{constraints.EqOpEqual, refValue}
}

// Op creates an Constraint which comprise of an operator and the
// reference value.
func Op[
	ValueT numeric,
](op constraints.EqOp, refValue ValueT) NumericConstraint[ValueT] {
	return &orderedOp[ValueT]{op, refValue}
}

// LessThan creates an Constraint which an instance will be
// declared as valid if its value is less than refValue.
func LessThan[
	ValueT typcon.Ordered,
](refValue ValueT) constraints.Constraint[ValueT] {
	return &orderedOp[ValueT]{ref: refValue, op: constraints.EqOpLess}
}

// LessThanOrEqualTo creates an Constraint which an instance will be
// declared as valid if its value is less than or equal to refValue.
func LessThanOrEqualTo[
	ValueT numeric,
](refValue ValueT) NumericConstraint[ValueT] {
	return &orderedOp[ValueT]{ref: refValue, op: constraints.EqOpLessOrEqual}
}

// GreaterThan creates an Constraint which an instance will be declared
// as valid if its value is greater than refValue.
func GreaterThan[
	ValueT numeric,
](refValue ValueT) NumericConstraint[ValueT] {
	return &orderedOp[ValueT]{ref: refValue, op: constraints.EqOpGreater}
}

// GreaterThanOrEqualTo creates an Constraint which an instance will be
// declared as valid if its value is greater than or equal to refValue.
func GreaterThanOrEqualTo[
	ValueT numeric,
](refValue ValueT) NumericConstraint[ValueT] {
	return &orderedOp[ValueT]{ref: refValue, op: constraints.EqOpGreaterOrEqual}
}

var (
	_ NumericConstraint[int64] = &orderedOp[int64]{}
)

type orderedOp[ValueT typcon.Ordered] struct {
	op  constraints.EqOp
	ref ValueT
}

func (c orderedOp[ValueT]) ConstraintDescription() string {
	return fmt.Sprintf(c.op.StringFormat(), c.ref)
}

func (c orderedOp[ValueT]) IsValid(v ValueT) bool {
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
	return false
}
