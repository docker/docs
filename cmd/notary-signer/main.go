// +build pkcs11

package main

import (
	"crypto/tls"
	"errors"
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/docker/distribution/health"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/signer/keydbstore"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/utils"
	"github.com/docker/notary/version"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"

	"github.com/Sirupsen/logrus"
	pb "github.com/docker/notary/proto"
)

const (
	jsonLogFormat   = "json"
	debugAddr       = "localhost:8080"
	envPrefix       = "NOTARY_SIGNER"
	defaultAliasEnv = "DEFAULT_ALIAS"
)

var (
	debug      bool
	logFormat  string
	configFile string
	mainViper  = viper.New()
)

func init() {
	utils.SetupViper(mainViper, envPrefix)
	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "show the version and exit")
	flag.StringVar(&logFormat, "logf", "json", "Set the format of the logs. Only 'json' and 'logfmt' are supported at the moment.")

	// this needs to be in init so that _ALL_ logs are in the correct format
	if logFormat == jsonLogFormat {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}
}

func passphraseRetriever(keyName, alias string, createNew bool, attempts int) (passphrase string, giveup bool, err error) {
	passphrase = mainViper.GetString(strings.ToUpper(alias))

	if passphrase == "" {
		return "", false, errors.New("expected env variable to not be empty: " + alias)
	}

	return passphrase, false, nil
}

// Reads the configuration file for storage setup, and sets up the cryptoservice
// mapping
func setUpCryptoservices(configuration *viper.Viper, allowedBackends []string) (
	signer.CryptoServiceIndex, error) {

	storeConfig, err := utils.ParseStorage(configuration, allowedBackends)
	if err != nil {
		return nil, err
	}

	var keyStore trustmanager.KeyStore
	if storeConfig.Backend == utils.MemoryBackend {
		keyStore = trustmanager.NewKeyMemoryStore(
			passphrase.ConstantRetriever("memory-db-ignore"))
	} else {
		defaultAlias := configuration.GetString("storage.default_alias")
		if defaultAlias == "" {
			// backwards compatibility - support this environment variable
			defaultAlias = configuration.GetString(defaultAliasEnv)
		}

		if defaultAlias == "" {
			return nil, fmt.Errorf("must provide a default alias for the key DB")
		}
		logrus.Debug("Default Alias: ", defaultAlias)

		dbStore, err := keydbstore.NewKeyDBStore(
			passphraseRetriever, defaultAlias, storeConfig.Backend, storeConfig.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to create a new keydbstore: %v", err)
		}

		health.RegisterPeriodicFunc(
			"DB operational", dbStore.HealthCheck, time.Second*60)
		keyStore = dbStore
	}

	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	cryptoServices := make(signer.CryptoServiceIndex)
	cryptoServices[data.ED25519Key] = cryptoService
	cryptoServices[data.ECDSAKey] = cryptoService
	return cryptoServices, nil
}

// set up the GRPC server
func setupGRPCServer(grpcAddr string, tlsConfig *tls.Config,
	cryptoServices signer.CryptoServiceIndex) (*grpc.Server, net.Listener, error) {

	//RPC server setup
	kms := &api.KeyManagementServer{CryptoServices: cryptoServices,
		HealthChecker: health.CheckStatus}
	ss := &api.SignerServer{CryptoServices: cryptoServices,
		HealthChecker: health.CheckStatus}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc server failed to listen on %s: %v",
			grpcAddr, err)
	}

	creds := credentials.NewTLS(tlsConfig)
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterKeyManagementServer(grpcServer, kms)
	pb.RegisterSignerServer(grpcServer, ss)

	return grpcServer, lis, nil
}

func setupHTTPServer(httpAddr string, tlsConfig *tls.Config,
	cryptoServices signer.CryptoServiceIndex) http.Server {

	return http.Server{
		Addr:      httpAddr,
		Handler:   api.Handlers(cryptoServices),
		TLSConfig: tlsConfig,
	}
}

func getAddrAndTLSConfig(configuration *viper.Viper) (string, string, *tls.Config, error) {
	tlsConfig, err := utils.ParseServerTLS(configuration, true)
	if err != nil {
		return "", "", nil, fmt.Errorf("unable to set up TLS: %s", err.Error())
	}

	grpcAddr := configuration.GetString("server.grpc_addr")
	if grpcAddr == "" {
		return "", "", nil, fmt.Errorf("grpc listen address required for server")
	}

	httpAddr := configuration.GetString("server.http_addr")
	if httpAddr == "" {
		return "", "", nil, fmt.Errorf("http listen address required for server")
	}

	return httpAddr, grpcAddr, tlsConfig, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if debug {
		go debugServer(debugAddr)
	}

	// when the signer starts print the version for debugging and issue logs later
	logrus.Infof("Version: %s, Git commit: %s", version.NotaryVersion, version.GitCommit)

	// parse viper config
	if err := utils.ParseViper(mainViper, configFile); err != nil {
		logrus.Fatal(err.Error())
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

	// parse server config
	httpAddr, grpcAddr, tlsConfig, err := getAddrAndTLSConfig(mainViper)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// setup the cryptoservices
	cryptoServices, err := setUpCryptoservices(mainViper,
		[]string{utils.MySQLBackend, utils.MemoryBackend})
	if err != nil {
		logrus.Fatal(err.Error())
	}

	grpcServer, lis, err := setupGRPCServer(grpcAddr, tlsConfig, cryptoServices)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	httpServer := setupHTTPServer(httpAddr, tlsConfig, cryptoServices)

	if debug {
		log.Println("RPC server listening on", grpcAddr)
		log.Println("HTTP server listening on", httpAddr)
	}

	go grpcServer.Serve(lis)
	err = httpServer.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal("HTTPS server failed to start:", err)
	}
}

func usage() {
	log.Println("usage:", os.Args[0], "<config>")
	flag.PrintDefaults()
}

// debugServer starts the debug server with pprof, expvar among other
// endpoints. The addr should not be exposed externally. For most of these to
// work, tls cannot be enabled on the endpoint, so it is generally separate.
func debugServer(addr string) {
	logrus.Infof("Debug server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatalf("error listening on debug interface: %v", err)
	}
}
