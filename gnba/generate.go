package gnba

import (
	"github.com/tympanix/master-2019/ltl"
)

// GenerateGNBA generates an GNBA from an LTL formula phi
func GenerateGNBA(phi ltl.Node) []*Node {

	// TODO:
	// - Find closure of phi (all subformulas of phi)
	// - Find AP (all tomic propositions of phi)
	// - Find elementary set from closure of phi
	closure := ltl.Closure(phi)
	aps := ltl.FindAtomicPropositions(phi)
	elemSets := ltl.FindElementarySets(closure)

	states := make([]*Node, len(elemSets))

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
			if s.shouldHaveEdgeTo(s2) {
				// TODO: add an edge here!
			}
		}
	}

	return nil
}
