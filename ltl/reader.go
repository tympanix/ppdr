package ltl

// Reader is a function for retrieving user identities
type Reader struct{}

// SameAs returns true if node is also true
func (r Reader) SameAs(node Node) bool {
	_, ok := node.(Reader)
	return ok
}

func (r Reader) Normalize() Node {
	return r
}

func (r Reader) Compile(m *RefTable) Node {
	return (*m)[UserRef]
}

func (r Reader) String() string {
	return "reader"
}

func (r Reader) Len() int {
	return 0
}

func (r Reader) Map(fn MapFunc) Node {
	return fn(r)
}

func (r Reader) Satisfied(res Resolver) bool {
	panic("reader can not be checked for satisfyability")
}
