package gnba

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/ltl"
)

// GNBA is a list of nodes
type GNBA []*State

func (g GNBA) String() string {
	var sb strings.Builder
	for _, s := range g {
		fmt.Fprintf(&sb, "%s\n", s.ElementarySet)

		for _, t := range s.Transitions {
			fmt.Fprintf(&sb, "\t%s\t->\t%s\n", t.Label, t.State.ElementarySet)
		}
	}

	return sb.String()
}

// State is a node in a GNBA
type State struct {
	ElementarySet ltl.Set
	Transitions   []Transition
	IsStartState  bool
	IsFinishState bool
}

// Has returns true if the gnba node has formula psi in the elementary set
func (n *State) Has(psi ltl.Node) bool {
	return n.ElementarySet.Contains(psi)
}

func (n *State) addTransition(node *State, label ltl.Set) {
	n.Transitions = append(n.Transitions, Transition{
		State: node,
		Label: label,
	})
}

func (n State) shouldHaveEdgeTo(node State, closure ltl.Set) bool {
	// case 1
	// n = B, node = B'
	for _, psi := range closure {
		if next, ok := psi.(ltl.Next); ok {
			if n.Has(next) != node.Has(next.ChildNode()) {
				return false
			}
		}
	}

	// case 2
	for _, psi := range closure {
		if until, ok := psi.(ltl.Until); ok {
			if n.Has(until) != (n.Has(until.RHSNode()) || (n.Has(until.LHSNode()) && node.Has(until))) {
				return false
			}
		}
	}

	return true
}
