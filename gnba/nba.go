package gnba

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/ltl"
)

// NBA is a structure for non-deterministic Büchi automatons
type NBA struct {
	States      []*State
	StartStates StateSet
	FinalStates StateSet
	Phi         ltl.Node
}

// NewNBA returns a new empty NBA
func NewNBA(phi ltl.Node) *NBA {
	return &NBA{
		States:      make([]*State, 0),
		StartStates: NewStateSet(),
		FinalStates: NewStateSet(),
		Phi:         phi,
	}
}

// AddState add one or more states to the NBA
func (n *NBA) AddState(states ...*State) {
	n.States = append(n.States, states...)
}

// AddInitialState add one or more initials states to the NBA
func (n *NBA) AddInitialState(states ...*State) {
	n.StartStates.Add(states...)
}

// AddAcceptanceState add one or more acceptance states to the NBA
func (n *NBA) AddAcceptanceState(states ...*State) {
	n.FinalStates.Add(states...)
}

// IsAcceptanceState returns true if state is an accepting state in the NBA
func (n *NBA) IsAcceptanceState(state *State) bool {
	return n.FinalStates.Contains(state)
}

// IsStartingState returns true if state is a starting state in the NBA
func (n *NBA) IsStartingState(state *State) bool {
	return n.StartStates.Contains(state)
}

func (n *NBA) String() string {
	var sb strings.Builder
	for _, s := range n.States {
		var prefix string
		if n.IsStartingState(s) {
			prefix = ">"
		}

		var suffix string
		if ok := n.IsAcceptanceState(s); ok {
			suffix = fmt.Sprintf("*")
		}

		fmt.Fprintf(&sb, "%s%s%s\n", prefix, s.ElementarySet, suffix)

		for _, t := range s.Transitions {
			fmt.Fprintf(&sb, "\t%s\t-->\t%s\n", t.Label, t.State.ElementarySet)
		}
	}

	return sb.String()
}
