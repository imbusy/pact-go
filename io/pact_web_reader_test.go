package io

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WebReader_ReadsPactWithoutAuth(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "../pact_examples/consumer-provider.json"
		b, err := ioutil.ReadFile(path)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		if _, err := w.Write(b); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
	}))
	defer s.Close()

	r := NewPactWebReader(s.URL, "", "")
	if f, err := r.Read(); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Error("expected valid pact file")
	}
}

func Test_WebReader_ReadsPactWithAuth(t *testing.T) {
	authSchme := "Basic"
	authVal := "some-encrypted-token"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorisation")
		if authHeader != fmt.Sprintf("%s %s", authSchme, authVal) {
			w.WriteHeader(http.StatusUnauthorized)
		}

		path := "../pact_examples/consumer-provider.json"
		b, err := ioutil.ReadFile(path)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		if _, err := w.Write(b); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
	}))
	defer s.Close()

	r := NewPactWebReader(s.URL, authSchme, authVal)
	if f, err := r.Read(); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Error("expected valid pact file")
	}
}

func Test_WebReader_ReturnsErrorWhenWebRequestFails(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer s.Close()

	r := NewPactWebReader(s.URL, "", "")
	if _, err := r.Read(); err == nil {
		t.Error("expected error")
	}
}
