package scanner

import (
	"fmt"
	"testing"

	"github.com/tympanix/ppdr/ltl/scanner/token"
)

func compareTokens(t *testing.T, s *Scanner, test scannerTest) {
	var i int
	var j int

	for {
		tk := s.NextToken()

		if tk.Kind() != test.tokens[i] {
			t.Fatal("unexpected token")
		}

		if tk.Kind() == token.AP {
			if tk.String() != test.names[j] {
				t.Error(fmt.Sprintf("token %s and %s mismatch", tk.String(), test.names[j]))
			}
			j++
		}

		if tk.Kind() == token.EOF {
			break
		}

		i++
	}

}

type scannerTest struct {
	tokens []token.Kind
	names  []string
}

func test(tokens []token.Kind, names []string) scannerTest {
	return scannerTest{
		tokens: tokens,
		names:  names,
	}
}

func TestScanner_1(t *testing.T) {

	tests := map[string]scannerTest{
		"a and b":                 {[]token.Kind{token.AP, token.AND, token.AP, token.EOF}, []string{"a", "b"}},
		"Oa or !b":                {[]token.Kind{token.NEXT, token.AP, token.OR, token.NOT, token.AP, token.EOF}, []string{"a", "b"}},
		"[]<> crit1 & []<> crit2": {[]token.Kind{token.ALWAYS, token.EVENTUALLY, token.AP, token.AND, token.ALWAYS, token.EVENTUALLY, token.AP, token.EOF}, []string{"crit1", "crit2"}},
		"!(a | b) and (b -> c)":   {[]token.Kind{token.NOT, token.LPAR, token.AP, token.OR, token.AP, token.RPAR, token.AND, token.LPAR, token.AP, token.IMPL, token.AP, token.RPAR, token.EOF}, []string{"a", "b", "b", "c"}},
		"true":                    {[]token.Kind{token.TRUE, token.EOF}, []string{}},
		"a = b":                   {[]token.Kind{token.AP, token.EQUALS, token.AP, token.EOF}, []string{"a", "b"}},
		"\"ok\"":                  {[]token.Kind{token.LITSTRING, token.EOF}, []string{}},
	}

	var i int

	for str, expect := range tests {
		name := fmt.Sprintf("test:%d", i)
		t.Run(name, func(t *testing.T) {
			s := NewFromString(str)
			compareTokens(t, s, expect)
		})
		i++
	}

}
