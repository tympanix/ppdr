package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/tympanix/ppdr/ltl/scanner/token"
)

var (
	symbols = map[rune]token.Kind{
		'O': token.NEXT,
		'U': token.UNTIL,
		'&': token.AND,
		'|': token.OR,
		'(': token.LPAR,
		')': token.RPAR,
	}
)

// Scanner is able to scan input files
type Scanner struct {
	r   *bufio.Reader
	buf bytes.Buffer
	i   int
}

// NewFromFile creates a new scanner from a file path
func NewFromFile(path string) (*Scanner, error) {

	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return &Scanner{
		r: bufio.NewReader(f),
	}, nil
}

// NewFromReader returns a new scanner from a io.Reader object
func NewFromReader(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

// NewFromString returns a new scanner from a string
func NewFromString(str string) *Scanner {
	return &Scanner{
		r: bufio.NewReader(strings.NewReader(str)),
	}
}

func (s *Scanner) next() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		return 0
	}
	s.buf.WriteRune(r)
	return r
}

func (s *Scanner) nextN(n int) {
	for i := 0; i < n; i++ {
		s.next()
	}
}

func (s *Scanner) has(r rune) bool {
	if s.peekRune() == r {
		s.next()
		return true
	}
	return false
}

func (s *Scanner) see(r rune) bool {
	return s.peekRune() == r
}

func (s *Scanner) clear() {
	s.buf.Reset()
}

func (s *Scanner) hasString(str string) bool {
	if s.peek(len(str)) == str {
		s.nextN(len(str))
		return true
	}
	return false
}

func (s *Scanner) hasLetter() bool {
	if unicode.IsLetter(s.peekRune()) {
		s.next()
		return true
	}
	return false
}

func (s *Scanner) hasDigit() bool {
	if unicode.IsNumber(s.peekRune()) {
		s.next()
		return true
	}
	return false
}

func (s *Scanner) discard() {
	s.r.Discard(1)
}

func (s *Scanner) rune() rune {
	return rune(s.buf.Bytes()[0])
}

func (s *Scanner) peekRune() rune {
	b, err := s.r.Peek(1)
	if err != nil {
		return 0
	}
	return rune(b[0])
}

func (s *Scanner) peek(n int) string {
	str, err := s.r.Peek(n)
	if err != nil {
		return ""
	}
	return string(str)
}

func (s *Scanner) get() string {
	str := s.buf.String()
	s.buf.Reset()
	return str
}

func (s *Scanner) newToken(kind token.Kind) *token.Token {
	return token.New(kind, s.get())
}

func (s *Scanner) unexpectedToken() {
	panic(fmt.Sprintf("unknown token: %s\n", s.get()))
}

func (s *Scanner) expect(r rune) {
	if !s.see(r) {
		s.unexpectedToken()
	}
	s.discard()
}

// NextToken retrieves the next token from the scanner
func (s *Scanner) NextToken() *token.Token {
	for {
		for unicode.IsSpace(s.peekRune()) {
			s.discard()
		}

		if s.hasString("[]") {
			return s.newToken(token.ALWAYS)
		} else if s.hasString("<>") {
			return s.newToken(token.EVENTUALLY)
		} else if s.hasString("&") || s.hasString("and") {
			return s.newToken(token.AND)
		} else if s.hasString("|") || s.hasString("or") {
			return s.newToken(token.OR)
		} else if s.hasString("->") {
			return s.newToken(token.IMPL)
		} else if s.hasString("true") {
			return s.newToken(token.TRUE)
		} else if s.hasString("false") {
			return s.newToken(token.FALSE)
		} else if s.hasString("self") {
			return s.newToken(token.SELF)
		} else if s.has('>') {
			if s.has('=') {
				return s.newToken(token.GTEQ)
			}
			return s.newToken(token.GT)
		} else if s.has('<') {
			if s.has('=') {
				return s.newToken(token.LTEQ)
			}
			return s.newToken(token.LT)
		} else if s.has('=') {
			return s.newToken(token.EQUALS)
		} else if s.has('O') {
			return s.newToken(token.NEXT)
		} else if s.has('U') {
			return s.newToken(token.UNTIL)
		} else if s.has('!') {
			if s.has('=') {
				return s.newToken(token.NEQ)
			}
			return s.newToken(token.NOT)
		} else if s.see('"') {
			s.discard()
			for s.hasLetter() || s.hasDigit() {
				// noop
			}
			s.expect('"')
			return s.newToken(token.LITSTRING)
		} else if s.hasDigit() {
			for s.hasDigit() {
				// noop
			}
			if s.has('.') {
				for s.hasDigit() {
					// noop
				}
			}
			return s.newToken(token.LITNUMBER)
		} else if s.hasLetter() {
			for s.hasLetter() || s.hasDigit() {
				// noop
			}
			return s.newToken(token.AP)
		} else if t, ok := symbols[s.peekRune()]; ok {
			s.next()
			return s.newToken(t)
		} else if s.has(0) {
			return s.newToken(token.EOF)
		} else {
			panic(fmt.Sprintf("unknown token: %c\n", s.next()))
		}
	}
}
