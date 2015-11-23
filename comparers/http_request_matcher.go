package comparers

import (
	"net/http"
)

func MatchRequest(expected, actual *http.Request) (bool, error) {
	if !methodMatches(expected.Method, actual.Method) ||
		!pathMatches(expected.URL.Path, actual.URL.Path) ||
		!queryMatches(expected.URL.Query(), actual.URL.Query()) ||
		!containsAllHeaders(expected.Header, actual.Header) {

		return false, nil
	}

	return bodyMatches(expected.Body, actual.Body)
}
