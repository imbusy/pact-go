package provider

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Method  string
	Path    string
	Query   string
	Headers http.Header
	HttpContent
}

func NewJsonRequest(method, path, query string, headers http.Header) *Request {
	if method == "" {
		//return error
	}

	if path == "" {
		//return error
	}

	return &Request{
		Method:      method,
		Path:        path,
		Query:       query,
		Headers:     headers,
		HttpContent: &jsonContent{},
	}
}

func (p *Request) ResetContent() {
	p.HttpContent = nil
}

func (p *Request) MarshalJSON() ([]byte, error) {
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

	if p.HttpContent != nil {
		if body := p.GetBody(); body != nil {
			obj["body"] = body
		}
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
