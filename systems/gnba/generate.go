package gnba

import (
	"github.com/tympanix/ppdr/debug"
	"github.com/tympanix/ppdr/ltl"
	"github.com/tympanix/ppdr/systems/ba"
)

// GenerateGNBA generates an GNBA from an LTL formula phi
func GenerateGNBA(phi ltl.Node) *GNBA {

	t := debug.NewTimer("gnba")

	defer func() {
		t.Stop()
	}()

	phi = phi.Normalize()
	closure := ltl.Closure(phi)
	elemSets := ltl.FindElementarySets(phi)

	gnba := NewGNBA(phi)

	for _, s := range elemSets {
		n := ba.State{
			ElementarySet: s,
			Transitions:   make([]ba.Transition, 0),
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
			set := ba.NewStateSet()
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
		intersec := s.ElementarySet.Intersection(gnba.AP)

		for _, s2 := range gnba.States {
			if s.ShouldHaveEdgeTo(*s2, closure) {
				s.AddTransitionFromSet(s2, intersec)
			}
		}
	}

	return gnba
}
