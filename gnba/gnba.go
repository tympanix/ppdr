package gnba

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/ltl"
)

// GNBA is a structure of a generalized non-deterministic Büchi automaton
type GNBA struct {
	States         []*State
	StartingStates StateSet
	FinalStates    []StateSet
}

// NewGNBA return a new empty GNBA
func NewGNBA() *GNBA {
	return &GNBA{
		States:         make([]*State, 0),
		StartingStates: NewStateSet(),
		FinalStates:    make([]StateSet, 0),
	}
}

// IsAcceptanceState return true if state is in any of the acceptance sets
func (g *GNBA) IsAcceptanceState(state *State) (int, bool) {
	for i, s := range g.FinalStates {
		if s.Contains(state) {
			return i, true
		}
	}
	return -1, false
}

// IsStartingState returns true if state is a starting state for the GNBA
func (g *GNBA) IsStartingState(state *State) bool {
	return g.StartingStates.Contains(state)
}

func (g GNBA) String() string {
	var sb strings.Builder
	for _, s := range g.States {
		var prefix string
		if g.IsStartingState(s) {
			prefix = ">"
		}

		var suffix string
		if i, ok := g.IsAcceptanceState(s); ok {
			suffix = fmt.Sprintf("{%d}", i)
		}

		fmt.Fprintf(&sb, "%s%s%s\n", prefix, s.ElementarySet, suffix)

		for _, t := range s.Transitions {
			fmt.Fprintf(&sb, "\t%s\t-->\t%s\n", t.Label, t.State.ElementarySet)
		}
	}

	return sb.String()
}

// State is a node in a GNBA
type State struct {
	ElementarySet ltl.Set
	Transitions   []Transition
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
	for psi := range closure {
		if next, ok := psi.(ltl.Next); ok {
			if n.Has(next) != node.Has(next.ChildNode()) {
				return false
			}
		}
	}

	// case 2
	for psi := range closure {
		if until, ok := psi.(ltl.Until); ok {
			if n.Has(until) != (n.Has(until.RHSNode()) || (n.Has(until.LHSNode()) && node.Has(until))) {
				return false
			}
		}
	}

	return true
}
