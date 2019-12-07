package ltl

import (
	"fmt"
	"strings"
)

// LitString is a string literal
type LitString struct {
	Str string
}

func (l LitString) Compile(m *RefTable) Node {
	panic(ErrCompile)
}

func (l LitString) Len() int {
	return 0
}

func (l LitString) Normalize() Node {
	return l
}

func (l LitString) Compare(n Node) (int, error) {
	if l2, ok := n.(LitString); ok {
		return strings.Compare(l.Str, l2.Str), nil
	}
	return 0, ErrNotComparable
}

func (l LitString) Map(fn MapFunc) Node {
	return fn(l)
}

func (l LitString) SameAs(n Node) bool {
	if l2, ok := n.(LitString); ok {
		return l.Str == l2.Str
	}
	return false
}

func (l LitString) String() string {
	return fmt.Sprintf("\"%v\"", l.Str)
}
