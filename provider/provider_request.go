package provider

import (
	"encoding/json"
	"net/http"
)

type ProviderRequest struct {
	Method  string
	Path    string
	Query   string
	Headers http.Header
	HttpContent
}

func NewJsonProviderRequest(method, path, query string, headers http.Header) *ProviderRequest {
	return &ProviderRequest{
		Method:      method,
		Path:        path,
		Query:       query,
		Headers:     headers,
		HttpContent: &jsonContent{},
	}
}

func (p *ProviderRequest) MarshalJSON() ([]byte, error) {
	body := p.GetBody()
	if body != nil {
		return json.Marshal(map[string]interface{}{
			"method":  p.Method,
			"path":    p.Path,
			"query":   p.Query,
			"headers": getHeaderWithSingleValues(p.Headers),
			"body":    body,
		})
	} else {
		return json.Marshal(map[string]interface{}{
			"method":  p.Method,
			"path":    p.Path,
			"query":   p.Query,
			"headers": getHeaderWithSingleValues(p.Headers),
		})
	}

}

func getHeaderWithSingleValues(headers http.Header) map[string]string {
	if headers == nil {
		return nil
	}

	h := make(map[string]string)
	for header, val := range headers {
		h[header] = val[0]
	}
	return h
}
