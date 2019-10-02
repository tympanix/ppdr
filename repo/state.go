package repo

import (
	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ts"
)

// State is a data node in the repo
type State struct {
	predicates   ltl.Set
	dependencies []ts.State
	confPolicies ltl.Set
}

// NewState returns a new empty state
func NewState(pred ...ltl.Node) *State {
	return &State{
		predicates:   ltl.NewSet(pred...),
		dependencies: make([]ts.State, 0),
		confPolicies: ltl.NewSet(),
	}
}

func (s *State) addDependency(state *State) {
	s.dependencies = append(s.dependencies, state)
}

// Predicates returns the set of predicates which hold in the state
func (s *State) Predicates() ltl.Set {
	return s.predicates
}

// Dependencies return a list of dependencies from this state
func (s *State) Dependencies() []ts.State {
	return s.dependencies
}

func (s *State) addConfPolicy(set ltl.Set) {
	s.confPolicies.AddSet(set)
}
