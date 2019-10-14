package ltl

import "fmt"

// LitBool is a string literal
type LitBool struct {
	Bool bool
}

func (l LitBool) Compile(m *RefTable) Node {
	panic(ErrCompile)
}

func (l LitBool) Len() int {
	return 0
}

func (l LitBool) Normalize() Node {
	return l
}

func (l LitBool) SameAs(n Node) bool {
	if l2, ok := n.(LitBool); ok {
		return l.Bool == l2.Bool
	}
	return false
}

func (l LitBool) Map(fn MapFunc) Node {
	return fn(l)
}

func (l LitBool) String() string {
	return fmt.Sprintf("%v", l.Bool)
}
