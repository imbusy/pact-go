package consumer

import (
	"fmt"
	"net/http"
	"testing"
)

var (
	mockServer *httpMockService
)

func init() {
	mockServer = newHttpMockService()
}

func Test_CreatesMockServer(t *testing.T) {
	interaction := getFakeInteraction()
	mockServer.RegisterInteraction(interaction)
	url := mockServer.Start()
	defer mockServer.Stop()

	client := &http.Client{}

	req, err := interaction.ToHttpRequest(url)
	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != interaction.Response.Status {
		t.Errorf("The response status is not correct: %s", resp.Status)
	}

	fmt.Println(resp.Status)

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
		Response: &ProviderResponse{Status: 200},
	}
}
