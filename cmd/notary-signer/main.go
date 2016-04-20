package main

import (
	"crypto/tls"
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary"
	"github.com/docker/notary/storage"
	"github.com/docker/notary/utils"
	"github.com/docker/notary/version"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const (
	jsonLogFormat   = "json"
	debugAddr       = "localhost:8080"
	envPrefix       = "NOTARY_SIGNER"
	defaultAliasEnv = "DEFAULT_ALIAS"
)

var (
	debug       bool
	logFormat   string
	configFile  string
	mainViper   = viper.New()
	doBootstrap bool
)

func init() {
	utils.SetupViper(mainViper, envPrefix)
	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "show the version and exit")
	flag.StringVar(&logFormat, "logf", "json", "Set the format of the logs. Only 'json' and 'logfmt' are supported at the moment.")
	flag.BoolVar(&doBootstrap, "bootstrap", false, "Do any necessary setup of configured backend storage services")

	// this needs to be in init so that _ALL_ logs are in the correct format
	if logFormat == jsonLogFormat {
		logrus.SetFormatter(new(logrus.JSONFormatter))
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
		[]string{notary.MySQLBackend, notary.MemoryBackend})
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

func bootstrap(s interface{}) error {
	store, ok := s.(storage.Bootstrapper)
	if !ok {
		return fmt.Errorf("Store does not support bootstrapping.")
	}
	return store.Bootstrap()
}
