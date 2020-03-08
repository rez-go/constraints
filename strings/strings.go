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

// Const creates a Constraint which will declare an instance as valid
// if its value matches refValue.
func Const(refValue string) Constraint {
	return &funcConstraint{
		fmt.Sprintf("const %q", refValue),
		func(v string) bool {
			return v == refValue
		}}
}

// Length returns a Constraint which an instance is
// valid against the constraint if its length is exactly as specified.
func Length(specifiedLength int64) Constraint {
	if specifiedLength < 0 {
		panic("specifiedLength must be zero or a positive integer")
	}
	return &lengthConstraint{
		fmt.Sprintf("length %d", specifiedLength),
		ints.Const(specifiedLength)}
}

// MaxLength returns a Constraint which will declare an instance as valid
// if its length is less than or equal to the maxLength.
func MaxLength(maxLength int64) Constraint {
	if maxLength < 0 {
		panic("maxLength must be zero or a positive integer")
	}
	return &lengthConstraint{
		fmt.Sprintf("max length %d", maxLength),
		ints.Max(maxLength)}
}

// MinLength returns a Constraint which will declare an instance as valid
// if its length is more than or equal to the minLength.
func MinLength(minLength int64) Constraint {
	if minLength < 0 {
		panic("minLength must be zero or a positive integer")
	}
	return &lengthConstraint{
		fmt.Sprintf("min length %d", minLength),
		ints.Min(minLength)}
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

func (c *funcConstraint) ValidOrError(v string) constraints.Error {
	if c != nil && c.IsValid(v) {
		return nil
	}
	return constraints.ViolationError(c)
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

func (c *targetOperandFuncConstraint) ValidOrError(v string) constraints.Error {
	if c != nil && c.IsValid(v) {
		return nil
	}
	return constraints.ViolationError(c)
}

var (
	_ Constraint         = &lengthConstraint{}
	_ constraints.Length = &lengthConstraint{}
)

type lengthConstraint struct {
	desc       string
	constraint ints.Constraint
}

func (w *lengthConstraint) ConstraintDescription() string {
	if w != nil {
		if w.desc != "" {
			return w.desc
		}
		if w.constraint != nil {
			return "length " + w.constraint.ConstraintDescription()
		}
	}
	return "string length constraint"
}

func (w *lengthConstraint) IsValid(v string) bool {
	return w != nil && w.constraint != nil && w.constraint.IsValid(int64(len(v)))
}

func (w *lengthConstraint) ValidOrError(v string) constraints.Error {
	if w == nil || !w.IsValid(v) {
		return constraints.ViolationError(w)
	}
	return nil
}
