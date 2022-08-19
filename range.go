package constraints

import (
	"fmt"

	typecons "golang.org/x/exp/constraints"
)

func Range[ValueT typecons.Ordered](min, max ValueT) Constraint[ValueT] {
	return &rangeConstraint[ValueT]{
		inclusive: true,
		min:       min,
		max:       max,
	}
}

type rangeConstraint[ValueT typecons.Ordered] struct {
	inclusive bool //TODO: inclusiveMin & inclusiveMax
	min       ValueT
	max       ValueT
}

var (
	_ Constraint[int] = rangeConstraint[int]{}
	_ Constraint[int] = &rangeConstraint[int]{}
)

func (rc rangeConstraint[ValueT]) ConstraintDescription() string {
	return fmt.Sprintf("from %v to %v", valueLiteralString(rc.min), valueLiteralString(rc.max))
}

func (rc rangeConstraint[ValueT]) IsValid(v ValueT) bool {
	if rc.inclusive {
		return v >= rc.min && v <= rc.max
	}
	return v > rc.min && v < rc.max
}
