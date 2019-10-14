package ltl

import "unsafe"

// Ptr is a logical true
type Ptr struct {
	Pointer unsafe.Pointer
}

// SameAs returns true if node is also true
func (p Ptr) SameAs(node Node) bool {
	if p2, ok := node.(Ptr); ok {
		return p.Pointer == p2.Pointer
	}
	return false
}

func (p Ptr) Normalize() Node {
	return p
}

func (p Ptr) Compile(m *RefTable) Node {
	return m.NewRef(p)
}

func (p Ptr) String() string {
	return "ptr"
}

func (p Ptr) Len() int {
	return 0
}

func (p Ptr) Map(fn MapFunc) Node {
	return fn(p)
}

func (p Ptr) Satisfied(r Resolver) bool {
	return r.Resolve("self").SameAs(p)
}
