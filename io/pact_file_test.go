package io

import "testing"

func Test_Validate_ValidFile(t *testing.T) {
	path := "../pact_examples/consumer-provider.json"
	p := readPactFile(t, path)

	if err := p.Validate(); err != nil {
		t.Error(err)
	}
}

func Test_Validate_MissingProvider(t *testing.T) {
	path := "./pactNoProviderSpec.json"
	p := readPactFile(t, path)

	expectError(t, p.Validate(), errEmptyProvider)
}

func Test_Validate_MissingConsumer(t *testing.T) {
	path := "./pactNoConsumerSpec.json"
	p := readPactFile(t, path)

	expectError(t, p.Validate(), errEmptyConsumer)
}

func Test_Validate_InvalidSpec(t *testing.T) {
	path := "./pactWrongSpec.json"
	p := readPactFile(t, path)

	expectError(t, p.Validate(), errIncompatiblePact)
}

func readPactFile(t *testing.T, path string) *PactFile {
	r := NewPactFileReader(path)
	p, err := r.Read()
	if err != nil {
		t.Error(err)
	}
	return p
}

func expectError(t *testing.T, actual, expected error) {
	if actual == nil {
		t.Error("expected an error")
	} else if actual != expected {
		t.Errorf("got %q error, expected %q", actual, expected)
	}
}
