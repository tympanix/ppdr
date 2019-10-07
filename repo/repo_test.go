package repo

import (
	"fmt"
	"testing"

	"github.com/tympanix/master-2019/ltl"
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

type exe struct {
	o op
	r result
}

func (e exe) String() string {
	return fmt.Sprintf("%v", e.o)
}

func TestRepo_one(t *testing.T) {

	r := NewRepo()

	s0 := NewState(ltl.AP{"test"})
	s1 := NewState()
	s2 := NewState()

	s1.addDependency(s0)
	s2.addDependency(s1)

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
	r := NewRepo()

	s0 := NewState(ltl.AP{"Charlie"})
	s1 := NewState(ltl.AP{"Charlie"})
	s2 := NewState(ltl.AP{"David"})
	s3 := NewState(ltl.AP{"Bob"})
	s4 := NewState(ltl.AP{"Mallory"})
	s5 := NewState(ltl.AP{"Alice"})
	s6 := NewState(ltl.AP{"Bob"})

	s1.addDependency(s0)
	s3.addDependency(s1)
	s4.addDependency(s1)
	s4.addDependency(s2)
	s5.addDependency(s3)
	s5.addDependency(s4)
	s6.addDependency(s2)

	phi1 := ltl.Impl{ltl.AP{"Bob"}, ltl.Next{ltl.AP{"Charlie"}}}
	phi2 := ltl.Always{ltl.Not{ltl.AP{"Mallory"}}}
	phi3 := ltl.Impl{ltl.AP{"Bob"}, ltl.Next{ltl.Always{ltl.AP{"Charlie"}}}}

	tests := []exe{
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

	tests := []exe{
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
				t.Errorf("operation failed: %v, error: %v", e, err)
			} else {
				t.Errorf("expected error on op: %v", e.o)
			}
		}
	}
}
