package gnba

// TransformGNBAtoNBA takes a GNBA and transforms it into an NBA
func TransformGNBAtoNBA(gnba *GNBA) *NBA {
	nba := NewNBA()

	// TODO: check if final states is the empty set
	// - Then all states are acceptance states

	// TODO: check if final states is a singleton set
	// - Then the GNBA can be considered directly as an NBA

	return nba
}

type renameTable map[*State]*State
