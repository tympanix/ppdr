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
