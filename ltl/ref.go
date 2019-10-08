package ltl

import "fmt"

// Ref is a reference to other propositional logic
type Ref struct {
	R int
}

func (r Ref) Compile(m *RefTable) Node {
	panic(ErrCompile)
}

func (r Ref) Len() int {
	return 0
}

func (r Ref) SameAs(n Node) bool {
	if r2, ok := n.(Ref); ok {
		return r.R == r2.R
	}
	return false
}

func (r Ref) Normalize() Node {
	return r
}

func (r Ref) String() string {
	return fmt.Sprintf("#%v", r.R)
}
