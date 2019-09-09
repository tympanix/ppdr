package ltl

import (
	"fmt"
	"testing"
)

func TestClosure(t *testing.T) {
	phi := Until{Next{AP{"a"}}, Not{AP{"b"}}}

	closure := Closure(phi)

	if len(closure) != 8 {
		t.Error("Expected length to be 8 but was ", len(closure))
	}

}

func ExampleFindElementarySets() {
	phi := Next{AP{"A"}}
	elemSets := FindElementarySets(Closure(phi))

	for _, s := range elemSets {
		fmt.Println(s)
	}

	// Output:
	// [OA !A]
	// [!OA !A]
	// [OA A]
	// [!OA A]
}
