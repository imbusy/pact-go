package pact

import (
	"github.com/SEEK-Jobs/pact-go/consumer"
	"github.com/SEEK-Jobs/pact-go/provider"
	"github.com/SEEK-Jobs/pact-go/writer"
)

//ProviderService - Interface to register and verify interaactions between consumer and service provider.
type ProviderService interface {
	Given(state string) ProviderService
	UponReceiving(description string) ProviderService
	With(request provider.ProviderRequest) ProviderService
	WillRespondWith(response provider.ProviderResponse) ProviderService
	ClearInteractions() ProviderService
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

func (p *mockProviderService) WillRespondWith(providerResponse provider.ProviderResponse) ProviderService {
	p.providerResponse = &providerResponse
	p.registerInteraction()
	return p
}

func (p *mockProviderService) ClearInteractions() ProviderService {
	p.service.ClearInteractions()
	p.resetTransientState()
	return p
}

func (p *mockProviderService) VerifyInteractions() error {
	return p.VerifyInteractions()
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

func (p *mockProviderService) registerInteraction() {
	interaction := consumer.NewInteraction(p.description, p.state, p.providerRequest, p.providerResponse)
	p.service.RegisterInteraction(interaction)
	p.resetTransientState()
}

func (p *mockProviderService) resetTransientState() {
	p.state = ""
	p.description = ""
	p.providerRequest = nil
	p.providerResponse = nil
}
