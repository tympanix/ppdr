package ltl

// User is a function for retrieving user identities
type User struct {
	Name string
}

// SameAs returns true if node is also true
func (u User) SameAs(node Node) bool {
	_, ok := node.(User)
	return ok
}

func (u User) Normalize() Node {
	return u
}

func (u User) Compile(m *RefTable) Node {
	panic("user can not be compiled")
}

func (u User) String() string {
	return "user"
}

func (u User) Len() int {
	return 0
}

func (u User) Map(fn MapFunc) Node {
	return fn(u)
}

func (u User) Satisfied(r Resolver) bool {
	panic("user can not be checked for satisfyability")
}
