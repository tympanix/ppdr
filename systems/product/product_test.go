package product

import (
	"testing"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ba"
	"github.com/tympanix/master-2019/systems/nba"
	"github.com/tympanix/master-2019/systems/ts"
)

func TestExample_4_64(t *testing.T) {
	tr, n := generateExample4_22()
	p := New(tr, n)
	if !p.Satisfy() {
		t.Error("Expected p to satisfy.")
	}
}

func generateExample4_22() (*ts.TS, *nba.NBA) {
	t := ts.New()
	n := nba.NewNBA()

	s1 := ts.NewState(ltl.AP{"red"})
	s2 := ts.NewState(ltl.AP{"green"})
	s1.AddDependency(s2)
	s2.AddDependency(s1)
	t.AddState(s1, s2)
	t.AddInitialState(s1)

	q0 := ba.NewState(ltl.NewSet(ltl.AP{"q0"}))
	q1 := ba.NewState(ltl.NewSet(ltl.AP{"q1"}))
	q2 := ba.NewState(ltl.NewSet(ltl.AP{"q2"}))
	q0.AddTransition(q0, ltl.True{})
	q0.AddTransition(q1, ltl.Not{ltl.AP{"green"}})
	q1.AddTransition(q1, ltl.Not{ltl.AP{"green"}})
	q1.AddTransition(q2, ltl.AP{"green"})
	q2.AddTransition(q2, ltl.True{})
	n.AddState(q0, q1, q2)
	n.AddInitialState(q0)
	n.AddAcceptanceState(q1)

	return t, n
}
