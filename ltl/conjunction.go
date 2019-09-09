package ltl

// Conjunction is the ltl structure for the logical implication
type Conjunction struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are implications and has identical sub-trees
func (c Conjunction) SameAs(node Node) bool {
	if c2, ok := node.(Conjunction); ok {
		return c.LHS.SameAs(c2.LHS) && c.RHS.SameAs(c2.RHS)
	}
	return false
}

func (c Conjunction) LHSNode() Node {
	return c.LHS
}

func (c Conjunction) RHSNode() Node {
	return c.RHS
}

func (c Conjunction) String() string {
	return binaryNodeString(c, "&")
}
