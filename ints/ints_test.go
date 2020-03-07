package ints

import (
	"testing"

	"github.com/rez-go/constraints"
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
	if err := eq0.ValidOrError(0); err != nil {
		t.Errorf(`err := eq0.ValidOrError(0); err != nil`)
	}
	err := eq0.ValidOrError(1)
	if err == nil {
		t.Errorf(`err must not be nil`)
	}
	if c := constraints.ViolatedConstraintFromError(err); c != eq0 {
		t.Errorf(`ViolatedConstraints must be eq0`)
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
	if Positive.ConstraintDescription() != "positive" {
		t.Errorf(`expecting "positive", got %q`,
			Positive.ConstraintDescription())
	}
	if Positive.IsValid(0) {
		t.Errorf(`Positive.IsValid(0)`)
	}
	if !Positive.IsValid(1) {
		t.Errorf(`!Positive.IsValid(1)`)
	}
	if Positive.IsValid(-1) {
		t.Errorf(`Positive.IsValid(-1)`)
	}
}

func TestEven(t *testing.T) {
	if Even.ConstraintDescription() != "even" {
		t.Errorf(`expecting "even", got %q`,
			Even.ConstraintDescription())
	}
	if !Even.IsValid(0) {
		t.Errorf(`!Even.IsValid(0)`)
	}
	if Even.IsValid(1) {
		t.Errorf(`Even.IsValid(1)`)
	}
}

func TestPowerOfTwo(t *testing.T) {
	if PowerOfTwo.ConstraintDescription() != "power of two" {
		t.Errorf(`expecting "power of two", got %q`,
			PowerOfTwo.ConstraintDescription())
	}
	if PowerOfTwo.IsValid(0) {
		t.Errorf(`PowerOfTwo.IsValid(0)`)
	}
	// 2^0 = 1
	if !PowerOfTwo.IsValid(1) {
		t.Errorf(`!PowerOfTwo.IsValid(1)`)
	}
	if !PowerOfTwo.IsValid(2) {
		t.Errorf(`!PowerOfTwo.IsValid(2)`)
	}
	if PowerOfTwo.IsValid(-2) {
		t.Errorf(`PowerOfTwo.IsValid(-2)`)
	}
	if !PowerOfTwo.IsValid(4) {
		t.Errorf(`!PowerOfTwo.IsValid(4)`)
	}
}
