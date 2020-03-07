package strings

import (
	"fmt"
	strlib "strings"

	"github.com/rez-go/constraints"
	"github.com/rez-go/constraints/ints"
)

// Constraint is an abstract type for string-related constraints.
type Constraint interface {
	constraints.Constraint
	IsValid(v string) bool
	ValidOrError(v string) constraints.Error
}

var (
	_ Constraint      = Set{}
	_ constraints.Set = Set{}
)

// Set defines a set of constraints. A value is considered valid
// if every constraint considers it as valid.
//
//TODO: Set with 'or' (instead of 'and')
type Set []Constraint

// ConstraintDescription conforms Constraint interface.
func (cs Set) ConstraintDescription() string {
	if cs != nil {
		descs := make([]string, 0, len(cs))
		for _, ci := range cs {
			descs = append(descs, ci.ConstraintDescription())
		}
		return strlib.Join(descs, ", ")
	}
	return "<empty constraint set>"
}

// ConstraintList conforms constraints.Set interface.
func (cs Set) ConstraintList() []constraints.Constraint {
	if cs != nil {
		copyConstraints := make([]constraints.Constraint, 0, len(cs))
		for _, c := range cs {
			copyConstraints = append(copyConstraints, c)
		}
		return copyConstraints
	}
	return nil
}

// IsValid conforms Constraint interface.
func (cs Set) IsValid(v string) bool {
	return cs != nil && cs.Validate(v) == nil
}

// ValidOrError conforms Constraint interface. If the value violated any
// or all of the constraints, this method will return a new instance
// of Set which contains only the violated constraints.
func (cs Set) ValidOrError(v string) constraints.Error {
	if cs != nil {
		violatedConstraints := cs.ValidateAll(v)
		if len(violatedConstraints) == 0 {
			return nil
		}
		return constraints.ViolationError(Set(violatedConstraints))
	}
	return constraints.ViolationError(Set{})
}

// Validate checks the value against the constraints. If the value is
// valid, this method will return nil. Otherwise, it'll return the first
// violated constraint.
func (cs Set) Validate(v string) Constraint {
	for _, ci := range cs {
		if !ci.IsValid(v) {
			return ci
		}
	}
	return nil
}

// ValidateAll validates the value against all constraints. It returns
// a list of violate constraints if any.
func (cs Set) ValidateAll(v string) (violated []Constraint) {
	for _, ci := range cs {
		if !ci.IsValid(v) {
			violated = append(violated, ci)
		}
	}
	return
}

var (
	// Empty is a constraint which value is considered valid if it's
	// an empty string.
	Empty Constraint = Func(
		"empty",
		func(v string) bool {
			return v == ""
		},
	)

	// NotEmpty is a constraint which value is considered valid if it's
	// not an empty string.
	NotEmpty Constraint = Func(
		"not empty",
		func(v string) bool {
			return v != ""
		},
	)

	// NotWhitespace is a constraint which will consider value is valid
	// if it contains not just whitespace.
	//
	// Note that empty string is considered as a non-whitespace string.
	NotWhitespace Constraint = Func(
		"not whitespace",
		func(v string) bool {
			return v == "" || strlib.TrimSpace(v) != ""
		},
	)
)

// Length returns a Constraint which a string instance is
// valid against the constraint if its length is exactly as specified.
func Length(specifiedLength int64) Constraint {
	if specifiedLength < 0 {
		panic("specifiedLength must be zero or a positive integer")
	}
	return &lengthOp{ints.Equals(specifiedLength)}
}

// MaxLength returns a Constraint which a string instance is
// valid against the constraint if its length is less than, or equal to,
// the reference value.
func MaxLength(maxLength int64) Constraint {
	if maxLength < 0 {
		panic("maxLength must be zero or a positive integer")
	}
	return &lengthOp{ints.Max(maxLength)}
}

// MinLength returns a Constraint which a string instance is
// valid against the constraint if its length is greater than, or equal to,
// the reference value.
func MinLength(minLength int64) Constraint {
	if minLength < 0 {
		panic("minLength must be zero or a positive integer")
	}
	return &lengthOp{ints.Min(minLength)}
}

// Const creates a Constraint which will declare a value as valid
// if it matches refValue.
func Const(refValue string) Constraint {
	return &constraintFunc{
		fmt.Sprintf("const %q", refValue),
		func(v string) bool {
			return v == refValue
		}}
}

// In creates a Constraint which will declare a value as valid
// if it matches one of the choices.
func In(choices ...string) Constraint {
	copyChoices := make([]string, len(choices))
	copy(copyChoices, choices)
	return &constraintFunc{
		fmt.Sprintf("in %v", copyChoices),
		func(v string) bool {
			for _, s := range copyChoices {
				if v == s {
					return true
				}
			}
			return false
		}}
}

var (
	_ Constraint = &lengthOp{}
)

type lengthOp struct {
	lengthCheck ints.Constraint
}

func (c *lengthOp) IsValid(v string) bool {
	return c.lengthCheck.IsValid(int64(len(v)))
}

// ValidOrError conforms Constraint interface.
func (c *lengthOp) ValidOrError(v string) constraints.Error {
	if c != nil && c.IsValid(v) {
		return nil
	}
	return constraints.ViolationError(c)
}

func (c *lengthOp) ConstraintDescription() string {
	return "length " + c.lengthCheck.ConstraintDescription()
}

// ValidatorFunc is an adapter to allow use of ordinary functions
// as a validator in constraints.
type ValidatorFunc func(s string) bool

// Func creates a constraint from a validator function.
//
// A good example is for defining UTF8 constraint:
//
// 	import "unicode/utf8"
//
// 	var mustUTF8 = Func("valid UTF-8", utf8.ValidString)
//
// Or regular expression string matcher:
//
// 	var myPattern = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
//  var myConstraint = Func("match my pattern", myPattern.MatchString)
//
func Func(desc string, fn func(v string) bool) Constraint {
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
	return "<unknown string constraint>"
}

func (c *constraintFunc) IsValid(v string) bool {
	if c != nil {
		return c.fn(v)
	}
	return false
}

func (c *constraintFunc) ValidOrError(v string) constraints.Error {
	if c != nil && c.IsValid(v) {
		return nil
	}
	return constraints.ViolationError(c)
}
