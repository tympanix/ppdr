package product

// Statespace is a set of Product states
type Statespace map[statespaceEntry]bool

type statespaceEntry struct {
	state *State
	tag   int
}

// NewStatespace returns a new (possibly empty) set of states
func NewStatespace() Statespace {
	space := make(Statespace)
	return space
}

// Size returns the number of elements in the set
func (s Statespace) Size() int {
	return len(s)
}

// Add adds one or more states to the set
func (s Statespace) Add(entries ...statespaceEntry) {
	for _, n := range entries {
		s[n] = true
	}
}

// Contains returns true if state is contained in the set
func (s Statespace) Contains(entry statespaceEntry) bool {
	_, ok := s[entry]
	return ok
}

// Get removes and returns a random state from the set
func (s Statespace) Get() *State {
	for k := range s {
		delete(s, k)
		return k.state
	}
	return nil
}
