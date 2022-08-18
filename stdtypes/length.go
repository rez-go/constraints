package stdtypes

import (
	"fmt"

	"github.com/rez-go/constraints"
)

// TODO: generic slices and others
type lenable interface {
	~string | []byte
}

// Length returns a Constraint that will declare that a value
// as valid if its length is exactly as specified.
//
// API status: experimental
func Length[ValueT lenable](specifiedLength int) constraints.Constraint[ValueT] {
	if specifiedLength < 0 {
		panic("specifiedLength must be zero or a positive integer")
	}
	return &lengthConstraint[ValueT]{min: specifiedLength, max: specifiedLength}
}

// Length returns a Constraint that will declare that a value
// as valid if its length is at most equal to maxLength.
//
// API status: experimental
func MaxLength[ValueT lenable](maxLength int) constraints.Constraint[ValueT] {
	if maxLength < 0 {
		panic("maxLength must be zero or a positive integer")
	}
	return &lengthConstraint[ValueT]{min: -1, max: maxLength}
}

// Length returns a Constraint that will declare that a value
// as valid if its length is at most least to maxLength.
//
// API status: experimental
func MinLength[ValueT lenable](minLength int) constraints.Constraint[ValueT] {
	if minLength < 0 {
		panic("minLength must be zero or a positive integer")
	}
	return &lengthConstraint[ValueT]{min: minLength, max: -1}
}

func LengthRange[ValueT lenable](min int, max int) constraints.Constraint[ValueT] {
	return &lengthConstraint[ValueT]{min: min, max: max}
}

// lengthConstraint defines exact length Constraint.
type lengthConstraint[ValueT lenable] struct {
	min int
	max int
}

var (
	_ StringConstraint = lengthConstraint[string]{}
)

// ConstraintDescription conforms constraints.Constraint interface.
func (c lengthConstraint[ValueT]) ConstraintDescription() string {
	if c.min == c.max {
		return fmt.Sprintf("length %d", c.min)
	}
	if c.min == -1 {
		return fmt.Sprintf("max length %d", c.max)
	}
	if c.max == -1 {
		return fmt.Sprintf("min length %d", c.min)
	}
	return fmt.Sprintf("length betwen %d and %d", c.min, c.max)
}

// IsValid conforms Constraint interface.
func (c lengthConstraint[ValueT]) IsValid(v ValueT) bool {
	if c.min == c.max {
		return len(v) == c.min
	}
	if c.min == -1 {
		return len(v) <= c.max
	}
	if c.max == -1 {
		return len(v) >= c.min
	}
	return len(v) >= c.min && len(v) <= c.max
}
