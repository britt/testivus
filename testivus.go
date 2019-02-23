// Package testivus adds disappointments to go test
package testivus

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
)

// Disappointments are all the ways your code has let you down without
// explicitly failing.
type Disappointments struct {
	sync.Mutex `json:"-"`
	Grievances map[string][]Disappointment `json:"grievances"`
}

// String renders a text representation of your disappointments for the
// airing of grievances.
func (d *Disappointments) String() string {
	d.Lock()
	defer d.Unlock()

	count, rows := d.summarize()
	if !testing.Verbose() {
		return fmt.Sprintf("I gotta lot of problems with you people! (%d disappointments)\n", count)
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "I gotta lot of problems with you people! (%d disappointments)\n", count)
	for _, r := range rows {
		fmt.Fprintf(w, "\t%s\t%d\t%s\n", r.Tag, r.Count, strings.Repeat("|", r.Count))
	}
	w.Flush()

	return buf.String()
}

type reportRow struct {
	Tag   string
	Count int
}

func (d *Disappointments) summarize() (count int, rows []reportRow) {
	countByTag := make(map[string]int)
	for _, v := range d.Grievances {
		count += len(v)
		for _, g := range v {
			for _, t := range g.Tags {
				countByTag[t] = countByTag[t] + 1
			}
		}
	}

	for t, c := range countByTag {
		rows = append(rows, reportRow{Tag: t, Count: c})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].Count > rows[j].Count
	})
	return count, rows
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
	fmt.Printf(d.String())
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
