package constraints

import "errors"

// An Error is a specialized error which describes constraint violation(s).
type Error[
	ValueT any,
] interface {
	error
	ViolatedConstraint() Constraint[ValueT]
}

// ViolationError creates a new constraint violation error.
func ViolationError[
	ValueT any,
](c Constraint[ValueT]) Error[ValueT] {
	return &requirementError[ValueT]{c}
}

// ViolatedConstraintFromError attempts to extract violated Constraint
// from an error.
func ViolatedConstraintFromError[
	ValueT any,
](err error) Constraint[ValueT] {
	if err != nil {
		if e, ok := err.(Error[ValueT]); ok {
			return e.ViolatedConstraint()
		}
		return ViolatedConstraintFromError[ValueT](errors.Unwrap(err))
	}
	return nil
}

var (
	_ Error[any] = &requirementError[any]{}
)

type requirementError[ValueT any] struct {
	violated Constraint[ValueT]
}

func (e *requirementError[ValueT]) Error() string {
	if e != nil {
		if c := e.violated; c != nil {
			return "required to be " +
				c.ConstraintDescription()
		}
		return "constraint violation: <undefined>"
	}
	return "unknown constraint violation"
}

func (e *requirementError[ValueT]) ViolatedConstraint() Constraint[ValueT] {
	if e != nil {
		return e.violated
	}
	return nil
}
