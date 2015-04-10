package main

import (
	"crypto/rand"
	"crypto/tls"
	_ "expvar"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/docker/distribution/health"
	"github.com/docker/vetinari/server/handlers"
	"github.com/docker/vetinari/utils"
	"github.com/gorilla/mux"
)

// ServerAddress is the secure server address to listen on
const ServerAddress = ":4443"

// DebugAddress is the debug server address to listen on
const DebugAddress = "localhost:8080"

var debug bool
var certFile, keyFile string

func init() {
	flag.StringVar(&certFile, "cert", "", "Intermediate certificates")
	flag.StringVar(&keyFile, "key", "", "Private key file")
	flag.BoolVar(&debug, "debug", false, "show the version and exit")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if DebugAddress != "" {
		go debugServer(DebugAddress)
	}

	if certFile == "" || keyFile == "" {
		usage()
		log.Fatalf("Certificate and key are mandatory")
	}

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA},
	}
	tlsConfig.Rand = rand.Reader

	hand := utils.RootHandlerFactory(&utils.InsecureAuthorizer{}, utils.ContextFactory)

	r := mux.NewRouter()
	// TODO (endophage): use correct regexes for image and tag names
	r.Methods("PUT").Path("/{imageName}/init").Handler(hand(handlers.GenKeysHandler, utils.SSCreate))
	r.Methods("GET").Path("/{imageName}/{tufFile}").Handler(hand(handlers.GetHandler, utils.SSNoAuth))
	r.Methods("DELETE").Path("/{imageName}/{tag}").Handler(hand(handlers.RemoveHandler, utils.SSDelete))
	r.Methods("POST").Path("/{imageName}/{tag}").Handler(hand(handlers.AddHandler, utils.SSUpdate))

	server := http.Server{
		Addr:      ServerAddress,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	if debug {
		log.Println("[Vetinari Server] : Listening on", ServerAddress)
	}

	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("[Vetinari Server] : Failed to start %s", err)
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
