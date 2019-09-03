package ltl

type Next struct {
	Child Node
}

func (next Next) SameAs(node Node) bool {
	if next2, ok := node.(Next); ok {
		return next.Child.SameAs(next2.Child)
	}
	return false
}

func (next Next) ChildNode() Node {
	return next.Child
}
