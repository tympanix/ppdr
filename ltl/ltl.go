package ltl

// Node is any node of an LTL formula
type Node interface {
	SameAs(Node) bool
}

type BinaryNode interface {
	LHSNode() Node
	RHSNode() Node
}

type UnaryNode interface {
	ChildNode() Node
}

// FindAtomicPropositions returns a list of all atomic propositions
// in the LTL formula starting from the node
func FindAtomicPropositions(node Node) []Node {
	ap := make([]Node, 0)
	return auxFindAtomicPropositions(node, ap)
}

func auxFindAtomicPropositions(node Node, acc []Node) []Node {
	if ap, ok := node.(AP); ok {
		return append(acc, ap)
	} else if unary, ok := node.(UnaryNode); ok {
		return auxFindAtomicPropositions(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		aps := auxFindAtomicPropositions(binary.LHSNode(), acc)
		return auxFindAtomicPropositions(binary.RHSNode(), aps)
	}
	panic("Unknown ltl node")
}
