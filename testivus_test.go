package testivus_test

import (
	"testing"

	"github.com/britt/testivus"
)

func TestMain(m *testing.M) {
	testivus.Run(m)
}

func TestTestivus(t *testing.T) {
	testivus.Grievance(t, "My son tells me your company stinks!")
	testivus.Grievance(t, "You're slow!", "speed")
	testivus.Grievance(t, "You're send too much data!", "speed", "download")
}
