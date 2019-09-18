package product

import (
	"github.com/tympanix/master-2019/gnba"
	"github.com/tympanix/master-2019/systems/ts"
)

// StateTuple is a tuple of states from the product transition system of TS and A
type StateTuple struct {
	StateTS  *ts.State
	StateNBA *gnba.State
}

// State is a state in the product transition system of TS and A
type State struct {
	StateTuple
	Transitions StateSet
	IsExpanded  bool
}

func newState(sTS *ts.State, sNBA *gnba.State) *State {
	return &State{
		StateTuple: StateTuple{
			StateTS:  sTS,
			StateNBA: sNBA,
		},
		Transitions: NewStateSet(),
		IsExpanded:  false,
	}
}

func (s *State) post(p *Product) StateSet {
	if !s.IsExpanded {
		return s.expand(p)
	}

	return s.Transitions
}

func (s *State) addTransition(s1 *State) {
	s.Transitions.Add(s1)
}

func (s *State) expand(p *Product) StateSet {
	for _, tTS := range s.StateTS.Dependencies {
		for _, tNBA := range s.StateNBA.Transitions {
			pNBA := tNBA.State
			if !tNBA.Label.Conflicts(tTS.Predicates) {
				sPrime := p.getOrAddState(tTS, pNBA)
				s.addTransition(sPrime)
			}
		}
	}
	s.IsExpanded = true
	return s.Transitions
}

func (s *State) unvisitedSucc(set StateSet, p *Product) *State {
	for s1 := range s.post(p) {
		if !set.Contains(s1) {
			return s1
		}
	}
	return nil
}

func (s *State) isFinalState(p *Product) bool {
	return p.NBA.IsAcceptanceState(s.StateNBA)
}
