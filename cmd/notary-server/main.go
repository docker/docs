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
	"github.com/docker/notary/server"
	"github.com/docker/notary/version"
)

// DebugAddress is the debug server address to listen on
const (
	jsonLogFormat = "json"
	DebugAddress  = "localhost:8080"
)

var (
	debug       bool
	logFormat   string
	configFile  string
	envPrefix   = "NOTARY_SERVER"
	doBootstrap bool
)

func init() {
	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable the debugging server on localhost:8080")
	flag.StringVar(&logFormat, "logf", "json", "Set the format of the logs. Only 'json' and 'logfmt' are supported at the moment.")
	flag.BoolVar(&doBootstrap, "bootstrap", false, "Do any necessary setup of configured backend storage services")

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

	ctx, serverConfig, err := parseServerConfig(configFile, health.RegisterPeriodicFunc)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	if doBootstrap {
		err = bootstrap(ctx)
	} else {
		logrus.Info("Starting Server")
		err = server.Run(ctx, serverConfig)
	}

	if err != nil {
		logrus.Fatal(err.Error())
	}
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
