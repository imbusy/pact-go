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
	*jsonContent
}

func NewProviderRequest(method, path, query string, headers http.Header) *ProviderRequest {
	return &ProviderRequest{
		Method:      method,
		Path:        path,
		Query:       query,
		Headers:     headers,
		jsonContent: &jsonContent{},
	}
}

func (p *ProviderRequest) MarshalJSON() ([]byte, error) {
	if len(p.data) > 0 {
		return json.Marshal(map[string]interface{}{
			"method":  p.Method,
			"path":    p.Path,
			"query":   p.Query,
			"headers": getHeaderWithSingleValues(p.Headers),
			"body":    p.data,
		})
	} else if len(p.sliceData) > 0 {
		return json.Marshal(map[string]interface{}{
			"method":  p.Method,
			"path":    p.Path,
			"query":   p.Query,
			"headers": getHeaderWithSingleValues(p.Headers),
			"body":    p.sliceData,
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
