package pact

import (
	"net/http"
	"testing"

	"github.com/SEEK-Jobs/pact-go/provider"
)

func Test_CanBuild(t *testing.T) {
	builder := NewConsumerPactBuilder(&Config{PactPath: "./pact_examples"}).
		ServiceConsumer("browser").
		HasPactWith("api")
	ps, _ := builder.GetMockProviderService()

	request := provider.NewJSONRequest("GET", "/user", "id=23", nil)
	header := make(http.Header)
	header.Add("content-type", "application/json")
	response := provider.NewJSONResponse(200, header)
	response.SetBody(`{ "id": 23, "firstName": "John", "lastName": "Doe" }`)

	if err := ps.Given("there is a user with id {23}").
		UponReceiving("get request for user with id {23}").
		With(*request).
		WillRespondWith(*response); err != nil {
		t.Error(err)
		t.FailNow()
	}

	request.Query = "id=200"
	response.Status = 404
	response.Headers = nil
	response.ResetContent()

	ps.Given("there is no user with id {200}").
		UponReceiving("get request for user with id {200}").
		With(*request).
		WillRespondWith(*response)

	if err := builder.Build(); err != nil {
		t.Error(err)
	}
}

func Test_CannotBuild_WhenThereIsNoConsumer(t *testing.T) {
	builder := NewConsumerPactBuilder(&Config{PactPath: "./"}).
		HasPactWith("serviceprovider")
	if err := builder.Build(); err != nil {
		if err != errInvalidConsumer {
			t.Error("expected invalid consumer error")
		}
	} else {
		t.Error("should not build without consumer")
	}
}

func Test_CannotBuild_WhenThereIsNoProvider(t *testing.T) {
	builder := NewConsumerPactBuilder(&Config{PactPath: "./"}).
		ServiceConsumer("consumer")
	if err := builder.Build(); err != nil {
		if err != errInvalidProvider {
			t.Error("expected invalid provider error")
		}
	} else {
		t.Error("should not build without consumer")
	}
}

func Test_CannotBuild_WhenPactCannotBePersisted(t *testing.T) {
	builder := NewConsumerPactBuilder(&Config{PactPath: "//3434"}).
		ServiceConsumer("consumer").HasPactWith("serviceProvider")
	if err := builder.Build(); err == nil {
		t.Error("should not build without consumer")
	}
}
