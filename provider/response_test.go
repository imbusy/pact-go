package provider

import (
	"bytes"
	"net/http"
	"testing"
)

func Test_Response_MarshalJSON_BodyShouldExist(t *testing.T) {
	header := make(http.Header)
	header.Add("content-type", "application/json")
	response := NewJSONResponse(200, header)
	response.SetBody(`[]`)

	result, _ := response.MarshalJSON()

	if !bytes.Contains(result, []byte(`"body"`)) {
		t.Error(t, "Response should contain body field.")
	}
}
