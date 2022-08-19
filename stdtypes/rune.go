package stdtypes

import (
	"fmt"
	"strconv"

	"github.com/rez-go/constraints"
)

type RuneConstraint = constraints.Constraint[rune]

var (
	RuneMatch = constraints.Match[rune]
	RuneOneOf = constraints.OneOf[rune]

	PrintableRune = constraints.Func(
		"printable rune",
		func(v rune) bool {
			return strconv.IsPrint(v)
		})
)

func RuneOneOfByString(allowedRunes string) RuneConstraint {
	return constraints.Func(
		fmt.Sprintf("rune from %q", allowedRunes), //TODO: splits
		func(v rune) bool {
			for _, cr := range allowedRunes {
				if v == cr {
					return true
				}
			}
			return false
		})
}

// RuneRange creates a RuneConstraint that declares a rune as valid
// if said rune is within the specified range.
//
// e.g., RuneRange('a', 'z') is valid for all lowercased latin alphabet runes.
var RuneRange = constraints.Range[rune]
