package constraints

import "strings"

func Any[ValueT any](
	constraints ...Constraint[ValueT],
) Constraint[ValueT] {
	return &anyConstraint[ValueT]{constraints: constraints}
}

type anyConstraint[ValueT any] struct {
	constraints []Constraint[ValueT]
}

func (ac anyConstraint[ValueT]) ConstraintDescription() string {
	if ac.constraints != nil {
		descs := make([]string, 0, len(ac.constraints))
		for _, ci := range ac.constraints {
			descs = append(descs, ci.ConstraintDescription())
		}
		return strings.Join(descs, " or ")
	}
	return ""
}

func (ac anyConstraint[ValueT]) IsValid(v ValueT) bool {
	if ac.constraints != nil {
		for _, c := range ac.constraints {
			if c.IsValid(v) {
				return true
			}
		}
	}
	return false
}
