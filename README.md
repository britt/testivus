# testivus

[![Documentation](https://godoc.org/github.com/britt/testivus?status.svg)](https://godoc.org/github.com/britt/testivus)

A test helper library for the rest of us. Let your code know how it disappoints you.

## Example Usage

```go
package testivus_test

import (
	"testing"

	"github.com/britt/testivus"
)

func TestMain(m *testing.M) {
	exitCode := testivus.Run(m)
	os.Exit(exitCode)
}

func TestTestivus(t *testing.T) {
	if 1 == 2 {
		t.Fail()
	}

	testivus.Grievance(t, "My son tells me your company stinks!")
	testivus.Grievance(t, "You're slow!", "speed").WithError(err)
	testivus.Grievance(t, "You're send too much data!", "speed", "download")
}
```

### Output

```
> go test -timeout 30s github.com/britt/testivus -v -testivus.outputfile testivus.json

=== RUN   TestTestivus
	DISAPPOINTMENT: My son tells me your company stinks!
	DISAPPOINTMENT: You're slow! (speed)
	DISAPPOINTMENT: You're send too much data! (speed, download)
--- PASS: TestTestivus (0.00s)
PASS

=== The airing of grievances:
I gotta lot of problems with you people! (3 disappointments)

By Tag:
 speed    2 ||
 download 1 |

By Error:
 timeout exceeded 1 |

ok  	github.com/britt/testivus	0.019s
```
