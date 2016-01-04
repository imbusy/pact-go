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
	PactUri(uri string, config *PactUriConfig) Verifier
	Verify() error
}

type Action func() error

type stateAction struct {
	setup    Action
	teardown Action
}

type pactFileVerfier struct {
	stateActions  map[string]*stateAction
	provider      string
	consumer      string
	pactUri       string
	pactUriConfig *PactUriConfig
	validator     consumerValidator
	config        *VerfierConfig
}

//NewPactFileVerifier creates a new pact verifier. The setup & teardown actions
//get executed before each interaction is verified.
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

//ServiceProvider provides the information needed to verify the interactions with service provider
func (v *pactFileVerfier) ServiceProvider(providerName string, c *http.Client, u *url.URL) Verifier {
	v.provider = providerName
	v.validator.ProviderService(c, u)
	return v
}

//ProviderState sets the setup and teardown action to be executed before a interaction with specific state gets verified
func (v *pactFileVerfier) ProviderState(state string, setup, teardown Action) Verifier {
	//sacrificed empty state validation in favour of chaining
	if state != "" {
		v.stateActions[state] = &stateAction{setup: setup, teardown: teardown}
	}
	return v
}

//HonoursPactWith consumer with which pact needs to be honoured
func (v *pactFileVerfier) HonoursPactWith(consumerName string) Verifier {
	v.consumer = consumerName
	return v
}

//PactUri sets the uri to get the pact file
func (v *pactFileVerfier) PactUri(uri string, config *PactUriConfig) Verifier {
	if config == nil {
		config = DefaultPactUriConfig
	}
	v.pactUriConfig = config
	v.pactUri = uri
	return v
}

//Verify verifies all the interactions of consumer against the provider
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
	var r io.PactReader
	if io.IsWebUri(v.pactUri) {
		r = io.NewPactWebReader(v.pactUri, v.pactUriConfig.AuthScheme, v.pactUriConfig.AuthValue)
	} else {
		r = io.NewPactFileReader(v.pactUri)
	}

	f, err := r.Read()
	if err != nil {
		return nil, err
	}

	if err := f.Validate(); err != nil {
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
