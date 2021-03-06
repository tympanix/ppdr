package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tympanix/ppdr/ltl"
	"github.com/tympanix/ppdr/ltl/scanner"
	"github.com/tympanix/ppdr/ltl/scanner/token"
)

// Parser parses the input program from a scanner
type Parser struct {
	s      *scanner.Scanner
	tokens []*token.Token
	prev   *token.Token
	i      int
}

// New return a new parser
func New(s *scanner.Scanner) *Parser {
	return &Parser{s: s}
}

func (p *Parser) pump(n int) {
	for len(p.tokens) < n {
		next := p.s.NextToken()
		p.tokens = append(p.tokens, next)
	}
}

func (p *Parser) current() *token.Token {
	p.pump(1)
	return p.tokens[0]
}

func (p *Parser) pop() {
	if len(p.tokens) > 0 {
		p.prev = p.tokens[0]
		p.tokens = p.tokens[1:]
	}
	p.pump(1)
}

func (p *Parser) peek(n int) *token.Token {
	p.pump(n + 1)
	return p.tokens[n]
}

func (p *Parser) have(t token.Kind) bool {
	p.pump(1)
	e := p.current()

	if e.Kind() == t {
		p.pop()
	}

	return e.Kind() == t
}

func (p *Parser) see(t token.Kind) bool {
	return p.current().Kind() == t
}

func (p *Parser) expect(t token.Kind) *token.Token {
	if !p.have(t) {
		panic(fmt.Sprintf("unexpected token: %s, expected: %s\n", p.current().String(), t.String()))
	}
	return p.last()
}

func (p *Parser) last() *token.Token {
	return p.prev
}

func (p *Parser) unexpectedToken() {
	panic(fmt.Sprintf("unexpected token: %v", p.current().Kind()))
}

// Parse parses the program
func (p *Parser) Parse() (exp ltl.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	exp = p.parseTopLevel()
	p.expect(token.EOF)
	return exp, nil
}

func (p *Parser) parseTopLevel() ltl.Node {
	return p.parseImpl()
}

func (p *Parser) parseImpl() ltl.Node {
	lhs := p.parseOr()

	if p.have(token.IMPL) {
		return ltl.Impl{lhs, p.parseImpl()}
	}
	return lhs
}

func (p *Parser) parseOr() ltl.Node {
	lhs := p.parseAnd()

	if p.have(token.OR) {
		return ltl.Or{lhs, p.parseOr()}
	}

	return lhs
}

func (p *Parser) parseAnd() ltl.Node {
	lhs := p.parseUntil()

	if p.have(token.AND) {
		return ltl.And{lhs, p.parseAnd()}
	}

	return lhs
}

func (p *Parser) parseUntil() ltl.Node {
	lhs := p.parseAtomic()

	if p.have(token.UNTIL) {
		return ltl.Until{lhs, p.parseUntil()}
	}

	return lhs
}

func (p *Parser) seeExpression() bool {
	if p.seeLiteral() || p.seeFunction() || p.see(token.AP) {
		switch p.peek(1).Kind() {
		case token.EQUALS, token.NEQ, token.GT, token.GTEQ, token.LT, token.LTEQ:
			return true
		}
	}

	return false
}

func (p *Parser) parseExpressionArg() ltl.Node {
	if p.seeLiteral() {
		return p.parseLiteral()
	} else if p.seeFunction() {
		return p.parseFunction()
	} else if p.have(token.AP) {
		return ltl.AP{p.last().String()}
	}
	p.unexpectedToken()
	return nil
}

func (p *Parser) parseExpression() ltl.Node {
	lhs := p.parseExpressionArg()
	if p.have(token.EQUALS) {
		return ltl.Equals{lhs, p.parseExpressionArg()}
	} else if p.have(token.NEQ) {
		return ltl.NotEqual{lhs, p.parseExpressionArg()}
	} else if p.have(token.GT) {
		return ltl.Greater{lhs, p.parseExpressionArg()}
	} else if p.have(token.GTEQ) {
		return ltl.GreaterEqual{lhs, p.parseExpressionArg()}
	} else if p.have(token.LT) {
		return ltl.Less{lhs, p.parseExpressionArg()}
	} else if p.have(token.LTEQ) {
		return ltl.LessEqual{lhs, p.parseExpressionArg()}
	}
	p.unexpectedToken()
	return nil
}

func (p *Parser) parseAtomic() ltl.Node {
	if p.have(token.ALWAYS) {
		return ltl.Always{p.parseAtomic()}
	} else if p.have(token.EVENTUALLY) {
		return ltl.Eventually{p.parseAtomic()}
	} else if p.seeExpression() {
		return p.parseExpression()
	} else if p.see(token.LPAR) {
		return p.parseParenthesis()
	} else if p.have(token.NOT) {
		return ltl.Not{p.parseAtomic()}
	} else if p.have(token.NEXT) {
		return ltl.Next{p.parseAtomic()}
	} else if p.have(token.SELF) {
		return ltl.Self{}
	} else if p.have(token.TRUE) {
		return ltl.True{}
	} else if p.have(token.AP) {
		return ltl.AP{p.last().String()}
	}
	p.unexpectedToken()
	return nil
}

func (p *Parser) seeFunction() bool {
	if p.peek(0).Kind() == token.AP && p.peek(1).Kind() == token.LPAR {
		return true
	}
	return false
}

func (p *Parser) parseFunction() ltl.Node {
	var l ltl.Node
	p.expect(token.AP)
	name := p.last().String()
	p.expect(token.LPAR)
	if !p.see(token.RPAR) {
		l = p.parseLiteral()
	}
	p.expect(token.RPAR)

	if name == "subject" && l == nil {
		return ltl.Subject{}
	} else if s, ok := l.(ltl.LitString); name == "user" && ok {
		return ltl.User{s.Str}
	}
	p.unexpectedToken()
	return nil
}

func (p *Parser) seeLiteral() bool {
	switch p.current().Kind() {
	case token.LITSTRING, token.LITNUMBER, token.FALSE, token.TRUE:
		return true
	}
	return false
}

func (p *Parser) parseLiteral() ltl.Node {
	if p.have(token.TRUE) {
		return ltl.LitBool{true}
	} else if p.have(token.FALSE) {
		return ltl.LitBool{false}
	} else if p.have(token.LITNUMBER) {
		f, err := strconv.ParseFloat(p.last().String(), 64)
		if err != nil {
			panic(err)
		}
		return ltl.LitNumber{f}
	} else if p.have(token.LITSTRING) {
		return ltl.LitString{p.last().String()}
	}
	p.unexpectedToken()
	return nil
}

func (p *Parser) parseParenthesis() ltl.Node {
	p.expect(token.LPAR)
	exp := p.parseTopLevel()
	p.expect(token.RPAR)
	return exp
}
