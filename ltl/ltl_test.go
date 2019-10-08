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

func TestClosure_three(t *testing.T) {
	phi := Negate(Always{Eventually{AP{"green"}}}.Normalize())

	closure := Closure(phi)

	if len(closure) != 8 {
		t.Error("Expected length to be 8 but was ", len(closure))
		t.Error(closure)
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
		Always{Or{Not{AP{"c1"}}, Not{AP{"c2"}}}}:       Not{Until{True{}, And{AP{"c1"}, AP{"c2"}}}},
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

func TestFormulaLength(t *testing.T) {

	tests := map[Node]int{
		AP{"a"}:                               0,
		And{AP{"a"}, AP{"b"}}:                 1,
		Not{AP{"a"}}:                          1,
		Until{And{AP{"a"}, AP{"b"}}, AP{"c"}}: 2,
		Next{Not{And{AP{"a"}, True{}}}}:       3,
		Eventually{Or{True{}, Not{AP{"a"}}}}:  3,
		Impl{Always{AP{"a"}}, AP{"b"}}:        2,
	}

	i := 0

	for n, l := range tests {
		name := fmt.Sprintf("test:%d", i)
		t.Run(name, func(t *testing.T) {
			if n.Len() != l {
				t.Errorf("expected formula length: %d, got %d", l, n.Len())
			}
		})
		i++
	}

}

func TestSatisfied_one(t *testing.T) {
	a := AP{"a"}
	b := AP{"b"}
	c := AP{"c"}
	d := AP{"d"}
	set := NewSet(a, b, c)

	tests := map[Node]bool{
		And{a, b}:           true,
		And{c, d}:           false,
		And{a, And{b, c}}:   true,
		Or{b, c}:            true,
		Or{a, d}:            true,
		Or{Or{a, b}, d}:     true,
		Impl{d, a}:          true,
		Impl{b, a}:          true,
		Impl{b, Impl{a, c}}: true,
		Impl{a, d}:          false,
		Impl{a, Impl{b, d}}: false,
		Not{a}:              false,
		Not{d}:              true,

		// Lazily evaluated examples
		Or{a, And{a, Next{b}}}: true,
		And{d, Or{a, Next{b}}}: false,
	}

	i := 0

	for k, v := range tests {
		name := fmt.Sprintf("test:%d", i)
		t.Run(name, func(t *testing.T) {
			s, err := Satisfied(k, set)
			if err != nil {
				t.Errorf("expected no errors from %v", k)
			}
			if s != v {
				t.Errorf("expected: %v\rgot: %v\n", v, s)
			}

		})

		i++
	}
}

func TestSatisfied_two(t *testing.T) {
	a := AP{"a"}
	b := AP{"b"}
	c := AP{"c"}
	set := NewSet(a, b, c)

	tests := []Node{
		Always{a},
		And{a, Always{b}},
		Eventually{a},
		Or{Eventually{a}, Eventually{b}},
		Next{a},
		Impl{a, Next{b}},
		Until{a, b},
		Or{Until{a, b}, c},
	}

	i := 0

	for _, v := range tests {
		name := fmt.Sprintf("test:%d", i)
		t.Run(name, func(t *testing.T) {
			_, err := Satisfied(v, set)
			if err != ErrNotPropositional {
				t.Errorf("expected: %v\rgot: %v\n", ErrNotPropositional, err)
			}

		})

		i++
	}
}

func TestCompile_one(t *testing.T) {

	tests := map[Node]Node{
		AP{"a"}:                         AP{"a"},
		Equals{AP{"a"}, LitString{"b"}}: Ref{0},
		And{AP{"a"}, Equals{AP{"a"}, LitString{"b"}}}: And{AP{"a"}, Ref{0}},

		// Multiple refereces
		Or{Equals{AP{"a"}, LitString{"b"}}, Equals{AP{"a"}, LitString{"b"}}}: Or{Ref{0}, Ref{1}},
	}

	for k, v := range tests {
		n, _, err := Compile(k)

		if err != nil {
			t.Error(err)
		}
		if !n.SameAs(v) {
			t.Errorf("extected %v, got %v", v, n)
		}
	}
}
