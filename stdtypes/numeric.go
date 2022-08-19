package stdtypes

import (
	typecons "golang.org/x/exp/constraints"

	"github.com/rez-go/constraints"
)

type numeric interface {
	typecons.Integer | typecons.Float
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
func Even[ValueT typecons.Integer]() NumericConstraint[ValueT] {
	return constraints.Func("even", EvenCheck[ValueT])
}

func EvenCheck[ValueT typecons.Integer](v ValueT) bool { return (v & 1) == 0 }

// Odd creates a constraint that will declare an instance as valid
// if its value is odd.
func Odd[ValueT typecons.Integer]() NumericConstraint[ValueT] {
	return constraints.Func("odd", OddCheck[ValueT])
}

func OddCheck[ValueT typecons.Integer](v ValueT) bool { return (v & 1) == 1 }

// PowerOfTwo creates a constraint that will declare an instance as valid
// if its value is power of two.
func PowerOfTwo[ValueT typecons.Integer]() NumericConstraint[ValueT] {
	return constraints.Func("power of two", PowerOfTwoCheck[ValueT])
}

func PowerOfTwoCheck[ValueT typecons.Integer](v ValueT) bool {
	return v > 0 && (v&(v-1)) == 0
}
