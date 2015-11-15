package consumer

import (
	"github.com/bennycao/pact-go/provider"
)

type ProviderService interface {
	Given(state string) ProviderService
	UponReceiving(description string) ProviderService
	With(request *provider.ProviderRequest) ProviderService
	WillRespondWith(response *provider.ProviderResponse) ProviderService
	ClearInteractions() ProviderService
	VerifyInteractions() error
}

//this can be private, needs to implement the above interface
//look into httptest package, we can create test server from there for a mock server

type MockProviderService struct {
	providerRequest  *provider.ProviderRequest
	providerResponse *provider.ProviderResponse
	state            string
	description      string
	service          *httpMockService
}

func NewMockProviderService(config *PactConfig) *MockProviderService {
	return &MockProviderService{service: newHttpMockService()}
}

func (p *MockProviderService) Given(state string) ProviderService {
	p.state = state
	return p
}

func (p *MockProviderService) UponReceiving(description string) ProviderService {
	p.description = description
	return p
}

func (p *MockProviderService) With(providerRequest *provider.ProviderRequest) ProviderService {
	p.providerRequest = providerRequest
	return p
}

func (p *MockProviderService) WillRespondWith(providerResponse *provider.ProviderResponse) ProviderService {
	p.providerResponse = providerResponse
	p.registerInteraction()
	p.resetTransientState()
	return p
}

func (p *MockProviderService) ClearInteractions() ProviderService {
	p.service.ClearInteractions()
	p.resetTransientState()
	return p
}

func (p *MockProviderService) VerifyInteractions() error {

	return nil
}

func (p *MockProviderService) registerInteraction() {
	interaction := NewInteraction(p.description, p.state, p.providerRequest, p.providerResponse)
	p.service.RegisterInteraction(interaction)
}

func (p *MockProviderService) resetTransientState() {
	p.state = ""
	p.description = ""
	p.providerRequest = nil
	p.providerResponse = nil
}
