package ltl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tympanix/master-2019/debug"
)

// Node is any node of an LTL formula
type Node interface {
	SameAs(Node) bool
	Normalize() Node
	String() string
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
func FindAtomicPropositions(node Node) Set {
	ap := make(Set, 0)
	return auxFindAtomicPropositions(node, ap)
}

func auxFindAtomicPropositions(node Node, acc Set) Set {
	if ap, ok := node.(AP); ok {
		return acc.Add(ap)
	} else if _, ok := node.(True); ok {
		return acc
	} else if unary, ok := node.(UnaryNode); ok {
		return auxFindAtomicPropositions(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		aps := auxFindAtomicPropositions(binary.LHSNode(), acc)
		return auxFindAtomicPropositions(binary.RHSNode(), aps)
	}
	panic(fmt.Errorf("unknown ltl node %v", node))
}

// Closure returns a list of all sub-nodes of a given node and the
// node itself.
func Closure(node Node) Set {
	closureTemp := make(Set, 0)
	return auxClosure(node, closureTemp)
}

func auxClosure(node Node, acc Set) Set {
	if ap, ok := node.(AP); ok {
		acc = addNegation(ap, acc)
		return acc.Add(ap)
	} else if t, ok := node.(True); ok {
		return acc.Add(t)
	} else if unary, ok := node.(UnaryNode); ok {
		acc.Add(unary)
		addNegation(unary, acc)
		return auxClosure(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		acc.Add(binary)
		addNegation(binary, acc)
		acc = auxClosure(binary.LHSNode(), acc)
		return auxClosure(binary.RHSNode(), acc)
	}
	panic(fmt.Sprintf("unknown ltl node %v", node))
}

// Function will add the negation of a node to an array, if the node
// ifself is not already a negation.
func addNegation(node Node, nodes Set) Set {
	if _, ok := node.(Not); !ok {
		nodes.Add(Not{node})
	}

	return nodes
}

// FindElementarySets finds all the elementary sets for a
// closure(phi).
func FindElementarySets(closure Set) []Set {
	t := debug.NewTimer("elemsets")

	defer func() {
		t.Stop()
	}()

	elementarySets := make([]Set, 0)
	powerSet := closure.PowerSet()

	for _, set := range powerSet {
		if set.IsElementary(closure) {
			elementarySets = append(elementarySets, set)
		}
	}

	sort.SliceStable(elementarySets, func(i, j int) bool {
		a := elementarySets[i].String()
		b := elementarySets[j].String()

		if len(a) != len(b) {
			return len(a) < len(b)
		}

		return strings.Compare(a, b) < 0
	})

	return elementarySets
}
