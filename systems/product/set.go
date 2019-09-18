package product

// StateSet is a set of Product states
type StateSet map[*State]bool

// NewStateSet returns a new (possibly empty) set of states
func NewStateSet(states ...*State) StateSet {
	set := make(StateSet)
	set.Add(states...)
	return set
}

// Size returns the number of elements in the set
func (s StateSet) Size() int {
	return len(s)
}

// Add adds one or more states to the set
func (s StateSet) Add(states ...*State) {
	for _, n := range states {
		s[n] = true
	}
}

// Contains returns true if state is contained in the set
func (s StateSet) Contains(state *State) bool {
	_, ok := s[state]
	return ok
}

// Get removes and returns a random state from the set
func (s StateSet) Get() *State {
	for k := range s {
		delete(s, k)
		return k
	}
	return nil
}
