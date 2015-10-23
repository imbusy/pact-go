package consumer

import ()

type ProviderService interface {
	Given(state string) ProviderService
	UponReceiving(description string) ProviderService
	With(request *ProviderRequest) ProviderService
	WillRespondWith(response *ProviderResponse)
	VerifyInteractions() error
}

//this can be private, needs to implement the above interface
//look into httptest package, we can create test server from there for a mock server

type MockProviderService struct {
	providerRequest  *ProviderRequest
	providerResponse *ProviderResponse
	state            string
	description      string
}

type ProviderRequest struct {
	method  string
	path    string
	query   string
	headers map[string]string
	body    string
}

type ProviderResponse struct {
	status  string
	headers map[string]string
	body    string
}

func NewMockProviderService(config *PactConfig) *MockProviderService {
	return &MockProviderService{}
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

func (p *MockProviderService) WillRespondWith(providerResponse *ProviderResponse) {
	p.providerResponse = providerResponse
	p.registerInteractions()
}

func (p *MockProviderService) VerifyInteractions() error {
	return nil
}

func (p *MockProviderService) registerInteractions() {

}
