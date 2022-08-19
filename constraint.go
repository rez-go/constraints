package constraints

import (
	"fmt"
	"strconv"
)

type ConstraintBase interface {
	// ConstraintDescription returns the description of the constraint, e.g.,
	// "5 or less".
	//
	// This is a porcelain interface.
	ConstraintDescription() string
}

// Constraint is base interface for all constraints.
type Constraint[ValueT any] interface {
	ConstraintBase

	IsValid(value ValueT) bool
}

//----

// ValidOrError tests the value v against constraint c. If the value is
// valid, it will return nil, otherwise, it'll return an error.
//
// If the constraint is a Set and the value violated any
// or all of the constraints, this method will return a Error which
// constraint contained is a new instance of Set contains only the
// violated constraints.
func ValidOrError[ValueT any](v ValueT, c Constraint[ValueT]) error {
	if c != nil {
		if cs, ok := c.(interface {
			ValidateAll(ValueT) []Constraint[ValueT]
		}); ok && cs != nil {
			violatedConstraints := cs.ValidateAll(v)
			if len(violatedConstraints) > 0 {
				return ViolationError[ValueT](Set(violatedConstraints...))
			}
			return nil
		}
		if c.IsValid(v) {
			return nil
		}
	}
	return ViolationError(c)
}

// Match creates a Constraint which will declare an instance as valid
// if its value matches refValue.
func Match[ValueT comparable](refValue ValueT) *matchConstraint[ValueT] {
	return &matchConstraint[ValueT]{refValue}
}

// matchConstraint defines constant value Constraint.
type matchConstraint[ValueT comparable] struct {
	refValue ValueT
}

var (
	_ Constraint[string] = matchConstraint[string]{}
)

// ConstraintDescription conforms constraints.Constraint interface.
func (c matchConstraint[ValueT]) ConstraintDescription() string {
	return fmt.Sprintf("match %s", valueLiteralString(c.refValue))
}

// IsValid conforms Constraint interface.
func (c matchConstraint[ValueT]) IsValid(v ValueT) bool {
	return v == c.refValue
}

func Negate[ValueT any](c Constraint[ValueT], descOverride string) Constraint[ValueT] {
	desc := descOverride
	if desc == "" {
		desc = "not " + c.ConstraintDescription()
	}
	return &constraintFunc[ValueT]{negate: true, desc: desc, fn: c.IsValid}
}

func valueLiteralString(v any) string {
	switch tv := v.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case rune:
		if strconv.IsPrint(tv) {
			return fmt.Sprintf("'%c'", v)
		}
		return fmt.Sprintf("0x%x", v)
	}
	return fmt.Sprintf("%v", v)
}
