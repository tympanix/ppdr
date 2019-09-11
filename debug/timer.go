package debug

import (
	"log"
	"time"
)

type Timer struct {
	name  string
	start time.Time
}

func NewTimer(name string) *Timer {
	return &Timer{
		name:  name,
		start: time.Now(),
	}
}

func (t *Timer) Stop() {
	t2 := time.Now()
	log.Printf("%s: %v\n", t.name, t2.Sub(t.start))
}
