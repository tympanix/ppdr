package ba

import (
	"fmt"

	"github.com/tympanix/ppdr/ltl"
)

// StateNamer gives human readable names to states
type StateNamer func(s *State) string

// RenameState returns the name of the state using the rename function
func (n StateNamer) RenameState(s *State) string {
	if n != nil {
		return n(s)
	}
	return fmt.Sprintf("%v", s.ElementarySet)
}

// NewStateNamerFromMap return a new state namer backed by a lookup table
func NewStateNamerFromMap(m map[string]ltl.Set) StateNamer {

	return StateNamer(func(s *State) string {

		for n, e := range m {
			if s.ElementarySet.ContainsAll(e) && s.ElementarySet.Size() == e.Size() {
				return n
			}
		}

		panic(fmt.Sprintf("no name for state %v", s))
	})
}
