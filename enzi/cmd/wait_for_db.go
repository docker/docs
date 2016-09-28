package cmd

import (
	"crypto/tls"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/schema"
)

// WaitForDB is the command to wait for all table replicas to be ready.
var WaitForDB = cli.Command{
	Name:   "wait-for-db",
	Usage:  "Wait for all DB table replicas to be ready",
	Action: runWaitForDB,
}

var waitFor string

func init() {
	WaitForDB.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "wait-for",
			Value:       schema.AllReplicasReady,
			Usage:       "table status to wait for",
			Destination: &waitFor,
		},
	}
}

func runWaitForDB(*cli.Context) error {
	tlsConfig := GetTLSConfig(tls.NoClientCert)

	log.Debugf("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	if err := schema.WaitForReadyTables(dbSession, schema.TableWaitTarget(waitFor)); err != nil {
		log.Fatalf("unable to wait for ready tables: %s", err)
	}

	log.Infof("Complete!")

	return nil
}
