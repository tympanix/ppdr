package parser

import (
	"errors"
	"fmt"

	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/ltl/scanner"
	"github.com/tympanix/master-2019/ltl/scanner/token"
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
		panic(fmt.Sprintf("expected token: %s\n", t.String()))
	}
	return p.last()
}

func (p *Parser) last() *token.Token {
	return p.prev
}

// Parse parses the program
func (p *Parser) Parse() (exp ltl.Node, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	exp = p.parseExpression()
	p.expect(token.EOF)
	return exp, nil
}

func (p *Parser) parseExpression() ltl.Node {
	return p.parseImpl()
}

func (p *Parser) parseImpl() ltl.Node {
	lhs := p.parseOr()

	for {
		if p.have(token.IMPL) {
			lhs = ltl.Impl{lhs, p.parseOr()}
		} else {
			break
		}
	}
	return lhs
}

func (p *Parser) parseOr() ltl.Node {
	lhs := p.parseAnd()

	for {
		if p.have(token.OR) {
			lhs = ltl.Or{lhs, p.parseAnd()}
		} else {
			break
		}
	}
	return lhs
}

func (p *Parser) parseAnd() ltl.Node {
	lhs := p.parseUntil()

	for {
		if p.have(token.AND) {
			lhs = ltl.And{lhs, p.parseUntil()}
		} else {
			break
		}
	}
	return lhs
}

func (p *Parser) parseUntil() ltl.Node {
	lhs := p.parseAtomic()

	for {
		if p.have(token.UNTIL) {
			lhs = ltl.Until{lhs, p.parseAtomic()}
		} else {
			break
		}
	}
	return lhs
}

func (p *Parser) parseAtomic() ltl.Node {
	if p.have(token.ALWAYS) {
		if p.see(token.LPAR) {
			return p.parseParenthesis()
		}
		return ltl.Always{p.parseAtomic()}
	} else if p.have(token.EVENTUALLY) {
		if p.see(token.LPAR) {
			return p.parseParenthesis()
		}
		return ltl.Eventually{p.parseAtomic()}
	} else if p.have(token.LPAR) {
		exp := p.parseExpression()
		p.expect(token.RPAR)
		return exp
	} else if p.have(token.NOT) {
		if p.have(token.LPAR) {
			return p.parseParenthesis()
		}
		p.expect(token.AP)
		return ltl.AP{p.last().String()}
	} else if p.have(token.NEXT) {
		if p.have(token.LPAR) {
			return p.parseParenthesis()
		}
		p.expect(token.AP)
		return ltl.AP{p.last().String()}
	} else if p.have(token.AP) {
		return ltl.AP{p.last().String()}
	} else {
		panic("unexpected token")
	}
}

func (p *Parser) parseParenthesis() ltl.Node {
	p.expect(token.LPAR)
	exp := p.parseExpression()
	p.expect(token.RPAR)
	return exp
}
