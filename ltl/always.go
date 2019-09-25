package ltl

import "fmt"

// Always is the structure for LTL Always
type Always struct {
	Child Node
}

// ChildNode returns the child node of tha Always op
func (a Always) ChildNode() Node {
	return a.Child
}

func (a Always) String() string {
	if _, ok := a.ChildNode().(BinaryNode); ok {
		return fmt.Sprintf("[](%v)", a.ChildNode())
	}
	return fmt.Sprintf("[]%v", a.ChildNode())
}

// SameAs returns true if both nodes are ltl Always and has
// identical sub-tree
func (a Always) SameAs(node Node) bool {
	if n, ok := node.(Always); ok {
		return a.ChildNode().SameAs(n.ChildNode())
	}
	return false
}

func (a Always) Normalize() Node {
	return Not{Eventually{Negate(a.ChildNode().Normalize())}.Normalize()}
}

func (a Always) Len() int {
	return 1 + a.Child.Len()
}
