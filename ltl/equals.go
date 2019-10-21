package ltl

import "fmt"

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

	if !e.isValidType(e.LHSNode()) {
		panic(fmt.Sprintf("equals lhs invalid type : %T", e.LHSNode()))
	}

	if !e.isValidType(e.RHSNode()) {
		panic(fmt.Sprintf("equals rhs invalid type : %T", e.RHSNode()))
	}

	return m.NewRef(e)
}

func (e Equals) isValidType(n Node) bool {
	switch n.(type) {
	// Base types
	case AP, Ptr:
		return true
	// Literals
	case LitBool, LitNumber, LitString:
		return true
	default:
		return false
	}
}

// Len returns the length of the equals operator and its children
func (e Equals) Len() int {
	return 1 + e.LHSNode().Len() + e.RHSNode().Len()
}

// Satisfied returns true if LHS and RHS are equal in type and value
func (e Equals) Satisfied(r Resolver) bool {
	var lhs Node = e.LHSNode()
	var rhs Node = e.RHSNode()

	if ap, ok := e.LHSNode().(AP); ok {
		lhs = r.Resolve(ap.Name)
	}

	if ap, ok := e.RHSNode().(AP); ok {
		rhs = r.Resolve(ap.Name)
	}

	if lhs == nil || rhs == nil {
		return false
	}

	return lhs.SameAs(rhs)
}

func (e Equals) Map(fn MapFunc) Node {
	return fn(Equals{e.LHSNode().Map(fn), e.RHSNode().Map(fn)})
}
