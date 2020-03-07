package constraints

import "errors"

// An Error is a specialized error which describes constraint violation(s).
type Error interface {
	error
	ViolatedConstraint() Constraint
}

// ViolationError creates a new constraint violation error.
func ViolationError(c Constraint) Error {
	return &requirementError{c}
}

// ViolatedConstraintFromError attempts to extract violated Constraint
// from an error.
func ViolatedConstraintFromError(err error) Constraint {
	if err != nil {
		if e, ok := err.(Error); ok {
			return e.ViolatedConstraint()
		}
		return ViolatedConstraintFromError(errors.Unwrap(err))
	}
	return nil
}

var (
	_ Error = &requirementError{}
)

type requirementError struct {
	violated Constraint
}

func (e *requirementError) Error() string {
	if e != nil {
		if c := e.violated; c != nil {
			return "required to be " +
				c.ConstraintDescription()
		}
		return "constraint violation: <undefined>"
	}
	return "unknown constraint violation"
}

func (e *requirementError) ViolatedConstraint() Constraint {
	if e != nil {
		return e.violated
	}
	return nil
}
