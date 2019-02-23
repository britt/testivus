package testivus_test

import (
	"testing"

	"github.com/britt/testivus"
)

func TestMain(m *testing.M) {
	testivus.Run(m)
}

func TestTestivus(t *testing.T) {
	if 1 == 2 {
		t.Fail()
	}

	testivus.Disappointment(t, "My son tells me your company stinks!")
}
