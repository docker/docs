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
)

var (
	debug      bool
	logFormat  string
	configFile string
	envPrefix  = "NOTARY_SERVER"
)

func init() {
	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable the debugging server on localhost:8080")
	flag.StringVar(&logFormat, "logf", "json", "Set the format of the logs. Only 'json' and 'logfmt' are supported at the moment.")

	// this needs to be in init so that _ALL_ logs are in the correct format
	if logFormat == jsonLogFormat {
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}
}

func parseServerConfig(configFilePath string, hRegister healthRegister) (context.Context, server.Config, error) {
	config := viper.New()
	utils.SetupViper(config, envPrefix)

	// parse viper config
	if err := utils.ParseViper(config, configFilePath); err != nil {
		return nil, server.Config{}, err
	}

	ctx := context.Background()

	// default is error level
	lvl, err := utils.ParseLogLevel(config, logrus.ErrorLevel)
	if err != nil {
		return nil, server.Config{}, err
	}
	logrus.SetLevel(lvl)

	// parse bugsnag config
	bugsnagConf, err := utils.ParseBugsnag(config)
	if err != nil {
		return ctx, server.Config{}, err
	}
	utils.SetUpBugsnag(bugsnagConf)

	trust, keyAlgo, err := getTrustService(config, client.NewNotarySigner, hRegister)
	if err != nil {
		return nil, server.Config{}, err
	}
	ctx = context.WithValue(ctx, "keyAlgorithm", keyAlgo)

	store, err := getStore(config, []string{utils.MySQLBackend, utils.MemoryBackend}, hRegister)
	if err != nil {
		return nil, server.Config{}, err
	}
	ctx = context.WithValue(ctx, "metaStore", store)

	currentCache, consistentCache, err := getCacheConfig(config)
	if err != nil {
		return nil, server.Config{}, err
	}

	httpAddr, tlsConfig, err := getAddrAndTLSConfig(config)
	if err != nil {
		return nil, server.Config{}, err
	}

	return ctx, server.Config{
		Addr:                         httpAddr,
		TLSConfig:                    tlsConfig,
		Trust:                        trust,
		AuthMethod:                   config.GetString("auth.type"),
		AuthOpts:                     config.Get("auth.options"),
		CurrentCacheControlConfig:    currentCache,
		ConsistentCacheControlConfig: consistentCache,
	}, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if debug {
		go debugServer(DebugAddress)
	}

	// when the server starts print the version for debugging and issue logs later
	logrus.Infof("Version: %s, Git commit: %s", version.NotaryVersion, version.GitCommit)

	ctx, serverConfig, err := parseServerConfig(configFile, health.RegisterPeriodicFunc)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	logrus.Info("Starting Server")
	err = server.Run(ctx, serverConfig)

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
