package product

import (
	"github.com/tympanix/master-2019/gnba"
	"github.com/tympanix/master-2019/ts"
)

// State is a state in the product transition system of TS and A
type State struct {
	StateTS  *ts.State
	StateNBA *gnba.State
}

type Product struct {
	States []*State
}

func (p *Product) AddState(s *State) {

}
