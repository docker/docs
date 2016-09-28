package main

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/cfssl/api"
	"github.com/cloudflare/cfssl/api/info"
	"github.com/cloudflare/cfssl/errors"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/ca/config"
)

type caServer struct {
	signer signer.Signer

	*http.ServeMux
}

func newCAServer(s signer.Signer) (*caServer, error) {
	server := &caServer{
		signer:   s,
		ServeMux: http.NewServeMux(),
	}

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	infoHandler, err := info.NewHandler(s)
	if err != nil {
		return nil, fmt.Errorf("unable to create info handler: %v", err)
	}

	server.Handle(config.APIInfoPath, infoHandler)
	server.HandleFunc(config.APISignPath, func(w http.ResponseWriter, r *http.Request) {
		response := server.handleSignRequest(r)
		if !response.Success {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode(response)
	})

	return server, nil
}

func hasControllerAuth(clientSubjectName pkix.Name) bool {
	if clientSubjectName.CommonName == "controller" {
		return true
	}

	// Newer controller certs with have an OU set to "ucp".
	for _, ou := range clientSubjectName.OrganizationalUnit {
		if ou == "ucp" {
			return true
		}
	}

	return false
}

func hasSwarmManagerAuth(clientSubjectName pkix.Name) bool {
	for _, ou := range clientSubjectName.OrganizationalUnit {
		if ou == "swarm-manager" {
			return true
		}
	}

	return false
}

// makeErrResponse is a helper function for creating api responses which
// represent an error.
func makeErrResponse(category errors.Category, reason errors.Reason, message string) *api.Response {
	err := errors.New(category, reason)
	errResponse := api.NewErrorResponse(message, err.ErrorCode)
	return &errResponse
}

// Decode the certificate signing request.
func decodeAndParseCSR(r io.Reader, s signer.Signer) (*signer.SignRequest, *pkix.Name, *api.Response) {
	var signReq signer.SignRequest
	if err := json.NewDecoder(r).Decode(&signReq); err != nil {
		return nil, nil, makeErrResponse(
			errors.APIClientError, errors.JSONError,
			fmt.Sprintf("unable to decode sign request JSON: %v", err),
		)
	}

	csrBlock, _ := pem.Decode([]byte(signReq.Request))
	if csrBlock == nil {
		return nil, nil, makeErrResponse(
			errors.CSRError, errors.DecodeFailed,
			"unable to decode csr PEM",
		)
	}

	if csrBlock.Type != "CERTIFICATE REQUEST" {
		return nil, nil, makeErrResponse(
			errors.CSRError, errors.BadRequest,
			"not a certificate or csr",
		)
	}

	// Get the subject name from the request.
	csr, err := signer.ParseCertificateRequest(s, csrBlock.Bytes)
	if err != nil {
		return nil, nil, makeErrResponse(
			errors.CSRError, errors.BadRequest,
			"unable to parse csr",
		)
	}

	requestSubject := local.PopulateSubjectFromCSR(signReq.Subject, csr.Subject)

	return &signReq, &requestSubject, nil
}

// Issue a new certificate from a given Certificate Signing Request.
func (s *caServer) handleSignRequest(r *http.Request) *api.Response {
	// Require client authentication via mutual TLS.
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		return makeErrResponse(
			errors.APIClientError, errors.AuthenticationFailure,
			"must authenticate sign request with mutual TLS",
		)
	}
	clientSubjectName := r.TLS.PeerCertificates[0].Subject

	// In order to be authorized, the client certificate subject must have
	// the common name "controller" OR must have an organizational unit of
	// "swarm-manager".
	isController := hasControllerAuth(clientSubjectName)
	isSwarmManager := hasSwarmManagerAuth(clientSubjectName)

	if !(isController || isSwarmManager) {
		return makeErrResponse(
			errors.APIClientError, errors.AuthenticationFailure,
			"must be authenticated as a UCP Controller or Swarm Manager",
		)
	}

	signRequest, requestSubject, errResponse := decodeAndParseCSR(r.Body, s.signer)
	if errResponse != nil {
		return errResponse
	}

	if isSwarmManager {
		// The client certificate must have an Org (the cluster ID).
		if len(clientSubjectName.Organization) == 0 {
			return makeErrResponse(
				errors.APIClientError, errors.AuthenticationFailure,
				"Swarm Manager client certificate does not specify cluster ID",
			)
		}

		clusterID := clientSubjectName.Organization[0]

		// The Org in the sign request subject should match the cluster
		// ID.
		if len(requestSubject.Organization) == 0 || requestSubject.Organization[0] != clusterID {
			return makeErrResponse(
				errors.CSRError, errors.BadRequest,
				"sign request subject org does not match cluster ID",
			)
		}
	}

	// Finally, sign the requested certificate.
	certificate, err := s.signer.Sign(*signRequest)
	if err != nil {
		return makeErrResponse(
			errors.APIClientError, errors.ServerRequestFailed,
			fmt.Sprintf("unable to sign requested certificate: %v", err),
		)
	}

	successResponse := api.NewSuccessResponse(
		map[string]string{
			"certificate": string(certificate),
		},
	)

	return &successResponse
}

func runServer(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if c.Bool("jsonlog") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	localSigner, err := local.NewSignerFromFile(c.String("ca"), c.String("ca-key"), config.Signing())
	if err != nil {
		return fmt.Errorf("unable to load local signer: %v", err)
	}

	serverCert, err := tls.LoadX509KeyPair(c.String("tls-cert"), c.String("tls-key"))
	if err != nil {
		return fmt.Errorf("unable to load TLS server key pair: %v", err)
	}

	clientCAPEM, err := ioutil.ReadFile(c.String("mutual-tls-ca"))
	if err != nil {
		return fmt.Errorf("unable to load TLS client CA certificate: %v", err)
	}

	clientCAPool := x509.NewCertPool()
	if !clientCAPool.AppendCertsFromPEM(clientCAPEM) {
		return fmt.Errorf("no valid client CA certificates found in %s", c.String("mutual-tls-ca"))
	}

	caServer, err := newCAServer(localSigner)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(c.String("address"), c.String("port")),
		Handler: caServer,
		TLSConfig: &tls.Config{
			ClientCAs:    clientCAPool,
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
		},
	}

	return httpServer.ListenAndServeTLS("", "")
}

func main() {
	app := cli.NewApp()
	app.Name = "ca"
	app.Usage = "Certificate Authority"

	serveCmd := cli.Command{
		Name:  "serve",
		Usage: "run a CA signing server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "address",
				Value: "localhost",
				Usage: "host address on which to listen for connections",
			},
			cli.StringFlag{
				Name:  "port",
				Value: "8080",
				Usage: "host port on which to listen for connections",
			},
			cli.StringFlag{
				Name:  "ca",
				Value: "/etc/ucp-ca/ca.pem",
				Usage: "filename of CA certificate",
			},
			cli.StringFlag{
				Name:  "ca-key",
				Value: "/etc/ucp-ca/ca-key.pem",
				Usage: "filename of CA private key",
			},
			cli.StringFlag{
				Name:  "tls-cert",
				Value: "/etc/ucp-ca/server.pem",
				Usage: "filename of CA server TLS certificate",
			},
			cli.StringFlag{
				Name:  "tls-key",
				Value: "/etc/ucp-ca/server-key.pem",
				Usage: "filename of CA server TLS key",
			},
			cli.StringFlag{
				Name:  "mutual-tls-ca",
				Value: "/etc/ucp-ca/client-ca.pem",
				Usage: "filename of CA server TLS client CA",
			},
			cli.StringFlag{
				Name:  "mutual-tls-cn",
				Value: "controller",
				Usage: "common name of UCP controller cert (deprecated)",
			},
			cli.StringFlag{
				Name:  "config",
				Value: "/etc/ucp-ca/config.json",
				Usage: "cfssl signing policy configuration file (deprecated)",
			},
			cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug logs",
			},
			cli.BoolFlag{
				Name:  "jsonlog",
				Usage: "format logs as JSON",
			},
		},
		Action: runServer,
	}

	app.Commands = []cli.Command{serveCmd}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
