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
		if p.isAccepting(v) {
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
	Space      Statespace
	S          *StateStack
	CycleFound bool
}

func newContext() *context {
	return &context{
		Space:      NewStatespace(),
		S:          NewStateStack(),
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

// HasAcceptingCycle returns true if the product has a cycle that goes through an acceptance state
func (p *Product) HasAcceptingCycle() bool {
	c := newContext()

	p.addInitialStates()

	for _, s := range p.InitialStates {
		if p.dfs(s, c) {
			return true
		}
	}
	return false
}

func (p *Product) dfs(s *State, c *context) bool {
	c.Space.Add(statespaceEntry{s, 0})
	c.S.Push(s)

	for t := range s.successors(p) {
		if !c.Space.Contains(statespaceEntry{t, 0}) {
			p.dfs(t, c)
		}
	}

	if p.isAccepting(s) {
		if p.ndfs(s, c) {
			return true
		}
	}

	c.S.Pop()

	return false
}

func (p *Product) ndfs(s *State, c *context) bool {
	c.Space.Add(statespaceEntry{s, 1})

	for t := range s.successors(p) {
		if !c.Space.Contains(statespaceEntry{t, 1}) {
			if p.ndfs(t, c) {
				return true
			}
		} else if c.S.Contains(t) {
			return true
		}
	}
	return false
}

func (p *Product) addInitialStates() {
	for _, s0 := range p.TS.InitialStates {
		for q0 := range p.NBA.StartStates {
			for _, t := range q0.Transitions {
				fmt.Println("t.label = ", t.Label)
				fmt.Println("s0.Predicates = ", s0.Predicates)
				if s0.Predicates.ContainsAny(t.Label) {
					q := t.State
					n := p.getOrAddState(s0, q)
					p.InitialStates = append(p.InitialStates, n)
				}
			}
		}
	}
}

func (p *Product) isAccepting(s *State) bool {
	return s.StateNBA.Has(p.NBA.Phi)
}
