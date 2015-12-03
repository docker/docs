package main

import (
	"crypto/tls"
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/health"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/token"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/signer/client"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"

	"github.com/docker/notary/server"
	"github.com/docker/notary/utils"
	"github.com/docker/notary/version"
	"github.com/spf13/viper"
)

// DebugAddress is the debug server address to listen on
const DebugAddress = "localhost:8080"

var (
	debug      bool
	configFile string
	envPrefix  = "NOTARY_SERVER"
	mainViper  = viper.New()
)

func init() {
	utils.SetupViper(mainViper, envPrefix)
	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable the debugging server on localhost:8080")
}

// get the address for the HTTP server, and parses the optional TLS
// configuration for the server - if no TLS configuration is specified,
// TLS is not enabled.
func getAddrAndTLSConfig(configuration *viper.Viper) (string, *tls.Config, error) {
	httpAddr := configuration.GetString("server.http_addr")
	if httpAddr == "" {
		return "", nil, fmt.Errorf("http listen address required for server")
	}

	tlsOpts, err := utils.ParseServerTLS(configuration, false)
	if err != nil {
		return "", nil, fmt.Errorf(err.Error())
	}
	// do not support this yet since the client doesn't have client cert support
	if tlsOpts != nil {
		tlsOpts.ClientCAFile = ""
		tlsConfig, err := utils.ConfigureServerTLS(tlsOpts)
		if err != nil {
			return "", nil, fmt.Errorf(
				"unable to set up TLS for server: %s", err.Error())
		}
		return httpAddr, tlsConfig, nil
	}
	return httpAddr, nil, nil
}

// sets up TLS for the GRPC connection to notary-signer
func grpcTLS(configuration *viper.Viper) (*tls.Config, error) {
	rootCA := utils.GetPathRelativeToConfig(configuration, "trust_service.tls_ca_file")
	serverName := configuration.GetString("trust_service.hostname")
	clientCert := utils.GetPathRelativeToConfig(configuration, "trust_service.tls_client_cert")
	clientKey := utils.GetPathRelativeToConfig(configuration, "trust_service.tls_client_key")

	if (clientCert == "" && clientKey != "") || (clientCert != "" && clientKey == "") {
		return nil, fmt.Errorf("Partial TLS configuration found. Either include both a client cert and client key file in the configuration, or include neither.")
	}

	tlsConfig, err := utils.ConfigureClientTLS(&utils.ClientTLSOpts{
		RootCAFile:     rootCA,
		ServerName:     serverName,
		ClientCertFile: clientCert,
		ClientKeyFile:  clientKey,
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
		// exactly unheatlthy, so always return healthy and just log an
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

func main() {
	flag.Usage = usage
	flag.Parse()

	if debug {
		go debugServer(DebugAddress)
	}

	// when the server starts print the version for debugging and issue logs later
	logrus.Infof("Version: %s, Git commit: %s", version.NotaryVersion, version.GitCommit)

	ctx := context.Background()

	filename := filepath.Base(configFile)
	ext := filepath.Ext(configFile)
	configPath := filepath.Dir(configFile)

	mainViper.SetConfigType(strings.TrimPrefix(ext, "."))
	mainViper.SetConfigName(strings.TrimSuffix(filename, ext))
	mainViper.AddConfigPath(configPath)

	err := mainViper.ReadInConfig()
	if err != nil {
		logrus.Error("Viper Error: ", err.Error())
		logrus.Error("Could not read config at ", configFile)
		os.Exit(1)
	}

	// default is error level
	lvl, err := utils.ParseLogLevel(mainViper, logrus.ErrorLevel)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.SetLevel(lvl)

	// parse bugsnag config
	bugsnagConf, err := utils.ParseBugsnag(mainViper)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	utils.SetUpBugsnag(bugsnagConf)

	trust, keyAlgo, err := getTrustService(mainViper,
		client.NewNotarySigner, health.RegisterPeriodicFunc)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	ctx = context.WithValue(ctx, "keyAlgorithm", keyAlgo)

	store, err := getStore(mainViper, []string{utils.MySQLBackend, utils.MemoryBackend})
	if err != nil {
		logrus.Fatal(err.Error())
	}
	ctx = context.WithValue(ctx, "metaStore", store)

	httpAddr, tlsConfig, err := getAddrAndTLSConfig(mainViper)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	logrus.Info("Starting Server")
	err = server.Run(
		ctx,
		httpAddr,
		tlsConfig,
		trust,
		mainViper.GetString("auth.type"),
		mainViper.Get("auth.options"),
	)

	logrus.Error(err.Error())
	return
}

func usage() {
	fmt.Println("usage:", os.Args[0])
	flag.PrintDefaults()
}

// debugServer starts the debug server with pprof, expvar among other
// endpoints. The addr should not be exposed externally. For most of these to
// work, tls cannot be enabled on the endpoint, so it is generally separate.
func debugServer(addr string) {
	logrus.Info("Debug server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatal("error listening on debug interface: ", err)
	}
}
