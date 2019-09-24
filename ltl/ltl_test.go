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

func TestClosure_two(t *testing.T) {
	phi := Eventually{AP{"green"}}.Normalize()

	closure := Closure(phi)

	if len(closure) != 6 {
		t.Error("Expected length to be 8 but was ", len(closure))
	}

	golden := []Node{
		phi,
		Negate(phi),
		True{},
		Negate(True{}),
		AP{"green"},
		Negate(AP{"green"}),
	}

	for _, v := range golden {
		if !closure.Contains(v) {
			t.Errorf("Expected closure to contain %s.", v)
		}
	}

}

func TestNormalize(t *testing.T) {

	tests := map[Node]Node{
		Always{AP{"a"}}:        Not{Until{True{}, Not{AP{"a"}}}},
		Impl{AP{"a"}, AP{"b"}}: Not{And{AP{"a"}, Not{AP{"b"}}}},
		And{Not{AP{"a"}}, Impl{Not{AP{"a"}}, AP{"b"}}}: And{Not{AP{"a"}}, Not{And{Not{AP{"a"}}, Not{AP{"b"}}}}},
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
	// [A, OA]
	// [!A, OA]
	// [!OA, A]
	// [!A, !OA]
}
