package gnba

// NBA is a structure for non-deterministic Büchi automatons
type NBA struct {
	States          []*State
	StartStates     StateSet
	FinishingStates StateSet
}

// NewNBA returns a new empty NBA
func NewNBA() *NBA {
	return &NBA{
		States:          make([]*State, 0),
		StartStates:     NewStateSet(),
		FinishingStates: NewStateSet(),
	}
}
