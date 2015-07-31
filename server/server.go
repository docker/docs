package server

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/registry/auth"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/signed"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/docker/notary/server/handlers"
	"github.com/docker/notary/utils"
)

func init() {
	data.SetDefaultExpiryTimes(
		map[string]int{
			"timestamp": 14,
		},
	)
}

// Run sets up and starts a TLS server that can be cancelled using the
// given configuration. The context it is passed is the context it should
// use directly for the TLS server, and generate children off for requests
func Run(ctx context.Context, addr, tlsCertFile, tlsKeyFile string, trust signed.CryptoService, authMethod string, authOpts interface{}) error {

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	var lsnr net.Listener
	lsnr, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	if tlsCertFile != "" && tlsKeyFile != "" {
		keypair, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
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

		logrus.Info("[Notary Server] Enabling TLS.")
		lsnr = tls.NewListener(lsnr, tlsConfig)
	} else if tlsCertFile != "" || tlsKeyFile != "" {
		return fmt.Errorf("Partial TLS configuration found. Either include both a cert and key file in the configuration, or include neither to disable TLS.")
	}

	var ac auth.AccessController
	if authMethod == "token" {
		authOptions, ok := authOpts.(map[string]interface{})
		if !ok {
			return fmt.Errorf("auth.options must be a map[string]interface{}")
		}
		ac, err = auth.GetAccessController(authMethod, authOptions)
		if err != nil {
			return err
		}
	}
	hand := utils.RootHandlerFactory(ac, ctx, trust)

	r := mux.NewRouter()
	r.Methods("GET").Path("/v2/").Handler(hand(handlers.MainHandler))
	r.Methods("POST").Path("/v2/{imageName:.*}/_trust/tuf/").Handler(hand(handlers.AtomicUpdateHandler, "push", "pull"))
	r.Methods("GET").Path("/v2/{imageName:.*}/_trust/tuf/{tufRole:(root|targets|snapshot)}.json").Handler(hand(handlers.GetHandler, "pull"))
	r.Methods("GET").Path("/v2/{imageName:.*}/_trust/tuf/timestamp.json").Handler(hand(handlers.GetTimestampHandler, "pull"))
	r.Methods("GET").Path("/v2/{imageName:.*}/_trust/tuf/timestamp.key").Handler(hand(handlers.GetTimestampKeyHandler, "push", "pull"))
	r.Methods("DELETE").Path("/v2/{imageName:.*}/_trust/tuf/").Handler(hand(handlers.DeleteHandler, "push", "pull"))
	r.Methods("GET", "POST", "PUT", "HEAD", "DELETE").Path("/{other:.*}").Handler(hand(utils.NotFoundHandler))
	svr := http.Server{
		Addr:    addr,
		Handler: r,
	}

	logrus.Info("[Notary Server] : Starting on ", addr)

	err = svr.Serve(lsnr)

	return err
}
