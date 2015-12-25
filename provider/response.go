package provider

import (
	"encoding/json"
	"errors"
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

//UnmarshalJSON cusotm json unmarshalling
func (p *Response) UnmarshalJSON(b []byte) error {
	var obj map[string]interface{}
	r := Response{httpContent: &jsonContent{}}

	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	if status, ok := obj["status"].(float64); ok { //default number deserialised as float64
		r.Status = int(status)
	} else {
		return errors.New("Could not unmarshal response, status value is either nil or not a int")
	}

	if headers, ok := obj["headers"].(map[string]string); ok {
		for key, val := range headers {
			r.Headers.Add(key, val)
		}
	}

	r.SetBody(obj["body"])
	*p = Response(r)
	return nil
}
