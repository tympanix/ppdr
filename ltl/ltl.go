package ltl

import (
	"fmt"
	"strings"
)

// Node is any node of an LTL formula
type Node interface {
	SameAs(Node) bool
}

// BinaryNode is an ltl node which has two child nodes
type BinaryNode interface {
	Node
	LHSNode() Node
	RHSNode() Node
}

// UnaryNode is an ltl node which only has one child
type UnaryNode interface {
	Node
	ChildNode() Node
}

func binaryNodeString(b BinaryNode, op string) string {
	var sb strings.Builder

	if _, ok := b.LHSNode().(BinaryNode); ok {
		fmt.Fprintf(&sb, "(%v)", b.LHSNode())
	} else {
		fmt.Fprint(&sb, b.LHSNode())
	}
	fmt.Fprintf(&sb, " %s ", op)
	if _, ok := b.RHSNode().(BinaryNode); ok {
		fmt.Fprintf(&sb, "(%v)", b.RHSNode())
	} else {
		fmt.Fprint(&sb, b.RHSNode())
	}
	return sb.String()
}

// Negate negates the ltl formula and removes double negations
func Negate(node Node) Node {
	if n, ok := node.(Not); ok {
		return n.ChildNode()
	}

	return Not{node}
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

// FindElementarySets finds all the elementary sets for a
// closure(phi).
func FindElementarySets(closure Set) []Set {
	elementarySets := make([]Set, 0)
	powerSet := closure.PowerSet()

	for _, set := range powerSet {
		if set.IsElementary(closure) {
			elementarySets = append(elementarySets, set)
		}
	}

	return elementarySets
}
