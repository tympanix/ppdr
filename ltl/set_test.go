package ltl

import (
	"testing"
)

func TestPowerSet(t *testing.T) {
	set := Set{AP{"A"}, AP{"B"}, AP{"C"}}

	powerSet := set.PowerSet()
	
	if len(powerSet) != 8 {
		t.Error("Expected length to be 8 but was ", len(powerSet))
	}

}
