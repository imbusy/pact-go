package provider

import (
	"net/http"
)

type ProviderResponse struct {
	Status  int
	Headers http.Header
	Body    string
}
