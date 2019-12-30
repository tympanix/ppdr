package ltl

import (
	"sort"
	"strings"

	"github.com/tympanix/ppdr/debug"
)

type esContext struct {
	sub     []Node
	es      *[]Set
	cur     Set
	closure Set
}

func (c *esContext) next() Node {
	return c.sub[0]
}

func (c *esContext) copy() *esContext {
	return &esContext{
		sub:     c.sub,
		es:      c.es,
		cur:     c.cur.Copy(),
		closure: c.closure,
	}
}

func (c *esContext) rec(n Node) *esContext {

	if !c.isConsistentAfter(n) {
		return nil
	}

	c1 := c.copy()

	c1.sub = c.sub[1:]
	c1.cur.Add(n)

	return c1
}

// isConsistent checks whether the context is still consistent with respect
// ro propositional logic and consistent with respect to the until operator
func (c *esContext) isConsistentAfter(n Node) bool {
	if c.cur.Contains(Negate(n)) {
		return false
	} else if a, ok := n.(And); ok {
		if !(c.cur.Contains(a.LHSNode()) && c.cur.Contains(a.RHSNode())) {
			return false
		}
	} else if u, ok := n.(Until); ok {
		if !(c.cur.Contains(u.RHSNode()) || c.cur.Contains(u.LHSNode())) {
			return false
		}
	} else if n, ok := n.(Not); ok {
		if u, ok := n.ChildNode().(Until); ok {
			if c.cur.Contains(u.RHSNode()) {
				return false
			}
		}
	}
	return true
}

// FindElementarySets return all sets of the closure which are elementary
func FindElementarySets(phi Node) []Set {
	t := debug.NewTimer("elemsets")

	defer func() {
		t.Stop()
	}()

	sub := Subformulae(phi).AsSlice()

	sort.SliceStable(sub, func(i, j int) bool {
		return sub[i].Len() < sub[j].Len()
	})

	es := make([]Set, 0)

	c := &esContext{
		sub:     sub,
		es:      &es,
		cur:     NewSet(),
		closure: Closure(phi),
	}

	auxFindElementarySets(c)

	sort.SliceStable(es, func(i, j int) bool {
		a := es[i].String()
		b := es[j].String()

		if len(a) != len(b) {
			return len(a) < len(b)
		}

		return strings.Compare(a, b) < 0
	})

	return es
}

func auxFindElementarySets(c *esContext) {

	if c == nil {
		return
	}

	if len(c.sub) == 0 {
		if c.cur.IsElementary(c.closure) {
			*(c.es) = append(*(c.es), c.cur)
		}
		return
	}

	if _, ok := c.next().(True); ok {
		auxFindElementarySets(c.rec(c.next()))
	} else {
		auxFindElementarySets(c.rec(c.next()))
		auxFindElementarySets(c.rec(Negate(c.next())))
	}

}
