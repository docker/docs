package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/distribution/context"
	workerapi "github.com/docker/orca/enzi/api/server/worker"
	"github.com/docker/orca/enzi/schema"
	"github.com/docker/orca/enzi/worker"
	"github.com/emicklei/go-restful"
)

// Worker is the command for running a worker server.
var Worker = cli.Command{
	Name:   "worker",
	Usage:  "Run a Worker Server",
	Action: runWorker,
}

var workerAddr string

func init() {
	Worker.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "addr",
			Value:       "worker.enzi",
			Usage:       "address (host[:port]) that this worker will advertise to API servers",
			Destination: &workerAddr,
		},
	}
}

const defaultWorkerPort = "12386"

func runWorker(*cli.Context) error {
	workerAddr, err := addrWithDefaultPort(workerAddr, defaultWorkerPort)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := GetTLSConfig(tls.RequireAndVerifyClientCert)

	log.Println("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	schemaMgr := schema.NewRethinkDBManager(dbSession)

	ctx := context.Background()

	w, err := worker.New(ctx, schemaMgr, workerAddr, "/work")
	if err != nil {
		context.GetLogger(ctx).Fatal(err)
	}

	// Run the worker in another goroutine.
	go w.Run()

	// Use this goroutine to run the webserver.
	serviceContainer := restful.NewContainer()
	serviceContainer.Add(workerapi.NewService(ctx, w, "/v0/").WebService)

	server := &http.Server{
		Addr:      ":4443",
		Handler:   serviceContainer,
		TLSConfig: tlsConfig,
		// There appears to be some bug with HTTP/2 so we explicitly
		// disable it for now.
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}

	log.Fatal(server.ListenAndServeTLS("", ""))

	return nil
}

// addrWithDefaultPort returns the given addr joined with the given defaultPort
// if it does not already specify a port. The given addr is parsed into its
// host and port components. If the addr fails to parse for reason other than
// having a missing port then that error is returned.
func addrWithDefaultPort(addr, defaultPort string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if addrErr, ok := err.(*net.AddrError); !(ok && addrErr.Err == "missing port in address") {
			return "", fmt.Errorf("unable to parse worker address: %s", err)
		}

		host, port = addr, defaultPort
	}

	return net.JoinHostPort(host, port), nil
}
