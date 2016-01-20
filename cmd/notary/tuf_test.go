package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenAuth(t *testing.T) {
	var (
		readOnly      bool
		baseTransport = &http.Transport{}
		gun           = "test"
	)
	require.Nil(t, tokenAuth("https://localhost:9999", baseTransport, gun, readOnly))
}

func NotFoundTestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func TestTokenAuthNon200Status(t *testing.T) {
	var (
		readOnly      bool
		baseTransport = &http.Transport{}
		gun           = "test"
	)
	s := httptest.NewServer(http.HandlerFunc(NotFoundTestHandler))
	defer s.Close()

	require.Nil(t, tokenAuth(s.URL, baseTransport, gun, readOnly))
}
