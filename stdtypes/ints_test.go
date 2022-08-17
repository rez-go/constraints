package stdtypes

import (
	"testing"
)

func TestEquals(t *testing.T) {
	eq0 := Equals(0)
	if eq0.ConstraintDescription() != "equals 0" {
		t.Errorf(`expecting "equals 0", got %q`,
			eq0.ConstraintDescription())
	}
	if !eq0.IsValid(0) {
		t.Errorf(`!eq0.IsValid(0)`)
	}
	if eq0.IsValid(1) {
		t.Errorf(`eq0.IsValid(1)`)
	}
	if eq0.IsValid(-1) {
		t.Errorf(`eq0.IsValid(-1)`)
	}
}

func TestLessThan(t *testing.T) {
	lt0 := LessThan(0)
	if lt0.ConstraintDescription() != "less than 0" {
		t.Errorf(`expecting "less than 0", got %q`,
			lt0.ConstraintDescription())
	}
	if lt0.IsValid(0) {
		t.Errorf(`lt0.IsValid(0)`)
	}
	if lt0.IsValid(1) {
		t.Errorf(`lt0.IsValid(1)`)
	}
	if !lt0.IsValid(-1) {
		t.Errorf(`!lt0.IsValid(-1)`)
	}
}

func TestPositive(t *testing.T) {
	positiveInt64 := Int64Positive
	if positiveInt64.ConstraintDescription() != "positive" {
		t.Errorf(`expecting "positive", got %q`,
			positiveInt64.ConstraintDescription())
	}
	if positiveInt64.IsValid(0) {
		t.Errorf(`Positive.IsValid(0)`)
	}
	if !positiveInt64.IsValid(1) {
		t.Errorf(`!Positive.IsValid(1)`)
	}
	if positiveInt64.IsValid(-1) {
		t.Errorf(`Positive.IsValid(-1)`)
	}
}

func TestEven(t *testing.T) {
	evenInt64 := Int64Even
	if evenInt64.ConstraintDescription() != "even" {
		t.Errorf(`expecting "even", got %q`,
			evenInt64.ConstraintDescription())
	}
	if !evenInt64.IsValid(0) {
		t.Errorf(`!evenInt64.IsValid(0)`)
	}
	if evenInt64.IsValid(1) {
		t.Errorf(`evenInt64.IsValid(1)`)
	}
}

func TestPowerOfTwo(t *testing.T) {
	powerOfTwoInt64 := PowerOfTwo[int64]()
	if powerOfTwoInt64.ConstraintDescription() != "power of two" {
		t.Errorf(`expecting "power of two", got %q`,
			powerOfTwoInt64.ConstraintDescription())
	}
	if powerOfTwoInt64.IsValid(0) {
		t.Errorf(`powerOfTwoInt64.IsValid(0)`)
	}
	// 2^0 = 1
	if !powerOfTwoInt64.IsValid(1) {
		t.Errorf(`!powerOfTwoInt64.IsValid(1)`)
	}
	if !powerOfTwoInt64.IsValid(2) {
		t.Errorf(`!powerOfTwoInt64.IsValid(2)`)
	}
	if powerOfTwoInt64.IsValid(-2) {
		t.Errorf(`powerOfTwoInt64.IsValid(-2)`)
	}
	if !powerOfTwoInt64.IsValid(4) {
		t.Errorf(`!powerOfTwoInt64.IsValid(4)`)
	}
}
