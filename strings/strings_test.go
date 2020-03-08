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
	if err := Empty.ValidOrError(""); err != nil {
		t.Errorf(`err := Empty.ValidOrError(""); err != nil`)
	}
	err := Empty.ValidOrError("empty")
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
	len10 := Length(10)
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
	len0 := Length(0)
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
	min0 := MinLength(0)
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
	min10 := MinLength(10)
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
	max0 := MaxLength(0)
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
	max10 := MaxLength(10)
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
	err := cs.ValidOrError("")
	if err.Error() != "required to be not empty" {
		t.Errorf(`expecting "required to be not empty", got %q`, err.Error())
	}
}

func TestSetMinhNotWhitespace(t *testing.T) {
	cs := Set{MinLength(5), NotWhitespace}
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
	err := cs.ValidOrError("   ")
	if err.Error() != "required to be min length 5, not whitespace" {
		t.Errorf(`expecting "required to be min length 5, not whitespace", got %q`, err.Error())
	}
}

func TestConst(t *testing.T) {
	emptyConst := Const("")
	if emptyConst.ConstraintDescription() != "const \"\"" {
		t.Errorf(`expecting "const \"\"", got %q`, emptyConst.ConstraintDescription())
	}
	if !emptyConst.IsValid("") {
		t.Errorf(`!emptyConst.IsValid("")`)
	}
	if emptyConst.IsValid("empty") {
		t.Errorf(`emptyConst.IsValid("empty")`)
	}
	helloConst := Const("hello")
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

func TestAny(t *testing.T) {
	emptyAny := Any()
	if emptyAny.ConstraintDescription() != "in []" {
		t.Errorf(`expecting "in []", got %q`, emptyAny.ConstraintDescription())
	}
	if emptyAny.IsValid("") {
		t.Errorf(`emptyIn.IsValid("")`)
	}
	if emptyAny.IsValid("nothing") {
		t.Errorf(`emptyIn.IsValid("nothing")`)
	}
	oneAny := Any("one")
	if oneAny.IsValid("") {
		t.Errorf(`oneIn.IsValid("")`)
	}
	if oneAny.IsValid("zero") {
		t.Errorf(`oneIn.IsValid("zero")`)
	}
	if !oneAny.IsValid("one") {
		t.Errorf(`!oneIn.IsValid("one")`)
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
