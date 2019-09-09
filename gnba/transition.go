package gnba

import (
	"github.com/tympanix/master-2019/ltl"
)

// Transition is a transition in a GNBA
type Transition struct {
	State *State
	Label ltl.Set
}
