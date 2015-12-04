package provider

import (
	"encoding/json"
	"net/http"
)

type ProviderResponse struct {
	Status  int
	Headers http.Header
	HttpContent
}

func NewJsonProviderResponse(status int, headers http.Header) *ProviderResponse {
	return &ProviderResponse{
		Status:      status,
		Headers:     headers,
		HttpContent: &jsonContent{},
	}
}

func (p *ProviderResponse) ResetContent() {
	p.HttpContent = nil
}

func (p *ProviderResponse) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{"status": p.Status}

	if p.Headers != nil {
		obj["headers"] = getHeaderWithSingleValues(p.Headers)
	}
	if p.HttpContent != nil {
		if body := p.GetBody(); body != nil {
			obj["body"] = body
		}
	}
	return json.Marshal(obj)
}
