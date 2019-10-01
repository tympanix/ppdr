package product

import (
	"fmt"
	"strings"
)

// StateStack is a stack of Product states
type StateStack struct {
	stack []*State
	set   map[*State]int
}

// NewStateStack returns a new (possibly empty) stack of states
func NewStateStack(states ...*State) *StateStack {
	stack := &StateStack{
		stack: make([]*State, 0),
		set:   make(map[*State]int),
	}
	stack.Push(states...)
	return stack
}

func (s *StateStack) String() string {
	var sb strings.Builder
	for _, v := range s.stack {
		fmt.Fprintf(&sb, "%v\n", v)
	}
	return sb.String()
}

func (s *StateStack) push(state *State) {
	s.stack = append(s.stack, state)
	s.set[state] = s.set[state] + 1
}

// Length returns the number of elements in the stack
func (s *StateStack) Length() int {
	return len(s.stack)
}

// Empty returns true if the stack is empty
func (s *StateStack) Empty() bool {
	return s.Length() == 0
}

// Push adds one or more states to the stack
func (s *StateStack) Push(states ...*State) {
	for _, state := range states {
		s.push(state)
	}
}

// Contains returns true if state is contained in the stack
func (s *StateStack) Contains(state *State) bool {
	_, ok := s.set[state]
	return ok
}

// Pop removes and returns the top state on the stack
func (s *StateStack) Pop() *State {

	state := s.Peek()
	s.stack = s.stack[:s.Length()-1]
	s.set[state] = s.set[state] - 1

	if s.set[state] <= 0 {
		delete(s.set, state)
	}

	return state
}

// Peek return the top state on the stack
func (s *StateStack) Peek() *State {
	return s.stack[s.Length()-1]
}
