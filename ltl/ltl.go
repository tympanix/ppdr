package ltl

// Node is any node of an LTL formula
type Node interface {
	SameAs(Node) bool
}

type BinaryNode interface {
	Node
	LHSNode() Node
	RHSNode() Node
}

type UnaryNode interface {
	Node
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

// Closure returns a list of all sub-nodes of a given node and the
// node itself.
func Closure(node Node) []Node {
	closureTemp := make([]Node, 0)
	closureWithDuplicates := auxClosure(node, closureTemp)
	return removeDuplicates(closureWithDuplicates)
}

func auxClosure(node Node, acc []Node) []Node {
	if ap, ok := node.(AP); ok {
		return append(acc, ap)
	} else if unary, ok := node.(UnaryNode); ok {
		acc = append(acc, unary)
		return auxClosure(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		acc = append(acc, binary)
		acc = auxClosure(binary.LHSNode(), acc)
		return auxClosure(binary.RHSNode(), acc)
	}
	panic("Unknown ltl node")
}

func removeDuplicates(nodes []Node) []Node {
	seen := make(map[Node]struct{}, len(nodes))
	j := 0
	for _, v := range nodes {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		nodes[j] = v
		j++
	}
	return nodes[:j]
}
