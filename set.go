package constraints

import "strings"

// A Set is a constraint which consisted of a set of constraints.
type ConstraintSet[
	ValueT any,
	ConstraintT Constraint[ValueT],
] interface {
	Constraint[ValueT]
	Validate(v ValueT) Constraint[ValueT]
	ValidateAll(v ValueT) (violated []Constraint[ValueT])
	ConstraintList() []ConstraintT
}

func Set[ValueT any](
	constraints ...Constraint[ValueT],
) ConstraintSet[ValueT, Constraint[ValueT]] {
	return &constraintSet[ValueT]{
		constraints: constraints, //copy?
	}
}

var (
	_ Constraint[string]                        = constraintSet[string]{}
	_ ConstraintSet[string, Constraint[string]] = constraintSet[string]{}
)

// constraintSet defines a set of constraints. A value is considered valid
// if every constraint considers the value as valid.
type constraintSet[ValueT any] struct {
	constraints []Constraint[ValueT]
}

// ConstraintDescription conforms Constraint interface.
func (cs constraintSet[ValueT]) ConstraintDescription() string {
	if cs.constraints != nil {
		descs := make([]string, 0, len(cs.constraints))
		for _, ci := range cs.constraints {
			descs = append(descs, ci.ConstraintDescription())
		}
		return strings.Join(descs, ", ")
	}
	return ""
}

// ConstraintList conforms Set interface.
func (cs constraintSet[ValueT]) ConstraintList() []Constraint[ValueT] {
	if cs.constraints != nil {
		copyConstraints := make([]Constraint[ValueT], 0, len(cs.constraints))
		copy(copyConstraints, cs.constraints)
		return copyConstraints
	}
	return nil
}

// IsValid conforms Constraint interface.
func (cs constraintSet[ValueT]) IsValid(v ValueT) bool {
	return cs.Validate(v) == nil
}

// Validate checks the value against the constraints. If the value is
// valid, this method will return nil. Otherwise, it'll return the first
// violated constraint.
func (cs constraintSet[ValueT]) Validate(v ValueT) Constraint[ValueT] {
	for _, ci := range cs.constraints {
		if !ci.IsValid(v) {
			return ci
		}
	}
	return nil
}

// ValidateAll validates the value against all constraints. It returns
// a list of violate constraints if any.
func (cs constraintSet[ValueT]) ValidateAll(v ValueT) (violated []Constraint[ValueT]) {
	for _, ci := range cs.constraints {
		if !ci.IsValid(v) {
			violated = append(violated, ci)
		}
	}
	return violated
}
