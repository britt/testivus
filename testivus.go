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

	"github.com/pkg/errors"
)

var reportFile = flag.String("testivus.outputfile", "", "write a detailed disappointment report to a file")

// Disappointments are all the ways your code has let you down without
// explicitly failing.
type Disappointments struct {
	sync.Mutex `json:"-"`
	Grievances map[string][]*Disappointment `json:"grievances"`
	Summary    Summary                      `json:"summary"`
}

// Summary is an aggregation of all your disappointments
type Summary struct {
	Total   int
	ByName  map[string]int
	ByTag   map[string]int
	ByError map[string]int

	nameRows  []reportRow
	tagRows   []reportRow
	errorRows []reportRow
}

// MarshalJSON renders the summary to JSON
func (s Summary) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"total":  s.Total,
		"byTag":  s.ByTag,
		"byName": s.ByName,
	}

	if len(s.ByError) > 0 {
		be := map[string]int{}
		for e, c := range s.ByError {
			be[e] = c
		}
		m["byError"] = be
	}

	return json.Marshal(m)
}

// String renders a text representation of your disappointments for the
// airing of grievances.
func (d *Disappointments) String() string {
	d.Lock()
	defer d.Unlock()

	s := d.summarize()
	if s.Total == 0 {
		return "No disapointments, you are truly master of your domain."
	} else if !testing.Verbose() {
		return fmt.Sprintf("I got a lot of problems with you people! (%d disappointments)\n", s.Total)
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "\n=== The airing of grievances:\n")
	fmt.Fprintf(w, "I got a lot of problems with you people! (%d disappointments)\n", s.Total)

	if len(s.tagRows) > 0 {
		fmt.Fprintf(w, "\nBy Tag:\n")
		for _, r := range s.tagRows {
			fmt.Fprintf(w, "\t%s\t%d\t%s\n", r.ID, r.Count, strings.Repeat("|", r.Count))
		}
	}
	w.Flush()

	if len(s.errorRows) > 0 {
		fmt.Fprintf(w, "\nBy Error:\n")
		for _, r := range s.errorRows {
			fmt.Fprintf(w, "\t%s\t%d\t%s\n", r.ID, r.Count, strings.Repeat("|", r.Count))
		}
	}
	w.Flush()

	fmt.Fprintf(w, "\nBy Test:\n")
	for _, r := range s.nameRows {
		fmt.Fprintf(w, "\t%s\t%d\t%s\n", r.ID, r.Count, strings.Repeat("|", r.Count))
	}
	fmt.Fprintf(w, "\n")
	w.Flush()

	return buf.String()
}

type reportRow struct {
	ID    string
	Count int
}

func (d *Disappointments) summarize() Summary {
	s := Summary{}
	count := 0

	// count grievances by tag
	countByTag := make(map[string]int)
	for _, v := range d.Grievances {
		count += len(v)
		for _, g := range v {
			for _, t := range g.Tags {
				countByTag[t] = countByTag[t] + 1
			}
		}
	}
	s.ByTag = countByTag
	for t, c := range countByTag {
		s.tagRows = append(s.tagRows, reportRow{ID: t, Count: c})
	}

	sort.SliceStable(s.tagRows, func(i, j int) bool {
		return s.tagRows[i].Count > s.tagRows[j].Count
	})

	s.Total = count

	// count grievances by name
	countByName := make(map[string]int)
	for _, v := range d.Grievances {
		count += len(v)
		for _, g := range v {
			countByName[g.Name] = countByName[g.Name] + 1
		}
	}
	s.ByName = countByName
	for t, c := range countByName {
		s.nameRows = append(s.nameRows, reportRow{ID: t, Count: c})
	}

	sort.SliceStable(s.nameRows, func(i, j int) bool {
		return s.nameRows[i].Count > s.nameRows[j].Count
	})

	// count grievances by error
	countByError := make(map[string]int)
	for _, v := range d.Grievances {
		for _, g := range v {
			if g.Error != nil {
				countByError[g.Error.Error()] = countByError[g.Error.Error()] + 1
			}
		}
	}
	s.ByError = countByError
	for e, c := range countByError {
		s.errorRows = append(s.errorRows, reportRow{ID: e, Count: c})
	}

	sort.SliceStable(s.errorRows, func(i, j int) bool {
		return s.errorRows[i].Count > s.errorRows[j].Count
	})

	return s
}

// Disappointment is how your code has disappointed you
type Disappointment struct {
	Message string   `json:"message"`
	Tags    []string `json:"tags"`
	Error   error    `json:"error"`
	Name    string   `json:"testName"`
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
		fmt.Println(errors.Wrap(err, "could not save report"))
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

	if reportFile != nil {
		// save output to file
		out, err := os.OpenFile(*reportFile, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer out.Close()

		err = json.NewEncoder(out).Encode(d)
		if err != nil {
			return err
		}
		out.Sync()
	}

	return nil
}

// Grievance registers a Disappointment
func Grievance(t *testing.T, msg string, tags ...string) *Disappointment {
	t.Helper()
	running.Lock()
	defer running.Unlock()

	g := &Disappointment{Name: t.Name(), Message: msg, Tags: tags}
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
