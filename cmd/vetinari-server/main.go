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
	"github.com/gorilla/mux"
)

const ADDR = ":4443"
const DEBUG_ADDR = "localhost:8080"

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

	if DEBUG_ADDR != "" {
		go debugServer(DEBUG_ADDR)
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

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.MainHandler)

	server := http.Server{
		Addr:      ADDR,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	if debug {
		log.Println("[Vetinari Server] : Listening on", ADDR)
	}

	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("[Vetinari Server] : Failed to start %s", err)
	}
}

func usage() {
	log.Println(os.Stderr, "usage:", os.Args[0], "<config>")
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
