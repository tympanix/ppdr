package product

import (
	"fmt"
	"strings"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/ba"
	"github.com/tympanix/master-2019/systems/nba"
	"github.com/tympanix/master-2019/systems/ts"
)

// Product is a product of a transition system and a NBA
type Product struct {
	States        map[StateTuple]*State
	InitialStates []*State
	TS            ts.TS
	NBA           *nba.NBA
	RefTable      ltl.RefTable
}

// StringWithRenamer strings the product using a naming function
func (p *Product) StringWithRenamer(r ba.StateNamer) string {
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

		fmt.Fprintf(&sb, "%s%s%s\n", prefix, v.StringWithRenamer(r), suffix)
		for s := range v.Transitions {
			fmt.Fprintf(&sb, "\t-->\t%s\n", s.StringWithRenamer(r))
		}
	}

	return sb.String()
}

func (p *Product) String() string {
	return p.StringWithRenamer(nil)
}

// New creates a new product.
func New(t ts.TS, n *nba.NBA, r ltl.RefTable) *Product {
	return &Product{
		States:   make(map[StateTuple]*State),
		TS:       t,
		NBA:      n,
		RefTable: r,
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

// Context contains information about a cycle check
type Context struct {
	Space      Statespace
	S          *StateStack
	T          *StateStack
	CycleFound bool
}

// TraceWithRenamer return the cycle trace using a naming function
func (c *Context) TraceWithRenamer(r ba.StateNamer) string {
	var sb strings.Builder
	for _, s := range c.S.stack {
		fmt.Fprintf(&sb, "%v\n", s.StringWithRenamer(r))
	}
	fmt.Fprintf(&sb, "begin cycle:\n")
	for _, s := range c.T.stack {
		fmt.Fprintf(&sb, "%v\n", s.StringWithRenamer(r))
	}
	return sb.String()
}

func newContext() *Context {
	return &Context{
		Space:      NewStatespace(),
		S:          NewStateStack(),
		T:          NewStateStack(),
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

func (p *Product) getOrAddState(sTS ts.State, sNBA *ba.State) *State {
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
func (p *Product) HasAcceptingCycle() (cycle *Context) {
	defer func() {
		if r := recover(); r != nil {
			if c, ok := r.(*Context); ok {
				cycle = c
			}
		}
	}()

	c := newContext()

	p.addInitialStates()

	for _, s := range p.InitialStates {
		p.dfs(s, c)
	}

	return nil
}

func (p *Product) dfs(s *State, c *Context) {
	c.Space.Add(statespaceEntry{s, 0})
	c.S.Push(s)

	for t := range s.successors(p) {
		if !c.Space.Contains(statespaceEntry{t, 0}) {
			p.dfs(t, c)
		}
	}

	if p.isAccepting(s) {
		c.T = NewStateStack()
		p.ndfs(s, c)
	}

	c.S.Pop()
}

func (p *Product) ndfs(s *State, c *Context) {
	c.Space.Add(statespaceEntry{s, 1})
	c.T.Push(s)

	for t := range s.successors(p) {
		if c.S.Contains(t) {
			c.T.Push(t)
			panic(c)
		} else if !c.Space.Contains(statespaceEntry{t, 1}) {
			p.ndfs(t, c)
		}
	}

	c.T.Pop()
}

func (p *Product) addInitialStates() {
	for _, s0 := range p.TS.InitialStates() {
		lf := p.NBA.AP.Intersection(s0.Predicates(p.NBA.AP, p.RefTable))
		for q0 := range p.NBA.StartStates {
			for _, t := range q0.Transitions {
				if t.Label.Equals(lf) {
					q := t.State
					n := p.getOrAddState(s0, q)
					p.InitialStates = append(p.InitialStates, n)
				}
			}
		}
	}
}

func (p *Product) isAccepting(s *State) bool {
	return p.NBA.IsAcceptanceState(s.StateNBA)
}
