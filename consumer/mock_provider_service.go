package consumer

import (
	"net/http"
)

type ProviderService interface {
	Given(state string) ProviderService
	UponReceiving(description string) ProviderService
	With(request *ProviderRequest) ProviderService
	WillRespondWith(response *ProviderResponse) ProviderService
	ClearInteractions() ProviderService
	VerifyInteractions() error
}

//this can be private, needs to implement the above interface
//look into httptest package, we can create test server from there for a mock server

type MockProviderService struct {
	providerRequest  *ProviderRequest
	providerResponse *ProviderResponse
	state            string
	description      string
	service          *httpMockService
}

type ProviderRequest struct {
	Method  string
	Path    string
	Query   string
	Headers http.Header
	Body    string
}

type ProviderResponse struct {
	Status  int
	Headers http.Header
	Body    string
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

func (p *MockProviderService) With(providerRequest *ProviderRequest) ProviderService {
	p.providerRequest = providerRequest
	return p
}

func (p *MockProviderService) WillRespondWith(providerResponse *ProviderResponse) ProviderService {
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
	interaction := NewInteraction(p.state, p.description, p.providerRequest, p.providerResponse)
	p.service.RegisterInteraction(interaction)
}

func (p *MockProviderService) resetTransientState() {
	p.state = ""
	p.description = ""
	p.providerRequest = nil
	p.providerResponse = nil
}
