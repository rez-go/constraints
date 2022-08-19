package stdtypes

import (
	"testing"

	"github.com/rez-go/constraints"
)

func TestEquals(t *testing.T) {
	eq0 := constraints.Match(0)
	assertEq(t, "match 0", eq0.ConstraintDescription())
	assertEq(t, true, eq0.IsValid(0))
	assertEq(t, false, eq0.IsValid(1))
	assertEq(t, false, eq0.IsValid(-1))
}

func TestLessThan(t *testing.T) {
	lt0 := constraints.LessThan(0)
	assertEq(t, "less than 0", lt0.ConstraintDescription())
	assertEq(t, false, lt0.IsValid(0))
	assertEq(t, false, lt0.IsValid(1))
	assertEq(t, true, lt0.IsValid(-1))
}

func TestPositive(t *testing.T) {
	positiveInt := IntPositive
	assertEq(t, "positive", positiveInt.ConstraintDescription())
	assertEq(t, false, positiveInt.IsValid(0))
	assertEq(t, true, positiveInt.IsValid(1))
	assertEq(t, false, positiveInt.IsValid(-1))
}

func TestNegative(t *testing.T) {
	negativeInt := IntNegative
	assertEq(t, "negative", negativeInt.ConstraintDescription())
	assertEq(t, false, negativeInt.IsValid(0))
	assertEq(t, false, negativeInt.IsValid(1))
	assertEq(t, true, negativeInt.IsValid(-1))
}

func TestEven(t *testing.T) {
	evenInt64 := Int64Even
	assertEq(t, "even", evenInt64.ConstraintDescription())
	assertEq(t, true, evenInt64.IsValid(0))
	assertEq(t, false, evenInt64.IsValid(1))
	assertEq(t, true, evenInt64.IsValid(2))
	assertEq(t, false, evenInt64.IsValid(-1))
	assertEq(t, true, evenInt64.IsValid(-2))
}

func TestNumericOdd(t *testing.T) {
	intOdd := Odd[int]()
	assertEq(t, "odd", intOdd.ConstraintDescription())
	assertEq(t, false, intOdd.IsValid(0))
	assertEq(t, true, intOdd.IsValid(1))
	assertEq(t, false, intOdd.IsValid(2))
	assertEq(t, true, intOdd.IsValid(-1))
	assertEq(t, false, intOdd.IsValid(-2))
}

func TestPowerOfTwo(t *testing.T) {
	powerOfTwoInt64 := PowerOfTwo[int64]()
	assertEq(t, "power of two", powerOfTwoInt64.ConstraintDescription())
	assertEq(t, false, powerOfTwoInt64.IsValid(0))
	assertEq(t, true, powerOfTwoInt64.IsValid(1))
	assertEq(t, true, powerOfTwoInt64.IsValid(2))
	assertEq(t, false, powerOfTwoInt64.IsValid(-2))
	assertEq(t, true, powerOfTwoInt64.IsValid(4))
}
