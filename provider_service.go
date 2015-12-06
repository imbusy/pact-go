package pact

import (
	"github.com/SEEK-Jobs/pact-go/consumer"
	"github.com/SEEK-Jobs/pact-go/provider"
	"github.com/SEEK-Jobs/pact-go/writer"
)

//ProviderService - Interface to register and verify interaactions between consumer and service provider.
type ProviderService interface {
	//Given the state exists
	Given(state string) ProviderService
	//UponReceiving this interaction with description
	UponReceiving(description string) ProviderService
	//With this request
	With(request provider.ProviderRequest) ProviderService
	//WillRespondWith with this response
	WillRespondWith(response provider.ProviderResponse) error
	//ClearInteractions clears all the registered interaactions
	ClearInteractions() ProviderService
	//VerifyInteractions checks if all the registered interactions have been verified
	VerifyInteractions() error
}

type mockProviderService struct {
	providerRequest  *provider.ProviderRequest
	providerResponse *provider.ProviderResponse
	state            string
	description      string
	service          *consumer.HTTPMockService
	config           *Config
}

func newMockProviderService(config *Config) *mockProviderService {
	return &mockProviderService{service: consumer.NewHTTPMockService(), config: config}
}

func (p *mockProviderService) Given(state string) ProviderService {
	p.state = state
	return p
}

func (p *mockProviderService) UponReceiving(description string) ProviderService {
	p.description = description
	return p
}

func (p *mockProviderService) With(providerRequest provider.ProviderRequest) ProviderService {
	p.providerRequest = &providerRequest
	return p
}

func (p *mockProviderService) WillRespondWith(providerResponse provider.ProviderResponse) error {
	p.providerResponse = &providerResponse
	if err := p.registerInteraction(); err != nil {
		return err
	}
	return nil
}

func (p *mockProviderService) ClearInteractions() ProviderService {
	p.service.ClearInteractions()
	p.resetTransientState()
	return p
}

func (p *mockProviderService) VerifyInteractions() error {
	return p.service.VerifyInteractions()
}

func (p *mockProviderService) start() string {
	return p.service.Start()
}

func (p *mockProviderService) stop() {
	p.ClearInteractions()
	p.service.Stop()
}

func (p *mockProviderService) persistPact(consumer, serviceProvider string) error {
	pact := writer.NewPactFile(consumer, serviceProvider, p.service.GetRegisteredInteractions())
	return writer.NewPactFileWriter(pact, p.config.PactPath).Write()
}

func (p *mockProviderService) registerInteraction() error {
	interaction, err := consumer.NewInteraction(p.description, p.state, p.providerRequest, p.providerResponse)
	if err != nil {
		return err
	}
	p.service.RegisterInteraction(interaction)
	p.resetTransientState()
	return nil
}

func (p *mockProviderService) resetTransientState() {
	p.state = ""
	p.description = ""
	p.providerRequest = nil
	p.providerResponse = nil
}
