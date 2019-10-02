package repo

import (
	"errors"

	"github.com/tympanix/master-2019/ltl"
)

// Repo is the data repository itself
type Repo struct {
	states map[*State]bool
	origin *State
}

// NewRepo returns a new empty repo
func NewRepo() *Repo {
	s0 := NewState()
	s0.addDependency(s0)

	r := &Repo{
		states: make(map[*State]bool),
		origin: s0,
	}

	return r
}

func (r *Repo) addState(state *State) {
	if len(state.Dependencies()) == 0 {
		state.addDependency(r.origin)
	}
	r.states[state] = true
}

// Query performs a lookup in the data repository with a integrity policy
func (r *Repo) Query(state *State, intr ltl.Node) (*State, error) {

	if _, ok := r.states[state]; !ok {
		return nil, errors.New("state not found in repo")
	}

	c := candidate{state}

	if !c.satisfiesFormula(intr) {
		return nil, errors.New("integrity not satisfied")
	}

	if !c.satisfiesConfPolicies() {
		return nil, errors.New("confidentiality not satisfied")
	}

	return state, nil

}

// Put adds a new data point to the repository with a confidentiality policy
func (r *Repo) Put(state *State) bool {
	if _, ok := r.states[state]; ok {
		return false
	}
	for _, d := range state.Dependencies() {
		if s, ok := d.(*State); ok {
			if _, ok := r.states[s]; ok {
				state.addConfPolicy(s.confPolicies)
			} else {
				panic("unknown dependency")
			}
		} else {
			panic("unknown repo state")
		}
	}

	c := candidate{state}

	if c.satisfiesConfPolicies() {
		r.addState(state)
		return true
	}
	return false
}
