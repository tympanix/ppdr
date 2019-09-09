package gnba

import (
	"github.com/tympanix/master-2019/ltl"
)

// GenerateGNBA generates an GNBA from an LTL formula phi
func GenerateGNBA(phi ltl.Node) GNBA {

	closure := ltl.Closure(phi)
	aps := ltl.FindAtomicPropositions(phi)
	elemSets := ltl.FindElementarySets(closure)

	states := make([]*Node, 0, len(elemSets))

	for _, s := range elemSets {
		n := Node{
			ElementarySet: s,
			Transitions:   make([]Transition, 0),
		}

		states = append(states, &n)
	}

	for _, s := range states {
		s.IsStartState = s.ElementarySet.Contains(phi)
	}

	for _, s := range states {
		intersec := s.ElementarySet.Intersection(aps)

		for _, s2 := range states {
			if s.shouldHaveEdgeTo(*s2, closure) {
				s.addTransition(s2, intersec)
			}
		}
	}

	return states
}
