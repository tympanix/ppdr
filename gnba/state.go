package gnba

import "github.com/tympanix/master-2019/ltl"

// State is a node in a GNBA
type State struct {
	ElementarySet ltl.Set
	Transitions   []Transition
}

// NewState returns a new state with no transitions with the given elementary set
func NewState(es ltl.Set) *State {
	return &State{
		ElementarySet: es,
		Transitions:   make([]Transition, 0),
	}
}

// Has returns true if the gnba node has formula psi in the elementary set
func (n *State) Has(psi ltl.Node) bool {
	return n.ElementarySet.Contains(psi)
}

// Copy makes a copy of the state using a rename tabe to rename transitions
func (n *State) Copy() *State {
	trns := make([]Transition, 0)

	for _, t := range n.Transitions {
		trns = append(trns, t)
	}

	return &State{
		ElementarySet: n.ElementarySet.Copy(),
		Transitions:   trns,
	}
}

// Rename renames all state transitions using a rename table
func (n *State) Rename(rt renameTable) {
	for i, t := range n.Transitions {
		n.Transitions[i] = t.Rename(rt)
	}
}

func (n *State) addTransition(node *State, label ltl.Set) {
	n.Transitions = append(n.Transitions, Transition{
		State: node,
		Label: label,
	})
}

func (n *State) shouldHaveEdgeTo(node State, closure ltl.Set) bool {
	// case 1
	// n = B, node = B'
	for psi := range closure {
		if next, ok := psi.(ltl.Next); ok {
			if n.Has(next) != node.Has(next.ChildNode()) {
				return false
			}
		}
	}

	// case 2
	for psi := range closure {
		if until, ok := psi.(ltl.Until); ok {
			if n.Has(until) != (n.Has(until.RHSNode()) || (n.Has(until.LHSNode()) && node.Has(until))) {
				return false
			}
		}
	}

	return true
}
