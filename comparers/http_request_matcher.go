package comparers

import (
	"net/http"
)

func MatchRequest(expected, actual *http.Request) (bool, error) {
	if !methodMatches(expected.Method, actual.Method) {
		return false, nil
	} else if !pathMatches(expected.URL.Path, actual.URL.Path) {
		return false, nil
	} else if res, _ := queryMatches(expected.URL.Query(), actual.URL.Query()); !res {
		return false, nil
	} else if res, _ := headerMatches(expected.Header, actual.Header); !res {
		return false, nil
	} else if res, _, err := bodyMatches(expected.Body, actual.Body); !res {
		return false, err
	}
	return true, nil
}
