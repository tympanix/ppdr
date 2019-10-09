package ts

import (
	"github.com/tympanix/master-2019/ltl"
)

// TS is an interface for the transition system itself
type TS interface {
	InitialStates() []State
}

// State is a state in a transition system
type State interface {
	Predicates(ap ltl.Set, t ltl.RefTable) ltl.Set
	Dependencies() []State
}
