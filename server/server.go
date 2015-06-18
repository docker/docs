package server

import (
	"crypto/rand"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/registry/auth"
	"github.com/endophage/gotuf/signed"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/docker/vetinari/config"
	"github.com/docker/vetinari/server/handlers"
	"github.com/docker/vetinari/utils"
)

type HTTPServer struct {
	http.Server
	conns map[net.Conn]struct{}
	id    string
}

func NewHTTPServer(s http.Server) *HTTPServer {
	return &HTTPServer{
		Server: s,
		conns:  make(map[net.Conn]struct{}),
		id:     uuid.New(),
	}
}

// Track connections for cleanup on shutdown.
func (svr *HTTPServer) ConnState(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		svr.conns[conn] = struct{}{}
	case http.StateClosed, http.StateHijacked:
		delete(svr.conns, conn)
	}
}

// This should only be called after closing the server's listeners.
func (svr *HTTPServer) TimeoutConnections() {
	time.Sleep(time.Second * 30)
	for conn, _ := range svr.conns {
		conn.Close()
	}
	logrus.Infof("[Vetinari] All connections closed for server %s", svr.id)
}

// Run sets up and starts a TLS server that can be cancelled using the
// given configuration. The context it is passed is the context it should
// use directly for the TLS server, and generate children off for requests
func Run(ctx context.Context, conf config.ServerConf, trust signed.CryptoService) error {

	// TODO: check validity of config
	return run(ctx, conf.Addr, conf.TLSCertFile, conf.TLSKeyFile, trust)
}

func run(ctx context.Context, addr, tlsCertFile, tlsKeyFile string, trust signed.CryptoService) error {

	keypair, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
	if err != nil {
		logrus.Errorf("[Vetinari] Error loading keys %s", err)
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

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	lsnr, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	tlsLsnr := tls.NewListener(lsnr, tlsConfig)

	var ac auth.AccessController = nil
	//ac, err := auth.GetAccessController("token", map[string]interface{}{})
	//if err != nil {
	//	return err
	//}
	hand := utils.RootHandlerFactory(ac, ctx, trust)

	r := mux.NewRouter()
	// TODO (endophage): use correct regexes for image and tag names
	r.Methods("GET").Path("/v2/{imageName:.*}/_trust/tuf/{tufRole:(root|targets|timestamp|snapshot)}.json").Handler(hand(handlers.GetHandler, "pull"))
	r.Methods("POST").Path("/v2/{imageName:.*}/_trust/tuf/{tufRole:(root|targets|timestamp|snapshot)}.json").Handler(hand(handlers.UpdateHandler, "push", "pull"))

	svr := NewHTTPServer(
		http.Server{
			Addr:    addr,
			Handler: r,
		},
	)

	logrus.Info("[Vetinari] : Listening on", addr)

	go stopWatcher(ctx, svr, lsnr, tlsLsnr)

	err = svr.Serve(tlsLsnr)

	return err
}

func stopWatcher(ctx context.Context, svr *HTTPServer, ls ...net.Listener) {
	doneChan := ctx.Done()
	<-doneChan
	logrus.Debug("[Vetinari] Received close signal")
	for _, l := range ls {
		l.Close()
	}
	svr.TimeoutConnections()
}
