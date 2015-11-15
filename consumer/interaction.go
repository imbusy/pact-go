package consumer

import (
	"bytes"
	"github.com/bennycao/pact-go/provider"
	"io"
	"net/http"
	"net/url"
)

type Interaction struct {
	State       string                    `json:"state"`
	Description string                    `json:"description"`
	Request     *prvider.ProviderRequest  `json:"request"`
	Response    *povider.ProviderResponse `json:"response"`
}

func NewInteraction(state string, description string, request *provider.ProviderRequest,
	response *provider.ProviderResponse) *Interaction {
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

	body, err := i.Request.GetBody()
	if err != nil {
		return nil, err
	}

	bodyReader := getReader(body)
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

	if body, err := i.Response.GetBody(); err != nil {
		return err
	} else if body != nil {
		if _, writeErr := w.Write(body); writeErr != nil {
			return writeErr
		}
	}
	return nil
}

func getReader(content []byte) io.Reader {
	if content != nil {
		return bytes.NewReader(content)
	}
	return nil
}
