package ts

import (
	"fmt"

	"github.com/tympanix/master-2019/ltl"
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
