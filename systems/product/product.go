package product

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/systems/ba"
	"github.com/tympanix/master-2019/systems/nba"
	"github.com/tympanix/master-2019/systems/ts"
)

// Product is a product of a transition system and a NBA
type Product struct {
	States        map[StateTuple]*State
	InitialStates []*State
	TS            *ts.TS
	NBA           *nba.NBA
}

func (p *Product) String() string {
	var sb strings.Builder
	for _, v := range p.States {
		var prefix string
		if p.isInitialState(v) {
			prefix = ">"
		}

		var suffix string
		if v.isFinalState(p) {
			suffix = fmt.Sprintf("*")
		}

		fmt.Fprintf(&sb, "%s%s%s\n", prefix, v, suffix)
		for s := range v.Transitions {
			fmt.Fprintf(&sb, "\t-->\t%s\n", s)
		}
	}

	return sb.String()
}

// New creates a new product.
func New(t *ts.TS, n *nba.NBA) *Product {
	return &Product{
		States: make(map[StateTuple]*State),
		TS:     t,
		NBA:    n,
	}
}

func (p *Product) isInitialState(s *State) bool {
	for _, s1 := range p.InitialStates {
		if s == s1 {
			return true
		}
	}
	return false
}

type context struct {
	R          StateSet
	U          *StateStack
	T          StateSet
	V          *StateStack
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

// AddState to the states of the product.
func (p *Product) AddState(s *State) *State {
	if s, ok := p.States[s.StateTuple]; ok {
		return s
	}

	p.States[s.StateTuple] = s
	return s
}

func (p *Product) getOrAddState(sTS *ts.State, sNBA *ba.State) *State {
	stateTuple := StateTuple{
		StateTS:  sTS,
		StateNBA: sNBA,
	}

	if s, ok := p.States[stateTuple]; ok {
		return s
	}

	s := newState(sTS, sNBA)

	p.AddState(s)

	return s
}

// Satisfy return true if a cycle exists in the product
func (p *Product) Satisfy() bool {
	p.addInitialStates()
	c := newContext()

	ir := NewStateSet(p.InitialStates...)
	for ir.Size() > 0 && !c.CycleFound {
		s := ir.Get()
		p.reachableCycle(s, c)
	}
	if !c.CycleFound {
		return true
	}
	return false // TODO: add counter example
}

func (p *Product) reachableCycle(s *State, c *context) {
	c.U.Push(s)
	c.R.Add(s)
	for ok := true; ok; ok = !(c.U.Empty() || c.CycleFound) {
		s1 := c.U.Peek()
		if s2 := s1.unvisitedSucc(c.R, p); s2 != nil {
			c.U.Push(s2)
			c.R.Add(s2)
		} else {
			c.U.Pop()
			if !s1.isFinalState(p) { // TODO: Check if negation is correct
				c.CycleFound = p.cycleCheck(s1, c)
			}
		}
	}
}

func (p *Product) cycleCheck(s *State, c *context) bool {
	// Reset V and T
	c.T = NewStateSet()
	c.V = NewStateStack()
	cycleFound := false

	c.V.Push(s)
	c.T.Add(s)

	for ok := true; ok; ok = !(c.V.Empty() || cycleFound) {
		s1 := c.V.Peek()
		if s1.post(p).Contains(s) {
			cycleFound = true
			c.V.Push(s)
		} else {
			if s2 := s1.unvisitedSucc(c.T, p); s2 != nil {
				c.V.Push(s2)
				c.T.Add(s2)
			} else {
				c.V.Pop()
			}
		}
	}
	return cycleFound
}

func (p *Product) addInitialStates() {
	for _, s0 := range p.TS.InitialStates {
		for q0 := range p.NBA.StartStates {
			for _, t := range q0.Transitions {
				if !t.Label.Conflicts(s0.Predicates) {
					q := t.State
					n := p.getOrAddState(s0, q)
					p.InitialStates = append(p.InitialStates, n)
				}
			}
		}
	}
}
