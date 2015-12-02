package consumer

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/bennycao/pact-go/provider"
)

func Test_MatchingInteractionFound_ReturnsCorrectResponse(t *testing.T) {
	mockHTTPServer := NewHTTPMockService()
	interaction := getFakeInteraction()
	mockHTTPServer.ClearInteractions()
	mockHTTPServer.RegisterInteraction(interaction)
	url := mockHTTPServer.Start()
	defer mockHTTPServer.Stop()

	client := &http.Client{}

	req, err := interaction.ToHttpRequest(url)
	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != interaction.Response.Status {
		t.Errorf("The response status is %s, expected %v", resp.Status, interaction.Response.Status)
		contents, _ := ioutil.ReadAll(resp.Body)
		t.Log(string(contents))

		t.FailNow()
	}

	defer resp.Body.Close()

	if expectedBody, err := interaction.Response.GetData(); err != nil {
		t.Error(err)
	} else {
		if actualBody, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error(err)
		} else {
			if bytes.Compare(expectedBody, actualBody) != 0 {
				t.Error("The response body does not match")
			} else {
				if err := mockHTTPServer.VerifyInteractions(); err != nil {
					t.Errorf("expected verfication to pass, got error: %s", err.Error())
				}
			}
		}
	}

}

func Test_MatchingInteractionNotFound_Returns404(t *testing.T) {
	mockHTTPServer := NewHTTPMockService()
	interaction := getFakeInteraction()

	url := mockHTTPServer.Start()
	defer mockHTTPServer.Stop()

	client := &http.Client{}

	req, err := interaction.ToHttpRequest(url)
	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("The response status is %s, expected %v", resp.Status, 404)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyText := strings.TrimSpace(string(body)) //had trim, not sure why there are trailing spaces

	if err != nil {
		t.Error(err)
	}

	if bodyText != errNotFound.Error() {
		t.Errorf("The expected response was '%s' but recieved '%s'", errNotFound.Error(), bodyText)
	}
}

func getFakeInteraction() *Interaction {
	header := make(http.Header)
	header.Add("content-type", "application/json")
	i := NewInteraction("description of the interaction",
		"some state",
		provider.NewJsonProviderRequest("GET", "/", "param=xyzmk", header),
		provider.NewJsonProviderResponse(201, header))
	i.Request.SetBody(`{ "firstName": "John", "lastName": "Doe" }`)

	return i
}
