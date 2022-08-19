package constraints

// ValidatorFunc is an adapter to allow use of ordinary functions
// as a validator in constraints.
type ValidatorFunc[ValueT any] func(i ValueT) bool

// Func creates a constraint from a validator function.
//
// A good example is for defining a constraint to ensure that provided
// string is a valid UTF8-encoded string:
//
//	import "unicode/utf8"
//
//	var mustUTF8 = Func("valid UTF-8", utf8.ValidString)
//
// Or regular expression string matcher:
//
//	var usernamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
//	var usernameConstraint = Func("username", usernamePattern.MatchString)
func Func[
	ValueT any,
](desc string, fn func(v ValueT) bool) Constraint[ValueT] {
	return &constraintFunc[ValueT]{
		negate: false,
		desc:   desc,
		fn:     fn,
	}
}

var (
	_ Constraint[int64] = &constraintFunc[int64]{}
	_ Constraint[int64] = constraintFunc[int64]{}
)

type constraintFunc[ValueT any] struct {
	negate bool
	desc   string
	fn     ValidatorFunc[ValueT]
}

func (c constraintFunc[ValueT]) ConstraintDescription() string {
	return c.desc
}

func (c constraintFunc[ValueT]) IsValid(v ValueT) bool {
	result := c.fn(v)
	if c.negate {
		return !result
	}
	return result
}
