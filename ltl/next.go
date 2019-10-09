package ltl

import (
	"fmt"
)

// Next is the ltl next operator
type Next struct {
	Child Node
}

// SameAs returns true if the other node is also a next operator
// and has an identical sub-tree
func (next Next) SameAs(node Node) bool {
	if next2, ok := node.(Next); ok {
		return next.Child.SameAs(next2.Child)
	}
	return false
}

// ChildNode returns the child of the next operator
func (next Next) ChildNode() Node {
	return next.Child
}

func (next Next) String() string {
	return fmt.Sprintf("O%v", next.ChildNode())
}

func (next Next) Normalize() Node {
	return Next{next.ChildNode().Normalize()}
}

func (next Next) Compile(m *RefTable) Node {
	return Next{next.ChildNode().Compile(m)}
}

func (next Next) Map(fn MapFunc) Node {
	return fn(Next{next.ChildNode().Map(fn)})
}

func (next Next) Len() int {
	return 1 + next.ChildNode().Len()
}
