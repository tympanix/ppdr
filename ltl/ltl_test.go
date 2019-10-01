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
