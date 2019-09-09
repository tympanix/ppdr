package gnba

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/ltl"
)

// GNBA is a list of nodes
type GNBA []*Node

func (g GNBA) String() string {
	var sb strings.Builder
	for _, s := range g {
		fmt.Fprintf(&sb, "%s\n", s.ElementarySet)

		for _, t := range s.Transitions {
			fmt.Fprintf(&sb, "\t%s\t->\t%s\n", t.Label, t.Node.ElementarySet)
		}
	}

	return sb.String()
}

// Node is a node in a GNBA
type Node struct {
	ElementarySet ltl.Set
	Transitions   []Transition
	IsStartState  bool
	IsFinishState bool
}

func (n *Node) addTransition(node *Node, label ltl.Set) {
	n.Transitions = append(n.Transitions, Transition{
		Node:  node,
		Label: label,
	})
}

func (n Node) shouldHaveEdgeTo(node Node, closure ltl.Set) bool {
	// case 1
	// n = B, node = B'
	for _, psi := range closure {
		if next, ok := psi.(ltl.Next); ok {
			if n.ElementarySet.Contains(next) != node.ElementarySet.Contains(next.ChildNode()) {
				return false
			}
		}
	}

	// case 2
	for _, e := range n.ElementarySet {
		if _, ok := e.(ltl.Until); ok {
			panic("TODO: implement ltl until rule for transitions")
			// TODO: implement ltl until rule for transitions
		}
	}

	return true
}
