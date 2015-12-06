package provider

import (
	"encoding/json"
	"net/http"
)

//Response provider response
type Response struct {
	Status  int
	Headers http.Header
	httpContent
}

//NewJSONResponse creates new response with body as json content
func NewJSONResponse(status int, headers http.Header) *Response {
	return &Response{
		Status:      status,
		Headers:     headers,
		httpContent: &jsonContent{},
	}
}

//ResetContent emoves an existing contentß
func (p *Response) ResetContent() {
	p.httpContent = nil
}

//MarshalJSON custom json marshaling
func (p *Response) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{"status": p.Status}

	if p.Headers != nil {
		obj["headers"] = getHeaderWithSingleValues(p.Headers)
	}
	if p.httpContent != nil {
		if body := p.GetBody(); body != nil {
			obj["body"] = body
		}
	}
	return json.Marshal(obj)
}
