package ltl

// Subject is a function for retrieving user identities
type Subject struct{}

// SameAs returns true if node is also true
func (r Subject) SameAs(node Node) bool {
	_, ok := node.(Subject)
	return ok
}

func (r Subject) Normalize() Node {
	return r
}

func (r Subject) Compile(m *RefTable) Node {
	return (*m)[UserRef]
}

func (r Subject) String() string {
	return "subject"
}

func (r Subject) Len() int {
	return 0
}

func (r Subject) Map(fn MapFunc) Node {
	return fn(r)
}

func (r Subject) Satisfied(res Resolver) bool {
	panic("subject can not be checked for satisfyability")
}
