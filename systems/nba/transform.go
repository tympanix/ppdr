package nba

import (
	"github.com/tympanix/master-2019/debug"
	"github.com/tympanix/master-2019/systems/ba"
	"github.com/tympanix/master-2019/systems/gnba"
)

// TransformGNBAtoNBA takes a GNBA and transforms it into an NBA
func TransformGNBAtoNBA(g *gnba.GNBA) *NBA {

	t := debug.NewTimer("transform")

	defer func() {
		t.Stop()
	}()

	// If final states is the empty set, all states are accepting
	if len(g.FinalStates) == 0 {
		copy := g.Copy()
		return &NBA{
			States:      copy.States,
			StartStates: copy.StartingStates,
			FinalStates: ba.NewStateSet(copy.States...),
			Phi:         copy.Phi,
			AP:          copy.AP,
		}
	}

	// If final states is a singleton set consider GNBA as NBA
	if len(g.FinalStates) == 1 {
		copy := g.Copy()
		return &NBA{
			States:      copy.States,
			StartStates: copy.StartingStates,
			FinalStates: copy.FinalStates[0],
			Phi:         copy.Phi,
			AP:          copy.AP,
		}
	}

	// Else, construct NBA by transformation
	copies := make([]*gnba.GNBA, 0, len(g.FinalStates))
	for range g.FinalStates {
		copies = append(copies, g.Copy())
	}

	// Make the i'th copy have the i'th acceptance set
	for i, c := range copies {
		c.FinalStates = []ba.StateSet{c.FinalStates[i]}
	}

	// Rearrange all transitions to point to the i'th+1 copy
	for i, c := range copies {
		for s := range c.FinalStates[0] {
			for j, t := range s.Transitions {
				next := copies[(i+1)%len(copies)]
				k := c.FindStateIndex(t.State)
				s.Transitions[j] = t.RenameTo(next.States[k])

			}
		}
	}

	return mergeCopiesToNBA(copies)
}

func mergeCopiesToNBA(copies []*gnba.GNBA) *NBA {
	nba := NewNBA(copies[0].Phi)

	states := make([]*ba.State, 0)

	for _, c := range copies {
		states = append(states, c.States...)
	}

	nba.States = states
	nba.StartStates = copies[0].StartingStates
	nba.FinalStates = copies[0].FinalStates[0]
	nba.AP = copies[0].AP

	return nba
}
