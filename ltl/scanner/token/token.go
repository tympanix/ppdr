//go:generate stringer -type Kind

package token

// New returns a new token with given type and textual represenetation
func New(kind Kind, repr string) *Token {
	return &Token{
		kind,
		repr,
	}
}

// Token represents a token kind with its textual representation
type Token struct {
	kind Kind
	repr string
}

// String returns the textual representation of the token
func (t Token) String() string {
	return t.repr
}

// Kind returns the kind of token
func (t Token) Kind() Kind {
	return t.kind
}

// Kind reprensents a token from the parser
type Kind int

const (
	EOF Kind = iota
	EQUALS
	GT
	GTEQ
	LT
	LTEQ
	NEQ
	AP
	EVENTUALLY
	ALWAYS
	UNTIL
	SELF
	NEXT
	AND
	OR
	IMPL
	LPAR
	RPAR
	NOT
	TRUE
	FALSE
	LITSTRING
	LITNUMBER
	LITBOOL
)
