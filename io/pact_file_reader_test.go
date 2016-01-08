package io

import "testing"

func Test_FileReader_ValidFile_ShouldReturnPactFile(t *testing.T) {
	path := "../pact_examples/consumer-provider.json"
	r := NewPactFileReader(path)

	if _, err := r.Read(); err != nil {
		t.Error(err)
	}
}

func Test_FileReader_BadFilePath_ShouldReturnError(t *testing.T) {
	path := "../badpath/nofile.json"
	r := NewPactFileReader(path)

	if _, err := r.Read(); err == nil {
		t.Error("expected error")
	}
}
