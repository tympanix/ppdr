package nba

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ba"
)

// NBA is a structure for non-deterministic Büchi automatons
type NBA struct {
	States      []*ba.State
	StartStates ba.StateSet
	FinalStates ba.StateSet
	Phi         ltl.Node
	AP          ltl.Set
}

// NewNBA returns a new empty NBA
func NewNBA(phi ltl.Node) *NBA {
	return &NBA{
		States:      make([]*ba.State, 0),
		StartStates: ba.NewStateSet(),
		FinalStates: ba.NewStateSet(),
		Phi:         phi,
	}
}

// AddState add one or more states to the NBA
func (n *NBA) AddState(states ...*ba.State) {
	n.States = append(n.States, states...)
}

// AddInitialState add one or more initials states to the NBA
func (n *NBA) AddInitialState(states ...*ba.State) {
	n.StartStates.Add(states...)
}

// AddAcceptanceState add one or more acceptance states to the NBA
func (n *NBA) AddAcceptanceState(states ...*ba.State) {
	n.FinalStates.Add(states...)
}

// IsAcceptanceState returns true if state is an accepting state in the NBA
func (n *NBA) IsAcceptanceState(state *ba.State) bool {
	return n.FinalStates.Contains(state)
}

// IsStartingState returns true if state is a starting state in the NBA
func (n *NBA) IsStartingState(state *ba.State) bool {
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
