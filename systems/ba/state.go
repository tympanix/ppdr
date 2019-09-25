package ba

import (
	"fmt"

	"github.com/tympanix/master-2019/ltl"
)

// State is a node in a GNBA
type State struct {
	ElementarySet ltl.Set
	Transitions   []Transition
}

func (s *State) String() string {
	return fmt.Sprint(s.ElementarySet)
}

// NewState returns a new state with no transitions with the given elementary set
func NewState(es ltl.Set) *State {
	return &State{
		ElementarySet: es,
		Transitions:   make([]Transition, 0),
	}
}

// AddTransition adds a transition from one state to another with labels for the transition
func (s *State) AddTransition(s1 *State, l ...ltl.Node) {
	s.Transitions = append(s.Transitions, Transition{
		State: s1,
		Label: ltl.NewSet(l...),
	})
}

// Has returns true if the gnba node has formula psi in the elementary set
func (s *State) Has(psi ltl.Node) bool {
	return s.ElementarySet.Contains(psi)
}

// Copy makes a copy of the state using a rename tabe to rename transitions
func (s *State) Copy() *State {
	trns := make([]Transition, 0)

	for _, t := range s.Transitions {
		trns = append(trns, t)
	}

	return &State{
		ElementarySet: s.ElementarySet.Copy(),
		Transitions:   trns,
	}
}

// Rename renames all state transitions using a rename table
func (s *State) Rename(rt RenameTable) {
	for i, t := range s.Transitions {
		s.Transitions[i] = t.Rename(rt)
	}
}

func (s *State) AddTransitionFromSet(state *State, label ltl.Set) {
	s.Transitions = append(s.Transitions, Transition{
		State: state,
		Label: label,
	})
}

func (s *State) ShouldHaveEdgeTo(state State, closure ltl.Set) bool {
	// case 1
	// s = B, state = B'
	for psi := range closure {
		if next, ok := psi.(ltl.Next); ok {
			if s.Has(next) != state.Has(next.ChildNode()) {
				return false
			}
		}
	}

	// case 2
	for psi := range closure {
		if until, ok := psi.(ltl.Until); ok {
			if s.Has(until) != (s.Has(until.RHSNode()) || (s.Has(until.LHSNode()) && state.Has(until))) {
				return false
			}
		}
	}

	return true
}
