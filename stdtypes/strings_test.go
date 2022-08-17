package stdtypes

import (
	"testing"

	"github.com/rez-go/constraints"
)

func TestEmpty(t *testing.T) {
	assertEq(t, "empty", EmptyString.ConstraintDescription())
	assertEq(t, true, EmptyString.IsValid(""))
	assertEq(t, false, EmptyString.IsValid("empty"))
	assertEq(t, nil, constraints.ValidOrError("", EmptyString))
	var err error = constraints.ValidOrError("empty", EmptyString)
	assertNeq(t, nil, err)
	assertEq(t, "required to be empty", err.Error())
	constraintError, ok := err.(constraints.Error[string])
	if !ok {
		t.Errorf(`err must be an implementation of constraints.Error`)
	}
	violatedConstraint := constraintError.ViolatedConstraint()
	if violatedConstraint != EmptyString {
		t.Errorf(`violatedConstraint must be Empty`)
	}
	violatedConstraint = constraints.ViolatedConstraintFromError[string](err)
	if violatedConstraint != EmptyString {
		t.Errorf(`violatedConstraint must be Empty`)
	}
}

func TestNotEmpty(t *testing.T) {
	if NonEmptyString.ConstraintDescription() != "non-empty" {
		t.Errorf(`expecting "non-empty", got %q`,
			NonEmptyString.ConstraintDescription())
	}
	if NonEmptyString.IsValid("") {
		t.Errorf(`NotEmpty.IsValid("")`)
	}
	if !NonEmptyString.IsValid("empty") {
		t.Errorf(`!NotEmpty.IsValid("empty")`)
	}
}

func TestNotWhitespace(t *testing.T) {
	if NonBlankString.ConstraintDescription() != "non-blank" {
		t.Errorf(`expecting "non-blank", got %q`,
			NonBlankString.ConstraintDescription())
	}
	if !NonBlankString.IsValid("") {
		t.Errorf(`!NotWhitespace.IsValid("")`)
	}
	if NonBlankString.IsValid("   ") {
		t.Errorf(`NotWhitespace.IsValid("   ")`)
	}
	if NonBlankString.IsValid("		") {
		t.Errorf(`NotWhitespace.IsValid("		")`)
	}
	if !NonBlankString.IsValid("whitespace") {
		t.Errorf(`!NotWhitespace.IsValid("whitespace")`)
	}
}

func TestLength10(t *testing.T) {
	len10 := StringLength(10)
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
	len0 := StringLength(0)
	assertEq(t, "length 0", len0.ConstraintDescription())
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
	min0 := StringMinLength(0)
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
	min10 := StringMinLength(10)
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
	max0 := StringMaxLength(0)
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
	max10 := StringMaxLength(10)
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
	cs := StringSet()
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
	cs := StringSet(NonEmptyString, NonBlankString)
	assertEq(t, "non-empty, non-blank", cs.ConstraintDescription())
	if cs.IsValid("") {
		t.Errorf(`cs.IsValid("")`)
	}
	if cs.IsValid("  ") {
		t.Errorf(`cs.IsValid("  ")`)
	}
	if !cs.IsValid("hello") {
		t.Errorf(`!cs.IsValid("hello")`)
	}
	if violatedConstraint := cs.Validate(""); violatedConstraint != NonEmptyString {
		t.Errorf(`expecting %#v, got %#v`, NonEmptyString, violatedConstraint)
	}
	if violatedConstraint := cs.Validate("  "); violatedConstraint != NonBlankString {
		t.Errorf(`expecting %#v, got %#v`, NonBlankString, violatedConstraint)
	}
	if violatedConstraints := cs.ValidateAll(""); len(violatedConstraints) != 1 || violatedConstraints[0] != NonEmptyString {
		t.Errorf("ValidateAll")
	}
	var err error = constraints.ValidOrError[string]("", cs)
	assertEq(t, "required to be non-empty", err.Error())
	cl := cs.ConstraintList()
	for i, cA := range cl {
		if cA != cl[i] {
			t.Errorf(`cs.ConstraintList()`)
		}
	}
}

func TestSetMinhNotWhitespace(t *testing.T) {
	cs := StringSet(StringMinLength(5), NonBlankString)
	assertEq(t, "min length 5, non-blank", cs.ConstraintDescription())
	if cs.IsValid("") {
		t.Errorf(`cs.IsValid("")`)
	}
	if cs.IsValid("hell") {
		t.Errorf(`cs.IsValid("hell")`)
	}
	if !cs.IsValid("hello") {
		t.Errorf(`!cs.IsValid("hello")`)
	}
	var err error = constraints.ValidOrError[string]("12345", cs)
	if err != nil {
		t.Errorf(`expecting <nil>, got %v`, err)
	}
	err = constraints.ValidOrError[string]("  ", cs)
	assertEq(t, "required to be min length 5, non-blank", err.Error())
}

func TestConst(t *testing.T) {
	emptyConst := constraints.Match("")
	if emptyConst.ConstraintDescription() != "match \"\"" {
		t.Errorf(`expecting "match \"\"", got %q`, emptyConst.ConstraintDescription())
	}
	if !emptyConst.IsValid("") {
		t.Errorf(`!emptyConst.IsValid("")`)
	}
	if emptyConst.IsValid("empty") {
		t.Errorf(`emptyConst.IsValid("empty")`)
	}
	helloConst := constraints.Match("hello")
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

func TestStringOneOf(t *testing.T) {
	emptyOneOf := constraints.OneOf[string]()
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
	oneOneOf := constraints.OneOf("one")
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
	oneOfOneTwoThree := constraints.OneOf("one", "two", "three")
	if oneOfOneTwoThree.ConstraintDescription() != "one of [one, two, three]" {
		t.Errorf(`expecting "one of [one, two, three]", got %q`,
			oneOfOneTwoThree.ConstraintDescription())
	}
}

func TestPrefix(t *testing.T) {
	emptyPrefix := StringPrefix("")
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
	noConsecutiveUnderscore := StringNoConsecutiveRune('_')
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
