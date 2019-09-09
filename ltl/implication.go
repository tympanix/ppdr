package ltl

// Implication is the ltl structure for the logical implication
type Implication struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are implications and has identical sub-trees
func (i Implication) SameAs(node Node) bool {
	if i2, ok := node.(Implication); ok {
		return i.LHS.SameAs(i2.LHS) && i.RHS.SameAs(i2.RHS)
	}
	return false
}

func (i Implication) LHSNode() Node {
	return i.LHS
}

func (i Implication) RHSNode() Node {
	return i.RHS
}

func (i Implication) String() string {
	return binaryNodeString(i, "->")
}

func (i Implication) Normalize() Node {
	return Disjunction{Negate(i.LHSNode().Normalize()), i.RHSNode().Normalize()}.Normalize()
}
