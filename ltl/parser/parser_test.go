package parser

import (
	"fmt"
	"testing"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/ltl/scanner"
)

func TestParser_1(t *testing.T) {

	tests := map[string]ltl.Node{
		"a or b":                    ltl.Or{ltl.AP{"a"}, ltl.AP{"b"}},
		"a -> b":                    ltl.Impl{ltl.AP{"a"}, ltl.AP{"b"}},
		"[]<> crit1 and []<> crit2": ltl.And{ltl.Always{ltl.Eventually{ltl.AP{"crit1"}}}, ltl.Always{ltl.Eventually{ltl.AP{"crit2"}}}},
	}

	var i int

	for str, golden := range tests {
		name := fmt.Sprintf("test:%d", i)

		t.Run(name, func(t *testing.T) {

			s := scanner.NewFromString(str)
			p := New(s)

			n, err := p.Parse()

			if err != nil {
				t.Fatal(err)
			}

			if !golden.SameAs(n) {
				t.Error(fmt.Printf("got %v, expected %v", n, golden))
			}

		})

		i++
	}
}
