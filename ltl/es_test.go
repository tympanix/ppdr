package ltl

import (
	"fmt"
	"testing"
)

// formula: !Oa
func TestElementarySets_one(t *testing.T) {
	phi := Not{Next{AP{"a"}}}
	elemSets := FindElementarySets(phi)

	golden := []Set{
		NewSet(Not{Next{AP{"a"}}}, AP{"a"}),
		NewSet(Not{Next{AP{"a"}}}, Not{AP{"a"}}),
		NewSet(Next{AP{"a"}}, AP{"a"}),
		NewSet(Next{AP{"a"}}, Not{AP{"a"}}),
	}

	compareElemSets(t, elemSets, golden)

}

// formula: true U !green
func TestElementarySets_two(t *testing.T) {
	phi := Until{True{}, Not{AP{"green"}}}
	elemSets := FindElementarySets(phi)

	golden := []Set{
		NewSet(phi, True{}, AP{"green"}),
		NewSet(phi, True{}, Not{AP{"green"}}),
		NewSet(Negate(phi), True{}, AP{"green"}),
	}

	compareElemSets(t, elemSets, golden)

}

// formula: a U (b and c)
func TestElementarySets_three(t *testing.T) {
	a := AP{"a"}
	b := AP{"b"}
	c := AP{"c"}
	and := And{b, c}

	phi := Until{a, and}
	elemSets := FindElementarySets(phi)

	golden := []Set{
		NewSet(phi, and, a, b, c),
		NewSet(phi, and, Not{a}, b, c),
		NewSet(phi, Not{and}, a, Not{b}, c),
		NewSet(phi, Not{and}, a, b, Not{c}),
		NewSet(phi, Not{and}, a, Not{b}, Not{c}),
		NewSet(Not{phi}, Not{and}, a, Not{b}, c),
		NewSet(Not{phi}, Not{and}, a, b, Not{c}),
		NewSet(Not{phi}, Not{and}, a, Not{b}, Not{c}),
		NewSet(Not{phi}, Not{and}, Not{a}, Not{b}, c),
		NewSet(Not{phi}, Not{and}, Not{a}, b, Not{c}),
		NewSet(Not{phi}, Not{and}, Not{a}, Not{b}, Not{c}),
	}

	compareElemSets(t, elemSets, golden)

}

// formula: true U (c1 and c2)
func TestElementarySets_four(t *testing.T) {
	c1 := AP{"c1"}
	c2 := AP{"c2"}
	and := And{c1, c2}

	phi := Until{True{}, and}
	elemSets := FindElementarySets(phi)

	golden := []Set{
		NewSet(Negate(phi), Not{and}, True{}, Not{c1}, c2),
		NewSet(Negate(phi), Not{and}, True{}, c1, Not{c2}),
		NewSet(Negate(phi), Not{and}, True{}, Not{c1}, Not{c2}),
		NewSet(phi, and, True{}, c1, c2),
		NewSet(phi, Not{and}, True{}, Not{c1}, c2),
		NewSet(phi, Not{and}, True{}, c1, Not{c2}),
		NewSet(phi, Not{and}, True{}, Not{c1}, Not{c2}),
	}

	compareElemSets(t, elemSets, golden)

}

func ExampleFindElementarySets_one() {
	phi := Next{AP{"A"}}
	elemSets := FindElementarySets(phi)

	for _, s := range elemSets {
		fmt.Println(s)
	}

	// Output:
	// [A, OA]
	// [!A, OA]
	// [!OA, A]
	// [!A, !OA]
}

func ExampleFindElementarySets_two() {
	phi := Until{True{}, Not{Until{True{}, AP{"green"}}}}
	elemSets := FindElementarySets(phi)

	for _, s := range elemSets {
		fmt.Println(s)
	}

	//Output:
	// [green, true, true U !(true U green), true U green]
	// [!green, true, true U !(true U green), true U green]
	// [!(true U !(true U green)), green, true, true U green]
	// [!(true U !(true U green)), !green, true, true U green]
	// [!(true U green), !green, true, true U !(true U green)]
}

func compareElemSets(t *testing.T, test []Set, golden []Set) {
	if len(test) != len(golden) {
		t.Fatalf("elementary sets length expected: %v, got: %v", len(golden), len(test))
	}
Elem:
	for _, g := range golden {
		for i, e := range test {
			if e.ContainsAll(g) && g.ContainsAll(e) && e.Size() == g.Size() {
				test = append(test[:i], test[i+1:]...)
				continue Elem
			}
		}
		t.Errorf("expected elementary set to contain: %v", g)
	}
}
