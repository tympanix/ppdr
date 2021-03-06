package ltl

import (
	"fmt"
)

// AP is an atomic poropostion in LTL
type AP struct {
	Name string
}

// SameAs returns true if two atomic propositions are the same
func (ap AP) SameAs(node Node) bool {
	if ap2, ok := node.(AP); ok {
		return ap.Name == ap2.Name
	}
	return false
}

func (ap AP) String() string {
	return fmt.Sprintf("%v", ap.Name)
}

func (ap AP) Normalize() Node {
	return ap
}

func (ap AP) Compile(m *RefTable) Node {
	return ap
}

func (ap AP) Len() int {
	return 0
}

func (a AP) Satisfied(r Resolver) bool {
	return r.ResolveBool(a.Name)
}

func (a AP) Map(fn MapFunc) Node {
	return fn(a)
}