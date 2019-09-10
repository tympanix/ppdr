package gnba

import (
	"fmt"

	"github.com/tympanix/master-2019/ltl"
)

func ExampleGenerateGNBA() {
	phi := ltl.Next{ltl.AP{"a"}}
	g := GenerateGNBA(phi)

	fmt.Println(g)
	//Output:
	// [!a, Oa]
	// 	[]	-->	[Oa, a]
	// 	[]	-->	[!Oa, a]
	// [!Oa, !a]
	// 	[]	-->	[!a, Oa]
	// 	[]	-->	[!Oa, !a]
	// [Oa, a]
	// 	[a]	-->	[Oa, a]
	// 	[a]	-->	[!Oa, a]
	// [!Oa, a]
	// 	[a]	-->	[!a, Oa]
	// 	[a]	-->	[!Oa, !a]
}

func ExampleGenerateGNBA_second() {
	phi := ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}}

	g := GenerateGNBA(phi)

	fmt.Println(g)

	//Output:
	// [!(a U b), !a, !b]
	// 	[]	-->	[!(a U b), !a, !b]
	// 	[]	-->	[a U b, a, !b]
	// 	[]	-->	[!(a U b), a, !b]
	// 	[]	-->	[a U b, !a, b]
	// 	[]	-->	[a U b, a, b]
	// [a, !b, a U b]
	// 	[a]	-->	[!b, a U b, a]
	// 	[a]	-->	[b, a U b, !a]
	// 	[a]	-->	[a U b, a, b]
	// [!(a U b), a, !b]
	// 	[a]	-->	[!(a U b), !a, !b]
	// 	[a]	-->	[!(a U b), a, !b]
	// [a U b, !a, b]
	// 	[b]	-->	[!(a U b), !a, !b]
	// 	[b]	-->	[!b, a U b, a]
	// 	[b]	-->	[!(a U b), a, !b]
	// 	[b]	-->	[a U b, !a, b]
	// 	[b]	-->	[a U b, a, b]
	// [b, a U b, a]
	// 	[a, b]	-->	[!(a U b), !a, !b]
	// 	[a, b]	-->	[a U b, a, !b]
	// 	[a, b]	-->	[!(a U b), a, !b]
	// 	[a, b]	-->	[a U b, !a, b]
	// 	[b, a]	-->	[a U b, a, b]
}
