package repo

import (
	"testing"

	"github.com/tympanix/master-2019/ltl"
)

func TestRepo_one(t *testing.T) {

	r := NewRepo()

	s0 := NewState(ltl.AP{"test"})
	s1 := NewState()
	s2 := NewState()

	s1.addDependency(s0)
	s2.addDependency(s1)

	if !r.Put(s0) {
		t.Errorf("could not add state %v", s0)
	}

	if !r.Put(s1) {
		t.Errorf("could not add state %v", s0)
	}

	if !r.Put(s2) {
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
