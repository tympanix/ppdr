package gnba

// StateSet is a set of GNBA states
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

// Copy returns a new set with pointers renamed by a renaming table
func (s StateSet) Copy(rt renameTable) StateSet {
	set := NewStateSet()
	for s := range s {
		set.Add(rt[s])
	}
	return set
}
