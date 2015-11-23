package comparers

import (
	"encoding/json"
	"io"
	"reflect"
)

func bodyMatches(expected, actual io.ReadCloser) (bool, error) {
	if expected == nil {
		return true, nil
	}

	var e, a interface{}
	decoder := json.NewDecoder(expected)
	err := decoder.Decode(&e)
	if err != nil {
		return false, err
	}

	decoder = json.NewDecoder(actual)
	err = decoder.Decode(&a)
	if err != nil {
		return false, err
	}

	expectedVal := e.(map[string]interface{})
	if expectedVal == nil {
		return true, nil
	}

	actualVal := a.(map[string]interface{})
	return reflect.DeepEqual(expectedVal, actualVal), nil
}
