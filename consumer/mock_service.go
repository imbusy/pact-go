package consumer

import (
	"fmt"
)

type MockService struct {
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
