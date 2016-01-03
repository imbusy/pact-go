package io

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type pactWebReader struct {
	url        string
	authScheme string
	authVal    string
}

func IsWebUri(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func NewPactWebReader(url, authScheme, authVal string) PactReader {
	return &pactWebReader{url: url, authScheme: authScheme, authVal: authVal}
}

func (p *pactWebReader) Read() (*PactFile, error) {
	req, err := http.NewRequest("GET", p.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	if p.authScheme != "" && p.authVal != "" {
		req.Header.Add("Authorisation", fmt.Sprintf("%s %s", p.authScheme, p.authVal))
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	} else if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get the pact file from %s, the response came back with %d status code", p.url, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var f PactFile
	if err := decoder.Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}
