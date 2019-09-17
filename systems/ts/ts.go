package ts

import (
	"github.com/tympanix/master-2019/ltl"
)

// State is a node in the transition system
type State struct {
	Predicates   ltl.Set
	Dependencies []*State
}

// TS is a transition system
type TS struct {
	States        []*State
	InitialStates []*State
}
