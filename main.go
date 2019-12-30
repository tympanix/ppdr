package main

import (
	"fmt"

	"github.com/tympanix/ppdr/debug"
	"github.com/tympanix/ppdr/ltl"
	"github.com/tympanix/ppdr/systems/gnba"
	"github.com/tympanix/ppdr/systems/nba"
)

func main() {
	phi := ltl.And{ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}}, ltl.Until{ltl.AP{"c"}, ltl.AP{"d"}}}
	g := gnba.GenerateGNBA(phi)
	nba.TransformGNBAtoNBA(g)

	debug.PrintMeasurements()
	fmt.Println("done.")

}
