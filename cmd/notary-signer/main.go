package main

import (
	_ "expvar"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/version"
	_ "github.com/go-sql-driver/mysql"
)

const (
	jsonLogFormat = "json"
	debugAddr     = "localhost:8080"
)

var (
	debug       bool
	logFormat   string
	configFile  string
	doBootstrap bool
)

func init() {
	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Show the version and exit")
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
		go debugServer(debugAddr)
	}

	// when the signer starts print the version for debugging and issue logs later
	logrus.Infof("Version: %s, Git commit: %s", version.NotaryVersion, version.GitCommit)

	signerConfig, err := parseSignerConfig(configFile)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	grpcServer, lis, err := setupGRPCServer(signerConfig.GRPCAddr, signerConfig.TLSConfig, signerConfig.CryptoServices)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	httpServer := setupHTTPServer(signerConfig.HTTPAddr, signerConfig.TLSConfig, signerConfig.CryptoServices)

	if debug {
		log.Println("RPC server listening on", signerConfig.GRPCAddr)
		log.Println("HTTP server listening on", signerConfig.HTTPAddr)
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
