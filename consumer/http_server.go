package consumer

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/SEEK-Jobs/pact-go/comparers"
)

type HTTPMockService struct {
	server                *httptest.Server
	interactions          []*Interaction
	requestedInteractions []*Interaction
}

var (
	errNotFound = errors.New("No matching interaction found.")
)

func NewHTTPMockService() *HTTPMockService {
	return &HTTPMockService{interactions: make([]*Interaction, 0)}
}

func (ms *HTTPMockService) Start() string {
	if ms.server == nil {
		ms.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			matchedInteraction, err := ms.findMatchingInteraction(r, ms.interactions)
			if matchedInteraction == nil && err == nil {
				http.Error(w, errNotFound.Error(), http.StatusNotFound)
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				matchedInteraction.WriteToHttpResponse(w)
				ms.requestedInteractions = append(ms.requestedInteractions, matchedInteraction)
			}
		}))
	}
	return ms.server.URL
}

func (ms *HTTPMockService) Stop() {
	if ms.server != nil {
		ms.server.Close()
		ms.server = nil
	}
}

func (ms *HTTPMockService) RegisterInteraction(interaction *Interaction) {
	ms.interactions = append(ms.interactions, interaction)
}

func (ms *HTTPMockService) ClearInteractions() {
	ms.interactions = nil
	ms.requestedInteractions = nil
	ms.interactions = make([]*Interaction, 0)
	ms.requestedInteractions = make([]*Interaction, 0)
}

//VerifyInteractions - Verfies if the registered interactions has been requested and verified
func (ms *HTTPMockService) VerifyInteractions() error {
	return verifyInteractions(ms.interactions, ms.requestedInteractions)
}

func (ms *HTTPMockService) GetRegisteredInteractions() []*Interaction {
	return ms.interactions
}

func (ms *HTTPMockService) findMatchingInteraction(r *http.Request, interactions []*Interaction) (*Interaction, error) {

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
