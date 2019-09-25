package gnba

import (
	"testing"

	"github.com/tympanix/master-2019/ltl"
)

func TestTransformGNBAToNBA(t *testing.T) {
	phi := ltl.And{ltl.Always{ltl.Eventually{ltl.AP{"crit1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"crit2"}}}}
	g := GenerateGNBA(phi)
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
	q0 := NewState(ltl.NewSet(ltl.AP{"q0"}))
	q1 := NewState(ltl.NewSet(ltl.AP{"q1"}))
	q2 := NewState(ltl.NewSet(ltl.AP{"q2"}))

	// Add transitions
	q0.addTransition(q1, ltl.NewSet(ltl.AP{"crit1"}))
	q0.addTransition(q2, ltl.NewSet(ltl.AP{"crit2"}))
	q0.addTransition(q0, ltl.NewSet(ltl.True{}))
	q1.addTransition(q0, ltl.NewSet(ltl.True{}))
	q2.addTransition(q0, ltl.NewSet(ltl.True{}))

	gnba := &GNBA{
		States:         []*State{q0, q1, q2},
		StartingStates: NewStateSet(q0),
		FinalStates:    []StateSet{NewStateSet(q1), NewStateSet(q2)},
	}

	nba := TransformGNBAtoNBA(gnba)

	if len(nba.States) != 2*len(gnba.States) {
		t.Error("nba does not have right amount of states")
	}
}
