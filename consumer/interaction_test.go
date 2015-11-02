package consumer

import (
	"net/url"
	"testing"
)

func Test_InvalidUrl_MappingToHttpRequestFails(t *testing.T) {
	request := &ProviderRequest{}
	interaction := NewInteraction("Some State", "description", request, nil)
	req, err := interaction.ToHttpRequest("bad.url")

	if err == nil || req != nil {
		t.Error("Should not parse invalid url.")
	}
	if req != nil {
		t.Error("Should not return request, as the url is invalid.")
	}
}

func Test_ValidInteraction_MapsToHttpRequest(t *testing.T) {
	interaction := getFakeInteraction()
	baseUrl, _ := url.Parse("http://localhost:52343/")
	req, err := interaction.ToHttpRequest(baseUrl.String())

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
