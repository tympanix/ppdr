package debug

import (
	"fmt"
	"strings"
	"time"
)

var (
	i       int
	entries []*Timer
)

// Timer contains information about a running timer
type Timer struct {
	name   string
	start  time.Time
	stop   time.Time
	indent int
}

func (t *Timer) duration() time.Duration {
	return t.stop.Sub(t.start)
}

// NewTimer returns a new timer for measuring execution time
func NewTimer(name string) *Timer {
	i++
	return &Timer{
		name:  name,
		start: time.Now(),
	}
}

// Stop stops the timer and print the result
func (t *Timer) Stop() {
	t.stop = time.Now()
	t.indent = i - 1
	entries = append(entries, t)
	i--
}

// PrintMeasurements prints all time measurements to stdout
func PrintMeasurements() {
	for i := range entries {
		entry := entries[len(entries)-i-1]
		fmt.Print(strings.Repeat("  ", entry.indent))
		fmt.Printf("%s: %v\n", entry.name, entry.duration())
	}
}
