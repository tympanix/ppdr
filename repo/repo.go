package repo

import (
	"errors"
	"unsafe"

	"github.com/tympanix/master-2019/ltl"
)

// Repo is the data repository itself
type Repo struct {
	states map[*State]bool
}

// NewRepo returns a new empty repo
func NewRepo() *Repo {
	r := &Repo{
		states: make(map[*State]bool),
	}

	return r
}

func (r *Repo) addState(states ...*State) {
	for _, s := range states {
		if len(s.Dependencies()) == 0 {
			s.addDependency(s)
		}
		r.states[s] = true
	}
}

// Query performs a lookup in the data repository with a integrity policy
func (r *Repo) Query(state *State, intr ltl.Node) (*State, error) {

	// Ensure that state actually exists in data repo
	if _, ok := r.states[state]; !ok {
		return nil, errors.New("state not found in repo")
	}

	// Make self predicates reference the argument state
	intr = ltl.RenameSelfPredicate(intr, unsafe.Pointer(state))

	// Ensure integrity policy is satisfied
	c := candidate{state}
	if !c.satisfiesFormula(intr) {
		return nil, errors.New("integrity not satisfied")
	}

	// Ensure conf policies are satsified
	if !c.satisfiesConfPolicies() {
		return nil, errors.New("confidentiality not satisfied")
	}

	// Accepted, return state
	return state, nil

}

// Put adds a new data point to the repository with a confidentiality policy
func (r *Repo) Put(state *State) error {
	// Make sure state doesn't already exist
	if _, ok := r.states[state]; ok {
		return errors.New("state already exists")
	}

	// Make self predicates reference this state
	state.replaceSelfReferences()

	// Add conf policies from dependents
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

	// Ensure conf policies are satisfied before submitting
	c := candidate{state}
	if !c.satisfiesConfPolicies() {
		return errors.New("confidentiality not satisfied")
	}

	// Submit state to data repo
	r.addState(state)
	return nil
}
