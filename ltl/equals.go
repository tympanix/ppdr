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

// LHSNode returns the LHS of the equals operator
func (e Equals) LHSNode() Node {
	return e.LHS
}

// RHSNode returns the RHS of the equals operator
func (e Equals) RHSNode() Node {
	return e.RHS
}

func (e Equals) String() string {
	return binaryNodeString(e, "=")
}

// Normalize for equals performs nothing
func (e Equals) Normalize() Node {
	return e
}

// Compile ensures that LHS is an AP and RHS is a literal
func (e Equals) Compile(m *RefTable) Node {
	if _, ok := e.LHSNode().(AP); !ok {
		panic(ErrCompile)
	}
	switch e.RHSNode().(type) {
	case LitString, LitBool, LitNumber:
		ref := m.NewRef(e)
		return ref
	default:
		panic(ErrCompile)
	}
}

// Len returns the length of the equals operator and its children
func (e Equals) Len() int {
	return 1 + e.LHSNode().Len() + e.RHSNode().Len()
}

// Satisfied returns true if LHS and RHS are equal in type and value
func (e Equals) Satisfied(r Resolver) bool {
	var ap AP
	var ok bool

	if ap, ok = e.LHSNode().(AP); !ok {
		return false
	}

	lhs := r.Resolve(ap.Name)

	if lhs == nil {
		return false
	}

	return e.RHSNode().SameAs(lhs)
}

func (e Equals) Map(fn MapFunc) Node {
	return fn(Equals{e.LHSNode().Map(fn), e.RHSNode().Map(fn)})
}