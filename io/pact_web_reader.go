package io

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type pactWebReader struct {
	url      string
	username string
	password string
}

func IsWebUri(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func NewPactWebReader(url, username, password string) PactReader {
	return &pactWebReader{url: url, username: username, password: password}
}

func (p *pactWebReader) Read() (*PactFile, error) {
	req, err := http.NewRequest("GET", p.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	if p.username != "" && p.password != "" {
		req.SetBasicAuth(p.username, p.password)
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
