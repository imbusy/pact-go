package consumer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Interaction struct {
	State       string
	Description string
	Request     *ProviderRequest
	Response    *ProviderResponse
}

func NewInteraction(state string, description string, request *ProviderRequest,
	response *ProviderResponse) *Interaction {
	return &Interaction{
		State:       state,
		Description: description,
		Request:     request,
		Response:    response,
	}
}

func (i *Interaction) ToHttpRequest(baseUrl string) (*http.Request, error) {
	u, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, err
	}
	u.Path = i.Request.Path
	u.RawQuery = i.Request.Query

	bodyReader := getBodyReader(i.Request.Body)
	req, err := http.NewRequest(i.Request.Method, u.String(), bodyReader)

	if err != nil {
		return nil, err
	}

	for header, val := range i.Request.Headers {
		req.Header.Add(header, val[0])
	}

	return req, nil
}

func (i *Interaction) WriteToHttpResponse(w http.ResponseWriter) error {
	w.WriteHeader(i.Response.Status)
	respHeader := w.Header()

	for header, val := range i.Response.Headers {
		respHeader.Add(header, val[0])
	}

	if i.Response.Body != "" {
		var body interface{}
		if err := json.Unmarshal([]byte(i.Response.Body), &body); err != nil {
			return err
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(body)
	}
	return nil
}

func getBodyReader(params string) io.Reader {
	if params != "" {
		return strings.NewReader(params)
	}
	return nil
}
