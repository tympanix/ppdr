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

	if s, err := ltl.Satisfied(phi, c.State.Predicates()); err == nil {
		return s
	}

	n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(ltl.Negate(phi)))
	p := product.New(c, n)
	return p.HasAcceptingCycle() == nil
}
