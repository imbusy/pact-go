package comparers

import (
	"net/url"
)

func pathMatches(expected, actual string) bool {
	return expected == actual
}

func queryMatches(expected, actual url.Values) bool {
	if len(expected) != len(actual) {
		return false
	}

	for key, evals := range expected {
		avals := actual[key]
		if len(evals) != len(avals) {
			return false
		}

		for i := 0; i < len(evals); i++ {
			if evals[i] != avals[i] {
				return false
			}
		}
	}
	return true
}
