package cmd

import (
	"crypto/tls"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/schema"
)

// SyncDB is the command for syncing RethinkDB table config.
var SyncDB = cli.Command{
	Name:   "sync-db",
	Usage:  "Setup Database",
	Action: runSyncDB,
}

var (
	numReplicas     int
	emergencyRepair bool
)

func init() {
	SyncDB.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "num-replicas, n",
			Usage:       "number of data replicas to configure (default is to use the number of detected DB cluster nodes)",
			Destination: &numReplicas,
		},
		cli.BoolFlag{
			Name:        "emergency-repair",
			Usage:       "perform emergency repair of all table replicas",
			Destination: &emergencyRepair,
		},
	}
}

func runSyncDB(*cli.Context) error {
	tlsConfig := GetTLSConfig(tls.NoClientCert)

	log.Println("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	log.Println("setting up database schema...")
	if err := schema.SetupDB(dbSession, uint(numReplicas), emergencyRepair); err != nil {
		log.Fatalf("unable to setup database schema: %s", err)
	}

	log.Println("complete!")

	return nil
}
