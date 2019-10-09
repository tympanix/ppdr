package repo

import (
	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/gnba"
	"github.com/tympanix/master-2019/systems/nba"
	"github.com/tympanix/master-2019/systems/product"
	"github.com/tympanix/master-2019/systems/ts"
)

type candidate struct {
	*State
}

func (c candidate) InitialStates() []ts.State {
	return []ts.State{c.State}
}

func (c candidate) satisfiesConfPolicies() bool {
	phi := c.State.confPolicies.Conjunction()
	return c.satisfiesFormula(phi)
}

func (c candidate) satisfiesFormula(phi ltl.Node) bool {
	if phi == nil {
		return true
	}

	phi2, table, err := ltl.Compile(phi)

	if err != nil {
		panic(err)
	}

	phi2 = phi2.Normalize()
	ap := ltl.FindAtomicPropositions(phi2)

	r := ltl.NewResolverFromSet(c.State.Predicates(ap, table))
	if s, err := ltl.Satisfied(phi2, r); err == nil {
		return s
	}

	n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(ltl.Negate(phi)))
	p := product.New(c, n)
	return p.HasAcceptingCycle() == nil
}
