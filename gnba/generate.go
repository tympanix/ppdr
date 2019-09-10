package gnba

import (
	"github.com/tympanix/master-2019/ltl"
)

// GenerateGNBA generates an GNBA from an LTL formula phi
func GenerateGNBA(phi ltl.Node) *GNBA {

	closure := ltl.Closure(phi)
	aps := ltl.FindAtomicPropositions(phi)
	elemSets := ltl.FindElementarySets(closure)

	gnba := NewGNBA()

	for _, s := range elemSets {
		n := State{
			ElementarySet: s,
			Transitions:   make([]Transition, 0),
		}

		gnba.States = append(gnba.States, &n)
	}

	// Find starting states
	for _, s := range gnba.States {
		if s.Has(phi) {
			gnba.StartingStates.Add(s)
		}
	}

	// Find acceptance sets
	for psi := range closure {
		if until, ok := psi.(ltl.Until); ok {
			set := NewStateSet()
			for _, s := range gnba.States {
				if !s.Has(until) || s.Has(until.RHSNode()) {
					set.Add(s)
				}
			}
			gnba.FinalStates = append(gnba.FinalStates, set)
		}
	}

	// Find transitions relations
	for _, s := range gnba.States {
		intersec := s.ElementarySet.Intersection(aps)

		for _, s2 := range gnba.States {
			if s.shouldHaveEdgeTo(*s2, closure) {
				s.addTransition(s2, intersec)
			}
		}
	}

	return gnba
}
