package main

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/health"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/token"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/signer/client"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// get the address for the HTTP server, and parses the optional TLS
// configuration for the server - if no TLS configuration is specified,
// TLS is not enabled.
func getAddrAndTLSConfig(configuration *viper.Viper) (string, *tls.Config, error) {
	httpAddr := configuration.GetString("server.http_addr")
	if httpAddr == "" {
		return "", nil, fmt.Errorf("http listen address required for server")
	}

	tlsConfig, err := utils.ParseServerTLS(configuration, false)
	if err != nil {
		return "", nil, fmt.Errorf(err.Error())
	}
	return httpAddr, tlsConfig, nil
}

// sets up TLS for the GRPC connection to notary-signer
func grpcTLS(configuration *viper.Viper) (*tls.Config, error) {
	rootCA := utils.GetPathRelativeToConfig(configuration, "trust_service.tls_ca_file")
	clientCert := utils.GetPathRelativeToConfig(configuration, "trust_service.tls_client_cert")
	clientKey := utils.GetPathRelativeToConfig(configuration, "trust_service.tls_client_key")

	if clientCert == "" && clientKey != "" || clientCert != "" && clientKey == "" {
		return nil, fmt.Errorf("either pass both client key and cert, or neither")
	}

	tlsConfig, err := tlsconfig.Client(tlsconfig.Options{
		CAFile:   rootCA,
		CertFile: clientCert,
		KeyFile:  clientKey,
	})
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to configure TLS to the trust service: %s", err.Error())
	}
	return tlsConfig, nil
}

// parses the configuration and returns a backing store for the TUF files
func getStore(configuration *viper.Viper, allowedBackends []string) (
	storage.MetaStore, error) {

	storeConfig, err := utils.ParseStorage(configuration, allowedBackends)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Using %s backend", storeConfig.Backend)

	if storeConfig.Backend == utils.MemoryBackend {
		return storage.NewMemStorage(), nil
	}

	store, err := storage.NewSQLStorage(storeConfig.Backend, storeConfig.Source)
	if err != nil {
		return nil, fmt.Errorf("Error starting DB driver: %s", err.Error())
	}
	health.RegisterPeriodicFunc(
		"DB operational", store.CheckHealth, time.Second*60)
	return store, nil
}

type signerFactory func(hostname, port string, tlsConfig *tls.Config) *client.NotarySigner
type healthRegister func(name string, checkFunc func() error, duration time.Duration)

// parses the configuration and determines which trust service and key algorithm
// to return
func getTrustService(configuration *viper.Viper, sFactory signerFactory,
	hRegister healthRegister) (signed.CryptoService, string, error) {

	switch configuration.GetString("trust_service.type") {
	case "local":
		logrus.Info("Using local signing service, which requires ED25519. " +
			"Ignoring all other trust_service parameters, including keyAlgorithm")
		return signed.NewEd25519(), data.ED25519Key, nil
	case "remote":
	default:
		return nil, "", fmt.Errorf(
			"must specify either a \"local\" or \"remote\" type for trust_service")
	}

	keyAlgo := configuration.GetString("trust_service.key_algorithm")
	if keyAlgo != data.ED25519Key && keyAlgo != data.ECDSAKey && keyAlgo != data.RSAKey {
		return nil, "", fmt.Errorf("invalid key algorithm configured: %s", keyAlgo)
	}

	clientTLS, err := grpcTLS(configuration)
	if err != nil {
		return nil, "", err
	}

	logrus.Info("Using remote signing service")

	notarySigner := sFactory(
		configuration.GetString("trust_service.hostname"),
		configuration.GetString("trust_service.port"),
		clientTLS,
	)

	minute := 1 * time.Minute
	hRegister(
		"Trust operational",
		// If the trust service fails, the server is degraded but not
		// exactly unhealthy, so always return healthy and just log an
		// error.
		func() error {
			err := notarySigner.CheckHealth(minute)
			if err != nil {
				logrus.Error("Trust not fully operational: ", err.Error())
			}
			return nil
		},
		minute)
	return notarySigner, keyAlgo, nil
}

// Gets the cache configuration for GET-ting current and checksummed metadata
// This is mainly the max-age (an integer in seconds, just like in the
// Cache-Control header) for consistent (content-addressable) downloads and
// current (latest version) downloads. The max-age must be between 0 and 31536000
// (one year in seconds, which is the recommended maximum time data is cached),
// else parsing will return an error.  A max-age of 0 will disable caching for
// that type of download (consistent or current).
func getCacheConfig(configuration *viper.Viper) (utils.CacheControlConfig, utils.CacheControlConfig, error) {
	var cccs []utils.CacheControlConfig
	types := []string{"current_metadata", "metadata_by_checksum"}

	for _, optionName := range types {
		m := configuration.GetString(fmt.Sprintf("caching.max_age.%s", optionName))
		if m == "" {
			continue
		}
		seconds, err := strconv.Atoi(m)
		if err != nil || seconds < 0 || seconds > maxMaxAge {
			return nil, nil, fmt.Errorf(
				"must specify a cache-control max-age between 0 and %v", maxMaxAge)
		}

		cccs = append(cccs, utils.NewCacheControlConfig(seconds, optionName == "current_metadata"))
	}
	return cccs[0], cccs[1], nil
}
