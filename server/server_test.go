package server

import (
	"net"
	"strings"
	"testing"

	_ "github.com/docker/distribution/registry/auth/silly"
	"github.com/docker/notary/tuf/signed"
	"golang.org/x/net/context"
)

func TestRunBadAddr(t *testing.T) {
	err := Run(
		context.Background(),
		"testAddr",
		"../fixtures/notary-server.crt",
		"../fixtures/notary-server.crt",
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
		"../fixtures/notary-server.crt",
		"../fixtures/notary-server.crt",
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
