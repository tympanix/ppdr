package ltl

// Impl is the ltl structure for the logical implication
type Impl struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are implications and has identical sub-trees
func (i Impl) SameAs(node Node) bool {
	if i2, ok := node.(Impl); ok {
		return i.LHS.SameAs(i2.LHS) && i.RHS.SameAs(i2.RHS)
	}
	return false
}

func (i Impl) LHSNode() Node {
	return i.LHS
}

func (i Impl) RHSNode() Node {
	return i.RHS
}

func (i Impl) String() string {
	return binaryNodeString(i, "->")
}

func (i Impl) Normalize() Node {
	return Or{Negate(i.LHSNode().Normalize()), i.RHSNode().Normalize()}.Normalize()
}

func (i Impl) Compile(m *RefTable) Node {
	return Impl{i.LHSNode().Compile(m), i.RHSNode().Compile(m)}
}

func (i Impl) Len() int {
	return 1 + i.LHSNode().Len() + i.RHSNode().Len()
}

func (i Impl) Satisfied(s Set) bool {
	if lhs, ok := i.LHSNode().(Satisfiable); ok {
		if rhs, ok := i.RHSNode().(Satisfiable); ok {
			return !lhs.Satisfied(s) || rhs.Satisfied(s)
		}
	}
	panic(ErrNotPropositional)
}
