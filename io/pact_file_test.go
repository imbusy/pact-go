package io

import "testing"

func Test_Validate_ValidFile(t *testing.T) {
	path := "../pact_examples/consumer-provider.json"
	r := NewPactFileReader(path)

	p, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	if err := p.Validate(); err != nil {
		t.Error(err)
	}
}
