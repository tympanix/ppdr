package repo

import (
	"fmt"
	"testing"

	"github.com/tympanix/ppdr/ltl"
)

type result int

const (
	OK  result = 0
	ERR result = iota
)

type op interface {
	Do(*Repo) error
}

type query struct {
	s   *State
	phi ltl.Node
}

func (q query) String() string {
	return fmt.Sprintf("query(%v, %v)", q.s, q.phi)
}

func (q query) Do(r *Repo) error {
	_, err := r.Query(q.s, q.phi)
	return err
}

type put struct {
	s *State
}

func (p put) String() string {
	return fmt.Sprintf("put(%v)", p.s)
}

func (p put) Do(r *Repo) error {
	return r.Put(p.s)
}

type user struct {
	id *Identity
}

func (u user) String() string {
	return fmt.Sprintf("user(%v)", u.id)
}

func (u user) Do(r *Repo) error {
	r.SetCurrentUser(u.id)
	return nil
}

type exe struct {
	o op
	r result
}

func (e exe) String() string {
	return fmt.Sprintf("%v", e.o)
}

func TestRepo_one(t *testing.T) {

	r := NewRepo()

	s0 := NewState("test", true)
	s1 := NewState()
	s2 := NewState()

	s1.addDependency(s0)
	s2.addDependency(s1)

	u := NewIdentity("john")
	r.SetCurrentUser(u)

	if r.Put(s0) != nil {
		t.Errorf("could not add state %v", s0)
	}

	if r.Put(s1) != nil {
		t.Errorf("could not add state %v", s0)
	}

	if r.Put(s2) != nil {
		t.Errorf("could not add state %v", s0)
	}

	var c *State
	var err error

	phi := ltl.Eventually{ltl.AP{"test"}}

	if c, err = r.Query(s2, phi); err != nil {
		t.Errorf("unexpected error on query: %v", s2)
	}

	if c != s2 {
		t.Errorf("unexpected state, got: %v, expected: %v", c, s2)
	}

	if c, err = r.Query(s2, ltl.Not{phi}); err == nil {
		t.Errorf("expected error on query: %v", s2)
	}

}

func TestRepo_two(t *testing.T) {
	alice := ltl.AP{"Alice"}
	bob := ltl.AP{"Bob"}
	charlie := ltl.AP{"Charlie"}
	david := ltl.AP{"David"}
	mallory := ltl.AP{"Mallory"}
	r := NewRepo()

	s0 := NewState(charlie)
	s1 := NewState(charlie)
	s2 := NewState(david)
	s3 := NewState(bob)
	s4 := NewState(mallory)
	s5 := NewState(alice)
	s6 := NewState(bob)

	s1.addDependency(s0)
	s3.addDependency(s1)
	s4.addDependency(s1)
	s4.addDependency(s2)
	s5.addDependency(s3)
	s5.addDependency(s4)
	s6.addDependency(s2)

	phi1 := ltl.Impl{bob, ltl.Next{charlie}}
	phi2 := ltl.Always{ltl.Not{mallory}}
	phi3 := ltl.Impl{bob, ltl.Next{ltl.Always{charlie}}}
	phi4 := alice
	phi5 := ltl.Or{ltl.Always{mallory}, charlie}

	u := NewIdentity("john")

	tests := []exe{
		// Set current user
		exe{user{u}, OK},

		// Add states
		exe{put{s0}, OK},
		exe{put{s1}, OK},
		exe{put{s2}, OK},
		exe{put{s3}, OK},
		exe{put{s4}, OK},
		exe{put{s5}, OK},
		exe{put{s6}, OK},

		// Bob -> O Charlie
		exe{query{s3, phi1}, OK},
		exe{query{s5, phi1}, OK},

		// [] !Mallory
		exe{query{s3, phi2}, OK},
		exe{query{s5, phi2}, ERR},
		exe{query{s4, phi2}, ERR},

		// Bob -> O[]Charlie
		exe{query{s5, phi3}, OK},
		exe{query{s3, phi3}, OK},
		exe{query{s6, phi3}, ERR},

		// Alice
		exe{query{s0, phi4}, ERR},
		exe{query{s1, phi4}, ERR},
		exe{query{s2, phi4}, ERR},
		exe{query{s3, phi4}, ERR},
		exe{query{s4, phi4}, ERR},
		exe{query{s5, phi4}, OK},
		exe{query{s6, phi4}, ERR},

		// []David \/ Charlie
		exe{query{s0, phi5}, OK},
		exe{query{s1, phi5}, OK},
		exe{query{s2, phi5}, ERR},
		exe{query{s3, phi5}, ERR},
		exe{query{s4, phi5}, ERR},
		exe{query{s5, phi5}, ERR},
		exe{query{s6, phi5}, ERR},
	}

	runTableTest(t, r, tests)

}

// -- Alice --
// -- Bob --
// -- Bob --
// - / - \ --
// Bob -- Charlie
func TestRepo_three(t *testing.T) {
	r := NewRepo()

	alice := ltl.AP{"Alice"}
	bob := ltl.AP{"Bob"}
	charlie := ltl.AP{"Charlie"}

	s0 := NewState(alice)
	s1 := NewState(bob)
	s2 := NewState(bob)

	s3 := NewState(bob)
	s4 := NewState(charlie)

	s1.addDependency(s0)
	s2.addDependency(s1)
	s3.addDependency(s2)
	s4.addDependency(s2)

	phi := ltl.Until{bob, alice}

	u := NewIdentity("john")

	tests := []exe{
		// Set current user
		exe{user{u}, OK},

		// Put
		exe{put{s0}, OK},
		exe{put{s1}, OK},
		exe{put{s2}, OK},
		exe{put{s3}, OK},
		exe{put{s4}, OK},

		// Bob U Alice
		exe{query{s3, phi}, OK},
		exe{query{s4, phi}, ERR},
	}

	runTableTest(t, r, tests)

}

func runTableTest(t *testing.T, r *Repo, tab []exe) {
	for _, e := range tab {
		err := e.o.Do(r)

		if (e.r == OK) != (err == nil) {
			if err != nil {
				t.Fatalf("operation failed: %v, error: %v", e, err)
			} else {
				t.Fatalf("expected error on op: %v", e.o)
			}
		}
	}
}

// s0 {value: "ok", number: true}
// s1 {value: true, number: 5}
// s2 {value: false}
func TestRepo_four(t *testing.T) {

	r := NewRepo()

	s0 := NewState("value", "ok", "number", true)
	s1 := NewState("value", true, "number", 5)
	s2 := NewState("value", false)

	s1.addDependency(s0)
	s2.addDependency(s1)

	eqOk := ltl.Equals{ltl.AP{"value"}, ltl.LitString{"ok"}}
	eq5 := ltl.Equals{ltl.AP{"number"}, ltl.LitNumber{5}}

	phi1 := ltl.Eventually{eqOk}
	phi2 := ltl.Always{eqOk}
	phi3 := ltl.Next{eqOk}
	phi4 := ltl.Until{ltl.Not{eqOk}, eqOk}
	phi5 := eqOk
	phi6 := eq5

	u := NewIdentity("john")

	tests := []exe{
		// Set current user
		exe{user{u}, OK},

		// Put
		exe{put{s0}, OK},
		exe{put{s1}, OK},
		exe{put{s2}, OK},

		// <>(value = "ok")
		exe{query{s0, phi1}, OK},
		exe{query{s1, phi1}, OK},
		exe{query{s2, phi1}, OK},

		// [](value = "ok")
		exe{query{s0, phi2}, OK},
		exe{query{s1, phi2}, ERR},
		exe{query{s2, phi2}, ERR},

		// O(value = "ok")
		exe{query{s0, phi3}, OK},
		exe{query{s1, phi3}, OK},
		exe{query{s2, phi3}, ERR},

		// !(value = "ok") U (value = "ok")
		exe{query{s0, phi4}, OK},
		exe{query{s1, phi4}, OK},
		exe{query{s2, phi4}, OK},

		// value = "ok"
		exe{query{s0, phi5}, OK},
		exe{query{s1, phi5}, ERR},
		exe{query{s2, phi5}, ERR},

		// number = 5
		exe{query{s0, phi6}, ERR},
		exe{query{s1, phi6}, OK},
		exe{query{s2, phi6}, ERR},
	}

	runTableTest(t, r, tests)
}

func TestRepo_five(t *testing.T) {

	r := NewRepo()

	s0 := NewState()
	s1 := NewState(ltl.AP{"a"})
	s1.AddPolicy(ltl.Until{ltl.AP{"a"}, ltl.Self{}})
	s2 := NewState(ltl.AP{"a"})
	s3 := NewState(ltl.AP{"a"})
	s4 := NewState()

	s1.addDependency(s0)
	s2.addDependency(s1)
	s3.addDependency(s2)
	s4.addDependency(s3)

	u := NewIdentity("john")

	tests := []exe{
		// Set current user
		exe{user{u}, OK},

		// Put states into repo
		exe{put{s0}, OK},
		exe{put{s1}, OK},
		exe{put{s2}, OK},
		exe{put{s3}, OK},
		exe{put{s4}, ERR},

		// Queries
		exe{query{s0, ltl.Always{ltl.Self{}}}, OK},
		exe{query{s1, ltl.Always{ltl.Self{}}}, ERR},
	}

	runTableTest(t, r, tests)
}

// Integrity:
//		author = john
//		author = jane
//		author = jack
//		author = subject()
func TestRepo_six(t *testing.T) {

	r := NewRepo()

	s0 := NewState()
	s1 := NewState()

	john := NewIdentity("john")
	jane := NewIdentity("jane")

	tests := []exe{
		exe{user{john}, OK},
		exe{put{s0}, OK},
		exe{user{jane}, OK},
		exe{put{s1}, OK},

		// Queries
		exe{query{s0, ltl.Equals{ltl.AP{"author"}, ltl.User{"john"}}}, OK},
		exe{query{s0, ltl.Equals{ltl.AP{"author"}, ltl.User{"jane"}}}, ERR},
		exe{query{s0, ltl.Equals{ltl.AP{"author"}, ltl.User{"jack"}}}, ERR},
		exe{query{s0, ltl.Equals{ltl.AP{"author"}, ltl.Subject{}}}, ERR},

		exe{query{s1, ltl.Equals{ltl.AP{"author"}, ltl.User{"john"}}}, ERR},
		exe{query{s1, ltl.Equals{ltl.AP{"author"}, ltl.User{"jane"}}}, OK},
		exe{query{s1, ltl.Equals{ltl.AP{"author"}, ltl.User{"jack"}}}, ERR},
		exe{query{s1, ltl.Equals{ltl.AP{"author"}, ltl.Subject{}}}, OK},
	}

	runTableTest(t, r, tests)
}

// Confidentiality: author = subject()
func TestRepo_seven(t *testing.T) {

	r := NewRepo()

	s0 := NewState()
	s0.AddPolicy(ltl.Equals{ltl.AP{"author"}, ltl.Subject{}})
	s1 := NewState()
	s2 := NewState()

	s1.addDependency(s0)
	s2.addDependency(s1)

	john := NewIdentity("john")
	jane := NewIdentity("jane")

	tests := []exe{
		exe{user{john}, OK},
		exe{put{s0}, OK},
		exe{user{jane}, OK},
		exe{put{s1}, OK},
		exe{put{s2}, OK},

		// Query as Jane
		exe{user{jane}, OK},
		exe{query{s0, ltl.True{}}, ERR},
		exe{query{s1, ltl.True{}}, OK},
		exe{query{s2, ltl.True{}}, OK},

		// Query as John
		exe{user{john}, OK},
		exe{query{s0, ltl.True{}}, OK},
		exe{query{s1, ltl.True{}}, ERR},
		exe{query{s2, ltl.True{}}, ERR},
	}

	runTableTest(t, r, tests)
}

// Confidentiality: self -> author == subject()
func TestRepo_eight(t *testing.T) {

	r := NewRepo()

	phi := ltl.Impl{ltl.Self{}, ltl.Equals{ltl.AP{"author"}, ltl.Subject{}}}
	s0 := NewState()
	s0.AddPolicy(phi)
	s1 := NewState()

	s1.addDependency(s0)

	john := NewIdentity("john")
	jane := NewIdentity("jane")

	tests := []exe{
		exe{user{john}, OK},
		exe{put{s0}, OK},
		exe{user{jane}, OK},
		exe{put{s1}, OK},

		// Query as Jane
		exe{user{jane}, OK},
		exe{query{s0, ltl.True{}}, ERR},
		exe{query{s1, ltl.True{}}, OK},

		// Query as John
		exe{user{john}, OK},
		exe{query{s0, ltl.True{}}, OK},
		exe{query{s1, ltl.True{}}, OK},
	}

	runTableTest(t, r, tests)
}

// Confidentiality: author = user("jane")
func TestRepo_nine(t *testing.T) {

	r := NewRepo()

	phi := ltl.Equals{ltl.User{"jane"}, ltl.AP{"author"}}
	s0 := NewState()
	s0.AddPolicy(phi)

	john := NewIdentity("john")
	jane := NewIdentity("jane")

	tests := []exe{
		exe{user{jane}, OK},
		exe{user{john}, OK},
		exe{put{s0}, ERR},

		exe{user{jane}, OK},
		exe{put{s0}, OK},

		exe{user{jane}, OK},
		exe{query{s0, ltl.True{}}, OK},
		exe{user{john}, OK},
		exe{query{s0, ltl.True{}}, OK},
	}

	runTableTest(t, r, tests)
}

// Confidentiality: subject() = user("jane")
func TestRepo_ten(t *testing.T) {

	r := NewRepo()

	phi0 := ltl.Equals{ltl.Subject{}, ltl.User{"jane"}}
	s0 := NewState()
	s0.AddPolicy(phi0)

	phi1 := ltl.Equals{ltl.Subject{}, ltl.User{"john"}}
	s1 := NewState()
	s1.AddPolicy(phi1)

	john := NewIdentity("john")
	jane := NewIdentity("jane")

	tests := []exe{
		exe{user{jane}, OK},
		exe{user{john}, OK},
		exe{put{s0}, ERR},
		exe{put{s1}, OK},

		exe{user{jane}, OK},
		exe{put{s0}, OK},

		exe{user{john}, OK},
		exe{query{s0, ltl.True{}}, ERR},
		exe{query{s1, ltl.True{}}, OK},

		exe{user{jane}, OK},
		exe{query{s0, ltl.True{}}, OK},
		exe{query{s1, ltl.True{}}, ERR},
	}

	runTableTest(t, r, tests)
}

func TestRepo_eleven(t *testing.T) {

	r := NewRepo()

	s0 := NewState(
		"str1", "bbb",
		"str2", "aaa",
		"num1", 7,
		"num2", 3.1415,
		"bol1", true,
		"bol2", false,
	)
	s1 := NewState()

	s1.addDependency(s0)

	str1 := ltl.AP{"str1"}
	str2 := ltl.AP{"str2"}
	num1 := ltl.AP{"num1"}
	num2 := ltl.AP{"num2"}
	bol1 := ltl.AP{"bol1"}
	bol2 := ltl.AP{"bol2"}
	author := ltl.AP{"author"}
	bob := ltl.User{"Bob"}

	alice := NewIdentity("Alice")

	tests := []exe{
		exe{user{alice}, OK},
		exe{put{s0}, OK},
		exe{put{s1}, OK},

		// Compare strings
		exe{query{s1, ltl.Eventually{ltl.Equals{str1, str2}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.NotEqual{str1, str2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.Greater{str1, str2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.GreaterEqual{str1, str2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.GreaterEqual{str1, str1}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.Less{str1, str2}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.LessEqual{str1, str2}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.LessEqual{str1, str1}}}, OK},

		// Compare numbers
		exe{query{s1, ltl.Eventually{ltl.Equals{num1, num2}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.NotEqual{num1, num2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.Greater{num1, num2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.GreaterEqual{num1, num2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.Less{num1, num2}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.LessEqual{num1, num2}}}, ERR},

		// Compare booleans
		exe{query{s1, ltl.Eventually{ltl.Equals{bol1, bol2}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.NotEqual{bol1, bol2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.Greater{bol1, bol2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.GreaterEqual{bol1, bol2}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.Less{bol1, bol2}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.LessEqual{bol1, bol2}}}, ERR},

		// Compare pointers
		exe{query{s1, ltl.Eventually{ltl.Equals{author, bob}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.NotEqual{author, bob}}}, OK},
		exe{query{s1, ltl.Eventually{ltl.Greater{author, bob}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.GreaterEqual{author, bob}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.Less{author, bob}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.LessEqual{author, bob}}}, ERR},

		// Compare different types
		exe{query{s1, ltl.Eventually{ltl.Greater{str1, num1}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.Greater{str1, bol1}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.Greater{num1, bol1}}}, ERR},

		exe{query{s1, ltl.Eventually{ltl.NotEqual{str1, num1}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.NotEqual{str1, bol1}}}, ERR},
		exe{query{s1, ltl.Eventually{ltl.NotEqual{num1, bol1}}}, ERR},
	}

	runTableTest(t, r, tests)
}
