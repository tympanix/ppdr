package parser

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/tympanix/ppdr/ltl"
	"github.com/tympanix/ppdr/ltl/scanner"
)

func TestParser(t *testing.T) {

	tests := map[string]ltl.Node{
		// Operators
		"a = b":    ltl.Equals{ltl.AP{"a"}, ltl.AP{"b"}},
		"a or b":   ltl.Or{ltl.AP{"a"}, ltl.AP{"b"}},
		"a -> b":   ltl.Impl{ltl.AP{"a"}, ltl.AP{"b"}},
		"a U b":    ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}},
		"Oa":       ltl.Next{ltl.AP{"a"}},
		"a and b":  ltl.And{ltl.AP{"a"}, ltl.AP{"b"}},
		"true U b": ltl.Until{ltl.True{}, ltl.AP{"b"}},
		"OOa":      ltl.Next{ltl.Next{ltl.AP{"a"}}},
		"<><>a":    ltl.Eventually{ltl.Eventually{ltl.AP{"a"}}},
		"[]<>a":    ltl.Always{ltl.Eventually{ltl.AP{"a"}}},
		"<>[]a":    ltl.Eventually{ltl.Always{ltl.AP{"a"}}},

		// Literals
		"false = a":             ltl.Equals{ltl.LitBool{false}, ltl.AP{"a"}},
		"b = true":              ltl.Equals{ltl.AP{"b"}, ltl.LitBool{true}},
		"3.1415 = pi":           ltl.Equals{ltl.LitNumber{3.1415}, ltl.AP{"pi"}},
		"5 = 7":                 ltl.Equals{ltl.LitNumber{5}, ltl.LitNumber{7}},
		"\"hello\" = \"world\"": ltl.Equals{ltl.LitString{"hello"}, ltl.LitString{"world"}},

		// Expressions
		"2 = 5":  ltl.Equals{ltl.LitNumber{2}, ltl.LitNumber{5}},
		"2 != 5": ltl.NotEqual{ltl.LitNumber{2}, ltl.LitNumber{5}},
		"2 > 5":  ltl.Greater{ltl.LitNumber{2}, ltl.LitNumber{5}},
		"2 >= 5": ltl.GreaterEqual{ltl.LitNumber{2}, ltl.LitNumber{5}},
		"2 < 5":  ltl.Less{ltl.LitNumber{2}, ltl.LitNumber{5}},
		"2 <= 5": ltl.LessEqual{ltl.LitNumber{2}, ltl.LitNumber{5}},

		// Functions
		"a = subject()":      ltl.Equals{ltl.AP{"a"}, ltl.Subject{}},
		"a = user(\"john\")": ltl.Equals{ltl.AP{"a"}, ltl.User{"john"}},

		// Negations
		"!Oa":        ltl.Not{ltl.Next{ltl.AP{"a"}}},
		"!(!(Oa))":   ltl.Not{ltl.Not{ltl.Next{ltl.AP{"a"}}}},
		"!!Oa":       ltl.Not{ltl.Not{ltl.Next{ltl.AP{"a"}}}},
		"!(a and b)": ltl.Not{ltl.And{ltl.AP{"a"}, ltl.AP{"b"}}},
		"!(a or b)":  ltl.Not{ltl.Or{ltl.AP{"a"}, ltl.AP{"b"}}},

		// Associativity
		"a and b U c":   ltl.And{ltl.AP{"a"}, ltl.Until{ltl.AP{"b"}, ltl.AP{"c"}}},
		"a U b U c":     ltl.Until{ltl.AP{"a"}, ltl.Until{ltl.AP{"b"}, ltl.AP{"c"}}},
		"a and b and c": ltl.And{ltl.AP{"a"}, ltl.And{ltl.AP{"b"}, ltl.AP{"c"}}},
		"a or b or c":   ltl.Or{ltl.AP{"a"}, ltl.Or{ltl.AP{"b"}, ltl.AP{"c"}}},
		"a -> b -> c":   ltl.Impl{ltl.AP{"a"}, ltl.Impl{ltl.AP{"b"}, ltl.AP{"c"}}},

		// Precedence
		"!a U b":                    ltl.Until{ltl.Not{ltl.AP{"a"}}, ltl.AP{"b"}},
		"a and b -> c":              ltl.Impl{ltl.And{ltl.AP{"a"}, ltl.AP{"b"}}, ltl.AP{"c"}},
		"a -> b or c and d":         ltl.Impl{ltl.AP{"a"}, ltl.Or{ltl.AP{"b"}, ltl.And{ltl.AP{"c"}, ltl.AP{"d"}}}},
		"a and b or c -> d":         ltl.Impl{ltl.Or{ltl.And{ltl.AP{"a"}, ltl.AP{"b"}}, ltl.AP{"c"}}, ltl.AP{"d"}},
		"[]<> crit1 and []<> crit2": ltl.And{ltl.Always{ltl.Eventually{ltl.AP{"crit1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"crit2"}}}},
		"<> green and <> red":       ltl.And{ltl.Eventually{ltl.AP{"green"}}, ltl.Eventually{ltl.AP{"red"}}},
		"a -> b = c":                ltl.Impl{ltl.AP{"a"}, ltl.Equals{ltl.AP{"b"}, ltl.AP{"c"}}},
		"[]a = b":                   ltl.Always{ltl.Equals{ltl.AP{"a"}, ltl.AP{"b"}}},
		"[]a = b and c":             ltl.And{ltl.Always{ltl.Equals{ltl.AP{"a"}, ltl.AP{"b"}}}, ltl.AP{"c"}},

		// Introduction
		"author != user(\"Mallory\") U author = user(\"Alice\")":                                  ltl.Until{ltl.NotEqual{ltl.AP{"author"}, ltl.User{"Mallory"}}, ltl.Equals{ltl.AP{"author"}, ltl.User{"Alice"}}},
		"author = user(\"Alice\") & [] author != user(\"Mallory\")":                               ltl.And{ltl.Equals{ltl.AP{"author"}, ltl.User{"Alice"}}, ltl.Always{ltl.NotEqual{ltl.AP{"author"}, ltl.User{"Mallory"}}}},
		"[] author != user(\"Mallory\") | author != user(\"Mallory\") U author = user(\"Alice\")": ltl.Or{ltl.Always{ltl.NotEqual{ltl.AP{"author"}, ltl.User{"Mallory"}}}, ltl.Until{ltl.NotEqual{ltl.AP{"author"}, ltl.User{"Mallory"}}, ltl.Equals{ltl.AP{"author"}, ltl.User{"Alice"}}}},
	}

	var i int

	for str, golden := range tests {
		name := fmt.Sprintf("test:%v", str)

		t.Run(name, func(t *testing.T) {

			s := scanner.NewFromString(str)
			p := New(s)

			n, err := p.Parse()

			if err != nil {
				t.Fatal(err)
			}

			if !golden.SameAs(n) {
				t.Error(fmt.Printf("got %v, expected %v\n", n, golden))
			}

		})

		i++
	}
}

func TestParserFailed(t *testing.T) {

	f, err := os.Open("./tests/fail.txt")

	if err != nil {
		t.Fatal(err)
	}

	r := bufio.NewScanner(f)

	for r.Scan() {
		sc := scanner.NewFromString(r.Text())
		p := New(sc)

		_, err := p.Parse()

		if err == nil {
			t.Error(fmt.Sprintf("expected failure when parsing: %s", r.Text()))
		}

	}
}

func TestParserSuccess(t *testing.T) {

	f, err := os.Open("./tests/pass.txt")

	if err != nil {
		t.Fatal(err)
	}

	r := bufio.NewScanner(f)

	for r.Scan() {
		sc := scanner.NewFromString(r.Text())
		p := New(sc)

		_, err := p.Parse()

		if err != nil {
			t.Error(fmt.Sprintf("expected failure when parsing: %s, %v", r.Text(), err))
		}

	}

}
