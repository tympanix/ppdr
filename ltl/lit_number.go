package ltl

import "fmt"

// LitNumber is a string literal
type LitNumber struct {
	Num float64
}

func (l LitNumber) Compile(m *RefTable) Node {
	panic(ErrCompile)
}

func (l LitNumber) Len() int {
	return 0
}

func (l LitNumber) Normalize() Node {
	return l
}

func (l LitNumber) Compare(n Node) (int, error) {
	if l2, ok := n.(LitNumber); ok {
		diff := l.Num - l2.Num
		if diff > 0 {
			return 1, nil
		} else if diff < 0 {
			return -1, nil
		}
		return 0, nil
	}
	return 0, ErrNotComparable
}

func (l LitNumber) Map(fn MapFunc) Node {
	return fn(l)
}

func (l LitNumber) SameAs(n Node) bool {
	if l2, ok := n.(LitNumber); ok {
		return l.Num == l2.Num
	}
	return false
}

func (l LitNumber) String() string {
	return fmt.Sprintf("%v", l.Num)
}
