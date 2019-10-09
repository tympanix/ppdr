package repo

import (
	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ts"
)

// Attrs is a map containing attributes for the state
type Attrs map[string]interface{}

// State is a data node in the repo
type State struct {
	dependencies []ts.State
	attributes   map[string]interface{}
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
			a[key] = vals[i+1]
			i += 2
		} else if ap, ok := vals[i].(ltl.AP); ok {
			a[ap.Name] = true
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

// Predicates returns the set of predicates which hold in the state
func (s *State) Predicates(ap ltl.Set, t ltl.RefTable) ltl.Set {
	preds := ltl.NewSet()
	for k := range ap {
		if a, ok := k.(ltl.AP); ok {
			if v, ok := s.attributes[a.Name]; ok {
				if b, ok := v.(bool); ok && b {
					preds.Add(k)
				}
			}
		}
	}
	return preds
}

// Dependencies return a list of dependencies from this state
func (s *State) Dependencies() []ts.State {
	return s.dependencies
}

func (s *State) addConfPolicy(set ltl.Set) {
	s.confPolicies.AddSet(set)
}
