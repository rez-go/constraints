package stdtypes

var (
	_ StringConstraint = targetOperandFuncConstraint[string]{}
)

type targetOperandFuncConstraint[ValueT any] struct {
	desc    string
	operand ValueT
	fn      func(target, operand ValueT) bool
}

func (c targetOperandFuncConstraint[ValueT]) ConstraintDescription() string {
	return c.desc
}

func (c targetOperandFuncConstraint[ValueT]) IsValid(v ValueT) bool {
	return c.fn != nil && c.fn(v, c.operand)
}
