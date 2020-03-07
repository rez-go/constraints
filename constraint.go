package constraints

// Constraint is base interface for all constraints. As Go doesn't support
// generics / contracts (yet), we are restrainting to define the
// interfaces which require types here. Ideally, we would like to have something
// like IsValid(type T)(v T) bool .
type Constraint interface {
	// ConstraintDescription returns the description of the constraint, e.g.,
	// "5 or less".
	//
	// This interface is a porcelain interface.
	ConstraintDescription() string
}

// A Set is a constraint which consisted of a set of constraints.
type Set interface {
	Constraint
	ConstraintList() []Constraint
}

// An EqOp specifies equality operator.
type EqOp int

// Supported equality operators.
const (
	EqOpEqual EqOp = iota
	EqOpNotEqual
	EqOpLess
	EqOpLessOrEqual
	EqOpGreater
	EqOpGreaterOrEqual
)

func (op EqOp) String() string {
	switch op {
	case EqOpEqual:
		return "equal"
	case EqOpNotEqual:
		return "not equal"
	case EqOpLess:
		return "less"
	case EqOpLessOrEqual:
		return "less or equal"
	case EqOpGreater:
		return "greater"
	case EqOpGreaterOrEqual:
		return "greater or equal"
	}
	return ""
}

// Symbol returns representative symbol of the operator.
func (op EqOp) Symbol() string {
	switch op {
	case EqOpEqual:
		return "="
	case EqOpNotEqual:
		return "≠"
	case EqOpLess:
		return "<"
	case EqOpLessOrEqual:
		return "≤"
	case EqOpGreater:
		return ">"
	case EqOpGreaterOrEqual:
		return "≥"
	}
	return "?"
}

// StringFormat returns a string which could be used to in *printf functions.
// The string format expects a value to be passed.
func (op EqOp) StringFormat() string {
	switch op {
	case EqOpEqual:
		return "equals %v"
	case EqOpNotEqual:
		return "not equal to %v"
	case EqOpLess:
		return "less than %v"
	case EqOpLessOrEqual:
		return "less than or equal to %v"
	case EqOpGreater:
		return "greater than %v"
	case EqOpGreaterOrEqual:
		return "greater than or equal to %v"
	}
	return "(%v)"
}
