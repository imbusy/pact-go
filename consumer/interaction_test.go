package consumer

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/imbusy/pact-go/provider"
)

func Test_InvalidUrl_MappingToHttpRequestFails(t *testing.T) {
	request := &provider.Request{}
	interaction, _ := NewInteraction("Some State", "description", request, nil)
	req, err := interaction.ToHTTPRequest("bad.url")

	if err == nil || req != nil {
		t.Error("Should not parse invalid url.")
	}
	if req != nil {
		t.Error("Should not return request, as the url is invalid.")
	}
}

func Test_ValidRequest_MapsToHttpRequest(t *testing.T) {
	interaction := getFakeInteraction()
	baseUrl, _ := url.Parse("http://localhost:52343/")
	req, err := interaction.ToHTTPRequest(baseUrl.String())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if req == nil {
		t.Error("Should return http request")
	}

	if req.URL.Scheme != baseUrl.Scheme || req.URL.Host != baseUrl.Host {
		t.Error("Url host not mapped correctly")
	}

	if req.URL.Path != interaction.Request.Path {
		t.Error("Url path not mapped correctly")
	}

	if req.URL.RawQuery != interaction.Request.Query {
		t.Error("Url query not mapped correctly")
	}
}

func Test_ValidResponse_WritesToHttpResponse(t *testing.T) {
	interaction := getFakeInteraction()
	rec := httptest.NewRecorder()

	interaction.WriteToHTTPResponse(rec)

	if rec.Code != interaction.Response.Status {
		t.Errorf("Expected status %v, but received %v", interaction.Response.Status, rec.Code)
	}

	respHeader := rec.Header()
	for header, val := range interaction.Response.Headers {
		if val[0] != respHeader.Get(header) {
			t.Errorf("Expected header %s to have %s, but received %s", header, val[0], respHeader.Get(header))
		}
	}

	expectedObj, err := interaction.Response.GetData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	actualObj := rec.Body.Bytes()

	if bytes.Compare(expectedObj, actualObj) != 0 {
		t.Error("Expected body is different from the received body")
	}
}

func getJsonObj(jsonText string) (interface{}, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(jsonText), &obj); err != nil {
		return nil, err
	}
	return obj, nil
}
