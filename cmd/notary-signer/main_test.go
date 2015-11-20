// +build pkcs11

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/notary/signer"
	"github.com/docker/notary/tuf/data"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	Cert = "../../fixtures/notary-signer.crt"
	Key  = "../../fixtures/notary-signer.key"
	Root = "../../fixtures/root-ca.crt"
)

// initializes a viper object with test configuration
func configure(jsonConfig string) *viper.Viper {
	config := viper.New()
	config.SetConfigType("json")
	config.ReadConfig(bytes.NewBuffer([]byte(jsonConfig)))
	return config
}

func TestGetAddrAndTLSConfigInvalidTLS(t *testing.T) {
	invalids := []string{
		`{"server": {"http_addr": ":1234", "grpc_addr": ":2345"}}`,
		`{"server": {
				"http_addr": ":1234",
				"grpc_addr": ":2345",
				"tls_cert_file": "nope",
				"tls_key_file": "nope"
		}}`,
	}
	for _, configJSON := range invalids {
		_, _, _, err := getAddrAndTLSConfig(configure(configJSON))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to set up TLS")
	}
}

func TestGetAddrAndTLSConfigNoGRPCAddr(t *testing.T) {
	_, _, _, err := getAddrAndTLSConfig(configure(fmt.Sprintf(`{
		"server": {
			"http_addr": ":1234",
			"tls_cert_file": "%s",
			"tls_key_file": "%s"
		}
	}`, Cert, Key)))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grpc listen address required for server")
}

func TestGetAddrAndTLSConfigNoHTTPAddr(t *testing.T) {
	_, _, _, err := getAddrAndTLSConfig(configure(fmt.Sprintf(`{
		"server": {
			"grpc_addr": ":1234",
			"tls_cert_file": "%s",
			"tls_key_file": "%s"
		}
	}`, Cert, Key)))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "http listen address required for server")
}

func TestGetAddrAndTLSConfigSuccess(t *testing.T) {
	httpAddr, grpcAddr, tlsConf, err := getAddrAndTLSConfig(configure(fmt.Sprintf(`{
		"server": {
			"http_addr": ":2345",
			"grpc_addr": ":1234",
			"tls_cert_file": "%s",
			"tls_key_file": "%s"
		}
	}`, Cert, Key)))
	assert.NoError(t, err)
	assert.Equal(t, ":2345", httpAddr)
	assert.Equal(t, ":1234", grpcAddr)
	assert.NotNil(t, tlsConf)
}

func TestSetupCryptoServicesNoDefaultAlias(t *testing.T) {
	tmpFile, err := ioutil.TempFile("/tmp", "sqlite3")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = setUpCryptoservices(
		configure(fmt.Sprintf(
			`{"storage": {"backend": "sqlite3", "db_url": "%s"}}`,
			tmpFile.Name())),
		[]string{"sqlite3"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must provide a default alias for the key DB")
}

func TestSetupCryptoServicesSuccess(t *testing.T) {
	tmpFile, err := ioutil.TempFile("/tmp", "sqlite3")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	cryptoServices, err := setUpCryptoservices(
		configure(fmt.Sprintf(
			`{"storage": {"backend": "sqlite3", "db_url": "%s"},
			"default_alias": "timestamp"}`,
			tmpFile.Name())),
		[]string{"sqlite3"})
	assert.NoError(t, err)
	assert.Len(t, cryptoServices, 2)

	edService, ok := cryptoServices[data.ED25519Key]
	assert.True(t, ok)

	ecService, ok := cryptoServices[data.ECDSAKey]
	assert.True(t, ok)

	assert.Equal(t, edService, ecService)
}

func TestSetupHTTPServer(t *testing.T) {
	httpServer := setupHTTPServer(":4443", nil, make(signer.CryptoServiceIndex))
	assert.Equal(t, ":4443", httpServer.Addr)
	assert.Nil(t, httpServer.TLSConfig)
}

func TestSetupGRPCServerInvalidAddress(t *testing.T) {
	_, _, err := setupGRPCServer("nope", nil, make(signer.CryptoServiceIndex))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "grpc server failed to listen on nope")
}

func TestSetupGRPCServerSuccess(t *testing.T) {
	tlsConf := tls.Config{InsecureSkipVerify: true}
	grpcServer, lis, err := setupGRPCServer(":7899", &tlsConf,
		make(signer.CryptoServiceIndex))
	defer lis.Close()
	assert.NoError(t, err)
	assert.Equal(t, "[::]:7899", lis.Addr().String())
	assert.Equal(t, "tcp", lis.Addr().Network())
	assert.NotNil(t, grpcServer)
}
