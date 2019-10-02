package systems

import (
	"fmt"
	"testing"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ba"
	"github.com/tympanix/master-2019/systems/gnba"
	"github.com/tympanix/master-2019/systems/nba"
	"github.com/tympanix/master-2019/systems/product"
	"github.com/tympanix/master-2019/systems/ts"
)

// An example of a simple traffix light which transitions between red and green
// light. It is checked whether certain properties are held for the traffic
// light (i.e. always eventually green)
func TestCycleDetection_one(t *testing.T) {
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
		ltl.Eventually{ltl.Always{ltl.Not{ltl.AP{"green"}}}}:                    true,
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

	checkFormulas(t, tr, tests, nil)

}

// formula: true U (c1 and c2)
func TestCycleDetection_two(t *testing.T) {
	tr := generateFigure_5_5()

	tests := map[ltl.Node]bool{
		ltl.Eventually{ltl.And{ltl.AP{"c1"}, ltl.AP{"c2"}}}:                                          true,
		ltl.Always{ltl.Or{ltl.Not{ltl.AP{"c1"}}, ltl.Not{ltl.AP{"c2"}}}}:                             false,
		ltl.Or{ltl.Always{ltl.Eventually{ltl.AP{"c1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"c2"}}}}:   false,
		ltl.And{ltl.Always{ltl.Eventually{ltl.AP{"c1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"c2"}}}}:  true,
		ltl.Impl{ltl.Always{ltl.Eventually{ltl.AP{"w1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"c1"}}}}: true,
		ltl.Impl{ltl.Always{ltl.Eventually{ltl.AP{"w2"}}}, ltl.Always{ltl.Eventually{ltl.AP{"c2"}}}}: true,
	}

	checkFormulas(t, tr, tests, nil)

}

func Example5_10() {
	c1 := ltl.AP{"c1"}
	c2 := ltl.AP{"c2"}
	and := ltl.And{c1, c2}

	phi := ltl.Until{ltl.True{}, and}

	names := map[string]ltl.Set{
		"B1": ltl.NewSet(ltl.Negate(phi), ltl.Not{and}, ltl.True{}, ltl.Not{c1}, c2),
		"B2": ltl.NewSet(ltl.Negate(phi), ltl.Not{and}, ltl.True{}, c1, ltl.Not{c2}),
		"B3": ltl.NewSet(ltl.Negate(phi), ltl.Not{and}, ltl.True{}, ltl.Not{c1}, ltl.Not{c2}),
		"B4": ltl.NewSet(phi, and, ltl.True{}, c1, c2),
		"B5": ltl.NewSet(phi, ltl.Not{and}, ltl.True{}, ltl.Not{c1}, c2),
		"B6": ltl.NewSet(phi, ltl.Not{and}, ltl.True{}, c1, ltl.Not{c2}),
		"B7": ltl.NewSet(phi, ltl.Not{and}, ltl.True{}, ltl.Not{c1}, ltl.Not{c2}),
	}
	r := ba.NewStateNamerFromMap(names)

	n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(phi))

	fmt.Println(n.StringWithRenamer(r))

	// Output:
	// >B4*
	// 	[c1, c2]	-->	B4
	// 	[c1, c2]	-->	B5
	// 	[c1, c2]	-->	B6
	// 	[c1, c2]	-->	B7
	// 	[c1, c2]	-->	B1
	// 	[c1, c2]	-->	B2
	// 	[c1, c2]	-->	B3
	// >B5
	// 	[c2]	-->	B4
	// 	[c2]	-->	B5
	// 	[c2]	-->	B6
	// 	[c2]	-->	B7
	// >B6
	// 	[c1]	-->	B4
	// 	[c1]	-->	B5
	// 	[c1]	-->	B6
	// 	[c1]	-->	B7
	// >B7
	// 	[]	-->	B4
	// 	[]	-->	B5
	// 	[]	-->	B6
	// 	[]	-->	B7
	// B1*
	// 	[c2]	-->	B1
	// 	[c2]	-->	B2
	// 	[c2]	-->	B3
	// B2*
	// 	[c1]	-->	B1
	// 	[c1]	-->	B2
	// 	[c1]	-->	B3
	// B3*
	// 	[]	-->	B1
	// 	[]	-->	B2
	// 	[]	-->	B3
}

// Figure 5.5 is a transition system of two processes and a critical section
// modelled by a mutex. Transition system model mutual exclusion.
func generateFigure_5_5() *ts.TS {
	t := ts.New()

	n1 := ltl.AP{"n1"}
	n2 := ltl.AP{"n2"}
	w1 := ltl.AP{"w1"}
	w2 := ltl.AP{"w2"}
	c1 := ltl.AP{"c1"}
	c2 := ltl.AP{"c2"}
	y := ltl.AP{"y"}

	s0 := ts.NewState(n1, n2, y)
	s1 := ts.NewState(w1, n2, y)
	s2 := ts.NewState(n1, w2, y)
	s3 := ts.NewState(c1, n2)
	s4 := ts.NewState(w1, w2)
	s5 := ts.NewState(n1, c2)
	s6 := ts.NewState(c1, w2)
	s7 := ts.NewState(w1, c2)

	t.AddInitialState(s0)

	s0.AddDependency(s1)
	s0.AddDependency(s2)
	s1.AddDependency(s3)
	s1.AddDependency(s4)
	s2.AddDependency(s4)
	s2.AddDependency(s5)
	s3.AddDependency(s0)
	s3.AddDependency(s6)
	s4.AddDependency(s6)
	s4.AddDependency(s7)
	s5.AddDependency(s0)
	s5.AddDependency(s7)
	s6.AddDependency(s2)
	s7.AddDependency(s1)

	return t
}

func checkFormulas(t *testing.T, tr *ts.TS, tests map[ltl.Node]bool, r ba.StateNamer) {
	i := 0

	for phi, cycle := range tests {
		name := fmt.Sprintf("test:%v", i)
		t.Run(name, func(t *testing.T) {
			n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(ltl.Negate(phi)))
			p := product.New(tr, n)

			context := p.HasAcceptingCycle()

			if (context != nil) != cycle {
				t.Errorf("expected cycle: %v, formula: %v", cycle, phi)
				t.Errorf("size of product: %v", len(p.States))

				if context != nil {
					fmt.Println(context.TraceWithRenamer(r))
				}
			}

		})
		i++
	}
}
