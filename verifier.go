package pact

import (
	"errors"
	"github.com/SEEK-Jobs/pact-go/io"
	"net/http"
	"net/url"
)

type Verifier interface {
	ProviderState(state string, setup, teardown Action) Verifier
	ServiceProvider(providerName string, c *http.Client, u *url.URL) Verifier
	HonoursPactWith(consumerName string) Verifier
	PactUri(uri string) Verifier
	Verify() error
}

type Action func() error

type stateAction struct {
	setup    Action
	teardown Action
}

type pactFileVerfier struct {
	stateActions map[string]*stateAction
	provider     string
	consumer     string
	pactUri      string
	validator    consumerValidator
	config       *VerfierConfig
}

func NewPactFileVerifier(setup, teardown Action, config *VerfierConfig) Verifier {
	if config == nil {
		config = DefaultVerifierConfig
	}

	return &pactFileVerfier{
		validator:    newConsumerValidator(setup, teardown, config.Logger),
		config:       config,
		stateActions: make(map[string]*stateAction),
	}
}

var (
	errEmptyProvider     = errors.New("Provider name cannot be empty, please provide a valid value using ServiceProvider function.")
	errEmptyConsumer     = errors.New("Consumer name cannot be empty, please provide a valid value using HonoursPactWith function.")
	errVerficationFailed = errors.New("Failed to verify the pact, please see the log for more details.")
)

func (v *pactFileVerfier) ServiceProvider(providerName string, c *http.Client, u *url.URL) Verifier {
	v.provider = providerName
	v.validator.ProviderService(c, u)
	return v
}

func (v *pactFileVerfier) ProviderState(state string, setup, teardown Action) Verifier {
	//sacrificed empty state validation in favour of chaining
	if state != "" {
		v.stateActions[state] = &stateAction{setup: setup, teardown: teardown}
	}
	return v
}

func (v *pactFileVerfier) HonoursPactWith(consumerName string) Verifier {
	v.consumer = consumerName
	return v
}

func (v *pactFileVerfier) PactUri(uri string) Verifier {
	v.pactUri = uri
	return v
}

func (v *pactFileVerfier) Verify() error {
	if err := v.verifyInternalState(); err != nil {
		return err
	}

	//get pact file
	f, err := v.getPactFile()
	if err != nil {
		return err
	}

	//validate interactions
	if ok, err := v.validator.Validate(f, v.stateActions); err != nil {
		return err
	} else if !ok {
		return errVerficationFailed
	}

	return nil
}

func (v *pactFileVerfier) getPactFile() (*io.PactFile, error) {
	r := io.NewPactFileReader(v.pactUri)
	f, err := r.Read()
	if err != nil {
		return nil, err
	} else if err := f.Validate(); err != nil {
		return nil, err
	}
	return f, nil
}

func (v *pactFileVerfier) verifyInternalState() error {
	if v.consumer == "" {
		return errEmptyConsumer
	}

	if v.provider == "" {
		return errEmptyProvider
	}
	return v.validator.CanValidate()
}
