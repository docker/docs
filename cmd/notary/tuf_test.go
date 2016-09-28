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
	auth, err := tokenAuth("https://localhost:9999", baseTransport, gun, readOnly)
	require.NoError(t, err)
	require.Nil(t, auth)
}

func StatusOKTestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("{}"))
}

func TestTokenAuth200Status(t *testing.T) {
	var (
		readOnly      bool
		baseTransport = &http.Transport{}
		gun           = "test"
	)
	s := httptest.NewServer(http.HandlerFunc(NotAuthorizedTestHandler))
	defer s.Close()

	auth, err := tokenAuth(s.URL, baseTransport, gun, readOnly)
	require.NoError(t, err)
	require.NotNil(t, auth)
}

func NotAuthorizedTestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(401)
}

func TestTokenAuth401Status(t *testing.T) {
	var (
		readOnly      bool
		baseTransport = &http.Transport{}
		gun           = "test"
	)
	s := httptest.NewServer(http.HandlerFunc(NotAuthorizedTestHandler))
	defer s.Close()

	auth, err := tokenAuth(s.URL, baseTransport, gun, readOnly)
	require.NoError(t, err)
	require.NotNil(t, auth)
}

func NotFoundTestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func TestTokenAuthNon200Non401Status(t *testing.T) {
	var (
		readOnly      bool
		baseTransport = &http.Transport{}
		gun           = "test"
	)
	s := httptest.NewServer(http.HandlerFunc(NotFoundTestHandler))
	defer s.Close()

	auth, err := tokenAuth(s.URL, baseTransport, gun, readOnly)
	require.NoError(t, err)
	require.Nil(t, auth)
}
