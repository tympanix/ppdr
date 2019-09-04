package ltl

import (
	"testing"
)

func TestClosure(t *testing.T) {
	phi := Until{Next{AP{"a"}}, Not{AP{"b"}}}

	closure := Closure(phi)

	if len(closure) != 5 {
		t.Error("Expected length to be 4.")
	}

	if _, ok := closure[0].(Until); !ok {
		t.Error("Expected \"Oa U !b\" but got: ", closure[0])
	}

	if _, ok := closure[1].(Next); !ok {
		t.Error("Expected \"Oa\" but got: ", closure[1])
	}

	if _, ok := closure[2].(AP); !ok {
		t.Error("Expected \"a\" but got: ", closure[2])
	}

	if _, ok := closure[3].(Not); !ok {
		t.Error("Expected \"!b\" but got: ", closure[3])
	}

	if _, ok := closure[4].(AP); !ok {
		t.Error("Expected \"b\" but got: ", closure[4])
	}

}
