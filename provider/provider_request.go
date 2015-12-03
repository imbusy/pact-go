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
	if method == "" {
		//return error
	}

	if path == "" {
		//return error
	}

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
	obj := map[string]interface{}{
		"method": p.Method,
		"path":   p.Path,
	}

	if p.Query != "" {
		obj["query"] = p.Query
	}

	if p.Headers != nil {
		obj["headers"] = getHeaderWithSingleValues(p.Headers)
	}

	if body != nil {
		obj["body"] = body
	}
	return json.Marshal(obj)
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
