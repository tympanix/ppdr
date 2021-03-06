package repo

import (
	"fmt"
	"unsafe"

	"github.com/tympanix/ppdr/ltl"
	"github.com/tympanix/ppdr/systems/ts"
)

// Attrs is a map containing attributes for the state
type Attrs map[string]ltl.Node

// State is a data node in the repo
type State struct {
	dependencies []ts.State
	attributes   Attrs
	confPolicies ltl.Set
}

// NewState returns a new empty state
func NewState(vals ...interface{}) *State {
	attr := parseAttributes(vals)

	return &State{
		dependencies: make([]ts.State, 0),
		attributes:   attr,
		confPolicies: ltl.NewSet(),
	}
}

func parseAttributes(vals []interface{}) Attrs {
	a := make(Attrs)
	i := 0

	for i < len(vals) {
		if key, ok := vals[i].(string); ok {
			a[key] = ltl.ValueToLiteral(vals[i+1])
			i += 2
		} else if ap, ok := vals[i].(ltl.AP); ok {
			a[ap.Name] = ltl.ValueToLiteral(true)
			i++
		} else {
			panic("could not parse attributes")
		}
	}

	return a
}

func (s *State) addDependency(state *State) {
	s.dependencies = append(s.dependencies, state)
}

func (s *State) resolver() ltl.Resolver {
	return ltl.NewResolverFromMap(s.attributes)
}

// Predicates returns the set of predicates which hold in the state
func (s *State) Predicates(ap ltl.Set, t ltl.RefTable) ltl.Set {
	preds := ltl.NewSet()
	for k := range ap {
		if a, ok := k.(ltl.AP); ok {
			// If node is an AP, return that AP if it is contained in the
			// attributes map as LitBool value
			if v, ok := s.attributes[a.Name]; ok {
				if b, ok := v.(ltl.LitBool); ok && b.Bool {
					preds.Add(k)
				}
			}
		} else if r, ok := k.(ltl.Ref); ok {
			// Else if AP is a reference, look up that reference in the
			// reference table and evaluate the expression
			exp, ok := t[r]

			if !ok {
				panic(fmt.Sprintf("unknown reference %v", r))
			}

			b, err := ltl.Satisfied(exp, s.resolver())

			if b && (err == nil) {
				preds.Add(r)
			}
		} else {
			panic(fmt.Sprintf("ltl node: %v, can not be evaluated as an atomic proposition", k))
		}
	}
	return preds
}

func (s *State) replaceSelfReferences() {
	set := ltl.NewSet()

	s.newAttrPtr("self", unsafe.Pointer(s))

	for p := range s.confPolicies {
		p1 := ltl.RenameSelfPredicate(p, s.getSelfAttr())
		set.Add(p1)
	}

	s.confPolicies = set
}

func (s *State) replaceUserPredicate(r *Repo) {
	set := ltl.NewSet()

	for p := range s.confPolicies {
		p1 := p.Map(r.replaceUserPredicate)
		set.Add(p1)
	}

	s.confPolicies = set
}

func (s *State) getSelfAttr() ltl.Node {
	return s.attributes["self"]
}

func (s *State) newAttrPtr(attr string, ptr unsafe.Pointer) ltl.Ptr {
	r := ltl.Ptr{
		Attr:    attr,
		Pointer: ptr,
	}
	s.attributes[r.Attr] = r
	return r
}

// Dependencies return a list of dependencies from this state
func (s *State) Dependencies() []ts.State {
	return s.dependencies
}

func (s *State) addConfPolicy(set ltl.Set) {
	s.confPolicies.AddSet(set)
}

func (s *State) AddPolicy(n ltl.Node) {
	s.confPolicies.Add(n)
}
