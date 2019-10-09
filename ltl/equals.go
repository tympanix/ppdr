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

func (e Equals) Compile(m *RefTable) Node {
	if _, ok := e.LHSNode().(AP); !ok {
		panic(ErrCompile)
	}
	if _, ok := e.RHSNode().(LitString); !ok {
		panic(ErrCompile)
	}
	ref := m.NewRef(e)
	return ref
}

func (e Equals) Len() int {
	return 1 + e.LHSNode().Len() + e.RHSNode().Len()
}

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

	if s, ok := e.RHSNode().(LitString); ok {
		if s2, ok := lhs.(string); ok {
			return s.Str == s2
		}
	}

	return false
}
