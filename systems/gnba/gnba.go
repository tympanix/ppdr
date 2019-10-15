package gnba

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ba"
)

// GNBA is a structure of a generalized non-deterministic Büchi automaton
type GNBA struct {
	States         []*ba.State
	StartingStates ba.StateSet
	FinalStates    []ba.StateSet
	Phi            ltl.Node
	AP             ltl.Set
}

// NewGNBA return a new empty GNBA
func NewGNBA(phi ltl.Node) *GNBA {
	return &GNBA{
		States:         make([]*ba.State, 0),
		StartingStates: ba.NewStateSet(),
		FinalStates:    make([]ba.StateSet, 0),
		Phi:            phi,
		AP:             ltl.FindAtomicPropositions(phi),
	}
}

// IsAcceptanceState return true if state is in any of the acceptance sets
func (g *GNBA) IsAcceptanceState(state *ba.State) (int, bool) {
	for i, s := range g.FinalStates {
		if s.Contains(state) {
			return i, true
		}
	}
	return -1, false
}

func (g *GNBA) getAcceptanceStateSets(state *ba.State) []int {
	f := make([]int, 0)
	for i, s := range g.FinalStates {
		if s.Contains(state) {
			f = append(f, i)
		}
	}
	return f
}

// IsStartingState returns true if state is a starting state for the GNBA
func (g *GNBA) IsStartingState(state *ba.State) bool {
	return g.StartingStates.Contains(state)
}

// HasState returns true if state is a part of the GNBA
func (g *GNBA) HasState(state *ba.State) bool {
	for _, s := range g.States {
		if s == state {
			return true
		}
	}
	return false
}

// FindStateIndex finds the index of the state in the GNBA structure
func (g *GNBA) FindStateIndex(state *ba.State) int {
	for i, s := range g.States {
		if s == state {
			return i
		}
	}
	return -1
}

// Copy creates a copy of the GNBA
func (g *GNBA) Copy() *GNBA {
	gnba := NewGNBA(g.Phi)

	var rt = make(ba.RenameTable)

	// Create a copy of each state and add to rename table
	for _, s := range g.States {
		copy := s.Copy()
		rt[s] = copy
		gnba.States = append(gnba.States, copy)
	}

	// Translate state transitions with renaming table
	for _, s := range gnba.States {
		s.Rename(rt)
	}

	// Copy and rename starting states
	for s := range g.StartingStates {
		gnba.StartingStates.Add(rt[s])
	}

	// Copy and rename acceptance set
	accSet := make([]ba.StateSet, 0)
	for _, s := range g.FinalStates {
		accSet = append(accSet, s.Copy(rt))
	}
	gnba.FinalStates = accSet

	return gnba
}

// StringWithRenamer strings the GNBA using a naming function
func (g GNBA) StringWithRenamer(r ba.StateNamer) string {
	var sb strings.Builder
	for _, s := range g.States {
		var prefix string
		if g.IsStartingState(s) {
			prefix = ">"
		}

		var suffix string
		if f := g.getAcceptanceStateSets(s); len(f) > 0 {
			suffix = fmt.Sprintf("{%d}", f)
		}

		fmt.Fprintf(&sb, "%s%s%s\n", prefix, r.RenameState(s), suffix)

		for _, t := range s.Transitions {
			fmt.Fprintf(&sb, "\t%s\t-->\t%s\n", t.Label, r.RenameState(t.State))
		}
	}

	return sb.String()
}

func (g GNBA) String() string {
	return g.StringWithRenamer(nil)
}
