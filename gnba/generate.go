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
	elemSets := generateSets(phi)
	aps := ltl.FindAtomicPropositions(phi)

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

func generateSets(phi ltl.Node) []ltl.Set {

	sets := make([]ltl.Set, 0)

	b1 := ltl.Set{
		ltl.AP{"a"},
		ltl.Next{ltl.AP{"a"}},
	}

	b2 := ltl.Set{
		ltl.AP{"a"},
		ltl.Not{ltl.Next{ltl.AP{"a"}}},
	}

	b3 := ltl.Set{
		ltl.Not{ltl.AP{"a"}},
		ltl.Next{ltl.AP{"a"}},
	}

	b4 := ltl.Set{
		ltl.Not{ltl.AP{"a"}},
		ltl.Not{ltl.Next{ltl.AP{"a"}}},
	}

	return append(sets, b1, b2, b3, b4)
}
