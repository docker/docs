package server

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/docker/distribution/registry/auth/silly"
	"github.com/docker/notary/tuf/signed"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestRunBadAddr(t *testing.T) {
	err := Run(
		context.Background(),
		"testAddr",
		nil,
		signed.NewEd25519(),
		"",
		nil,
	)
	if err == nil {
		t.Fatal("Passed bad addr, Run should have failed")
	}
}

func TestRunReservedPort(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())

	err := Run(
		ctx,
		"localhost:80",
		nil,
		signed.NewEd25519(),
		"",
		nil,
	)

	if _, ok := err.(*net.OpError); !ok {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "bind: permission denied") {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
}

func TestMetricsEndpoint(t *testing.T) {
	handler := RootHandler(nil, context.Background(), signed.NewEd25519())
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/_notary_server/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
