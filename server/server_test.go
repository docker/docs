package server

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/endophage/go-tuf/signed"
	"golang.org/x/net/context"

	"github.com/docker/vetinari/config"
)

func TestRunBadCerts(t *testing.T) {
	err := Run(
		context.Background(),
		config.ServerConf{},
		signed.NewEd25519(),
	)
	if err == nil {
		t.Fatal("Passed empty certs, Run should have failed")
	}
}

func TestRunBadAddr(t *testing.T) {
	config := config.ServerConf{
		Addr:        "testAddr",
		TLSCertFile: "../fixtures/ca.pem",
		TLSKeyFile:  "../fixtures/ca-key.pem",
	}
	err := Run(context.Background(), config, signed.NewEd25519())
	if err == nil {
		t.Fatal("Passed bad addr, Run should have failed")
	}
}

func TestRunReservedPort(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())

	config := config.ServerConf{
		Addr:        "localhost:80",
		TLSCertFile: "../fixtures/vetinari.pem",
		TLSKeyFile:  "../fixtures/vetinari.key",
	}

	err := Run(ctx, config, signed.NewEd25519())

	if _, ok := err.(*net.OpError); !ok {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "bind: permission denied") {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
}

func TestRunGoodCancel(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	config := config.ServerConf{
		Addr:        "localhost:8002",
		TLSCertFile: "../fixtures/vetinari.pem",
		TLSKeyFile:  "../fixtures/vetinari.key",
	}

	go func() {
		time.Sleep(time.Second * 3)
		cancelFunc()
	}()

	err := Run(ctx, config, signed.NewEd25519())

	if _, ok := err.(*net.OpError); !ok {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "use of closed network connection") {
		t.Fatalf("Received unexpected err: %s", err.Error())
	}
}
