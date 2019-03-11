package testivus_test

import (
	"errors"
	"testing"

	"github.com/britt/testivus"
)

func TestMain(m *testing.M) {
	testivus.Run(m)
}

func TestTestivus(t *testing.T) {
	testivus.Grievance(t, "My son tells me your company stinks!")
	testivus.Grievance(t, "You're slow!", "speed").WithError(errors.New("timeout exceeded"))
	testivus.Grievance(t, "You're send too much data!", "speed", "download")
}
