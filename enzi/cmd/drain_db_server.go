package cmd

import (
	"crypto/tls"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/schema"
)

// DrainDBServer is the command for safely removing a RethinkDB server from
// the cluster.
var DrainDBServer = cli.Command{
	Name:   "drain-db-server",
	Usage:  "Prepare a db server for removal from the cluster",
	Action: runDrainDBServer,
}

var (
	hostname string
	port     int
)

func init() {
	DrainDBServer.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "hostname, H",
			Value:       "localhost",
			Usage:       "canonical hostname for the server",
			Destination: &hostname,
		},
		cli.IntFlag{
			Name:        "cluster-port, p",
			Value:       29105,
			Usage:       "canonical port on which the server listens for intra-cluster connections",
			Destination: &port,
		},
	}
}

func runDrainDBServer(*cli.Context) error {
	tlsConfig := GetTLSConfig(tls.NoClientCert)

	log.Println("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	log.Println("draining db server...")
	if err := schema.DrainServer(dbSession, hostname, uint(port)); err != nil {
		log.Fatalf("unable to drain database server: %s", err)
	}

	log.Println("complete!")

	return nil
}
