package main

import (
	_ "expvar"
	"flag"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/docker/vetinari/config"
	"github.com/docker/vetinari/server"
)

// DebugAddress is the debug server address to listen on
const DebugAddress = "localhost:8080"

var debug bool
var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "show the version and exit")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if DebugAddress != "" {
		go debugServer(DebugAddress)
	}

	conf, err := config.Load(configFile)
	if err != nil {
		// TODO: log and exit
	}

	ctx := context.Background()
	childCtx, _ := context.WithCancel(ctx)
	server.Run(childCtx, conf)

}

func usage() {
	log.Println("usage:", os.Args[0], "<config>")
	flag.PrintDefaults()
}

// debugServer starts the debug server with pprof, expvar among other
// endpoints. The addr should not be exposed externally. For most of these to
// work, tls cannot be enabled on the endpoint, so it is generally separate.
func debugServer(addr string) {
	log.Println("[Vetinari Debug Server] server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[Vetinari Debug Server] error listening on debug interface: %v", err)
	}
}
