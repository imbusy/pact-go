package comparers

import (
	"bytes"
	"github.com/imbusy/pact-go/provider"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type matchResponseTest struct {
	exp       *provider.Response
	act       *http.Response
	diffCount int
	diffPath  string
}

var (
	matchResponseTestData = []*matchResponseTest{
		&matchResponseTest{buildTestProviderResponse(200, nil, ""), &http.Response{StatusCode: 201}, 1, "[\"status\"]"},
		&matchResponseTest{buildTestProviderResponse(200, http.Header{"Content-Type": {"application/json"}}, ""), &http.Response{StatusCode: 200}, 1, "[\"header\"]"},
		&matchResponseTest{buildTestProviderResponse(200, nil, `{"name":"John Doe"}`), &http.Response{StatusCode: 200}, 1, "[\"body\"]"},
		&matchResponseTest{buildTestProviderResponse(200, http.Header{"Content-Type": {"application/json"}}, `{"name":"John Doe"}`), buildTestHttpResponse(200, http.Header{"Content-Type": {"application/json"}}, `{"name":"John Doe"}`), 0, ""},
	}
)

func buildTestProviderResponse(status int, h http.Header, body string) *provider.Response {
	pr := provider.NewJSONResponse(200, h)
	pr.SetBody(body)
	return pr
}

func buildTestHttpResponse(status int, h http.Header, body string) *http.Response {
	return &http.Response{
		StatusCode:    status,
		Header:        h,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
	}
}

func Test_MatchResponse_Scenarios(t *testing.T) {
	for _, test := range matchResponseTestData {
		diff, err := MatchResponse(test.exp, test.act)
		if err != nil {
			t.Error(err)
		}
		if len(diff) != test.diffCount {
			t.Errorf("expected diff count to be %d, but actual count id %d", test.diffCount, len(diff))
		} else if len(diff) > 0 && !strings.Contains(diff.Error(), test.diffPath) {
			t.Errorf("expected diff at %s, but got %s", test.diffPath, diff.Error())
		}
	}
}
