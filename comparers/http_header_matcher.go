package comparers

import (
	"github.com/bennycao/pact-go/diff"
)

func headerMatches(expected, actual map[string][]string) (bool, diff.Differences) {
	return diff.DeepDiff(expected, actual, nil)
}
