package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/health"
	"github.com/docker/notary/signer/client"
	"golang.org/x/net/context"

	"github.com/docker/notary/server"
	"github.com/docker/notary/utils"
	"github.com/docker/notary/version"
	"github.com/spf13/viper"
)

// DebugAddress is the debug server address to listen on
const (
	jsonLogFormat = "json"
	DebugAddress  = "localhost:8080"
	// This is the generally recommended maximum age for Cache-Control headers
	// (one year, in seconds, since one year is forever in terms of internet
	// content)
	maxMaxAge = 60 * 60 * 24 * 365
)

var (
	debug      bool
	logFormat  string
	configFile string
	envPrefix  = "NOTARY_SERVER"
	mainViper  = viper.New()
)

func init() {
	utils.SetupViper(mainViper, envPrefix)
	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable the debugging server on localhost:8080")
	flag.StringVar(&logFormat, "logf", "json", "Set the format of the logs. Only 'json' and 'logfmt' are supported at the moment.")

	// this needs to be in init so that _ALL_ logs are in the correct format
	if logFormat == jsonLogFormat {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}
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

	currentCache, consistentCache, err := getCacheConfig(mainViper)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	httpAddr, tlsConfig, err := getAddrAndTLSConfig(mainViper)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	logrus.Info("Starting Server")
	err = server.Run(
		ctx,
		server.Config{
			Addr:                         httpAddr,
			TLSConfig:                    tlsConfig,
			Trust:                        trust,
			AuthMethod:                   mainViper.GetString("auth.type"),
			AuthOpts:                     mainViper.Get("auth.options"),
			CurrentCacheControlConfig:    currentCache,
			ConsistentCacheControlConfig: consistentCache,
		},
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
	logrus.Infof("Debug server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatalf("error listening on debug interface: %v", err)
	}
}
