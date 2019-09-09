package ltl

// True is a logical true
type True struct{}

// SameAs returns true if node is also true
func (t True) SameAs(node Node) bool {
	_, ok := node.(True)
	return ok
}

func (t True) Normalize() Node {
	return t
}

func (t True) String() string {
	return "true"
}
