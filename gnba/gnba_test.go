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
	// [Oa !a]
	// 	[]	->	[Oa a]
	// 	[]	->	[!Oa a]
	// [!Oa !a]
	// 	[]	->	[Oa !a]
	// 	[]	->	[!Oa !a]
	// [Oa a]
	// 	[a]	->	[Oa a]
	// 	[a]	->	[!Oa a]
	// [!Oa a]
	// 	[a]	->	[Oa !a]
	// 	[a]	->	[!Oa !a]
}
