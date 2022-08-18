package constraints

import (
	"fmt"
	"strings"
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

// A Set is a constraint which consisted of a set of constraints.
type ConstraintSet[
	ValueT any,
	ConstraintT Constraint[ValueT],
] interface {
	Constraint[ValueT]
	Validate(v ValueT) Constraint[ValueT]
	ValidateAll(v ValueT) (violated []Constraint[ValueT])
	ConstraintList() []ConstraintT
}

// An EqOp specifies equality operator.
type EqOp int

// Supported equality operators.
const (
	EqOpEqual EqOp = iota
	EqOpNotEqual
	EqOpLess
	EqOpLessOrEqual
	EqOpGreater
	EqOpGreaterOrEqual
)

func (op EqOp) String() string {
	switch op {
	case EqOpEqual:
		return "equal"
	case EqOpNotEqual:
		return "not equal"
	case EqOpLess:
		return "less"
	case EqOpLessOrEqual:
		return "less or equal"
	case EqOpGreater:
		return "greater"
	case EqOpGreaterOrEqual:
		return "greater or equal"
	}
	return ""
}

// Symbol returns representative symbol of the operator.
func (op EqOp) Symbol() string {
	switch op {
	case EqOpEqual:
		return "="
	case EqOpNotEqual:
		return "≠"
	case EqOpLess:
		return "<"
	case EqOpLessOrEqual:
		return "≤"
	case EqOpGreater:
		return ">"
	case EqOpGreaterOrEqual:
		return "≥"
	}
	return "?"
}

// StringFormat returns a string which could be used to in *printf functions.
// The string format expects a value to be passed.
func (op EqOp) StringFormat() string {
	switch op {
	case EqOpEqual:
		return "equals %v"
	case EqOpNotEqual:
		return "not equal to %v"
	case EqOpLess:
		return "less than %v"
	case EqOpLessOrEqual:
		return "less than or equal to %v"
	case EqOpGreater:
		return "greater than %v"
	case EqOpGreaterOrEqual:
		return "greater than or equal to %v"
	}
	return "(%v)"
}

//----

func Set[ValueT any](
	constraints ...Constraint[ValueT],
) ConstraintSet[ValueT, Constraint[ValueT]] {
	// copy?
	return constraintSet[ValueT](constraints)
}

var (
	_ Constraint[string]                        = constraintSet[string]{}
	_ ConstraintSet[string, Constraint[string]] = constraintSet[string]{}
)

// constraintSet defines a set of constraints. A value is considered valid
// if every constraint considers it as valid.
//
// TODO: Set with 'or' (instead of 'and')
type constraintSet[ValueT any] []Constraint[ValueT]

// ConstraintDescription conforms Constraint interface.
func (cs constraintSet[ValueT]) ConstraintDescription() string {
	if cs != nil {
		descs := make([]string, 0, len(cs))
		for _, ci := range cs {
			descs = append(descs, ci.ConstraintDescription())
		}
		return strings.Join(descs, ", ")
	}
	return ""
}

// ConstraintList conforms Set interface.
func (cs constraintSet[ValueT]) ConstraintList() []Constraint[ValueT] {
	if cs != nil {
		copyConstraints := make([]Constraint[ValueT], 0, len(cs))
		for _, c := range cs {
			copyConstraints = append(copyConstraints, c)
		}
		return copyConstraints
	}
	return nil
}

// IsValid conforms Constraint interface.
func (cs constraintSet[ValueT]) IsValid(v ValueT) bool {
	return cs.Validate(v) == nil
}

// Validate checks the value against the constraints. If the value is
// valid, this method will return nil. Otherwise, it'll return the first
// violated constraint.
func (cs constraintSet[ValueT]) Validate(v ValueT) Constraint[ValueT] {
	for _, ci := range cs {
		if !ci.IsValid(v) {
			return ci
		}
	}
	return nil
}

// ValidateAll validates the value against all constraints. It returns
// a list of violate constraints if any.
func (cs constraintSet[ValueT]) ValidateAll(v ValueT) (violated []Constraint[ValueT]) {
	for _, ci := range cs {
		if !ci.IsValid(v) {
			violated = append(violated, ci)
		}
	}
	return
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
				return ViolationError[ValueT](constraintSet[ValueT](violatedConstraints))
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

// OneOf creates a Constraint which will declare a value as valid
// if it matches one of the options.
//
// API status: experimental
func OneOf[ValueT comparable](options ...ValueT) *oneOfConstraint[ValueT] {
	copies := make([]ValueT, len(options))
	copy(copies, options)
	return &oneOfConstraint[ValueT]{negate: false, options: options}
}

func NoneOf[ValueT comparable](options ...ValueT) *oneOfConstraint[ValueT] {
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

// ValidatorFunc is an adapter to allow use of ordinary functions
// as a validator in constraints.
type ValidatorFunc[ValueT any] func(i ValueT) bool

// Func creates a constraint from a validator function.
//
// A good example is for defining a constraint to ensure that provided
// string is a valid UTF8-encoded string:
//
//	import "unicode/utf8"
//
//	var mustUTF8 = Func("valid UTF-8", utf8.ValidString)
//
// Or regular expression string matcher:
//
//	var usernamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
//	var usernameConstraint = Func("username", usernamePattern.MatchString)
func Func[
	ValueT any,
](desc string, fn func(v ValueT) bool) Constraint[ValueT] {
	return &constraintFunc[ValueT]{
		negate: false,
		desc:   desc,
		fn:     fn,
	}
}

var (
	_ Constraint[int64] = &constraintFunc[int64]{}
	_ Constraint[int64] = constraintFunc[int64]{}
)

type constraintFunc[ValueT any] struct {
	negate bool
	desc   string
	fn     ValidatorFunc[ValueT]
}

func (c constraintFunc[ValueT]) ConstraintDescription() string {
	return c.desc
}

func (c constraintFunc[ValueT]) IsValid(v ValueT) bool {
	result := c.fn(v)
	if c.negate {
		return !result
	}
	return result
}

func Negate[ValueT any](c Constraint[ValueT], descOverride string) Constraint[ValueT] {
	desc := descOverride
	if desc == "" {
		desc = "not " + c.ConstraintDescription()
	}
	return &constraintFunc[ValueT]{negate: true, desc: desc, fn: c.IsValid}
}

func valueLiteralString(v any) string {
	switch v.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	}
	return fmt.Sprintf("%v", v)
}
