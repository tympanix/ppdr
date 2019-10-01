package gnba

import (
	"fmt"
	"testing"

	"github.com/tympanix/master-2019/ltl"
)

func ExampleGenerateGNBA() {
	phi := ltl.Next{ltl.AP{"a"}}
	g := GenerateGNBA(phi)

	fmt.Println(g)
	//Output:
	// >[Oa, a]
	// 	[a]	-->	[Oa, a]
	// 	[a]	-->	[!Oa, a]
	// [!Oa, a]
	// 	[a]	-->	[!a, Oa]
	// 	[a]	-->	[!Oa, !a]
	// >[!a, Oa]
	// 	[]	-->	[Oa, a]
	// 	[]	-->	[!Oa, a]
	// [!Oa, !a]
	// 	[]	-->	[!a, Oa]
	// 	[]	-->	[!Oa, !a]
}

func ExampleGenerateGNBA_second() {
	phi := ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}}

	g := GenerateGNBA(phi)

	fmt.Println(g)

	//Output:
	// >[a, a U b, b]{[0]}
	// 	[a, b]	-->	[a, a U b, b]
	// 	[a, b]	-->	[!a, a U b, b]
	// 	[a, b]	-->	[!b, a, a U b]
	// 	[a, b]	-->	[!(a U b), !b, a]
	// 	[a, b]	-->	[!(a U b), !a, !b]
	// >[!a, a U b, b]{[0]}
	// 	[b]	-->	[a, a U b, b]
	// 	[b]	-->	[!a, a U b, b]
	// 	[b]	-->	[!b, a, a U b]
	// 	[b]	-->	[!(a U b), !b, a]
	// 	[b]	-->	[!(a U b), !a, !b]
	// >[!b, a, a U b]
	// 	[a]	-->	[a, a U b, b]
	// 	[a]	-->	[!a, a U b, b]
	// 	[a]	-->	[!b, a, a U b]
	// [!(a U b), !b, a]{[0]}
	// 	[a]	-->	[!(a U b), !b, a]
	// 	[a]	-->	[!(a U b), !a, !b]
	// [!(a U b), !a, !b]{[0]}
	// 	[]	-->	[a, a U b, b]
	// 	[]	-->	[!a, a U b, b]
	// 	[]	-->	[!b, a, a U b]
	// 	[]	-->	[!(a U b), !b, a]
	// 	[]	-->	[!(a U b), !a, !b]
}

func TestCopyGNBA(t *testing.T) {
	tests := []ltl.Node{
		ltl.Next{ltl.AP{"a"}},
		ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}},
	}

	for i, phi := range tests {

		name := fmt.Sprintf("test:%d", i)

		t.Run(name, func(t *testing.T) {
			checkGNBAEquality(t, phi)
		})

	}
}

func checkGNBAEquality(t *testing.T, phi ltl.Node) {
	g1 := GenerateGNBA(phi)
	g2 := g1.Copy()

	// Check string representation is the same
	if g1.String() != g2.String() {
		t.Error("gnba copy is not equal to original one")
	}

	// Check states are disjoint
	for _, s := range g2.States {
		if g1.HasState(s) {
			t.Error("gnba copy states are not disjoint")
		}
	}

	// Check transitions are disjoint
	for _, s := range g2.States {
		for _, tr := range s.Transitions {
			if tr.State == nil {
				t.Error("gnba copy has nil transition")
			}
			if g1.HasState(tr.State) {
				t.Error("gnba copies are not disjoint")
			}
		}
	}

	// Check starting states
	if g1.StartingStates.Size() != g2.StartingStates.Size() {
		t.Error("gnba copies do not have equal number of starting states")
	}
	for s := range g2.StartingStates {
		if g1.StartingStates.Contains(s) || g1.HasState(s) {
			t.Error("gnba copy starting states are not disjoint")
		}
	}

	// Check acceptance sets
	if len(g2.FinalStates) != len(g1.FinalStates) {
		t.Error("gnba copy number of acceptance sets not equal")
	}
	for _, set := range g2.FinalStates {
		for s := range set {
			if _, ok := g1.IsAcceptanceState(s); ok {
				t.Error("gnba acceptance sets are not disjoint")
			}
			if g1.HasState(s) {
				t.Error("gnba acceptence sets are not disjoint")
			}
		}

	}
}
