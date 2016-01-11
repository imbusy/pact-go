package io

import (
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
	username := "a-pact-user"
	password := "a-secret-password"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualUsername, actualPassword, ok := r.BasicAuth()

		if !ok || username != actualUsername || password != actualPassword {
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

	r := NewPactWebReader(s.URL, username, password)
	if f, err := r.Read(); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Error("expected valid pact file")
	}

	r = NewPactWebReader(s.URL, username, "incorrect password")
	if _, err := r.Read(); err == nil {
		t.Error("expected 401 error")
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
