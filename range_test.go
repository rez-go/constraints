package constraints

import "testing"

func TestRangeSimple(t *testing.T) {
	c := Range(0, 10)
	assertEq(t, "from 0 to 10", c.ConstraintDescription())
	assertEq(t, true, c.IsValid(1))
	assertEq(t, true, c.IsValid(0))
	assertEq(t, false, c.IsValid(-1))
	assertEq(t, false, c.IsValid(11))
}

func TestRangeOne(t *testing.T) {
	c := Range(1, 1)
	assertEq(t, "from 1 to 1", c.ConstraintDescription())
	assertEq(t, true, c.IsValid(1))
	assertEq(t, false, c.IsValid(0))
}
