package ltl

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/tympanix/ppdr/debug"
)

// Set is a set of formulas in LTL.
type Set map[Node]bool

// NewSet returns a new empty set
func NewSet(nodes ...Node) Set {
	set := make(Set)
	set.Add(nodes...)
	return set
}

// SetSlice is a list of nodes (for all elements in the set)
type SetSlice []Node

func (a SetSlice) Len() int {
	return len(a)
}

func (a SetSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SetSlice) Less(i, j int) bool {
	return strings.Compare(a[i].String(), a[j].String()) < 0
}

func (s Set) String() string {
	var sb strings.Builder

	sb.WriteString("[")

	i := 0

	elems := s.AsSlice()
	sort.Sort(elems)

	for _, e := range elems {
		fmt.Fprint(&sb, e)

		if i < s.Size()-1 {
			sb.WriteString(", ")
		}

		i++
	}

	sb.WriteString("]")

	return sb.String()
}

// Contains returns true if the elementary sets contains phi.
func (s Set) Contains(phi Node) bool {
	_, ok := s[phi]
	return ok
}

// Size returns the size of a set.
func (s Set) Size() int {
	return len(s)
}

// Add adds an element to a set.
func (s *Set) Add(node ...Node) Set {
	for _, n := range node {
		(*s)[n] = true
	}
	return *s
}

// AddSet adds all elements from the argument set to this one
func (s *Set) AddSet(set Set) {
	for e := range set {
		s.Add(e)
	}
}

// Copy returns a new set which is identical to the original one
func (s *Set) Copy() Set {
	return NewSet(s.AsSlice()...)
}

// ContainsAll return true if all elements are contained in the set.
func (s Set) ContainsAll(set Set) bool {
	for e := range set {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

// Conjunction returns a single formula which is a conjunction of all formulae
// in the set
func (s Set) Conjunction() Node {

	if s.Size() == 0 {
		return nil
	}

	if s.Size() == 1 {
		return s.AsSlice()[0]
	}

	l := s.AsSlice()
	conj := And{l[0], l[1]}

	for _, e := range l[2:] {
		conj = And{e, conj}
	}

	return conj
}

// ContainsAny return true if any element is contained in the set.
func (s Set) ContainsAny(set Set) bool {
	for e := range set {
		if s.Contains(e) {
			return true
		}
	}
	return false
}

// Equals returns true on set equality (i.e. same length and same elements)
func (s Set) Equals(set Set) bool {
	if s.Size() != set.Size() {
		return false
	}

	for e := range s {
		if !set.Contains(e) {
			return false
		}
	}

	return true
}

// Intersection find the intersection LTL nodes from another set.
func (s Set) Intersection(set Set) Set {
	res := make(Set, 0)
	for e := range s {
		for e2 := range set {
			if e.SameAs(e2) {
				res.Add(e)
			}
		}
	}
	return res
}

// AsSlice collects all values in the set as an array. Operation is expensive
func (s Set) AsSlice() SetSlice {
	slice := make([]Node, 0)
	for n := range s {
		slice = append(slice, n)
	}
	return slice
}

// PowerSet returns an array (set) containing all possible
// subsets for a set.
func (s Set) PowerSet() []Set {
	t := debug.NewTimer("powerset")

	powerSet := make([]Set, 0)
	powerSetSize := int(math.Pow(2, float64(s.Size())))
	values := s.AsSlice()

	for i := 0; i < powerSetSize; i++ {
		subset := make(Set, 0)
		for j := range values {
			if (i & (1 << uint(j))) > 0 {
				subset.Add(values[j])
			}
		}
		powerSet = append(powerSet, subset)
	}

	defer func() {
		t.Stop()
	}()

	return powerSet

}

// IsElementary returns true if the set is elementary
func (s Set) IsElementary(closure Set) bool {
	if !s.isMaximal(closure) {
		return false
	}

	if !s.isLocallyConsistent(closure) {
		return false
	}

	if !s.isConsistent(closure) {
		return false
	}

	return true
}

func (s Set) isConsistent(closure Set) bool {
	// Case 1
	for psi := range closure {
		if a, ok := psi.(And); ok {
			if s.Contains(a) != (s.Contains(a.LHSNode()) && s.Contains(a.RHSNode())) {
				return false
			}
		}
	}

	// Case 2
	for psi := range closure {
		if s.Contains(psi) {
			if s.Contains(Not{psi}) {
				return false
			}
		}
	}

	// Case 3
	if closure.Contains(True{}) {
		if !s.Contains(True{}) {
			return false
		}
	}

	return true
}

func (s Set) isLocallyConsistent(closure Set) bool {
	// Case 1
	for psi := range closure {
		if until, ok := psi.(Until); ok {
			if s.Contains(until.RHSNode()) {
				if !s.Contains(until) {
					return false
				}
			}
		}
	}

	// Case 2
	for psi := range closure {
		if until, ok := psi.(Until); ok {
			if s.Contains(until) && !s.Contains(until.RHSNode()) {
				if !s.Contains(until.LHSNode()) {
					return false
				}
			}
		}
	}

	return true
}

func (s Set) isMaximal(closure Set) bool {
	// Case 1
	for psi := range closure {
		if !s.Contains(psi) {
			if !s.Contains(Negate(psi)) {
				return false
			}
		}
	}

	return true
}
