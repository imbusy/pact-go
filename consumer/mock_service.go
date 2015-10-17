package consumer

import (
	"fmt"
)

type MockProviderService interface {
	Given(state string) MockProviderService
	UponReceiving(description string) MockProviderService
	With(request *ProviderRequest) MockProviderService
	WillRespondWith(response *ProviderResponse)
	VerifyInteractions() error
}

//this can be private, needs to implement the above functionality
type mockService struct {
	providerRequest  *ProviderRequest
	ProviderResponse *ProviderResponse
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
