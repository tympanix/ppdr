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
func Closure(node Node) Set {
	closureTemp := make(Set, 0)
	closureWithDuplicates := auxClosure(node, closureTemp)
	return removeDuplicates(closureWithDuplicates)
}

func auxClosure(node Node, acc Set) Set {
	if ap, ok := node.(AP); ok {
		acc = addNegation(ap, acc)
		return append(acc, ap)
	} else if unary, ok := node.(UnaryNode); ok {
		acc = append(acc, unary)
		acc = addNegation(unary, acc)
		return auxClosure(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		acc = append(acc, binary)
		acc = addNegation(binary, acc)
		acc = auxClosure(binary.LHSNode(), acc)
		return auxClosure(binary.RHSNode(), acc)
	}
	panic("Unknown ltl node")
}

// Function will add the negation of a node to an array, if the node
// ifself is not already a negation.
func addNegation(node Node, nodes Set) Set {
	if _, ok := node.(Not); !ok {
		return append(nodes, Not{node})
	}

	return nodes
}

func removeDuplicates(nodes Set) Set {
	seen := make(map[Node]struct{}, len(nodes))
	i := 0
	for _, node := range nodes {
		if _, ok := seen[node]; ok {
			continue
		}
		seen[node] = struct{}{}
		nodes[i] = node
		i++
	}
	return nodes[:i]
}
