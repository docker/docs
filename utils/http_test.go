package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/endophage/gotuf/signed"

	"github.com/docker/notary/errors"
)

func MockContextHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	return nil
}

func MockBetterErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) *errors.HTTPError {
	return &errors.HTTPError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       9999,
		Err:        fmt.Errorf("TestError"),
	}
}

func TestRootHandlerFactory(t *testing.T) {
	hand := RootHandlerFactory(nil, context.Background(), &signed.Ed25519{})
	handler := hand(MockContextHandler)
	if _, ok := interface{}(handler).(http.Handler); !ok {
		t.Fatalf("A rootHandler must implement the http.Handler interface")
	}

	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, received %d", res.StatusCode)
	}
}

//func TestRootHandlerUnauthorized(t *testing.T) {
//	hand := RootHandlerFactory(nil, context.Background(), &signed.Ed25519{})
//	handler := hand(MockContextHandler)
//
//	ts := httptest.NewServer(handler)
//	defer ts.Close()
//
//	res, err := http.Get(ts.URL)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if res.StatusCode != http.StatusUnauthorized {
//		t.Fatalf("Expected 401, received %d", res.StatusCode)
//	}
//}

func TestRootHandlerError(t *testing.T) {
	hand := RootHandlerFactory(nil, context.Background(), &signed.Ed25519{})
	handler := hand(MockBetterErrorHandler)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected 500, received %d", res.StatusCode)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := strings.Trim(string(content), "\r\n\t ")
	if contentStr != "9999: TestError" {
		t.Fatalf("Error Body Incorrect: `%s`", content)
	}
}
