package ltl

// Resolver is an interface for resolving values
type Resolver interface {
	Resolve(string) interface{}
	ResolveBool(string) bool
}

// Satisfiable is an interface that determines if a node is satisfiable
type Satisfiable interface {
	Satisfied(Resolver) bool
}

type mapResolver map[string]interface{}

func (m mapResolver) Resolve(s string) interface{} {
	return m[s]
}

func (m mapResolver) ResolveBool(s string) bool {
	if v, ok := m[s]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// NewResolverFromMap return a new resolver which is based on a lookup table
func NewResolverFromMap(m map[string]interface{}) Resolver {
	return mapResolver(m)
}

type setResolver struct {
	Set
}

func (s setResolver) Resolve(str string) interface{} {
	return s.ResolveBool(str)
}

func (s setResolver) ResolveBool(str string) bool {
	return s.Contains(AP{str})
}

// NewResolverFromSet return a new resolver which is based on a LTL set
func NewResolverFromSet(s Set) Resolver {
	return setResolver{s}
}
