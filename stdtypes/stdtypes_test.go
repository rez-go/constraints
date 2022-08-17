package stdtypes

import (
	"reflect"
	"testing"
)

func assertEq(t *testing.T, expected interface{}, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\tExpected: %#v\n\tActual: %#v", expected, actual)
	}
}

func assertNeq(t *testing.T, referenceValue interface{}, actualValue interface{}) {
	t.Helper()
	if reflect.DeepEqual(referenceValue, actualValue) {
		t.Fatalf("\n\tNot expected: %#v\n\tActual: %#v", referenceValue, actualValue)
	}
}
