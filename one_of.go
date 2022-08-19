package constraints

import (
	"fmt"
	"strings"
)

// OneOf creates a Constraint which will declare a value as valid
// if it matches one of the options.
//
// API status: experimental
func OneOf[ValueT comparable](options ...ValueT) Constraint[ValueT] {
	copies := make([]ValueT, len(options))
	copy(copies, options)
	return &oneOfConstraint[ValueT]{negate: false, options: options}
}

func NoneOf[ValueT comparable](options ...ValueT) Constraint[ValueT] {
	copies := make([]ValueT, len(options))
	copy(copies, options)
	return &oneOfConstraint[ValueT]{negate: true, options: options}
}

// oneOfConstraint defines choice-based Constraint.
type oneOfConstraint[ValueT comparable] struct {
	negate  bool
	options []ValueT
}

var (
	_ Constraint[string] = oneOfConstraint[string]{}
)

// ConstraintDescription conforms constraints.Constraint interface.
func (c oneOfConstraint[ValueT]) ConstraintDescription() string {
	var opts []string
	for _, o := range c.options {
		opts = append(opts, fmt.Sprintf("%v", o))
	}
	opt := "[" + strings.Join(opts, ", ") + "]"
	if c.negate {
		return fmt.Sprintf("none of %v", opt)
	}
	return fmt.Sprintf("one of %v", opt)
}

// IsValid conforms Constraint interface.
func (c oneOfConstraint[ValueT]) IsValid(v ValueT) bool {
	for _, s := range c.options {
		if s == v {
			return !c.negate
		}
	}
	return c.negate
}
