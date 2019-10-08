package ltl

// Equals is the ltl structure for the logical equals symbol
type Equals struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are equals and has identical sub-trees
func (e Equals) SameAs(node Node) bool {
	if e2, ok := node.(Equals); ok {
		return e.LHS.SameAs(e2.LHS) && e.RHS.SameAs(e2.RHS)
	}
	return false
}

func (e Equals) LHSNode() Node {
	return e.LHS
}

func (e Equals) RHSNode() Node {
	return e.RHS
}

func (e Equals) String() string {
	return binaryNodeString(e, "=")
}

func (e Equals) Normalize() Node {
	return e
}

func (e Equals) Len() int {
	return 1 + e.LHSNode().Len() + e.RHSNode().Len()
}
