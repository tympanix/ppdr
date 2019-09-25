package nba

import (
	"testing"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ba"
	"github.com/tympanix/master-2019/systems/gnba"
)

func TestTransformGNBAToNBA(t *testing.T) {
	phi := ltl.And{ltl.Always{ltl.Eventually{ltl.AP{"crit1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"crit2"}}}}
	g := gnba.GenerateGNBA(phi)
	nba := TransformGNBAtoNBA(g)

	if 4*len(g.States) != len(nba.States) {
		t.Error("nba does not have the right amount of states")
	}

	if len(g.FinalStates[0]) != len(nba.FinalStates) {
		t.Error("nba has incorrect number of accepting states")
	}
}

func TestTransformExample4_57(t *testing.T) {
	// Add states
	q0 := ba.NewState(ltl.NewSet(ltl.AP{"q0"}))
	q1 := ba.NewState(ltl.NewSet(ltl.AP{"q1"}))
	q2 := ba.NewState(ltl.NewSet(ltl.AP{"q2"}))

	// Add transitions
	q0.AddTransition(q2, ltl.AP{"crit2"})
	q0.AddTransition(q1, ltl.AP{"crit1"})
	q0.AddTransition(q0, ltl.True{})
	q1.AddTransition(q0, ltl.True{})
	q2.AddTransition(q0, ltl.True{})

	gnba := &gnba.GNBA{
		States:         []*ba.State{q0, q1, q2},
		StartingStates: ba.NewStateSet(q0),
		FinalStates:    []ba.StateSet{ba.NewStateSet(q1), ba.NewStateSet(q2)},
	}

	nba := TransformGNBAtoNBA(gnba)

	if len(nba.States) != 2*len(gnba.States) {
		t.Error("nba does not have right amount of states")
	}
}
