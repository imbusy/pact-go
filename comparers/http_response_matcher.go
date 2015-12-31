package comparers

import (
	"bytes"
	"github.com/SEEK-Jobs/pact-go/diff"
	"github.com/SEEK-Jobs/pact-go/provider"
	"io"
	"net/http"
)

func MatchResponse(expected *provider.Response, actual *http.Response) (diff.Differences, error) {
	diffs := make(diff.Differences, 0)

	b, err := expected.GetData()
	if err != nil {
		return nil, err
	}
	var expBody io.Reader
	if len(b) > 0 {
		expBody = bytes.NewReader(b)
	}
	if res, sDiff := diff.DeepDiff(expected.Status, actual.StatusCode,
		&diff.DiffConfig{AllowUnexpectedKeys: true, RootPath: "[\"status\"]"}); !res {
		diffs = append(diffs, sDiff...)
	} else if res, hDiff := headerMatches(expected.Headers, actual.Header); !res {
		diffs = append(diffs, hDiff...)
	} else if res, bDiff, err := bodyMatches(expBody, actual.Body); err != nil {
		return nil, err
	} else if !res {
		diffs = append(diffs, bDiff...)
	}

	return diffs, nil
}
