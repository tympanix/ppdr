package ltl

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// ErrNotPropositional is an error for nodes not supporting propositional logic
var ErrNotPropositional = errors.New("not propositional logic")

// ErrCompile is an error for when compilations has failed
var ErrCompile = errors.New("compile error")

// Node is any node of an LTL formula
type Node interface {
	SameAs(Node) bool
	Normalize() Node
	Compile(*RefTable) Node
	Map(MapFunc) Node
	Len() int
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

// RefTable references other propositional logic
type RefTable map[Ref]Node

// MapFunc is a function which maps one node to another
type MapFunc func(Node) Node

// NewRef adds a new reference to the reference table
func (r *RefTable) NewRef(n Node) Ref {
	ref := Ref{len(*r)}
	(*r)[ref] = n
	return ref
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
	} else if r, ok := node.(Ref); ok {
		return acc.Add(r)
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
	sub := Subformulae(node)

	for s := range sub {
		sub.Add(Negate(s))
	}

	return sub
}

// Subformulae return all subformulaes of a LTL formula
func Subformulae(node Node) Set {
	closureTemp := make(Set, 0)
	return auxSubformulae(node, closureTemp)
}

func auxSubformulae(node Node, acc Set) Set {
	if ap, ok := node.(AP); ok {
		return acc.Add(ap)
	} else if r, ok := node.(Ref); ok {
		return acc.Add(r)
	} else if p, ok := node.(Ptr); ok {
		return acc.Add(p)
	} else if t, ok := node.(True); ok {
		return acc.Add(t)
	} else if unary, ok := node.(UnaryNode); ok {
		acc.Add(unary)
		return auxSubformulae(unary.ChildNode(), acc)
	} else if binary, ok := node.(BinaryNode); ok {
		acc.Add(binary)
		acc = auxSubformulae(binary.LHSNode(), acc)
		return auxSubformulae(binary.RHSNode(), acc)
	}
	panic(fmt.Sprintf("unknown ltl node %v", node))
}

// Satisfied returns true if the set satisfies phi
func Satisfied(phi Node, r Resolver) (sat bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if r == ErrNotPropositional {
				err = ErrNotPropositional
			} else {
				panic(r)
			}
		}
	}()
	if p, ok := phi.(Satisfiable); ok {
		return p.Satisfied(r), nil
	}
	return false, ErrNotPropositional
}

// Compile runs through each node and substitutes unwanted propositional logic
// with references which can then later be evaluated
func Compile(phi Node) (n Node, t RefTable, err error) {
	defer func() {
		if r := recover(); r != nil {
			if r == ErrCompile {
				err = ErrCompile
			}
		}
	}()

	t = make(RefTable)
	return phi.Compile(&t), t, nil
}

// ValueToLiteral return the value as an LTL literal node
func ValueToLiteral(value interface{}) Node {
	switch v := value.(type) {
	case string:
		return LitString{v}
	case bool:
		return LitBool{v}
	case int:
		return LitNumber{float64(v)}
	case int64:
		return LitNumber{float64(v)}
	case float32:
		return LitNumber{float64(v)}
	case float64:
		return LitNumber{float64(v)}
	}
	panic("unsupported literal type")
}

// RenameSelfPredicate renames all self predicates
func RenameSelfPredicate(phi Node, ptr unsafe.Pointer) Node {
	return phi.Map(func(n Node) Node {
		if _, ok := n.(Self); ok {
			return Ptr{ptr}
		}
		return n
	})
}
