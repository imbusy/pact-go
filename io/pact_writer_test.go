package io

import (
	"net/http"
	"testing"

	"github.com/SEEK-Jobs/pact-go/consumer"
	"github.com/SEEK-Jobs/pact-go/provider"
)

func Test_Writer_InvalidPath_ShouldThrowError(t *testing.T) {
	pact := NewPactFile("consumer", "provider", nil)
	writer := NewPactFileWriter(pact, "//g34/example")

	if err := writer.Write(); err == nil {
		t.Error("expected error as file path is invalid")
	}
}

func Test_Writer_ValidPact_ShouldWritePactFile(t *testing.T) {
	var interactions []*consumer.Interaction
	interactions = append(interactions, getFakeInteraction())

	pact := NewPactFile("consumer", "provider", interactions)
	writer := NewPactFileWriter(pact, "../pact_examples")

	if err := writer.Write(); err != nil {
		t.Error(err)
	}

}

func getFakeInteraction() *consumer.Interaction {
	header := make(http.Header)
	header.Add("content-type", "application/json")
	i, _ := consumer.NewInteraction("description of the interaction",
		"some state",
		provider.NewJSONRequest("POST", "/", "param=xyzmk", header),
		provider.NewJSONResponse(201, header))

	i.Request.SetBody(`{ "firstName": "John", "lastName": "Doe" }`)
	i.Response.SetBody(`{"result": true}`)
	return i
}
