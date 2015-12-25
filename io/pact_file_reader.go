package io

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type PactReader interface {
	Read() (*PactFile, error)
}

type pactFileReader struct {
	filePath string
}

var errIncompatiblePact = fmt.Errorf("Incompatible pact specification! We only support version %s.", pactSpecification)

func NewPactFileReader(filePath string) PactReader {
	return &pactFileReader{filePath: filePath}
}

func (r *pactFileReader) Read() (f *PactFile, err error) {
	var b []byte
	f = &PactFile{}
	if b, err = ioutil.ReadFile(r.filePath); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, f); err != nil {
		return nil, err
	}

	if f.Metadata.PactSpecification != pactSpecification {
		return nil, errIncompatiblePact
	}
	return f, nil
}
