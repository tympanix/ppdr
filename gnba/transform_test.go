package gnba

import (
	"testing"

	"github.com/tympanix/master-2019/ltl"
)

func TestTransformGNBAToNBA(t *testing.T) {
	phi := ltl.Or{ltl.Always{ltl.Eventually{ltl.AP{"crit1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"crit2"}}}}
	t.Fatal(len(ltl.Closure(phi)))
	g := GenerateGNBA(phi)
	t.Errorf("\n%v", g)
}
