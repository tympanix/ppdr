package product

// StateStack is a stack of Product states
type StateStack []*State

// NewStateStack returns a new (possibly empty) stack of states
func NewStateStack(states ...*State) StateStack {
	stack := make(StateStack, 0)
	stack.Push(states...)
	return stack
}

// Length returns the number of elements in the stack
func (s StateStack) Length() int {
	return len(s)
}

// Empty returns true if the stack is empty
func (s StateStack) Empty() bool {
	return s.Length() == 0
}

// Push adds one or more states to the stack
func (s *StateStack) Push(states ...*State) {
	*s = append(*s, states...)
}

// Pop removes and returns the top state on the stack
func (s *StateStack) Pop() *State {

	state := s.Peek()

	*s = (*s)[:s.Length()-1]

	return state
}

// Peek return the top state on the stack
func (s *StateStack) Peek() *State {

	state := (*s)[s.Length()-1]

	return state
}

// Contains returns true if it contains the state
func (s *StateStack) Contains(state *State) bool {
	// TODO: Optimise this
	for _, s1 := range *s {
		if s1 == state {
			return true
		}
	}
	return false
}
