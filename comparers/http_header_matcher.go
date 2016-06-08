package comparers

import (
	"github.com/imbusy/pact-go/diff"
)

func headerMatches(expected, actual map[string][]string) (bool, diff.Differences) {
	if expected == nil {
		return true, nil
	}
	return diff.DeepDiff(expected, actual, &diff.DiffConfig{AllowUnexpectedKeys: true, RootPath: "[\"header\"]"})
}
