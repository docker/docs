package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/docker/distribution/registry/auth/silly"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/tuf/data"
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

// GetKeys supports only the timestamp and snapshot key endpoints
func TestGetKeysEndpoint(t *testing.T) {
	ctx := context.WithValue(
		context.Background(), "metaStore", storage.NewMemStorage())
	ctx = context.WithValue(ctx, "keyAlgorithm", data.ED25519Key)

	handler := RootHandler(nil, ctx, signed.NewEd25519())
	ts := httptest.NewServer(handler)
	defer ts.Close()

	rolesToStatus := map[string]int{
		data.CanonicalTimestampRole: http.StatusOK,
		data.CanonicalSnapshotRole:  http.StatusOK,
		data.CanonicalTargetsRole:   http.StatusNotFound,
		data.CanonicalRootRole:      http.StatusNotFound,
		"somerandomrole":            http.StatusNotFound,
	}

	for role, expectedStatus := range rolesToStatus {
		res, err := http.Get(
			fmt.Sprintf("%s/v2/gun/_trust/tuf/%s.key", ts.URL, role))
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	}
}
