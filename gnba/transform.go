package gnba

// TransformGNBAtoNBA takes a GNBA and transforms it into an NBA
func TransformGNBAtoNBA(gnba *GNBA) *NBA {

	// If final states is the empty set, all states are accepting
	if len(gnba.FinalStates) == 0 {
		copy := gnba.Copy()
		return &NBA{
			States:      copy.States,
			StartStates: copy.StartingStates,
			FinalStates: NewStateSet(copy.States...),
		}
	}

	// If final states is a singleton set consider GNBA as NBA
	if len(gnba.FinalStates) == 1 {
		copy := gnba.Copy()
		return &NBA{
			States:      copy.States,
			StartStates: copy.StartingStates,
			FinalStates: copy.FinalStates[0],
		}
	}

	// Else, construct NBA by transformation
	copies := make([]*GNBA, 0, len(gnba.FinalStates))
	for range gnba.FinalStates {
		copies = append(copies, gnba.Copy())
	}

	// Make the i'th copy have the i'th acceptance set
	for i, c := range copies {
		c.FinalStates = []StateSet{c.FinalStates[i]}
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

func mergeCopiesToNBA(copies []*GNBA) *NBA {
	nba := NewNBA()

	states := make([]*State, 0)

	for _, c := range copies {
		states = append(states, c.States...)
	}

	nba.States = states
	nba.StartStates = copies[0].StartingStates
	nba.FinalStates = copies[0].FinalStates[0]

	return nba
}

type renameTable map[*State]*State
