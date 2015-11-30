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

func (p *ProviderResponse) MarshalJSON() ([]byte, error) {
	body := p.GetBody()

	if body != nil {
		return json.Marshal(map[string]interface{}{
			"status":  p.Status,
			"headers": getHeaderWithSingleValues(p.Headers),
			"body":    body,
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"status":  p.Status,
			"headers": getHeaderWithSingleValues(p.Headers),
		})
	}

}
