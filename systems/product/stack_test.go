package product

import "testing"

func TestStack(t *testing.T) {

	s := NewStateStack()

	s1 := newState(nil, nil)

	s.Push(s1)
	if !s.Contains(s1) {
		t.Errorf("does not contain %v", s1)
	}

	s.Push(s1)
	if s.Length() != 2 {
		t.Errorf("wrong length, got %d, want %d", s.Length(), 2)
	}

	s.Pop()
	if !s.Contains(s1) {
		t.Errorf("does not contain %v", s1)
	}

	s.Pop()
	if s.Length() != 0 {
		t.Errorf("wrong length, got %d, want %d", s.Length(), 0)
	}

	if s.Contains(s1) {
		t.Errorf("does not contain %v", s1)
	}

	if len(s.set) != 0 {
		t.Errorf("expected set to be empty")
	}
}
