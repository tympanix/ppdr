package gnba

import (
	"github.com/tympanix/master-2019/ltl"
)

// Transition is a transition in a GNBA
type Transition struct {
	State *State
	Label ltl.Set
}

// Copy returns a copy of the transition with renaming of nodes
func (t Transition) Copy(rt renameTable) Transition {
	return Transition{
		State: rt[t.State],
		Label: t.Label.Copy(),
	}
}

// Rename renames the transition to point to a new state using a rename table
func (t Transition) Rename(rt renameTable) Transition {
	return Transition{
		State: rt[t.State],
		Label: t.Label,
	}
}

// RenameTo returns a new transitions where the destination state of the transition
// is changed to new one
func (t Transition) RenameTo(state *State) Transition {
	return Transition{
		State: state,
		Label: t.Label,
	}
}
