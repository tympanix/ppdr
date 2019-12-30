package repo

import (
	"unsafe"

	"github.com/tympanix/ppdr/ltl"
	"github.com/tympanix/ppdr/systems/gnba"
	"github.com/tympanix/ppdr/systems/nba"
	"github.com/tympanix/ppdr/systems/product"
	"github.com/tympanix/ppdr/systems/ts"
)

type candidate struct {
	*State
	*Repo
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

	// Rename subject with current user
	phi = phi.Map(func(n ltl.Node) ltl.Node {
		if _, ok := n.(ltl.Subject); ok {
			return ltl.Ptr{
				Attr:    "subject",
				Pointer: unsafe.Pointer(c.Repo.currentUser),
			}
		}
		return n
	})

	phi, table, err := ltl.Compile(phi)

	if err != nil {
		panic(err)
	}

	phi = phi.Normalize()
	ap := ltl.FindAtomicPropositions(phi)

	r1 := ltl.NewResolverFromMap(c.State.attributes)
	r2 := ltl.NewResolverFromSet(c.State.Predicates(ap, table))
	r := ltl.NewResolverCombined(r1, r2)

	if s, err := ltl.Satisfied(phi, r); err == nil {
		return s
	}

	n := nba.TransformGNBAtoNBA(gnba.GenerateGNBA(ltl.Negate(phi)))
	p := product.New(c, n, table)
	return p.HasAcceptingCycle() == nil
}
