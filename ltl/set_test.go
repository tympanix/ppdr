package ltl

import (
	"fmt"
	"testing"
)

func TestContains(t *testing.T) {
	apA := AP{"A"}
	apB := AP{"B"}
	notAPC := Not{AP{"C"}}
	aUntilB := Until{apA, apB}
	set := NewSet(apA, apB, notAPC, aUntilB)

	if !set.Contains(apA) {
		t.Errorf("Set did not contain %+v; expected to do", apA)
	}
	if !set.Contains(apB) {
		t.Errorf("Set did not contain %+v; expected to do", apB)
	}
	if !set.Contains(notAPC) {
		t.Errorf("Set did not contain %+v; expected to do", notAPC)
	}
	if !set.Contains(aUntilB) {
		t.Errorf("Set did not contain %+v; expected to do", aUntilB)
	}
}

func TestSize(t *testing.T) {
	apA := AP{"A"}
	apB := AP{"B"}
	notAPC := Not{AP{"C"}}
	aUntilB := Until{apA, apB}
	set := NewSet(apA, apB, notAPC, aUntilB)

	if set.Size() != 4 {
		t.Errorf("set.Size() = %d; want 4", set.Size())
	}
}

func TestAdd(t *testing.T) {
	apA := AP{"A"}
	set := NewSet()

	if set.Contains(apA) {
		t.Errorf("Set do contain %+v; expected not to do", apA)
	}

	set.Add(apA)

	if !set.Contains(apA) {
		t.Errorf("Set did not contain %+v; expected to do", apA)
	}
}

func TestPowerSet(t *testing.T) {
	set := NewSet(AP{"A"}, AP{"B"}, AP{"C"})

	powerSet := set.PowerSet()

	if len(powerSet) != 8 {
		t.Error("Expected length to be 8 but was ", len(powerSet))
	}
}

func ExampleSet_PowerSet() {
	set := NewSet(AP{"A"}, AP{"B"}, AP{"C"}, Until{AP{"A"}, AP{"C"}})

	powerSet := set.SortedPowerSet()

	for _, e := range powerSet {
		fmt.Println(e)
	}

	// Output:
	// []
	// [A]
	// [B]
	// [C]
	// [A, B]
	// [A, C]
	// [B, C]
	// [A U C]
	// [A, B, C]
	// [A U C, B]
	// [A U C, C]
	// [A, A U C]
	// [A U C, B, C]
	// [A, A U C, B]
	// [A, A U C, C]
	// [A, A U C, B, C]
}
