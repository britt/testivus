package testivus_test

import (
	"errors"
	"os"
	"testing"

	"github.com/britt/testivus"
)

func TestMain(m *testing.M) {
	code := testivus.Run(m)
	os.Exit(code)
}

func TestTestivus(t *testing.T) {
	testivus.Grievance(t, "My son tells me your company stinks!")
	testivus.Grievance(t, "You're slow!", "speed").WithError(errors.New("timeout exceeded"))
	testivus.Grievance(t, "You're send too much data!", "speed", "download")
}
