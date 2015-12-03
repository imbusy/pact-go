package pact

import (
	"net/http"
	"testing"

	"github.com/bennycao/pact-go/provider"
)

func Test_CanBuild(t *testing.T) {
	builder := NewConsumerPactBuilder(&Config{PactPath: "./pact_examples"}).
		ServiceConsumer("browser").
		HasPactWith("api")
	ps, _ := builder.GetMockProviderService()

	request := provider.NewJsonProviderRequest("GET", "/user", "id=23", nil)
	header := make(http.Header)
	header.Add("content-type", "application/json")
	response := provider.NewJsonProviderResponse(200, header)
	response.SetBody(`{ "id": 23, "firstName": "John", "lastName": "Doe" }`)

	ps.Given("there is a user with id {23}").
		UponReceiving("get request for user with id {23}").
		With(*request).
		WillRespondWith(*response)

	request.Query = "id=200"
	response.Status = 404
	response.Headers = nil
	response.Clear()

	ps.Given("there is no user with id {200}").
		UponReceiving("get request for user with id {200}").
		With(*request).
		WillRespondWith(*response)

	if err := builder.Build(); err != nil {
		t.Error(err)
	}
}
