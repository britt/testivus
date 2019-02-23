# testivus
A test library for the rest of us. Let your code know how it disappoints you.

> I got a lotta problems with you people!

## Example Usage

```go
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

	testivus.Grievance(t, "My son tells me your company stinks!")
	testivus.Grievance(t, "You're slow!", "speed")
	testivus.Grievance(t, "You're send too much data!", "speed", "download")
}
```