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

	s0 := NewState("test", true)
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

func TestRepo_four(t *testing.T) {

	r := NewRepo()

	s0 := NewState("value", "ok")
	s1 := NewState("value", true)
	s2 := NewState()

	s1.addDependency(s0)
	s2.addDependency(s1)

	eq := ltl.Equals{ltl.AP{"value"}, ltl.LitString{"ok"}}

	phi1 := ltl.Eventually{eq}
	phi2 := ltl.Always{eq}
	phi3 := ltl.Next{eq}
	phi4 := ltl.Until{ltl.Not{eq}, eq}
	phi5 := eq

	tests := []exe{
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
	}

	runTableTest(t, r, tests)
}
