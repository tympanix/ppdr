package gnba

import (
	"testing"

	"github.com/tympanix/master-2019/ltl"
)

func TestTransformGNBAToNBA(t *testing.T) {
	phi := ltl.And{ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}}, ltl.Until{ltl.AP{"b"}, ltl.AP{"a"}}}
	g := GenerateGNBA(phi)dd
	nba := TransformGNBAtoNBA(g)

	if 2*len(g.States) != len(nba.States) {
		t.Error("nba does not have the right amount of states")
	}

	if len(g.FinalStates[0]) != len(nba.FinalStates) {
		t.Error("nba has incorrect number of accepting states")
	}
}
