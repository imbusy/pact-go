package consumer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/SEEK-Jobs/pact-go/provider"
)

func Test_ShouldVerify_MissingInteractions(t *testing.T) {
	i, _ := NewInteraction("description", "state", provider.NewJSONRequest("GET", "/", "", nil), provider.NewJSONResponse(200, nil))
	registered := []*Interaction{
		i,
	}
	requested := []*Interaction{}
	if err := verifyInteractions(registered, requested); err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf(errIntMissingUsage, i.Description, i.State)) {
			t.Errorf("expected message to contain: %s, actual message: %s", fmt.Sprintf(errIntMissingUsage, i.Description, i.State), err.Error())
		}
	} else if err == nil {
		t.Error("expected missing interaction error")
	}
}

func Test_ShouldVerify_MultipleSameInteractionCalls(t *testing.T) {
	i, _ := NewInteraction("description", "state", provider.NewJSONRequest("GET", "/", "", nil), provider.NewJSONResponse(200, nil))
	registered := []*Interaction{i}
	requested := []*Interaction{i, i, i}
	if err := verifyInteractions(registered, requested); err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf(errIntMultipleCalls, i.Description, i.State, len(requested))) {
			t.Errorf("expected message to contain: %s, actual message: %s",
				fmt.Sprintf(errIntMultipleCalls, i.Description, i.State, len(requested)), err.Error())
		}
	} else if err == nil {
		t.Error("expected multiple calls interaction error")
	}
}

func Test_ShouldVerify_UnexpectedInteractionVerfications(t *testing.T) {
	i, _ := NewInteraction("description", "state", provider.NewJSONRequest("GET", "/", "", nil), provider.NewJSONResponse(200, nil))
	registered := []*Interaction{}
	requested := []*Interaction{i}
	if err := verifyInteractions(registered, requested); err != nil {
		if !strings.Contains(err.Error(), errIntShouldNotBeCalled) {
			t.Errorf("expected message to contain: %s, actual message: %s", errIntShouldNotBeCalled, err.Error())
		}
	} else if err == nil {
		t.Error("expected unexpected calls interaction error")
	}
}

func Test_ShouldVerify_Interactions(t *testing.T) {
	i, _ := NewInteraction("description", "state", provider.NewJSONRequest("GET", "/", "", nil), provider.NewJSONResponse(200, nil))
	registered := []*Interaction{i}
	requested := []*Interaction{i}
	if err := verifyInteractions(registered, requested); err != nil {
		t.Errorf("expected verfication to succed, got error: %s", err.Error())
	}
}
