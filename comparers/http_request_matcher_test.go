package comparers

import (
	_ "encoding/json"
	"net/http"
	"strings"
	"testing"
)

func Test_MethodIsDifferent_WillNotMatch(t *testing.T) {
	a, _ := http.NewRequest("GET", "", nil)
	b, _ := http.NewRequest("POST", "", nil)

	result, err := MatchRequest(a, b)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result {
		t.Error("The request should not match")
	}
}

func Test_UrlIsDifferent_WillNotMatch(t *testing.T) {
	a, _ := http.NewRequest("GET", "", nil)
	b, _ := http.NewRequest("GET", "/", nil)

	result, err := MatchRequest(a, b)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result {
		t.Error("The request should not match")
	}
}

func Test_ExpectedNoBodyButActualRequestHasBody_WillMatch(t *testing.T) {
	a, _ := http.NewRequest("GET", "/test", nil)
	b, _ := http.NewRequest("GET", "/test", strings.NewReader(`{"name": "John"}`))

	result, err := MatchRequest(a, b)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !result {
		t.Error("The request should match")
	}
}

func Test_BodyIsDifferent_WillNotMatch(t *testing.T) {
	a, err := http.NewRequest("GET", "/test", strings.NewReader(`{"name": "John", "age": 12 }`))
	b, err := http.NewRequest("GET", "/test", strings.NewReader(`{"name": "John"}`))

	result, err := MatchRequest(a, b)
	if result {
		t.Error("The request should not match")
	}

	if err != nil {
		t.Error(err)
	}

}

func Test_HeadersAreMissing_WillNotMatch(t *testing.T) {
	a, err := http.NewRequest("GET", "/test", nil)
	b, err := http.NewRequest("GET", "/test", nil)
	a.Header = make(http.Header)

	a.Header.Add("content-type", "application/json")

	result, err := MatchRequest(a, b)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result {
		t.Error("The request should not match")
	}
}

func Test_HeadersAreDifferent_WillNotMatch(t *testing.T) {
	a, err := http.NewRequest("GET", "/test", nil)
	b, err := http.NewRequest("GET", "/test", nil)
	a.Header = make(http.Header)
	b.Header = make(http.Header)

	a.Header.Add("content-type", "application/json")
	b.Header.Add("content-type", "text/plain")

	result, err := MatchRequest(a, b)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if result {
		t.Error("The request should not match")
	}
}

func Test_AllHeadersFound_WillMatch(t *testing.T) {
	a, err := http.NewRequest("GET", "/test", nil)
	b, err := http.NewRequest("GET", "/test", nil)
	a.Header = make(http.Header)
	b.Header = make(http.Header)

	a.Header.Add("content-type", "application/json")
	b.Header.Add("content-type", "application/json")
	b.Header.Add("extra-header", "value")

	result, err := MatchRequest(a, b)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !result {
		t.Error("The request should match")
	}
}
