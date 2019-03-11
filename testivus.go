// Package testivus adds disappointments to go test
package testivus

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"

	"upspin.io/errors"
)

var reportFile string

func init() {
	flag.StringVar(&reportFile, "testivus.outputfile", "", "write a detailed disappointment report to a file")
}

// Disappointments are all the ways your code has let you down without
// explicitly failing.
type Disappointments struct {
	sync.Mutex `json:"-"`
	Grievances map[string][]*Disappointment `json:"grievances"`
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
	fmt.Fprintf(w, "\n=== The airing of grievances:\n")
	fmt.Fprintf(w, "I gotta lot of problems with you people! (%d disappointments)\n", count)
	for _, r := range rows {
		fmt.Fprintf(w, "\t%s\t%d\t%s\n", r.Tag, r.Count, strings.Repeat("|", r.Count))
	}
	fmt.Fprintf(w, "\n")
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
	Error   error    `json:"error"`
}

func (d Disappointment) String() string {
	if len(d.Tags) == 0 {
		return d.Message
	}

	t := strings.Join(d.Tags, ", ")
	return fmt.Sprintf("%s (%s)", d.Message, t)
}

// WithMessage sets the message on the disappointment
func (d *Disappointment) WithMessage(msg string) *Disappointment {
	d.Message = msg
	return d
}

// WithError adds an error to the disappointment
func (d *Disappointment) WithError(err error) *Disappointment {
	d.Error = err
	return d
}

// WithTags appends the given tags to the disappointment
func (d *Disappointment) WithTags(tags ...string) *Disappointment {
	d.Tags = append(d.Tags, tags...)
	return d
}

var running *Disappointments

// Run can be used in place of TestMain to allow disappointment reporting
func Run(m *testing.M) {
	flag.Parse()
	running = New(m)
	code := m.Run()
	err := Report(running)
	if err != nil {
		fmt.Println(errors.E(err, "could not save report"))
		os.Exit(1)
	}
	os.Exit(code)
}

// New creates a new set of disappointments.
// Use this only if you need a custom TestMain. Otherwise you should just use Run.
func New(m *testing.M) *Disappointments {
	return &Disappointments{Grievances: make(map[string][]*Disappointment)}
}

// Report airs your grievances and shows a report of your disappointments.
// Use this only if you need a custom TestMain. Otherwise you should just use Run.
func Report(d *Disappointments) error {
	fmt.Printf(d.String())

	if reportFile != "" {
		// save output to file
		out, err := os.OpenFile(reportFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer out.Close()

		err = json.NewEncoder(out).Encode(d)
		if err != nil {
			return err
		}
	}

	return nil
}

// Grievance registers a Disappointment
func Grievance(t *testing.T, msg string, tags ...string) *Disappointment {
	t.Helper()
	running.Lock()
	defer running.Unlock()

	g := &Disappointment{Message: msg, Tags: tags}
	if testing.Verbose() {
		fmt.Println("\tDISAPPOINTMENT:", g)
	}

	v, ok := running.Grievances[t.Name()]
	if !ok {
		running.Grievances[t.Name()] = []*Disappointment{g}
		return g
	}

	v = append(v, g)
	running.Grievances[t.Name()] = v
	return g
}
