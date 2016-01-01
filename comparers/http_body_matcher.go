package comparers

import (
	"encoding/json"
	"github.com/SEEK-Jobs/pact-go/diff"
	"io"
)

func bodyMatches(expected, actual io.Reader) (bool, diff.Differences, error) {
	if expected == nil {
		return true, nil, nil
	}

	var e, a interface{}
	decoder := json.NewDecoder(expected)
	err := decoder.Decode(&e)
	if err != nil {
		return false, nil, err
	}

	if actual != nil {
		decoder = json.NewDecoder(actual)
		err = decoder.Decode(&a)
		if err != nil {
			return false, nil, err
		}
	}

	if result, diffs := diff.DeepDiff(e, a, &diff.DiffConfig{AllowUnexpectedKeys: true, RootPath: "[\"body\"]"}); result {
		return result, nil, nil
	} else {
		return result, diffs, nil
	}
}
