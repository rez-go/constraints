package testing

import (
	"fmt"
	"reflect"
	"testing"
)

type T = testing.T

func AssertEq(t *testing.T, expected interface{}, actual interface{}, extra ...any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		var formatted string
		if len(extra) > 0 {
			if msg, ok := (extra[0]).(string); ok {
				if len(extra) > 1 {
					formatted = fmt.Sprintf(msg, extra[1:]...)
				} else {
					formatted = msg
				}
			}
		}
		if formatted != "" {
			t.Fatalf("\n%s\n\tExpected: %#v\n\tActual: %#v", formatted, expected, actual)
		} else {
			t.Fatalf("\n\tExpected: %#v\n\tActual: %#v", expected, actual)
		}
	}
}

func AssertNeq(t *testing.T, referenceValue interface{}, actualValue interface{}) {
	t.Helper()
	if reflect.DeepEqual(referenceValue, actualValue) {
		t.Fatalf("\n\tNot expected: %#v\n\tActual: %#v", referenceValue, actualValue)
	}
}
