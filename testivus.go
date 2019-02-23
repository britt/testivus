// Package testivus adds disappointments to go test
package testivus

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

// Disappointments are all the ways your code has let you down without
// explicitly failing.
type Disappointments struct {
	sync.Mutex
	grievances map[string][]string
}

var running *Disappointments

// Run can be used in place of TestMain to allow disappointment reporting
func Run(m *testing.M) {
	running = New(m)
	code := m.Run()
	Report(running)
	os.Exit(code)
}

// New creates a new set of disappointments.
// Use this if you need a custom TestMain. Otherwise you should just use Run.
func New(m *testing.M) *Disappointments {
	return &Disappointments{grievances: make(map[string][]string)}
}

// Report airs your grievances and shows a report of your disappointments.
func Report(d *Disappointments) {
	running.Lock()
	defer running.Unlock()
	var count int
	for _, v := range running.grievances {
		count += len(v)
	}
	fmt.Printf("I gotta lot of problems with you people! (%d disappointments)\n", count)
}

// Disappointment registers a disappointment
func Disappointment(t *testing.T, msg string) {
	running.Lock()
	defer running.Unlock()

	fmt.Println("\tDISAPPOINTMENT:", t.Name(), msg)
	v, ok := running.grievances[t.Name()]
	if !ok {
		running.grievances[t.Name()] = []string{msg}
		return
	}

	v = append(v, msg)
	running.grievances[t.Name()] = v
	return
}
