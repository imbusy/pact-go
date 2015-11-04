package comparers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func MatchRequest(expected, actual *http.Request) (bool, error) {
	if !methodMatches(expected.Method, actual.Method) ||
		!urlMatches(expected.URL, actual.URL) ||
		!containsAllHeaders(expected.Header, actual.Header) {

		return false, nil
	}

	return bodyMatches(expected.Body, actual.Body)
}

func methodMatches(expected, actual string) bool {
	return strings.EqualFold(expected, actual)
}

func urlMatches(expected, actual *url.URL) bool {

	return strings.EqualFold(expected.Path, actual.Path) && strings.EqualFold(expected.RawQuery, actual.RawQuery)
}

func containsAllHeaders(expected, actual map[string][]string) bool {
	if len(expected) > len(actual) {
		return false
	}

	for key, val := range expected {
		if !strings.EqualFold(val[0], actual[key][0]) {
			return false
		}
	}
	return true
}

func bodyMatches(expected, actual io.ReadCloser) (bool, error) {
	if expected == nil {
		return true, nil
	}

	var e, a interface{}
	decoder := json.NewDecoder(expected)
	err := decoder.Decode(&e)
	if err != nil {
		return false, err
	}

	decoder = json.NewDecoder(actual)
	err = decoder.Decode(&a)
	if err != nil {
		return false, err
	}

	expectedVal := e.(map[string]interface{})
	if expectedVal == nil {
		return true, nil
	}

	actualVal := a.(map[string]interface{})
	return reflect.DeepEqual(expectedVal, actualVal), nil
}
