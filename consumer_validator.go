package pact

import (
	"errors"
	"fmt"
	"github.com/imbusy/pact-go/comparers"
	"github.com/imbusy/pact-go/consumer"
	"github.com/imbusy/pact-go/diff"
	"github.com/imbusy/pact-go/io"
	"github.com/imbusy/pact-go/util"

	"net/http"
	"net/url"
)

type consumerValidator interface {
	ProviderService(c *http.Client, u *url.URL)
	CanValidate() error
	Validate(f *io.PactFile, states map[string]*stateAction) (bool, error)
}

var (
	errNilProviderClient        = errors.New("Provider http client cannot be nil, please provide a valid value using ServiceProvider function.")
	errNilProviderURL           = errors.New("Provider url cannot be nil, please provide a valid value using ServiceProvider function.")
)

type pactValidator struct {
	c        *http.Client
	u        *url.URL
	setup    Action
	teardown Action
	l        util.Logger
}

func newConsumerValidator(setup, teardown Action, l util.Logger) consumerValidator {
	return &pactValidator{setup: setup, teardown: teardown, l: l}
}

func (v *pactValidator) CanValidate() error {
	if v.c == nil {
		return errNilProviderClient
	} else if v.u == nil {
		return errNilProviderURL
	}
	return nil
}

func (v *pactValidator) ProviderService(c *http.Client, u *url.URL) {
	v.c = c
	v.u = u
}

func (v *pactValidator) Validate(p *io.PactFile, s map[string]*stateAction) (bool, error) {
	isValid := true

	for _, i := range p.Interactions {
		//default setup
		if err := v.executeAction(v.setup); err != nil {
			return false, err
		}

		//state setup
		var sa *stateAction

		if i.State != "" {
			if sa = s[i.State]; sa == nil {
				continue
			} else if err := v.executeAction(sa.setup); err != nil {
				return false, err
			}
		}

		//interaction validation
		if ok, err := v.validateInteraction(i); err != nil {
			return false, err
		} else if !ok {
			isValid = false
		}

		//state teardown
		if sa != nil {
			if err := v.executeAction(sa.teardown); err != nil {
				return false, err
			}
		}

		//default teardown
		if err := v.executeAction(v.teardown); err != nil {
			return false, err
		}

	}
	return isValid, nil
}

func (v *pactValidator) validateInteraction(i *consumer.Interaction) (bool, error) {
	req, err := i.ToHTTPRequest(v.u.String())
	if err != nil {
		return false, err
	}
	resp, err := v.c.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return false, err
	} else if diffs, err := comparers.MatchResponse(i.Response, resp); err != nil {
		return false, err
	} else if len(diffs) > 0 {
		diff.FormatDiff(diffs, v.l, fmt.Sprintf("The response for state '%s' did not match, the differences are below:", i.State))
		return false, nil
	}
	return true, nil
}

func (v *pactValidator) executeAction(a Action) error {
	if a != nil {
		if err := a(); err != nil {
			return err
		}
	}
	return nil
}
