// Package strings contains constraints for strings.
//
// TODO: case-folding / caseless
package stdtypes

import (
	"fmt"
	strlib "strings"

	"github.com/rez-go/constraints"
)

var StringSet = constraints.Set[string]

// StringConstraint is an abstract type for string-related constraints.
type StringConstraint = constraints.Constraint[string]

// Built-in non-parametric constraints.
var (
	StringLength    = Length[string]
	StringMinLength = MinLength[string]
	StringMaxLength = MaxLength[string]

	// EmptyString is a constraint which value is considered valid if it's
	// an empty string.
	EmptyString StringConstraint = constraints.Func(
		"empty",
		func(v string) bool { return v == "" })

	// NonEmptyString is a constraint which value is considered valid if it's
	// not an empty string.
	NonEmptyString StringConstraint = constraints.Func(
		"non-empty",
		func(v string) bool { return v != "" })

	// NonBlankString is a constraint which will consider value is valid
	// if it contains not just whitespace.
	//
	// Note that an empty string is considered as a non-whitespace string.
	NonBlankString StringConstraint = constraints.Func(
		"non-blank",
		func(v string) bool {
			return v == "" || strlib.TrimSpace(v) != ""
		},
	)
)

// StringNoConsecutiveRune creates a Constraint which will declare a string
// as valid if it doesn't containt any conscutive of r.
func StringNoConsecutiveRune(r rune) StringConstraint {
	return constraints.Func(
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
		})
}

// StringPrefix creates a Constraint which an instance will be declared
// as valid if its value is prefixed with the specified prefix.
func StringPrefix(prefix string) StringConstraint {
	return &targetOperandFuncConstraint[string]{
		fmt.Sprintf("prefix %q", prefix),
		prefix,
		strlib.HasPrefix}
}

// StringSuffix creates a Constraint which an instance will be declared
// as valid if its value is suffixed with the specified suffix.
func StringSuffix(suffix string) StringConstraint {
	return &targetOperandFuncConstraint[string]{
		fmt.Sprintf("suffix %q", suffix),
		suffix,
		strlib.HasSuffix}
}
