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

func Test(t *testing.T) {
	tr := ts.New()

	s1 := ts.NewState(ltl.AP{"red"})
	s2 := ts.NewState(ltl.AP{"green"})
	s1.AddDependency(s2)
	s2.AddDependency(s1)
	tr.AddState(s1, s2)
	tr.AddInitialState(s1)

	phi := ltl.Eventually{ltl.AP{"green"}}

	n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(ltl.Negate(phi)))

	p := product.New(tr, n)

	t.Error(p.HasAcceptingCycle())

	t.Error(p)
}

func ExampleASD() {
	tr := ts.New()

	s1 := ts.NewState(ltl.AP{"red"})
	s2 := ts.NewState(ltl.AP{"green"})
	s1.AddDependency(s2)
	s2.AddDependency(s1)
	tr.AddState(s1, s2)
	tr.AddInitialState(s1)

	phi := ltl.Negate(ltl.Eventually{ltl.AP{"green"}}).Normalize()

	n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(phi))

	//fmt.Println(ltl.FindElementarySets(ltl.Closure(phi)))
	//fmt.Println(gnba.GenerateGNBA(phi))

	//fmt.Println(n)

	p := product.New(tr, n)
	fmt.Println(p.HasAcceptingCycle())
	fmt.Println(p)

	// Output:
	// false
}
