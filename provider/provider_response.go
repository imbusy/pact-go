package provider

import (
	"encoding/json"
	"net/http"
)

type ProviderResponse struct {
	Status  int
	Headers http.Header
	*jsonContent
}

func NewProviderResponse(status int, headers http.Header) *ProviderResponse {
	return &ProviderResponse{
		Status:      status,
		Headers:     headers,
		jsonContent: &jsonContent{},
	}
}

func (p *ProviderResponse) MarshalJSON() ([]byte, error) {
	if len(p.data) > 0 {
		return json.Marshal(map[string]interface{}{
			"status":  p.Status,
			"headers": p.Headers,
			"body":    p.data,
		})
	} else if len(p.sliceData) > 0 {
		return json.Marshal(map[string]interface{}{
			"status":  p.Status,
			"headers": p.Headers,
			"body":    p.sliceData,
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"status":  p.Status,
			"headers": p.Headers,
		})
	}

}
