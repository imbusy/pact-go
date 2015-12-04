package comparers

import (
	"github.com/SEEK-Jobs/pact-go/diff"
)

func headerMatches(expected, actual map[string][]string) (bool, diff.Differences) {
	return diff.DeepDiff(expected, actual, nil)
}
