package provider

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int
	Headers http.Header
	HttpContent
}

func NewJsonResponse(status int, headers http.Header) *Response {
	return &Response{
		Status:      status,
		Headers:     headers,
		HttpContent: &jsonContent{},
	}
}

func (p *Response) ResetContent() {
	p.HttpContent = nil
}

func (p *Response) MarshalJSON() ([]byte, error) {
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
