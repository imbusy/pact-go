package diff

import (
	"github.com/SEEK-Jobs/pact-go/util"
)

const (
	redColor = "\x1b[31m"
	reset    = "\x1b[0m"
)

var (
	diffMisMatchAt = "mismatch at %s: %s"
	diffExpected   = "expected:"
	diffVal        = "\t%#v"
	diffActual     = "recieved:"
	diffHeading    = "%s%s%s"
)

func FormatDiff(diffs Differences, l util.Logger, heading string) {
	if heading != "" {
		l.Printf(diffHeading, redColor, heading, reset)
	}

	for _, d := range diffs {
		l.Printf(diffMisMatchAt, d.path, d.how)
		l.Println(diffExpected)
		l.Printf(diffVal, interfaceOf(d.v1))
		l.Println(diffActual)
		l.Printf(diffVal, interfaceOf(d.v2))
	}
}
