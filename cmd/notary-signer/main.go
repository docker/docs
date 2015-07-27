package main

import (
	"crypto/rand"
	"crypto/tls"
	"database/sql"
	"errors"
	_ "expvar"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	_ "github.com/docker/distribution/health"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	_ "github.com/go-sql-driver/mysql"
	"github.com/miekg/pkcs11"

	pb "github.com/docker/notary/proto"
)

const (
	_Addr            = ":4444"
	_RpcAddr         = ":7899"
	_DebugAddr       = "localhost:8080"
	_DBType          = "mysql"
	_EnvPrefix       = "NOTARY_SIGNER"
	_DefaultAliasEnv = _EnvPrefix + "_DEFAULT_ALIAS"
)

var debug bool
var certFile, keyFile, pkcs11Lib, pin, dbURL string

func init() {
	flag.StringVar(&certFile, "cert", "", "Intermediate certificates")
	flag.StringVar(&keyFile, "key", "", "Private key file")
	flag.StringVar(&dbURL, "dburl", "", "URL of the database")
	flag.StringVar(&pkcs11Lib, "pkcs11", "", "enables HSM mode and uses the provided pkcs11 library path")
	flag.StringVar(&pin, "pin", "", "the PIN to use for the HSM")
	flag.BoolVar(&debug, "debug", false, "show the version and exit")
}

func passphraseRetriever(keyName, alias string, createNew bool, attempts int) (passphrase string, giveup bool, err error) {
	envVar := _EnvPrefix + "_" + strings.ToUpper(alias)
	passphrase = os.Getenv(envVar)

	if passphrase == "" {
		return "", false, errors.New("expected env variable to not be empty: " + envVar)
	}

	return passphrase, false, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if _DebugAddr != "" {
		go debugServer(_DebugAddr)
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

	cryptoServices := make(signer.CryptoServiceIndex)

	if pkcs11Lib != "" {
		if pin == "" {
			log.Fatalf("Using PIN is mandatory with pkcs11")
		}

		ctx, session := SetupHSMEnv(pkcs11Lib)

		defer cleanup(ctx, session)

		cryptoServices[data.RSAKey] = api.NewRSAHardwareCryptoService(ctx, session)
	}

	dbSQL, err := sql.Open(_DBType, dbURL)
	if err != nil {
		log.Fatalf("failed to open the database: %s, %v", dbURL, err)
	}

	keyStore, err := trustmanager.NewKeyDBStore(passphraseRetriever, _DefaultAliasEnv, _DBType, dbSQL)
	if err != nil {
		log.Fatalf("failed to create a new keydbstore: %v", err)
	}
	cryptoService := cryptoservice.NewCryptoService("", keyStore)

	cryptoServices[data.ED25519Key] = cryptoService
	cryptoServices[data.ECDSAKey] = cryptoService

	//RPC server setup
	kms := &api.KeyManagementServer{CryptoServices: cryptoServices}
	ss := &api.SignerServer{CryptoServices: cryptoServices}

	grpcServer := grpc.NewServer()
	pb.RegisterKeyManagementServer(grpcServer, kms)
	pb.RegisterSignerServer(grpcServer, ss)

	lis, err := net.Listen("tcp", _RpcAddr)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to generate credentials %v", err)
	}
	go grpcServer.Serve(creds.NewListener(lis))

	//HTTP server setup
	server := http.Server{
		Addr:      _Addr,
		Handler:   api.Handlers(cryptoServices),
		TLSConfig: tlsConfig,
	}

	if debug {
		log.Println("[Notary-signer RPC Server] : Listening on", _RpcAddr)
		log.Println("[Notary-signer Server] : Listening on", _Addr)
	}

	err = server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("[Notary-signer Server] : Failed to start %s", err)
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
	log.Println("[Notary-signer Debug Server] server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[Notary-signer Debug Server] error listening on debug interface: %v", err)
	}
}

// SetupHSMEnv is a method that depends on the existences
func SetupHSMEnv(libraryPath string) (*pkcs11.Ctx, pkcs11.SessionHandle) {
	p := pkcs11.New(libraryPath)

	if p == nil {
		log.Fatalf("Failed to init library")
	}

	if err := p.Initialize(); err != nil {
		log.Fatalf("Initialize error %s\n", err.Error())
	}

	slots, err := p.GetSlotList(true)
	if err != nil {
		log.Fatalf("Failed to list HSM slots %s", err)
	}
	// Check to see if we got any slots from the HSM.
	if len(slots) < 1 {
		log.Fatalln("No HSM Slots found")
	}

	// CKF_SERIAL_SESSION: TRUE if cryptographic functions are performed in serial with the application; FALSE if the functions may be performed in parallel with the application.
	// CKF_RW_SESSION: TRUE if the session is read/write; FALSE if the session is read-only
	session, err := p.OpenSession(slots[0], pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		log.Fatalf("Failed to Start Session with HSM %s", err)
	}

	// (diogo): Configure PIN from config file
	if err = p.Login(session, pkcs11.CKU_USER, pin); err != nil {
		log.Fatalf("User PIN %s\n", err.Error())
	}

	return p, session
}

func cleanup(ctx *pkcs11.Ctx, session pkcs11.SessionHandle) {
	ctx.Destroy()
	ctx.Finalize()
	ctx.CloseSession(session)
	ctx.Logout(session)
}
