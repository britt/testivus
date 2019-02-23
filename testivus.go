// Package testivus adds disappointments to go test
package testivus

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

// Disappointments are all the ways your code has let you down without
// explicitly failing.
type Disappointments struct {
	sync.Mutex `json:"-"`
	Grievances map[string][]Disappointment `json:"grievances"`
}

// Disappointment is how your code has disappointed you
type Disappointment struct {
	Message string   `json:"message"`
	Tags    []string `json:"tags"`
}

func (d Disappointment) String() string {
	if len(d.Tags) == 0 {
		return d.Message
	}

	t := strings.Join(d.Tags, ", ")
	return fmt.Sprintf("%s (%s)", d.Message, t)
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
// Use this only if you need a custom TestMain. Otherwise you should just use Run.
func New(m *testing.M) *Disappointments {
	return &Disappointments{Grievances: make(map[string][]Disappointment)}
}

// Report airs your grievances and shows a report of your disappointments.
// Use this only if you need a custom TestMain. Otherwise you should just use Run.
func Report(d *Disappointments) {
	running.Lock()
	defer running.Unlock()
	var count int
	for _, v := range running.Grievances {
		count += len(v)
	}
	fmt.Printf("I gotta lot of problems with you people! (%d disappointments)\n", count)
}

// Grievance registers a Disappointment
func Grievance(t *testing.T, msg string, tags ...string) {
	running.Lock()
	defer running.Unlock()

	g := Disappointment{Message: msg, Tags: tags}
	if testing.Verbose() {
		fmt.Println("\tDISAPPOINTMENT:", g)
	}

	v, ok := running.Grievances[t.Name()]
	if !ok {
		running.Grievances[t.Name()] = []Disappointment{g}
		return
	}

	v = append(v, g)
	running.Grievances[t.Name()] = v
	return
}
