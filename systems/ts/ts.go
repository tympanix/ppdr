package ts

import (
	"fmt"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ba"
)

// State is a node in the transition system
type State struct {
	Predicates   ltl.Set
	Dependencies []*State
}

func (s *State) String() string {
	return fmt.Sprint(s.Predicates)
}

// NewState creates and return a new state
func NewState(p ...ltl.Node) *State {
	return &State{
		Predicates:   ltl.NewSet(p...),
		Dependencies: make([]*State, 0),
	}
}

// AddDependency adds an dependency/transition to another state
func (s *State) AddDependency(s1 *State) {
	s.Dependencies = append(s.Dependencies, s1)
}

// ShouldHaveTransitionTo return true if TS state should have a transition
// to an NBA state with the given labeling function (ltl set)
func (s *State) ShouldHaveTransitionTo(t ba.Transition, lf ltl.Set) bool {
	if lf.Size() == 0 && t.Label.Size() == 0 {
		return true
	}

	if t.Label.Size() != 0 && lf.Size() == 0 {
		return false
	}

	return t.Label.ContainsAll(lf)
	//return lf.ContainsAny(t.Label)
}

// TS is a transition system
type TS struct {
	States        []*State
	InitialStates []*State
}

// New creates a new transition system, TS
func New() *TS {
	return &TS{
		States:        make([]*State, 0),
		InitialStates: make([]*State, 0),
	}
}

// AddState add one or more states to the transition system
func (t *TS) AddState(states ...*State) {
	t.States = append(t.States, states...)
}

// AddInitialState add one or more states as initial states to the transition system
func (t *TS) AddInitialState(states ...*State) {
	t.InitialStates = append(t.InitialStates, states...)
}
