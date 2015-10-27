package consumer

import (
	"fmt"
	"github.com/bennycao/pact-go/comparers"
	"net/http"
	"net/http/httptest"
)

type httpMockService struct {
	server       *httptest.Server
	interactions []*Interaction
}

func newHttpMockService() *httpMockService {
	return &httpMockService{interactions: make([]*Interaction, 0)}
}

func (ms *httpMockService) Start() string {

	ms.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		matchedInteraction := ms.findMatchingInteraction(r, ms.interactions)
		if matchedInteraction == nil {
			//return error
			http.NotFound(w, r)
		} else {

			w.WriteHeader(matchedInteraction.Response.Status)
		}
	}))

	return ms.server.URL
}

func (ms *httpMockService) Stop() {
	ms.server.Close()
}

func (ms *httpMockService) RegisterInteraction(interaction *Interaction) {
	ms.interactions = append(ms.interactions, interaction)
}

func (ms *httpMockService) ClearInteractions() {
	ms.interactions = nil
	ms.interactions = make([]*Interaction, 0)
}

func (ms *httpMockService) findMatchingInteraction(r *http.Request, interactions []*Interaction) *Interaction {

	for i := range interactions {
		req, _ := interactions[i].ToHttpRequest(ms.server.URL)

		result, _ := comparers.MatchRequest(req, r)
		fmt.Println(result)

		if result {
			fmt.Println("found match")

			return interactions[i]
		}
	}
	return nil
}
