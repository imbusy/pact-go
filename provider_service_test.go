package pact

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/SEEK-Jobs/pact-go/provider"
)

func Test_CanReigsterInteraction_WithValidData(t *testing.T) {
	ps := newMockProviderService(&Config{})

	header := make(http.Header)
	header.Add("content-type", "payload/nuclear")
	request := provider.NewJSONRequest("POST", "/luke", "action=attack", header)
	request.SetBody(`{ "simulation": false, "target": "Death Star" }`)

	response := provider.NewJSONResponse(200, nil)

	if err := ps.Given("Force is strong with Luke Skywalker").
		UponReceiving("Destroy death star").
		With(*request).
		WillRespondWith(*response); err != nil {
		t.Error(err)
	}
}

func Test_CannotReigsterInteraction_WithInvalidData(t *testing.T) {
	ps := newMockProviderService(&Config{})

	request := provider.NewJSONRequest("POST", "/luke", "action=attack", nil)
	response := provider.NewJSONResponse(200, nil)

	if err := ps.Given("Force is strong with Luke Skywalker").
		With(*request).
		WillRespondWith(*response); err == nil {
		t.Error("Should not be able to register interaction with empty description")
	}
}

func Test_CanResetTransientState_AfterRegistration(t *testing.T) {
	ps := newMockProviderService(&Config{})

	header := make(http.Header)
	header.Add("content-type", "payload/nuclear")
	request := provider.NewJSONRequest("POST", "/luke", "action=attack", header)
	request.SetBody(`{ "simulation": false, "target": "Death Star" }`)

	response := provider.NewJSONResponse(200, nil)

	if err := ps.Given("Force is strong with Luke Skywalker").
		UponReceiving("Destroy death star").
		With(*request).
		WillRespondWith(*response); err != nil {
		t.Error(err)
	}

	if ps.state != "" || ps.description != "" || ps.providerRequest != nil || ps.providerResponse != nil {
		t.Error("Provider services transient state is not cleared after interaction registration")
	}
}

func Test_CanClearInteractions(t *testing.T) {
	ps := newMockProviderService(&Config{})

	request := provider.NewJSONRequest("POST", "/luke", "action=attack", nil)
	response := provider.NewJSONResponse(200, nil)

	if err := ps.Given("Force is strong with Luke Skywalker").
		UponReceiving("Destroy death star").
		With(*request).
		WillRespondWith(*response); err != nil {
		t.Error(err)
	}

	ps.ClearInteractions()

	if interactions := ps.service.GetRegisteredInteractions(); len(interactions) > 0 {
		t.Error("Interactions have not been cleared")
	}
	if ps.state != "" || ps.description != "" || ps.providerRequest != nil || ps.providerResponse != nil {
		t.Error("Provider services transient state is not cleared after interaction registration")
	}
}

func Test_CanVerifyInteractions(t *testing.T) {
	ps := newMockProviderService(&Config{})

	request := provider.NewJSONRequest("POST", "/luke", "action=attack", nil)
	response := provider.NewJSONResponse(200, nil)

	if err := ps.Given("Force is strong with Luke Skywalker").
		UponReceiving("Destroy death star").
		With(*request).
		WillRespondWith(*response); err != nil {
		t.Error(err)
	}

	url := ps.start()
	defer ps.stop()

	client := &http.Client{}
	if req, err := http.NewRequest(request.Method, fmt.Sprintf("%s%s?%s", url, request.Path, request.Query), nil); err != nil {
		t.Error(err)
		t.FailNow()
	} else if _, err := client.Do(req); err != nil {
		t.Error(err)
		t.FailNow()
	} else if err := ps.VerifyInteractions(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
