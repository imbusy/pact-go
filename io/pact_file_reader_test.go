package io

import "testing"

func Test_Reader_ValidFile_ShouldReturnPactFile(t *testing.T) {
	path := "../pact_examples/consumer-provider.json"
	r := NewPactFileReader(path)

	if _, err := r.Read(); err != nil {
		t.Error(err)
	}
}

func Test_Reader_BadFilePath_ShouldReturnError(t *testing.T) {
	path := "../badpath/nofile.json"
	r := NewPactFileReader(path)

	if _, err := r.Read(); err == nil {
		t.Error("expected error")
	}
}

func Test_Reader_InvalidSpec_ShouldReturnError(t *testing.T) {
	path := "./pactWrongSpec.json"
	r := NewPactFileReader(path)

	if _, err := r.Read(); err == nil {
		t.Error("expected not supported pact error")
	} else if err != errIncompatiblePact {
		t.Error("got the wrong error")
	}
}
