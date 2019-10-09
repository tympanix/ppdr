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

func (u Until) Compile(m *RefTable) Node {
	return Until{u.LHSNode().Compile(m), u.RHSNode().Compile(m)}
}

func (u Until) Len() int {
	return 1 + u.LHSNode().Len() + u.RHSNode().Len()
}

func (u Until) Filter(fn MapFunc) Node {
	fn(Until{fn(u.LHSNode()), fn(u.RHSNode())})
}
