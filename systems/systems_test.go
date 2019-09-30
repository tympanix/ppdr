package systems

import (
	"fmt"
	"testing"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/gnba"
	"github.com/tympanix/master-2019/systems/nba"
	"github.com/tympanix/master-2019/systems/product"
	"github.com/tympanix/master-2019/systems/ts"
)

func TestCycleDetection(t *testing.T) {
	tr := ts.New()

	s1 := ts.NewState(ltl.AP{"red"})
	s2 := ts.NewState(ltl.AP{"green"})
	s1.AddDependency(s2)
	s2.AddDependency(s1)
	tr.AddState(s1, s2)
	tr.AddInitialState(s1)

	tests := map[ltl.Node]bool{
		ltl.Eventually{ltl.AP{"green"}}:                                         false,
		ltl.Eventually{ltl.AP{"yellow"}}:                                        true,
		ltl.Always{ltl.AP{"green"}}:                                             true,
		ltl.Always{ltl.Eventually{ltl.AP{"green"}}}:                             false,
		ltl.Next{ltl.AP{"green"}}:                                               false,
		ltl.Next{ltl.AP{"red"}}:                                                 true,
		ltl.Next{ltl.Next{ltl.AP{"red"}}}:                                       false,
		ltl.Next{ltl.Next{ltl.AP{"green"}}}:                                     true,
		ltl.Impl{ltl.AP{"red"}, ltl.Next{ltl.AP{"green"}}}:                      false,
		ltl.Impl{ltl.AP{"red"}, ltl.Next{ltl.AP{"red"}}}:                        true,
		ltl.And{ltl.Eventually{ltl.AP{"red"}}, ltl.Eventually{ltl.AP{"green"}}}: false,
		ltl.And{ltl.Always{ltl.AP{"red"}}, ltl.Always{ltl.AP{"green"}}}:         true,
		ltl.Always{ltl.Or{ltl.AP{"red"}, ltl.AP{"green"}}}:                      false,
	}

	i := 0

	for phi, cycle := range tests {
		name := fmt.Sprintf("test:%v", i)
		t.Run(name, func(t *testing.T) {
			n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(ltl.Negate(phi)))
			p := product.New(tr, n)

			if p.HasAcceptingCycle() != cycle {
				t.Errorf("expected cycle: %v, formula: %v", cycle, phi)
			}
		})
		i++
	}
}
