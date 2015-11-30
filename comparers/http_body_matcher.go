package comparers

import (
	"encoding/json"
	"github.com/bennycao/pact-go/diff"
	"io"
)

func bodyMatches(expected, actual io.ReadCloser) (bool, diff.Differences, error) {
	if expected == nil {
		return true, nil, nil
	}

	var e, a interface{}
	decoder := json.NewDecoder(expected)
	err := decoder.Decode(&e)
	if err != nil {
		return false, nil, err
	}

	decoder = json.NewDecoder(actual)
	err = decoder.Decode(&a)
	if err != nil {
		return false, nil, err
	}

	if result, diffs := diff.DeepDiff(e, a, diff.DefaultConfig); result {
		return result, nil, nil
	} else {
		return result, diffs, nil
	}
}
