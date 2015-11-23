package comparers

import (
	"strings"
)

func methodMatches(expected, actual string) bool {
	return strings.EqualFold(expected, actual)
}
