package ltl

import (
	"fmt"
)

// Not is negation in LTL
type Not struct {
	Child Node
}

func (not Not) SameAs(node Node) bool {
	if not2, ok := node.(Not); ok {
		return not.Child.SameAs(not2.Child)
	}
	return false
}

func (not Not) ChildNode() Node {
	return not.Child
}

func (not Not) Compile(m *RefTable) Node {
	return Not{not.ChildNode().Compile(m)}
}

func (not Not) String() string {
	if _, ok := not.Child.(BinaryNode); ok {
		return fmt.Sprintf("!(%v)", not.ChildNode())
	}
	return fmt.Sprintf("!%v", not.ChildNode())
}

func (not Not) Normalize() Node {
	return Negate(not.ChildNode().Normalize())
}

func (not Not) Len() int {
	return 1 + not.ChildNode().Len()
}

func (not Not) Map(fn MapFunc) Node {
	return fn(Not{not.Child.Map(fn)})
}

func (n Not) Satisfied(r Resolver) bool {
	if child, ok := n.ChildNode().(Satisfiable); ok {
		return !child.Satisfied(r)
	}
	panic(ErrNotPropositional)
}
