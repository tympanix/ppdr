package gnba

// StateSet is a set of GNBA states
type StateSet map[*State]bool

// NewStateSet returns a new (possibly empty) set of states
func NewStateSet(states ...*State) StateSet {
	set := make(StateSet)
	set.Add(states...)
	return set
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
