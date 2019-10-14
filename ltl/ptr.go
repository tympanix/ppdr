package ltl

import "unsafe"

// Ptr is a logical true
type Ptr struct {
	Pointer unsafe.Pointer
}

// SameAs returns true if node is also true
func (t Ptr) SameAs(node Node) bool {
	if p2, ok := node.(Ptr); ok {
		return t.Pointer == p2.Pointer
	}
	return false
}

func (t Ptr) Normalize() Node {
	return t
}

func (t Ptr) Compile(m *RefTable) Node {
	return t
}

func (t Ptr) String() string {
	return "ptr"
}

func (t Ptr) Len() int {
	return 0
}

func (t Ptr) Map(fn MapFunc) Node {
	return fn(t)
}

func (t Ptr) Satisfied(r Resolver) bool {
	return r.Resolve("self").SameAs(t)
}
