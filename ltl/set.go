package ltl

import (
	"math"
)

// Set is a set of formulas in LTL.
type Set []Node

// Contains returns true if the elementary sets contains phi.
func (s Set) Contains(phi Node) bool {
	for _, f := range s {
		if f.SameAs(phi) {
			return true
		}
	}
	return false
}

// Size returns the size of a set.
func (s Set) Size() int {
	return len(s)
}

// Add adds an element to a set.
func (s Set) Add(node Node) Set {
	return append(s, node)
}

// ContainsAll return true if all elements are contained in the set.
func (s Set) ContainsAll(set Set) bool {
	for _, e := range set {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

// Intersection find the intersection LTL nodes from another set.
func (s Set) Intersection(set Set) Set {
	res := make([]Node, 0)
	for _, e := range s {
		for _, e2 := range set {
			if e.SameAs(e2) {
				res = append(res, e)
			}
		}
	}
	return res
}

// PowerSet returns an array (set) containing all possible
// subsets for a set.
func (s Set) PowerSet() []Set {
	powerSet := make([]Set, 0)

	powerSetSize := int(math.Pow(2, float64(s.Size())))

	for i := 0; i < powerSetSize; i++ {
		subset := make(Set, 0)
		for j := range s {
			if (i & (1 << j)) > 0 {
				subset = subset.Add(s[j])
			}
		}
		powerSet = append(powerSet, subset)
	}

	return powerSet

}

// IsElementary returns true if the set is elementary.
func (s Set) IsElementary() bool {
	return s.isConsistent() && s.isLocallyConsistent() && s.isMaximal()
}

func (s Set) isConsistent() bool {
	return false
}

func (s Set) isLocallyConsistent() bool {
	return false
}

func (s Set) isMaximal() bool {
	return false
}
