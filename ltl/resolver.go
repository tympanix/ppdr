package ltl

// Resolver is an interface for resolving values
type Resolver interface {
	Resolve(string) Node
	ResolveBool(string) bool
	ResolveRef(Ref) bool
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

func (m mapResolver) ResolveRef(r Ref) bool {
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

func (s setResolver) ResolveRef(r Ref) bool {
	return s.Contains(r)
}

// NewResolverFromSet return a new resolver which is based on a LTL set
func NewResolverFromSet(s Set) Resolver {
	return setResolver{s}
}

type combinedResolver struct {
	r1 Resolver
	r2 Resolver
}

func (c combinedResolver) Resolve(str string) Node {
	if v := c.r1.Resolve(str); v != nil {
		return v
	}
	return c.r2.Resolve(str)
}

func (c combinedResolver) ResolveBool(str string) bool {
	return c.r1.ResolveBool(str) || c.r2.ResolveBool(str)
}

func (c combinedResolver) ResolveRef(r Ref) bool {
	return c.r1.ResolveRef(r) || c.r2.ResolveRef(r)
}

// NewResolverCombined return a new resolver composed by combining two
// individual resolvers. The first resolver, r1, has priority
func NewResolverCombined(r1 Resolver, r2 Resolver) Resolver {
	return combinedResolver{
		r1: r1,
		r2: r2,
	}
}
