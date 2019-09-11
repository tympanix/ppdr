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
