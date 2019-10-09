package ltl

// Resolver is an interface for resolving values
type Resolver interface {
	Resolve(string) Node
	ResolveBool(string) bool
}

// Satisfiable is an interface that determines if a node is satisfiable
type Satisfiable interface {
	Satisfied(Resolver) bool
}

type mapResolver map[string]Node

func (m mapResolver) Resolve(s string) Node {
	return m[s]
}

func (m mapResolver) ResolveBool(s string) bool {
	if v, ok := m[s]; ok {
		if b, ok := v.(LitBool); ok {
			return b.Bool
		}
	}
	return false
}

// NewResolverFromMap return a new resolver which is based on a lookup table
func NewResolverFromMap(m map[string]Node) Resolver {
	return mapResolver(m)
}

type setResolver struct {
	Set
}

func (s setResolver) Resolve(str string) Node {
	return LitBool{s.ResolveBool(str)}
}

func (s setResolver) ResolveBool(str string) bool {
	return s.Contains(AP{str})
}

// NewResolverFromSet return a new resolver which is based on a LTL set
func NewResolverFromSet(s Set) Resolver {
	return setResolver{s}
}
