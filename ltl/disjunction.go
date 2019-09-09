package ltl

// Disjunction is the ltl structure for the logical implication
type Disjunction struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are implications and has identical sub-trees
func (d Disjunction) SameAs(node Node) bool {
	if c2, ok := node.(Disjunction); ok {
		return d.LHS.SameAs(c2.LHS) && d.RHS.SameAs(c2.RHS)
	}
	return false
}

func (d Disjunction) LHSNode() Node {
	return d.LHS
}

func (d Disjunction) RHSNode() Node {
	return d.RHS
}

func (d Disjunction) String() string {
	return binaryNodeString(d, "|")
}

// Normalize rewrites the disjunctions to a conjunction
func (d Disjunction) Normalize() Node {
	return Not{Conjunction{Negate(d.LHSNode().Normalize()), Negate(d.RHSNode().Normalize())}}
}
