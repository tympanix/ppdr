package gnba

import "github.com/tympanix/master-2019/ltl"

// Node is a node in a GNBA
type Node struct {
	ElementarySet ltl.Set
	Transitions   []Transition
	IsStartState  bool
	IsFinishState bool
}

func (n Node) shouldHaveEdgeTo(node Node) bool {
	for _, e := range n.ElementarySet {
		if next, ok := e.(ltl.Next); ok {
			if !node.ElementarySet.Contains(next.ChildNode()) {
				return false
			}
		}
	}

	for _, e := range n.ElementarySet {
		if until, ok := e.(ltl.Until); ok {
			// TODO: implement ltl until rule for transitions
		}
	}

	return true
}
