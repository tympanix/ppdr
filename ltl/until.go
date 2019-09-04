package ltl

import (
	"fmt"
)

type Until struct {
	LHS, RHS Node
}

func (until Until) SameAs(node Node) bool {
	if until2, ok := node.(Until); ok {
		return until.LHS.SameAs(until2.LHS) && until.RHS.SameAs(until2.RHS)
	}
	return false
}

func (until Until) LHSNode() Node {
	return until.LHS
}

func (until Until) RHSNode() Node {
	return until.RHS
}

func (until Until) String() string {
	return fmt.Sprintf("%v U %v", until.LHSNode(), until.RHSNode())
}
