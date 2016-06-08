package consumer

import (
	"errors"
	"fmt"
	"github.com/imbusy/pact-go/comparers"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type HTTPMockService struct {
	server                *httptest.Server
	interactions          []*Interaction
	inScopeInteractions   []*Interaction
	requestedInteractions []*Interaction
}

var (
	errNotFound = errors.New("No matching interaction found.")
)

const (
	errDuplicateInteractionInScopeMsg = "An interaction already exists with the description '%s' and provider state '%s' in this test. Please supply a different description or provider state."
	errDuplicateInteractionMsg        = "An interaction registered by another test already exists with the description '%s' and provider state '%s', however the interaction does not match exactly. Please supply a different description or provider state. Alternatively align this interaction to match the duplicate exactly."
)

func NewHTTPMockService() *HTTPMockService {
	return &HTTPMockService{interactions: make([]*Interaction, 0)}
}

func (ms *HTTPMockService) Start() string {
	if ms.server == nil {
		ms.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			matchedInteraction, err := ms.findMatchingInteractionInScope(r)
			if matchedInteraction == nil && err == nil {
				http.Error(w, errNotFound.Error(), http.StatusNotFound)
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				matchedInteraction.WriteToHTTPResponse(w)
				ms.requestedInteractions = append(ms.requestedInteractions, matchedInteraction)
			}
		}))
	}
	return ms.server.URL
}

func (ms *HTTPMockService) Stop() {
	ms.interactions = nil
	if ms.server != nil {
		ms.server.Close()
		ms.server = nil
	}
}

func (ms *HTTPMockService) RegisterInteraction(interaction *Interaction) error {
	if si := findSimilarInteraction(ms.inScopeInteractions, interaction); si != nil {
		return fmt.Errorf(errDuplicateInteractionInScopeMsg, interaction.Description, interaction.State)
	}

	if si := findSimilarInteraction(ms.interactions, interaction); si == nil {
		ms.interactions = append(ms.interactions, interaction)
	} else if !reflect.DeepEqual(si, interaction) {
		return errors.New(fmt.Sprintf(errDuplicateInteractionMsg, interaction.Description, interaction.State))
	}
	ms.inScopeInteractions = append(ms.inScopeInteractions, interaction)
	return nil
}

func (ms *HTTPMockService) ClearInteractions() {
	ms.inScopeInteractions = nil
	ms.requestedInteractions = nil
	ms.inScopeInteractions = make([]*Interaction, 0)
	ms.requestedInteractions = make([]*Interaction, 0)
}

func (ms *HTTPMockService) IsTestScopeClear() bool {
	return len(ms.inScopeInteractions) == 0 && len(ms.requestedInteractions) == 0
}

//VerifyInteractions - Verfies if the registered interactions has been requested and verified
func (ms *HTTPMockService) VerifyInteractions() error {
	return verifyInteractions(ms.inScopeInteractions, ms.requestedInteractions)
}

func (ms *HTTPMockService) GetRegisteredInteractions() []*Interaction {
	return ms.interactions
}

func findSimilarInteraction(src []*Interaction, i *Interaction) *Interaction {
	for _, o := range src {
		if i.IsSimilar(o) {
			return o
		}
	}
	return nil
}

func (ms *HTTPMockService) findMatchingInteractionInScope(r *http.Request) (*Interaction, error) {

	for i := range ms.inScopeInteractions {
		req, err := ms.inScopeInteractions[i].ToHTTPRequest(ms.server.URL)

		if err != nil {
			return nil, err
		}

		result, err := comparers.MatchRequest(req, r)

		if err != nil {
			return nil, err
		}

		if result {
			return ms.inScopeInteractions[i], nil
		}
	}
	return nil, nil
}
