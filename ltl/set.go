package ltl

// Set is a set of formulas in LTL
type Set []Node

// Contains returns true if the elementary sets contains phi
func (s Set) Contains(phi Node) bool {
	for _, f := range s {
		if f.SameAs(phi) {
			return true
		}
	}
	return false
}

// ContainsAll return true if all elements are contained in the set
func (s Set) ContainsAll(set Set) bool {
	for _, e := range set {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

// Intersection find the intersection LTL nodes from another set
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

// IsElementary returns true if the set is elementary
func (s Set) IsElementary() bool {
	// TODO: Implement logic for elemetary sets
	return false
}
