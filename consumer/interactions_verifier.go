package consumer

import (
	"fmt"
	"strings"
)

var (
	errIntMultipleCalls      = "The interaction with description '%s' and provider state '%s', was used %d time/s by the test."
	errIntMissingUsage       = "The interaction with description '%s' and provider state '%s', was not verified."
	errIntShouldNotBeCalled  = "No interactions were registered, however the mock provider service was called."
	errIntVerificationFailed = "Verifying - actual interactions do not match expected interactions\n%s"
)

func verifyInteractions(registered, requested []*Interaction) error {
	var msgs []string

	if (registered == nil || len(registered) == 0) && (requested != nil && len(requested) > 0) {
		msgs = append(msgs, errIntShouldNotBeCalled)
	}

	for _, val := range registered {
		if result, count := contains(requested, val); !result {
			msgs = append(msgs, fmt.Sprintf(errIntMissingUsage, val.Description, val.State))
		} else if count > 1 {
			msgs = append(msgs, fmt.Sprintf(errIntMultipleCalls, val.Description, val.State, count))
		}
	}

	if len(msgs) > 0 {
		return fmt.Errorf(errIntVerificationFailed, strings.Join(msgs, ", "))
	}
	return nil
}

func contains(s []*Interaction, e *Interaction) (bool, int) {
	var count int
	for _, a := range s {
		if a == e {
			count++
		}
	}
	return (count > 0), count
}
