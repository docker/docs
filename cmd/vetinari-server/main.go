package main

import (
	_ "expvar"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

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

	ctx := context.Background()

	conf, err := parseConfig(configFile)
	if err != nil {
		log.Fatalf("Error parsing config: %s", err.Error())
	}

	sigHup := make(chan os.Signal)
	sigTerm := make(chan os.Signal)

	signal.Notify(sigHup, syscall.SIGHUP)
	signal.Notify(sigTerm, syscall.SIGTERM)

	for {
		log.Println("[Vetinari] Starting Server")
		childCtx, cancel := context.WithCancel(ctx)
		go server.Run(childCtx, conf)

		for {
			select {
			// On a sighup we cancel and restart a new server
			// with updated config
			case <-sigHup:
				log.Printf("[Vetinari] Server restart requested. Attempting to parse config at %s", configFile)
				conf, err = parseConfig(configFile)
				if err != nil {
					log.Printf("[Vetinari] Unable to parse config. Old configuration will keep running. Parse Err: %s", err.Error())
					continue
				} else {
					cancel()
					log.Println("[Vetinari] Stopping server for restart")
					break
				}
			// On sigkill we cancel and shutdown
			case <-sigTerm:
				cancel()
				log.Println("[Vetinari] Shutting Down Hard")
				os.Exit(0)
			}
		}
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
	log.Println("[Vetinari Debug Server] server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[Vetinari Debug Server] error listening on debug interface: %v", err)
	}
}

func parseConfig(path string) (*config.Configuration, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	return config.Load(file)
}
