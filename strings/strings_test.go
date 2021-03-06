package strings

import (
	"testing"

	"github.com/rez-go/constraints"
)

func TestEmpty(t *testing.T) {
	if Empty.ConstraintDescription() != "empty" {
		t.Errorf(`expecting "empty", got %q`,
			Empty.ConstraintDescription())
	}
	if !Empty.IsValid("") {
		t.Errorf(`!Empty.IsValid("")`)
	}
	if Empty.IsValid("empty") {
		t.Errorf(`Empty.IsValid("empty")`)
	}
	if err := ValidOrError("", Empty); err != nil {
		t.Errorf(`err := ValidOrError("", Empty); err != nil`)
	}
	err := ValidOrError("empty", Empty)
	if err == nil {
		t.Errorf(`err must not be nil`)
	}
	if err.Error() != "required to be empty" {
		t.Errorf(`expecting "required to be empty", got %q`, err.Error())
	}
	constraintError, ok := err.(constraints.Error)
	if !ok {
		t.Errorf(`err must be an implementation of constraints.Error`)
	}
	violatedConstraint := constraintError.ViolatedConstraint()
	if violatedConstraint != Empty {
		t.Errorf(`violatedConstraint must be Empty`)
	}
	violatedConstraint = constraints.ViolatedConstraintFromError(err)
	if violatedConstraint != Empty {
		t.Errorf(`violatedConstraint must be Empty`)
	}
}

func TestNotEmpty(t *testing.T) {
	if NotEmpty.ConstraintDescription() != "not empty" {
		t.Errorf(`expecting "not empty", got %q`,
			NotEmpty.ConstraintDescription())
	}
	if NotEmpty.IsValid("") {
		t.Errorf(`NotEmpty.IsValid("")`)
	}
	if !NotEmpty.IsValid("empty") {
		t.Errorf(`!NotEmpty.IsValid("empty")`)
	}
}

func TestNotWhitespace(t *testing.T) {
	if NotWhitespace.ConstraintDescription() != "not whitespace" {
		t.Errorf(`expecting "not whitespace", got %q`,
			NotWhitespace.ConstraintDescription())
	}
	if !NotWhitespace.IsValid("") {
		t.Errorf(`!NotWhitespace.IsValid("")`)
	}
	if NotWhitespace.IsValid("   ") {
		t.Errorf(`NotWhitespace.IsValid("   ")`)
	}
	if NotWhitespace.IsValid("		") {
		t.Errorf(`NotWhitespace.IsValid("		")`)
	}
	if !NotWhitespace.IsValid("whitespace") {
		t.Errorf(`!NotWhitespace.IsValid("whitespace")`)
	}
}

func TestLength10(t *testing.T) {
	len10 := NewLength(10)
	if !len10.IsValid("1234567890") {
		t.Errorf(`!len10.IsValid("1234567890")`)
	}
	if len10.IsValid("") {
		t.Errorf(`len10.IsValid("")`)
	}
	if len10.IsValid("123456789") {
		t.Errorf(`!len10.IsValid("123456789")`)
	}
	if len10.IsValid("12345678901") {
		t.Errorf(`!len10.IsValid("12345678901")`)
	}
}

func TestLength0(t *testing.T) {
	len0 := NewLength(0)
	if len0.ConstraintDescription() != "length 0" {
		t.Errorf(`expecting "length 0", got %q`,
			len0.ConstraintDescription())
	}
	if !len0.IsValid("") {
		t.Errorf(`len0.IsValid("")`)
	}
	if len0.IsValid("1234567890") {
		t.Errorf(`!len0.IsValid("1234567890")`)
	}
	if len0.IsValid("12345678901") {
		t.Errorf(`!len0.IsValid("12345678901")`)
	}
}

func TestMinLength0(t *testing.T) {
	min0 := NewMinLength(0)
	if !min0.IsValid("1234567890") {
		t.Errorf(`!min0.IsValid("1234567890")`)
	}
	if !min0.IsValid("") {
		t.Errorf(`!min0.IsValid("")`)
	}
	if !min0.IsValid("12345678901") {
		t.Errorf(`!min0.IsValid("12345678901")`)
	}
}

func TestMinLength10(t *testing.T) {
	min10 := NewMinLength(10)
	if min10.IsValid("") {
		t.Errorf(`min10.IsValid("")`)
	}
	if !min10.IsValid("1234567890") {
		t.Errorf(`!min10.IsValid("1234567890")`)
	}
	if min10.IsValid("123456789") {
		t.Errorf(`min10.IsValid("123456789")`)
	}
	if !min10.IsValid("12345678901") {
		t.Errorf(`!min10.IsValid("12345678901")`)
	}
}

func TestMaxLength0(t *testing.T) {
	max0 := NewMaxLength(0)
	if max0.ConstraintDescription() != "max length 0" {
		t.Errorf(`expecting "max length 0", got %q`,
			max0.ConstraintDescription())
	}
	if !max0.IsValid("") {
		t.Errorf(`!max0.IsValid("")`)
	}
	if max0.IsValid("1234567890") {
		t.Errorf(`max0.IsValid("1234567890")`)
	}
	if max0.IsValid("12345678901") {
		t.Errorf(`max0.IsValid("12345678901")`)
	}
}

func TestMaxLength10(t *testing.T) {
	max10 := NewMaxLength(10)
	if !max10.IsValid("") {
		t.Errorf(`!max10.IsValid("")`)
	}
	if !max10.IsValid("1234567890") {
		t.Errorf(`!max10.IsValid("1234567890")`)
	}
	if !max10.IsValid("123456789") {
		t.Errorf(`!max10.IsValid("123456789")`)
	}
	if max10.IsValid("12345678901") {
		t.Errorf(`max10.IsValid("12345678901")`)
	}
}

func TestSetEmpty(t *testing.T) {
	cs := Set{}
	if cs.ConstraintDescription() != "" {
		t.Errorf(`expecting "", got %q`,
			cs.ConstraintDescription())
	}
	if !cs.IsValid("") {
		t.Errorf(`!cs.IsValid("")`)
	}
	if !cs.IsValid("hello") {
		t.Errorf(`!cs.IsValid("hello")`)
	}
	cl := cs.ConstraintList()
	if len(cl) != 0 {
		t.Errorf(`expecting empty constraint list`)
	}
}

func TestSetBasic(t *testing.T) {
	cs := Set{NotEmpty, NotWhitespace}
	if cs.ConstraintDescription() != "not empty, not whitespace" {
		t.Errorf(`expecting "not empty, not whitespace", got %q`,
			cs.ConstraintDescription())
	}
	if cs.IsValid("") {
		t.Errorf(`cs.IsValid("")`)
	}
	if cs.IsValid("  ") {
		t.Errorf(`cs.IsValid("  ")`)
	}
	if !cs.IsValid("hello") {
		t.Errorf(`!cs.IsValid("hello")`)
	}
	if violatedConstraint := cs.Validate(""); violatedConstraint != NotEmpty {
		t.Errorf(`expecting %#v, got %#v`, NotEmpty, violatedConstraint)
	}
	if violatedConstraint := cs.Validate("  "); violatedConstraint != NotWhitespace {
		t.Errorf(`expecting %#v, got %#v`, NotWhitespace, violatedConstraint)
	}
	if violatedConstraints := cs.ValidateAll(""); len(violatedConstraints) != 1 || violatedConstraints[0] != NotEmpty {
		t.Errorf("ValidateAll")
	}
	err := ValidOrError("", cs)
	if err.Error() != "required to be not empty" {
		t.Errorf(`expecting "required to be not empty", got %q`, err.Error())
	}
	cl := cs.ConstraintList()
	for i, cA := range cl {
		if cA != cs[i] {
			t.Errorf(`cs.ConstraintList()`)
		}
	}
}

func TestSetMinhNotWhitespace(t *testing.T) {
	cs := Set{NewMinLength(5), NotWhitespace}
	if cs.ConstraintDescription() != "min length 5, not whitespace" {
		t.Errorf(`expecting "min length 5, not whitespace", got %q`,
			cs.ConstraintDescription())
	}
	if cs.IsValid("") {
		t.Errorf(`cs.IsValid("")`)
	}
	if cs.IsValid("hell") {
		t.Errorf(`cs.IsValid("hell")`)
	}
	if !cs.IsValid("hello") {
		t.Errorf(`!cs.IsValid("hello")`)
	}
	err := ValidOrError("12345", cs)
	if err != nil {
		t.Errorf(`expecting <nil>, got %v`, err)
	}
	err = ValidOrError("  ", cs)
	if err.Error() != "required to be min length 5, not whitespace" {
		t.Errorf(`expecting "required to be min length 5, not whitespace", got %q`, err.Error())
	}
}

func TestConst(t *testing.T) {
	emptyConst := NewConst("")
	if emptyConst.ConstraintDescription() != "const \"\"" {
		t.Errorf(`expecting "const \"\"", got %q`, emptyConst.ConstraintDescription())
	}
	if !emptyConst.IsValid("") {
		t.Errorf(`!emptyConst.IsValid("")`)
	}
	if emptyConst.IsValid("empty") {
		t.Errorf(`emptyConst.IsValid("empty")`)
	}
	helloConst := NewConst("hello")
	if helloConst.IsValid("") {
		t.Errorf(`helloConst.IsValid("")`)
	}
	if helloConst.IsValid("empty") {
		t.Errorf(`helloConst.IsValid("empty")`)
	}
	if !helloConst.IsValid("hello") {
		t.Errorf(`!helloConst.IsValid("hello")`)
	}
}

func TestOneOf(t *testing.T) {
	emptyOneOf := NewOneOf()
	if emptyOneOf.ConstraintDescription() != "one of []" {
		t.Errorf(`expecting "one of []", got %q`,
			emptyOneOf.ConstraintDescription())
	}
	if emptyOneOf.IsValid("") {
		t.Errorf(`emptyOneOf.IsValid("")`)
	}
	if emptyOneOf.IsValid("nothing") {
		t.Errorf(`emptyOneOf.IsValid("nothing")`)
	}
	oneOneOf := NewOneOf("one")
	if oneOneOf.ConstraintDescription() != "one of [one]" {
		t.Errorf(`expecting "one of [one]", got %q`,
			oneOneOf.ConstraintDescription())
	}
	if oneOneOf.IsValid("") {
		t.Errorf(`oneOneOf.IsValid("")`)
	}
	if oneOneOf.IsValid("zero") {
		t.Errorf(`oneOneOf.IsValid("zero")`)
	}
	if !oneOneOf.IsValid("one") {
		t.Errorf(`!oneOneOf.IsValid("one")`)
	}
	oneOfOneTwoThree := NewOneOf("one", "two", "three")
	if oneOfOneTwoThree.ConstraintDescription() != "one of [one, two, three]" {
		t.Errorf(`expecting "one of [one, two, three]", got %q`,
			oneOfOneTwoThree.ConstraintDescription())
	}
}

func TestPrefix(t *testing.T) {
	emptyPrefix := Prefix("")
	if emptyPrefix.ConstraintDescription() != "prefix \"\"" {
		t.Errorf(`expecting "prefix \"\"", got %q`, emptyPrefix.ConstraintDescription())
	}
	if !emptyPrefix.IsValid("") {
		t.Errorf(`!emptyPrefix.IsValid("")`)
	}
	if !emptyPrefix.IsValid("no prefix") {
		t.Errorf(`!emptyPrefix.IsValid("no prefix")`)
	}
}

func TestNoConsecutiveRune(t *testing.T) {
	noConsecutiveUnderscore := NoConsecutiveRune('_')
	if noConsecutiveUnderscore.ConstraintDescription() != "no consecutive '_'" {
		t.Errorf(`expecting "no consecutive '_'", got %q`,
			noConsecutiveUnderscore.ConstraintDescription())
	}
	if !noConsecutiveUnderscore.IsValid("") {
		t.Errorf(`!noConsecutiveUnderscore.IsValid("")`)
	}
	if !noConsecutiveUnderscore.IsValid("_") {
		t.Errorf(`!noConsecutiveUnderscore.IsValid("_")`)
	}
	if noConsecutiveUnderscore.IsValid("__") {
		t.Errorf(`noConsecutiveUnderscore.IsValid("__")`)
	}
	if !noConsecutiveUnderscore.IsValid("_._") {
		t.Errorf(`!noConsecutiveUnderscore.IsValid("_._")`)
	}
}
