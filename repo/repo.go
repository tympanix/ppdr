package repo

import (
	"crypto/rand"
	"errors"
	"unsafe"

	"github.com/tympanix/master-2019/ltl"
)

// Identity is the identity of a user
type Identity struct {
	name string
	id   []byte
}

// NewIdentity generates a new identity
func NewIdentity(name string) *Identity {
	id := make([]byte, 4)
	_, err := rand.Read(id)
	if err != nil {
		panic(err)
	}
	return &Identity{
		name: name,
		id:   id,
	}
}

// Repo is the data repository itself
type Repo struct {
	states      map[*State]bool
	users       map[string]*Identity
	currentUser *Identity
}

// NewRepo returns a new empty repo
func NewRepo() *Repo {
	r := &Repo{
		states: make(map[*State]bool),
		users:  make(map[string]*Identity),
	}

	return r
}

// SetCurrentUser changes the current user
func (r *Repo) SetCurrentUser(user *Identity) {
	r.users[user.name] = user
	r.currentUser = user
}

func (r *Repo) addState(states ...*State) {
	for _, s := range states {
		if len(s.Dependencies()) == 0 {
			s.addDependency(s)
		}
		r.states[s] = true
	}
}

// RenameUserPredicate renames all user predicates
func (r *Repo) renameUserPredicate(phi ltl.Node) ltl.Node {
	return phi.Map(func(n ltl.Node) ltl.Node {
		if u, ok := n.(ltl.User); ok {
			return ltl.Ptr{
				Attr:    "user",
				Pointer: unsafe.Pointer(r.users[u.Name]),
			}
		}
		return n
	})
}

func (r *Repo) getUserPredicate() ltl.Ptr {
	return ltl.Ptr{
		Attr:    "user",
		Pointer: unsafe.Pointer(r.currentUser),
	}
}

// Query performs a lookup in the data repository with a integrity policy
func (r *Repo) Query(state *State, intr ltl.Node) (*State, error) {
	// Ensure that current user is set
	if r.currentUser == nil {
		panic("current user must be set")
	}

	// Ensure that state actually exists in data repo
	if _, ok := r.states[state]; !ok {
		return nil, errors.New("state not found in repo")
	}

	// Make user predicate reference appropriate identities
	intr = r.renameUserPredicate(intr)

	// Make self predicates reference the argument state
	intr = ltl.RenameSelfPredicate(intr, state.getSelfAttr())

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
	// Ensure that the user is set
	if r.currentUser == nil {
		panic("current user must be set")
	}

	// Make sure state doesn't already exist
	if _, ok := r.states[state]; ok {
		return errors.New("state already exists")
	}

	// Save attribute author in this state
	state.newAttrPtr("author", unsafe.Pointer(r.currentUser))

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
