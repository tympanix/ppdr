package product

import (
	"github.com/tympanix/master-2019/gnba"
	"github.com/tympanix/master-2019/systems/ts"
)

// State is a state in the product transition system of TS and A
type State struct {
	StateTS     *ts.State
	StateNBA    *gnba.State
	Transitions []*State
	IsExpanded  bool
}

func newState(sTS *ts.State, sNBA *gnba.State) *State {
	return &State{
		StateTS:     sTS,
		StateNBA:    sNBA,
		Transitions: make([]*State, 0),
		IsExpanded:  false,
	}
}

func (s *State) post(p *Product) StateSet {
	if !s.IsExpanded {
		s.expand(p)
	}
	return NewStateSet()
}

func (s *State) addTransition(s1 *State) {
	s.Transitions = append(s.Transitions, s1)
}

func (s *State) expand(p *Product) {
	for _, tTS := range s.StateTS.Dependencies {
		for _, tNBA := range s.StateNBA.Transitions {
			pNBA := tNBA.State
			if !tNBA.Label.Conflicts(tTS.Predicates) {
				sPrime := newState(tTS, pNBA)
				p.AddState(sPrime)
				s.addTransition(sPrime)
			}
		}

	}

	s.IsExpanded = true
}

// Product is a product of a transition system and a NBA
type Product struct {
	States        []*State
	InitialStates []*State
	TS            *ts.TS
	NBA           *gnba.NBA
}

// AddState to the states of the product.
func (p *Product) AddState(s *State) {
	p.States = append(p.States, s)
}

// New creates a new product.
func New(t *ts.TS, n *gnba.NBA) *Product {
	return &Product{
		States: make([]*State, 0),
		TS:     t,
		NBA:    n,
	}
}

type context struct {
	R          StateSet
	U          StateStack
	T          StateSet
	V          StateStack
	CycleFound bool
}

func newContext() *context {
	return &context{
		R:          NewStateSet(),
		U:          NewStateStack(),
		T:          NewStateSet(),
		V:          NewStateStack(),
		CycleFound: false,
	}
}

// HasCycle return true if a cycle exists in the product
func (p *Product) HasCycle() bool {
	// TODO: Detect cycles in the product of TS and A
	p.addInitialStates()
	c := newContext()

	ir := NewStateSet(p.InitialStates...)
	for ir.Size() > 0 || !c.CycleFound {
		s := ir.Get()
		reachableCycle(s, c)
	}
	if !c.CycleFound {
		return true
	} else {
		return false // TODO: add counter example
	}
}

func reachableCycle(s *State, c *context) {
	c.U.Push(s)
	c.R.Add(s)
	for ok := true; ok; ok = !(c.U.Empty() || c.CycleFound) {
		s1 := c.U.Peek()
		if true {

		} else {

		}
	}
}

func (p *Product) addInitialStates() {
	for _, s0 := range p.TS.InitialStates {
		for q0 := range p.NBA.StartStates {
			for _, t := range q0.Transitions {
				if !t.Label.Conflicts(s0.Predicates) {
					q := t.State
					n := &State{
						StateTS:  s0,
						StateNBA: q,
					}
					p.InitialStates = append(p.InitialStates, n)
				}
			}
		}
	}
}
