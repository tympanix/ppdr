package ltl

import "fmt"

// Greater is the ltl structure for the logical equals symbol
type Greater struct {
	LHS, RHS Node
}

// SameAs returns true if both nodes are equals and has identical sub-trees
func (e Greater) SameAs(node Node) bool {
	if e2, ok := node.(Greater); ok {
		return e.LHS.SameAs(e2.LHS) && e.RHS.SameAs(e2.RHS)
	}
	return false
}

// LHSNode returns the LHS of the equals operator
func (e Greater) LHSNode() Node {
	return e.LHS
}

// RHSNode returns the RHS of the equals operator
func (e Greater) RHSNode() Node {
	return e.RHS
}

func (e Greater) String() string {
	return binaryNodeString(e, ">")
}

// Normalize for equals performs nothing
func (e Greater) Normalize() Node {
	return e
}

// Compile ensures that LHS is an AP and RHS is a literal
func (e Greater) Compile(m *RefTable) Node {

	if !e.isValidType(e.LHSNode()) {
		panic(fmt.Sprintf("greater lhs invalid type : %T", e.LHSNode()))
	}

	if !e.isValidType(e.RHSNode()) {
		panic(fmt.Sprintf("greater rhs invalid type : %T", e.RHSNode()))
	}

	return m.NewRef(e)
}

func (e Greater) isValidType(n Node) bool {
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
func (e Greater) Len() int {
	return 1 + e.LHSNode().Len() + e.RHSNode().Len()
}

// Satisfied returns true if LHS and RHS are equal in type and value
func (e Greater) Satisfied(r Resolver) bool {
	lhs := e.LHSNode()
	rhs := e.RHSNode()

	if ap, ok := e.LHSNode().(AP); ok {
		lhs = r.Resolve(ap.Name)
	}

	if ap, ok := e.RHSNode().(AP); ok {
		rhs = r.Resolve(ap.Name)
	}

	if lhs == nil || rhs == nil {
		return false
	}

	if e, ok := lhs.(Comparable); ok {
		if d, err := e.Compare(rhs); err == nil {
			return d > 0
		}
	}
	return false
}

// Map maps the grater operator with a function
func (e Greater) Map(fn MapFunc) Node {
	return fn(Greater{e.LHSNode().Map(fn), e.RHSNode().Map(fn)})
}
