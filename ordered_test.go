package constraints

import "testing"

func TestMinIntSimple(t *testing.T) {
	var c Constraint[int] = Min(5)
	assertEq(t, "min 5", c.ConstraintDescription())
	assertEq(t, true, c.IsValid(5))
	assertEq(t, true, c.IsValid(10))
	assertEq(t, false, c.IsValid(4))
	assertEq(t, false, c.IsValid(-1))
}

func TestMaxIntSimple(t *testing.T) {
	var c Constraint[int] = Max(5)
	assertEq(t, "max 5", c.ConstraintDescription())
	assertEq(t, true, c.IsValid(5))
	assertEq(t, false, c.IsValid(10))
	assertEq(t, true, c.IsValid(4))
	assertEq(t, true, c.IsValid(-1))
}
