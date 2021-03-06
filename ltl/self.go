package ltl

// Self is a predicate for the state to which a policy is applied
type Self struct{}

// SameAs returns true if node is also true
func (t Self) SameAs(node Node) bool {
	_, ok := node.(Self)
	return ok
}

func (t Self) Normalize() Node {
	return t
}

func (t Self) Compile(m *RefTable) Node {
	panic("self can not be compiled")
}

func (t Self) String() string {
	return "self"
}

func (t Self) Len() int {
	return 0
}

func (t Self) Map(fn MapFunc) Node {
	return fn(t)
}

func (t Self) Satisfied(r Resolver) bool {
	panic("self can not be checked for satisfyability")
}
