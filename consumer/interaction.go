package consumer

import (
	"bytes"
	"errors"
	"github.com/imbusy/pact-go/provider"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Interaction struct {
	State       string             `json:"provider_state,omitempty"`
	Description string             `json:"description"`
	Request     *provider.Request  `json:"request"`
	Response    *provider.Response `json:"response"`
}

var (
	errEmptyDescription = errors.New("Cannot register interaction with empty description.")
	errNilRequest       = errors.New("Cannot register interaction with nil request")
	errNilResponse      = errors.New("Cannot register interaction with nil response")
)

func NewInteraction(description string, state string, request *provider.Request,
	response *provider.Response) (*Interaction, error) {

	if description == "" {
		return nil, errEmptyDescription
	} else if request == nil {
		return nil, errNilRequest
	} else if response == nil {
		return nil, errNilResponse
	}

	return &Interaction{
		State:       state,
		Description: description,
		Request:     request,
		Response:    response,
	}, nil
}

func (i *Interaction) ToHTTPRequest(baseURL string) (*http.Request, error) {
	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = i.Request.Path
	u.RawQuery = i.Request.Query

	body, err := i.Request.GetData()
	if err != nil {
		return nil, err
	}

	bodyReader := getReader(body)
	req, err := http.NewRequest(strings.ToUpper(i.Request.Method), u.String(), bodyReader)

	if err != nil {
		return nil, err
	}

	for header, val := range i.Request.Headers {
		req.Header.Add(header, val[0])
	}

	return req, nil
}

func (i *Interaction) IsSimilar(to *Interaction) bool {
	return strings.EqualFold(i.Description, to.Description) &&
		strings.EqualFold(i.State, to.State)
}

func (i *Interaction) WriteToHTTPResponse(w http.ResponseWriter) error {
	w.WriteHeader(i.Response.Status)
	respHeader := w.Header()

	for header, val := range i.Response.Headers {
		respHeader.Add(header, val[0])
	}

	if body, err := i.Response.GetData(); err != nil {
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
