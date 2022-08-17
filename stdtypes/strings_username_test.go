package stdtypes

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	. "github.com/rez-go/constraints"
)

func ExampleStringUsername() {

	var (
		usernameMinLength = StringMinLength(6)
		usernameMaxLength = StringMaxLength(32)
		// All these constraints in this example could be declared as a
		// single regular expression pattern, but we are trying to design
		// a mechanism which is more readable, constructed of smaller, clear
		// rules rather than putting all the rules into a complex pattern.
		//
		// We might want to find a way to declare this kind of constraint,
		// but we will limit how far we will go before we reinvent regular
		// expression.
		usernameAllowedCharacters = Func(
			`allowed characters are A to Z, 0 to 9 and underscore`,
			regexp.MustCompile(`^[A-Za-z0-9_]+$`).MatchString)
		// Name the variables based on the semantic rather than the
		// description of the constraint (what a constraint is for, rather than
		// what the constraint does). For the description of the constraint,
		// we put it into the constraint instance itself.
		//
		// The variable name and the description are good if they sound good if
		// we merge them:
		// "username [starts with] a letter".
		usernameStartsWith = Func(
			`starts with a letter`,
			func(v string) bool {
				if v != "" {
					r, _ := utf8.DecodeRuneInString(v)
					return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
				}
				return false
			})
		usernameEndsWith = Func(
			`ends with anything but underscore`,
			func(v string) bool {
				if v != "" {
					r, _ := utf8.DecodeLastRuneInString(v)
					return r != '_'
				}
				return false
			})
		usernameNoConsecutiveUnderscore = StringNoConsecutiveRune('_')
	)
	var usernameConstraints = Set(
		usernameMinLength,
		usernameMaxLength,
		usernameAllowedCharacters,
		usernameStartsWith,
		usernameEndsWith,
		usernameNoConsecutiveUnderscore,
	)
	fmt.Printf("%s\n", usernameConstraints.ConstraintDescription())
	// Output: min length 6, max length 32, allowed characters are A to Z, 0 to 9 and underscore, starts with a letter, ends with anything but underscore, no consecutive '_'
}
