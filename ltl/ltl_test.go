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

func TestNormalize(t *testing.T) {

	tests := map[Node]Node{
		Always{AP{"a"}}:               Not{Until{True{}, Not{AP{"a"}}}},
		Implication{AP{"a"}, AP{"b"}}: Not{Conjunction{AP{"a"}, Not{AP{"b"}}}},
		Conjunction{Not{AP{"a"}}, Implication{Not{AP{"a"}}, AP{"b"}}}: Conjunction{Not{AP{"a"}}, Not{Conjunction{Not{AP{"a"}}, Not{AP{"b"}}}}},
	}

	i := 0

	for k, v := range tests {
		name := fmt.Sprintf("test:%d", i)
		t.Run(name, func(t *testing.T) {
			if !k.Normalize().SameAs(v) {
				t.Errorf("expected: %v\rgot: %v\n", v, k.Normalize())
			}
		})

		i++
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
