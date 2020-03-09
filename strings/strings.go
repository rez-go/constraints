// Package strings contains constraints for strings.
//
// TODO: case-folding / caseless
package strings

import (
	"fmt"
	strlib "strings"

	"github.com/rez-go/constraints"
)

// Constraint is an abstract type for string-related constraints.
type Constraint interface {
	constraints.Constraint
	IsValid(v string) bool
}

// ValidOrError tests the value v against constraint c. If the value is
// valid, it will return nil, otherwise, it'll return an error.
//
// If the constraint is a Set and the value violated any
// or all of the constraints, this method will return a Error which
// constraint contained is a new instance of Set contains only the
// violated constraints.
func ValidOrError(v string, c Constraint) constraints.Error {
	if c != nil {
		if cs, ok := c.(interface{ ValidateAll(string) []Constraint }); ok && cs != nil {
			violatedConstraints := cs.ValidateAll(v)
			if len(violatedConstraints) > 0 {
				return constraints.ViolationError(Set(violatedConstraints))
			}
			return nil
		}
		if c.IsValid(v) {
			return nil
		}
	}
	return constraints.ViolationError(c)
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
		func(v string) bool { return v == "" })

	// NotEmpty is a constraint which value is considered valid if it's
	// not an empty string.
	NotEmpty Constraint = Func(
		"not empty",
		func(v string) bool { return v != "" })

	// NotWhitespace is a constraint which will consider value is valid
	// if it contains not just whitespace.
	//
	// Note that an empty string is considered as a non-whitespace string.
	NotWhitespace Constraint = Func(
		"not whitespace",
		func(v string) bool {
			return v == "" || strlib.TrimSpace(v) != ""
		},
	)
)

// NewConst creates a Constraint which will declare an instance as valid
// if its value matches refValue.
func NewConst(refValue string) Constraint { return &Const{refValue} }

// Const defines constant value Constraint.
type Const struct {
	refValue string
}

var (
	_ Constraint = &Const{}
)

// ConstraintDescription conforms constraints.Constraint interface.
func (c *Const) ConstraintDescription() string {
	if c != nil {
		return fmt.Sprintf("const %q", c.refValue)
	}
	return "const constraint"
}

// IsValid conforms Constraint interface.
func (c *Const) IsValid(v string) bool {
	return c != nil && v == c.refValue
}

// NewLength returns a Constraint which a string instance will be
// declared as valid if its length is exactly as specified.
//
// API status: experimental
func NewLength(specifiedLength int64) Constraint {
	if specifiedLength < 0 {
		panic("specifiedLength must be zero or a positive integer")
	}
	return &Length{specifiedLength}
}

// Length defines exact length Constraint.
type Length struct {
	refValue int64
}

var (
	_ Constraint = &Length{}
)

// ConstraintDescription conforms constraints.Constraint interface.
func (c *Length) ConstraintDescription() string {
	if c != nil {
		return fmt.Sprintf("length %d", c.refValue)
	}
	return "length constraint"
}

// IsValid conforms Constraint interface.
func (c *Length) IsValid(v string) bool {
	return c != nil && int64(len(v)) == c.refValue
}

// NewMaxLength returns a Constraint which a string instance will be
// declared as valid if its length is at most equal to maxLength.
//
// API status: experimental
func NewMaxLength(maxLength int64) Constraint {
	if maxLength < 0 {
		panic("maxLength must be zero or a positive integer")
	}
	return &MaxLength{maxLength}
}

// MaxLength defines minimum length Constraint.
type MaxLength struct {
	refValue int64
}

var (
	_ Constraint = &MaxLength{}
)

// ConstraintDescription conforms constraints.Constraint interface.
func (c *MaxLength) ConstraintDescription() string {
	if c != nil {
		return fmt.Sprintf("max length %d", c.refValue)
	}
	return "max length constraint"
}

// IsValid conforms Constraint interface.
func (c *MaxLength) IsValid(v string) bool {
	return c != nil && int64(len(v)) <= c.refValue
}

// NewMinLength returns a Constraint which a string instance will be
// declared as valid if its length is at most least to maxLength.
//
// API status: experimental
func NewMinLength(minLength int64) Constraint {
	if minLength < 0 {
		panic("minLength must be zero or a positive integer")
	}
	return &MinLength{minLength}
}

// MinLength defines minimum length Constraint.
type MinLength struct {
	refValue int64
}

var (
	_ Constraint = &MinLength{}
)

// ConstraintDescription conforms constraints.Constraint interface.
func (c *MinLength) ConstraintDescription() string {
	if c != nil {
		return fmt.Sprintf("min length %d", c.refValue)
	}
	return "min length constraint"
}

// IsValid conforms Constraint interface.
func (c *MinLength) IsValid(v string) bool {
	return c != nil && int64(len(v)) >= c.refValue
}

// NoConsecutiveRune creates a Constraint which will declare a string
// as valid if it doesn't containt any conscutive of r.
func NoConsecutiveRune(r rune) Constraint {
	return &funcConstraint{
		fmt.Sprintf("no consecutive '%c'", r),
		func(v string) bool {
			lastRuneMatched := false
			for _, ir := range v {
				if ir == r {
					if lastRuneMatched {
						return false
					}
					lastRuneMatched = true
				} else {
					lastRuneMatched = false
				}
			}
			return true
		}}
}

// Prefix creates a Constraint which an instance will be declared
// as valid if its value is prefixed with the specified prefix.
func Prefix(prefix string) Constraint {
	return &targetOperandFuncConstraint{
		fmt.Sprintf("prefix %q", prefix),
		prefix,
		strlib.HasPrefix}
}

// Suffix creates a Constraint which an instance will be declared
// as valid if its value is suffixed with the specified suffix.
func Suffix(suffix string) Constraint {
	return &targetOperandFuncConstraint{
		fmt.Sprintf("suffix %q", suffix),
		suffix,
		strlib.HasSuffix}
}

// OneOf creates a Constraint which will declare a value as valid
// if it matches one of the choices.
//
// API status: experimental
func OneOf(choices ...string) Constraint {
	copyChoices := make([]string, len(choices))
	copy(copyChoices, choices)
	return &funcConstraint{
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

// ValidatorFunc is an adapter to allow use of ordinary functions
// as a validator in constraints.
type ValidatorFunc func(s string) bool

// Func creates a constraint from a validator function.
//
// A good example is for defining a constraint to ensure that provided
// string is a valid UTF8-encoded string:
//
// 	import "unicode/utf8"
//
// 	var mustUTF8 = Func("valid UTF-8", utf8.ValidString)
//
// Or regular expression string matcher:
//
//  var usernamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
//  var usernameConstraint = Func("username", usernamePattern.MatchString)
//
func Func(desc string, fn func(v string) bool) Constraint {
	return &funcConstraint{desc, fn}
}

var (
	_ Constraint = &funcConstraint{}
)

type funcConstraint struct {
	desc string
	fn   ValidatorFunc
}

func (c *funcConstraint) ConstraintDescription() string {
	if c != nil {
		return c.desc
	}
	return "string constraint"
}

func (c *funcConstraint) IsValid(v string) bool {
	if c != nil {
		return c.fn(v)
	}
	return false
}

var (
	_ Constraint = &targetOperandFuncConstraint{}
)

type targetOperandFuncConstraint struct {
	desc    string
	operand string
	fn      func(target, operand string) bool
}

func (c *targetOperandFuncConstraint) ConstraintDescription() string {
	if c != nil {
		return c.desc
	}
	return "string constraint"
}

func (c *targetOperandFuncConstraint) IsValid(v string) bool {
	return c != nil && c.fn != nil && c.fn(v, c.operand)
}
