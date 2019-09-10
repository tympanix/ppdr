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
	// [Oa, a]
	// 	[a]	-->	[Oa, a]
	// 	[a]	-->	[!Oa, a]
	// [!Oa, a]
	// 	[a]	-->	[!a, Oa]
	// 	[a]	-->	[!Oa, !a]
	// [!a, Oa]
	// 	[]	-->	[Oa, a]
	// 	[]	-->	[!Oa, a]
	// [!Oa, !a]
	// 	[]	-->	[!a, Oa]
	// 	[]	-->	[!Oa, !a]
}

func ExampleGenerateGNBA_second() {
	phi := ltl.Until{ltl.AP{"a"}, ltl.AP{"b"}}

	g := GenerateGNBA(phi)

	fmt.Println(g)

	//Output:
	// [a, a U b, b]
	// 	[a, b]	-->	[a, a U b, b]
	// 	[a, b]	-->	[!a, a U b, b]
	// 	[a, b]	-->	[!b, a, a U b]
	// 	[a, b]	-->	[!(a U b), !b, a]
	// 	[a, b]	-->	[!(a U b), !a, !b]
	// [!a, a U b, b]
	// 	[b]	-->	[a, a U b, b]
	// 	[b]	-->	[!a, a U b, b]
	// 	[b]	-->	[!b, a, a U b]
	// 	[b]	-->	[!(a U b), !b, a]
	// 	[b]	-->	[!(a U b), !a, !b]
	// [!b, a, a U b]
	// 	[a]	-->	[a, a U b, b]
	// 	[a]	-->	[!a, a U b, b]
	// 	[a]	-->	[!b, a, a U b]
	// [!(a U b), !b, a]
	// 	[a]	-->	[!(a U b), !b, a]
	// 	[a]	-->	[!(a U b), !a, !b]
	// [!(a U b), !a, !b]
	// 	[]	-->	[a, a U b, b]
	// 	[]	-->	[!a, a U b, b]
	// 	[]	-->	[!b, a, a U b]
	// 	[]	-->	[!(a U b), !b, a]
	// 	[]	-->	[!(a U b), !a, !b]
}
