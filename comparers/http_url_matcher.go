package comparers

import (
	"github.com/imbusy/pact-go/diff"
	"net/url"
)

func pathMatches(expected, actual string) bool {
	if expected != actual {
		return false
	}
	return true
}

func queryMatches(expected, actual url.Values) (bool, diff.Differences) {
	return diff.DeepDiff(expected, actual, &diff.DiffConfig{AllowUnexpectedKeys: false})
}
