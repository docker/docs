package main

import (
	"database/sql"
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	_ "github.com/docker/distribution/registry/auth/token"
	"github.com/endophage/gotuf/signed"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"

	"github.com/docker/vetinari/config"
	"github.com/docker/vetinari/server"
	"github.com/docker/vetinari/server/version"
	"github.com/docker/vetinari/signer"
)

// DebugAddress is the debug server address to listen on
const DebugAddress = "localhost:8080"

var debug bool
var configFile string

func init() {
	// Set default logging level to Error
	logrus.SetLevel(logrus.ErrorLevel)

	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable the debugging server on localhost:8080")
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
		logrus.Fatal("Error parsing config: ", err.Error())
		return // not strictly needed but let's be explicit
	}
	if conf.Logging.Level > 0 {
		logrus.SetLevel(logrus.Level(conf.Logging.Level))
	}

	sigHup := make(chan os.Signal)
	sigTerm := make(chan os.Signal)

	signal.Notify(sigHup, syscall.SIGHUP)
	signal.Notify(sigTerm, syscall.SIGTERM)

	var trust signed.CryptoService
	if conf.TrustService.Type == "remote" {
		logrus.Info("[Vetinari Server] : Using remote signing service")
		trust = signer.NewRufusSigner(conf.TrustService.Hostname, conf.TrustService.Port, conf.TrustService.TLSCAFile)
	} else {
		logrus.Info("[Vetinari] : Using local signing service")
		trust = signed.NewEd25519()
	}

	db, err := sql.Open("mysql", "dockercondemo:dockercondemo@tcp(mysql:3306)/dockercondemo")
	if err != nil {
		logrus.Fatal("Error starting DB driver: ", err.Error())
		return // not strictly needed but let's be explicit
	}
	ctx = context.WithValue(ctx, "versionStore", version.NewVersionDB(db))
	for {
		logrus.Info("[Vetinari] Starting Server")
		childCtx, cancel := context.WithCancel(ctx)
		go server.Run(childCtx, conf.Server, trust)

		for {
			select {
			// On a sighup we cancel and restart a new server
			// with updated config
			case <-sigHup:
				logrus.Infof("[Vetinari] Server restart requested. Attempting to parse config at %s", configFile)
				conf, err = parseConfig(configFile)
				if err != nil {
					logrus.Infof("[Vetinari] Unable to parse config. Old configuration will keep running. Parse Err: %s", err.Error())
					continue
				} else {
					cancel()
					logrus.Info("[Vetinari] Stopping server for restart")
					break
				}
			// On sigkill we cancel and shutdown
			case <-sigTerm:
				cancel()
				logrus.Info("[Vetinari] Shutting Down Hard")
				os.Exit(0)
			}
		}
	}
}

func usage() {
	fmt.Println("usage:", os.Args[0], "<config>")
	flag.PrintDefaults()
}

// debugServer starts the debug server with pprof, expvar among other
// endpoints. The addr should not be exposed externally. For most of these to
// work, tls cannot be enabled on the endpoint, so it is generally separate.
func debugServer(addr string) {
	logrus.Info("[Vetinari Debug Server] server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatal("[Vetinari Debug Server] error listening on debug interface: ", err)
	}
}

func parseConfig(path string) (*config.Configuration, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		logrus.Error("Failed to open configuration file located at: ", path)
		return nil, err
	}

	return config.Load(file)
}
