package ltl

import "fmt"

// Eventually is the structure for LTL eventually
type Eventually struct {
	Child Node
}

// ChildNode returns the child node of the eventually op
func (e Eventually) ChildNode() Node {
	return e.Child
}

func (e Eventually) String() string {
	if _, ok := e.ChildNode().(BinaryNode); ok {
		return fmt.Sprintf("<>(%v)", e.ChildNode())
	}
	return fmt.Sprintf("<>%v", e.ChildNode())
}

// SameAs returns true if both nodes are ltl eventually and has
// identical sub-tree
func (e Eventually) SameAs(node Node) bool {
	if n, ok := node.(Eventually); ok {
		return e.ChildNode().SameAs(n.ChildNode())
	}
	return false
}

func (e Eventually) Normalize() Node {
	return Until{True{}, e.ChildNode().Normalize()}
}

func (e Eventually) Compile(m *RefTable) Node {
	return Eventually{e.ChildNode().Compile(m)}
}

func (e Eventually) Len() int {
	return 1 + e.ChildNode().Len()
}
