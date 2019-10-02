package ts

import (
	"github.com/tympanix/master-2019/ltl"
)

// TS is an interface for the transition system itself
type TS interface {
	States() []State
	InitialStates() []State
}

// State is a state in a transition system
type State interface {
	Predicates() ltl.Set
	Dependencies() []State
}
