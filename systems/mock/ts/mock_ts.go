package ts

import (
	"fmt"

	"github.com/tympanix/ppdr/ltl"
	"github.com/tympanix/ppdr/systems/ts"
)

// State is a node in the transition system
type State struct {
	predicates   ltl.Set
	dependencies []ts.State
}

// Dependencies return the list of dependencies from this state
func (s *State) Dependencies() []ts.State {
	return s.dependencies
}

// Predicates returns the set of predicates which hold in this state
func (s *State) Predicates(ap ltl.Set, t ltl.RefTable) ltl.Set {
	return s.predicates
}

func (s *State) String() string {
	return fmt.Sprint(s.predicates)
}

// NewState creates and return a new state
func NewState(p ...ltl.Node) *State {
	return &State{
		predicates:   ltl.NewSet(p...),
		dependencies: make([]ts.State, 0),
	}
}

// AddDependency adds an dependency/transition to another state
func (s *State) AddDependency(s1 *State) {
	s.dependencies = append(s.dependencies, s1)
}

// TS is a transition system
type TS struct {
	states        []ts.State
	initialStates []ts.State
}

// New creates a new transition system, TS
func New() *TS {
	return &TS{
		states:        make([]ts.State, 0),
		initialStates: make([]ts.State, 0),
	}
}

// States returns the list of states in the transition system
func (t *TS) States() []ts.State {
	return t.states
}

// InitialStates return the list of initial states in the transition system
func (t *TS) InitialStates() []ts.State {
	return t.initialStates
}

// AddState add one or more states to the transition system
func (t *TS) AddState(s ...ts.State) {
	t.states = append(t.states, s...)
}

// AddInitialState add one or more states as initial states to the transition system
func (t *TS) AddInitialState(s ...ts.State) {
	t.initialStates = append(t.initialStates, s...)
}
