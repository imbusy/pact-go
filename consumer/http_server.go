package consumer

import (
	"errors"
	"github.com/bennycao/pact-go/comparers"
	"net/http"
	"net/http/httptest"
)

type httpMockService struct {
	server       *httptest.Server
	interactions []*Interaction
}

var (
	notFoundError = errors.New("No matching interaction found.")
)

func newHttpMockService() *httpMockService {
	return &httpMockService{interactions: make([]*Interaction, 0)}
}

func (ms *httpMockService) Start() string {
	if ms.server == nil {
		ms.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			matchedInteraction, err := ms.findMatchingInteraction(r, ms.interactions)
			if matchedInteraction == nil && err == nil {
				http.Error(w, notFoundError.Error(), http.StatusNotFound)
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				matchedInteraction.WriteToHttpResponse(w)
			}
		}))
	}
	return ms.server.URL
}

func (ms *httpMockService) Stop() {
	if ms.server != nil {
		ms.server.Close()
	}
}

func (ms *httpMockService) RegisterInteraction(interaction *Interaction) {
	ms.interactions = append(ms.interactions, interaction)
}

func (ms *httpMockService) ClearInteractions() {
	ms.interactions = nil
	ms.interactions = make([]*Interaction, 0)
}

func (ms *httpMockService) findMatchingInteraction(r *http.Request, interactions []*Interaction) (*Interaction, error) {

	for i := range interactions {
		req, err := interactions[i].ToHttpRequest(ms.server.URL)

		if err != nil {
			return nil, err
		}

		result, err := comparers.MatchRequest(req, r)

		if err != nil {
			return nil, err
		}

		if result {
			return interactions[i], nil
		}
	}
	return nil, nil
}
