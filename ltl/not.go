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

func (not Not) String() string {
	return fmt.Sprintf("!%v", not.ChildNode())
}
