package main

import (
	"fmt"

	"github.com/tympanix/master-2019/debug"
	"github.com/tympanix/master-2019/ltl"
	"github.com/tympanix/master-2019/systems/gnba"
	"github.com/tympanix/master-2019/systems/nba"
)

func main() {
	phi := ltl.And{ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}}, ltl.Until{ltl.AP{"c"}, ltl.AP{"d"}}}
	g := gnba.GenerateGNBA(phi)
	nba.TransformGNBAtoNBA(g)

	debug.PrintMeasurements()
	fmt.Println("done.")

}
