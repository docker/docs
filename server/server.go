package server

import (
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/endophage/go-tuf/signed"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/docker/vetinari/config"
	"github.com/docker/vetinari/server/handlers"
	"github.com/docker/vetinari/utils"
)

// Run sets up and starts a TLS server that can be cancelled using the
// given configuration. The context it is passed is the context it should
// use directly for the TLS server, and generate children off for requests
func Run(ctx context.Context, conf *config.Configuration) error {

	var trust signed.TrustService
	if conf.TrustService.Type == "remote" {
		log.Println("[Vetinari Server] : Using remote signing service")
		trust = newRufusSigner(conf.TrustService.Hostname, conf.TrustService.Port)
	} else {
		log.Println("[Vetinari Server] : Using local signing service")
		trust = signed.NewEd25519()
	}

	keypair, err := tls.LoadX509KeyPair(conf.Server.TLSCertFile, conf.Server.TLSKeyFile)
	if err != nil {
		return err
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
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{keypair},
		Rand:         rand.Reader,
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", conf.Server.Addr)
	if err != nil {
		return err
	}
	lsnr, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	tlsLsnr := tls.NewListener(lsnr, tlsConfig)

	// This is a basic way to shutdown the running listeners.
	// A more complete implementation would ensure each individual connection
	// gets cleaned up.
	go func() {
		doneChan := ctx.Done()
		<-doneChan
		// TODO: log that we received close signal
		lsnr.Close()
		tlsLsnr.Close()
	}()


	hand := utils.RootHandlerFactory(&utils.InsecureAuthorizer{}, utils.NewContext, trust)

	r := mux.NewRouter()
	// TODO (endophage): use correct regexes for image and tag names
	r.Methods("PUT").Path("/{imageName}/init").Handler(hand(handlers.GenKeysHandler, utils.SSCreate))
	r.Methods("GET").Path("/{imageName}/{tufFile}").Handler(hand(handlers.GetHandler, utils.SSNoAuth))
	r.Methods("DELETE").Path("/{imageName}:{tag}").Handler(hand(handlers.RemoveHandler, utils.SSDelete))
	r.Methods("POST").Path("/{imageName}:{tag}").Handler(hand(handlers.AddHandler, utils.SSUpdate))

	server := http.Server{
		Addr:    conf.Server.Addr,
		Handler: r,
	}

	log.Println("[Vetinari Server] : Listening on", conf.Server.Addr)

	err = server.Serve(tlsLsnr)

	return err
}
