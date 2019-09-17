package ts

import (
	"fmt"

	"github.com/tympanix/master-2019/gnba"
	"github.com/tympanix/master-2019/ltl"
)

// State is a node in the transition system
type State struct {
	Predicates ltl.Set
	Dependency []*State
}

// TS is a transition system
type TS struct {
	States        []*State
	StartingState []*State
}

// ValidateFormula return true if forumla phi holds in the transition system
func (t *TS) ValidateFormula(phi ltl.Node) bool {
	// TODO: Negate formula and create NBA
	nba := gnba.TransformGNBAtoNBA(gnba.GenerateGNBA(ltl.Negate(phi)))
	fmt.Println(nba)

	// TODO: Explore and find !phi states on-the-fly
	ts := &TS{}

	// TODO: Detect cycles in the product TS

	return true
}
