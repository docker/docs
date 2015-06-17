package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/endophage/gotuf/signed"

	"github.com/docker/vetinari/utils"
)

func TestMainHandlerGet(t *testing.T) {
	hand := utils.RootHandlerFactory(nil, context.Background(), &signed.Ed25519{})
	handler := hand(MainHandler)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	_, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Received error on GET /: %s", err.Error())
	}
}

func TestMainHandlerNotGet(t *testing.T) {
	hand := utils.RootHandlerFactory(nil, context.Background(), &signed.Ed25519{})
	handler := hand(MainHandler)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Head(ts.URL)
	if err != nil {
		t.Fatalf("Received error on GET /: %s", err.Error())
	}
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404, received %d", res.StatusCode)
	}
}
