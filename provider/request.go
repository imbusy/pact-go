package provider

import (
	"encoding/json"
	"errors"
	"net/http"
)

//Request provider request
type Request struct {
	Method  string
	Path    string
	Query   string
	Headers http.Header
	httpContent
}

//NewJSONRequest creates new http request with content body as json
func NewJSONRequest(method, path, query string, headers http.Header) *Request {
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
		httpContent: &jsonContent{},
	}
}

//ResetContent removes an existing contentß
func (p *Request) ResetContent() {
	p.httpContent = nil
}

//MarshalJSON custom json marshaling
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

	if p.httpContent != nil {
		if body := p.GetBody(); body != nil {
			obj["body"] = body
		}
	}

	return json.Marshal(obj)
}

//UnmarshalJSON cusotm json unmarshalling
func (p *Request) UnmarshalJSON(b []byte) error {
	var obj map[string]interface{}
	r := Request{httpContent: &jsonContent{}}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}

	if method, ok := obj["method"].(string); ok {
		r.Method = method
	} else {
		return errors.New("Could not unmarshal request, method value is either nil or not a string")
	}

	if path, ok := obj["path"].(string); ok {
		r.Path = path
	} else {
		return errors.New("Could not unmarshal request, path value is either nil or not a string")
	}

	if query, ok := obj["query"].(string); ok {
		r.Query = query
	}

	if headers, ok := obj["headers"].(map[string]interface{}); ok {
		r.Headers = make(http.Header)
		for key, val := range headers {
			if str, ok := val.(string); ok {
				r.Headers.Add(key, str)
			}
		}
	}

	r.SetBody(obj["body"])
	*p = Request(r)
	return nil
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
