package ltl

type Until struct {
	LHS, RHS Node
}

func (u Until) SameAs(node Node) bool {
	if u2, ok := node.(Until); ok {
		return u.LHS.SameAs(u2.LHS) && u.RHS.SameAs(u2.RHS)
	}
	return false
}

func (u Until) LHSNode() Node {
	return u.LHS
}

func (u Until) RHSNode() Node {
	return u.RHS
}

func (u Until) String() string {
	return binaryNodeString(u, "U")
}

func (u Until) Normalize() Node {
	return Until{u.LHSNode().Normalize(), u.RHSNode().Normalize()}
}

func (u Until) Len() int {
	return 1 + u.LHSNode().Len() + u.RHSNode().Len()
}
