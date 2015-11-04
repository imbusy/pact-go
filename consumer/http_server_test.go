package consumer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_MatchingInteractionFound_ReturnsCorrectResponse(t *testing.T) {
	mockHttpServer := newHttpMockService()
	interaction := getFakeInteraction()
	mockHttpServer.ClearInteractions()
	mockHttpServer.RegisterInteraction(interaction)
	url := mockHttpServer.Start()
	defer mockHttpServer.Stop()

	client := &http.Client{}

	req, err := interaction.ToHttpRequest(url)
	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != interaction.Response.Status {
		t.Errorf("The response status is %s, expected %v", resp.Status, interaction.Response.Status)
		t.FailNow()
	}

	defer resp.Body.Close()
	var expectedBody, actualBody interface{}

	if err := json.Unmarshal([]byte(interaction.Response.Body), &expectedBody); err != nil {
		t.Error(err)
		t.FailNow()
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&actualBody); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expectedBody, actualBody) {
		t.Error("The response body is does not match")
	}
}

func Test_MatchingInteractionNotFound_Returns404(t *testing.T) {
	mockHttpServer := newHttpMockService()
	interaction := getFakeInteraction()

	url := mockHttpServer.Start()
	defer mockHttpServer.Stop()

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

	if bodyText != notFoundError.Error() {
		t.Errorf("The expected response was '%s' but recieved '%s'", notFoundError.Error(), bodyText)
	}
}

func getFakeInteraction() *Interaction {
	header := make(http.Header)
	header.Add("content-type", "application/json")
	return &Interaction{
		Request: &ProviderRequest{
			Method:  "GET",
			Path:    "/",
			Query:   "param=xyzmk",
			Body:    `{ "firstName": "John", "lastName": "Doe" }`,
			Headers: header,
		},
		Response: &ProviderResponse{
			Status:  201,
			Headers: header,
			Body:    `{"result": true}`,
		},
	}
}
