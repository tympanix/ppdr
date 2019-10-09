package ltl

// Or is the ltl structure for the logical implication
type Or struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are implications and has identical sub-trees
func (d Or) SameAs(node Node) bool {
	if c2, ok := node.(Or); ok {
		return d.LHS.SameAs(c2.LHS) && d.RHS.SameAs(c2.RHS)
	}
	return false
}

func (d Or) LHSNode() Node {
	return d.LHS
}

func (d Or) RHSNode() Node {
	return d.RHS
}

func (d Or) String() string {
	return binaryNodeString(d, "|")
}

// Normalize rewrites the disjunctions to a conjunction
func (d Or) Normalize() Node {
	return Not{And{Negate(d.LHSNode().Normalize()), Negate(d.RHSNode().Normalize())}}
}

func (d Or) Compile(m *RefTable) Node {
	return Or{d.LHSNode().Compile(m), d.RHSNode().Compile(m)}
}

func (d Or) Len() int {
	return 1 + d.LHSNode().Len() + d.RHSNode().Len()
}

func (d Or) Satisfied(r Resolver) bool {
	if lhs, ok := d.LHSNode().(Satisfiable); ok {
		if rhs, ok := d.RHSNode().(Satisfiable); ok {
			return lhs.Satisfied(r) || rhs.Satisfied(r)
		}
	}
	panic(ErrNotPropositional)
}

func (d Or) Map(fn MapFunc) Node {
	return fn(Or{fn(d.LHSNode()), fn(d.RHSNode())})
}