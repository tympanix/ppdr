package product

import (
	"fmt"

	"github.com/tympanix/master-2019/systems/ba"
	"github.com/tympanix/master-2019/systems/ts"
)

// StateTuple is a tuple of states from the product transition system of TS and A
type StateTuple struct {
	StateTS  ts.State
	StateNBA *ba.State
}

func (s *StateTuple) String() string {
	return fmt.Sprintf("<%v, %v>", s.StateTS, s.StateNBA)
}

// State is a state in the product transition system of TS and A
type State struct {
	StateTuple
	Transitions StateSet
	IsExpanded  bool
}

// StringWithRenamer strings the state using a rename function
func (s *State) StringWithRenamer(r ba.StateNamer) string {
	return fmt.Sprintf("<%v, %v>", s.StateTS, r.RenameState(s.StateNBA))
}

func (s *State) String() string {
	return s.StringWithRenamer(nil)
}

func newState(sTS ts.State, sNBA *ba.State) *State {
	return &State{
		StateTuple: StateTuple{
			StateTS:  sTS,
			StateNBA: sNBA,
		},
		Transitions: NewStateSet(),
		IsExpanded:  false,
	}
}

func (s *State) successors(p *Product) StateSet {
	if !s.IsExpanded {
		return s.expand(p)
	}

	return s.Transitions
}

func (s *State) addTransition(s1 *State) {
	s.Transitions.Add(s1)
}

func (s *State) expand(p *Product) StateSet {
	for _, tTS := range s.StateTS.Dependencies() {
		lf := p.NBA.AP.Intersection(tTS.Predicates())
		for _, tNBA := range s.StateNBA.Transitions {
			pNBA := tNBA.State
			if tNBA.Label.Equals(lf) {
				sPrime := p.getOrAddState(tTS, pNBA)
				s.addTransition(sPrime)
			}
		}
	}
	s.IsExpanded = true
	return s.Transitions
}
