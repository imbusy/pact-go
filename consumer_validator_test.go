package pact

import (
	"errors"
	"github.com/imbusy/pact-go/consumer"
	"github.com/imbusy/pact-go/io"
	"github.com/imbusy/pact-go/provider"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Validator_IsInAStateToValidate(t *testing.T) {
	v := newConsumerValidator(nil, nil, DefaultLogger)

	if err := v.CanValidate(); err == nil || err != errNilProviderClient {
		t.Errorf("expected %s", errNilProviderClient)
	}

	v.ProviderService(&http.Client{}, nil)

	if err := v.CanValidate(); err == nil || err != errNilProviderURL {
		t.Errorf("expected %s", errNilProviderURL)
	}
}

func Test_Validator_ReturnsErrorWhenRequestCreationFails(t *testing.T) {
	v := newConsumerValidator(nil, nil, DefaultLogger)
	interaction, _ := consumer.NewInteraction("description", "state", provider.NewJSONRequest("Get", "/", "", nil), provider.NewJSONResponse(200, nil))
	f := io.NewPactFile("consumer", "provider", []*consumer.Interaction{interaction})
	sa := &stateAction{setup: nil, teardown: nil}

	v.ProviderService(&http.Client{}, &url.URL{})
	if _, err := v.Validate(f, map[string]*stateAction{"state": sa}); err == nil {
		t.Errorf("expected error whilst creating the request")
	}
}

func Test_Validator_ReturnsErrorWhenRequestFails(t *testing.T) {
	v := newConsumerValidator(nil, nil, DefaultLogger)
	interaction, _ := consumer.NewInteraction("description", "state", provider.NewJSONRequest("Get", "/", "", nil), provider.NewJSONResponse(200, nil))
	f := io.NewPactFile("consumer", "provider", []*consumer.Interaction{interaction})
	sa := &stateAction{setup: nil, teardown: nil}
	u, _ := url.Parse("http://localhost:54322")

	v.ProviderService(&http.Client{}, u)
	if _, err := v.Validate(f, map[string]*stateAction{"state": sa}); err == nil {
		t.Errorf("expected error whilst making request")
	}
}

func Test_Validator_ReturnsErrorFromResponseMatcher(t *testing.T) {
	v := newConsumerValidator(nil, nil, DefaultLogger)
	r := provider.NewJSONResponse(200, nil)
	r.SetBody(`{"name":"John Doe"}`)
	interaction, _ := consumer.NewInteraction("description", "state", provider.NewJSONRequest("Get", "/", "", nil), r)
	f := io.NewPactFile("consumer", "provider", []*consumer.Interaction{interaction})
	sa := &stateAction{setup: nil, teardown: nil}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(interaction.Response.Status)
		w.Write([]byte("bad json"))
	}))
	defer s.Close()
	u, _ := url.Parse(s.URL)

	v.ProviderService(&http.Client{}, u)
	if _, err := v.Validate(f, map[string]*stateAction{"state": sa}); err == nil {
		t.Errorf("expected error from response matcher")
	}
}

func Test_Validator_ExecutesSetupsAndTeardowns(t *testing.T) {
	i := 1
	v := newConsumerValidator(func() error {
		if i != 1 {
			t.Errorf("Expected this action to be called at %v position but is at %d", 1, i)
		} else {
			i++
		}
		return nil
	}, func() error {
		if i != 4 {
			t.Errorf("Expected this action to be called at %d position but is at %d", 4, i)
		}
		return nil
	}, DefaultLogger)

	sa := &stateAction{setup: func() error {
		if i != 2 {
			t.Errorf("Expected this action to be called at %d position but is at %d", 2, i)
		} else {
			i++
		}
		return nil
	}, teardown: func() error {
		if i != 3 {
			t.Errorf("Expected this action to be called at %d position but is at %d", 3, i)
		} else {
			i++
		}
		return nil
	}}

	interaction, _ := consumer.NewInteraction("description", "state", provider.NewJSONRequest("Get", "/", "", nil), provider.NewJSONResponse(200, nil))
	f := io.NewPactFile("consumer", "provider", []*consumer.Interaction{interaction})

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(interaction.Response.Status)
	}))
	defer s.Close()
	u, _ := url.Parse(s.URL)

	v.ProviderService(&http.Client{}, u)
	if res, err := v.Validate(f, map[string]*stateAction{"state": sa}); err != nil {
		t.Error(err)
	} else if !res {
		t.Error("Validation Failed")
	} else if i != 4 {
		t.Error("Setup and teardown actions were not called correctly")
	}
}

func Test_Validator_ThrowsErrorFromSetupsAndTeardowns(t *testing.T) {
	testErr := errors.New("action error")
	fn := func() error {
		return testErr
	}

	sa := &stateAction{setup: nil, teardown: nil}

	interaction, _ := consumer.NewInteraction("description", "state", provider.NewJSONRequest("Get", "/", "", nil), provider.NewJSONResponse(200, nil))
	f := io.NewPactFile("consumer", "provider", []*consumer.Interaction{interaction})

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(interaction.Response.Status)
	}))
	defer s.Close()

	u, _ := url.Parse(s.URL)

	//test setup action for every interaction
	v := newConsumerValidator(fn, nil, DefaultLogger)
	v.ProviderService(&http.Client{}, u)
	if _, err := v.Validate(f, map[string]*stateAction{"state": sa}); err == nil {
		t.Errorf("expected %s", testErr)
	} else if err != testErr {
		t.Errorf("expected %s, got %s", testErr, err)
	}

	//test teardown action for every interaction
	v = newConsumerValidator(nil, fn, DefaultLogger)
	v.ProviderService(&http.Client{}, u)
	if _, err := v.Validate(f, map[string]*stateAction{"state": sa}); err == nil {
		t.Errorf("expected %s", testErr)
	} else if err != testErr {
		t.Errorf("expected %s, got %s", testErr, err)
	}

	//test setup action for specific interaction
	sa = &stateAction{setup: fn, teardown: nil}
	v = newConsumerValidator(nil, fn, DefaultLogger)
	v.ProviderService(&http.Client{}, u)
	if _, err := v.Validate(f, map[string]*stateAction{"state": sa}); err == nil {
		t.Errorf("expected %s", testErr)
	} else if err != testErr {
		t.Errorf("expected %s, got %s", testErr, err)
	}

	//test teardown action for every interaction
	sa = &stateAction{setup: nil, teardown: fn}
	v = newConsumerValidator(nil, fn, DefaultLogger)
	v.ProviderService(&http.Client{}, u)
	if _, err := v.Validate(f, map[string]*stateAction{"state": sa}); err == nil {
		t.Errorf("expected %s", testErr)
	} else if err != testErr {
		t.Errorf("expected %s, got %s", testErr, err)
	}
}
