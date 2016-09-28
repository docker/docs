package cmd

import (
	"crypto/tls"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/schema"
)

// CleanupDBJob is the worker subcommand for cleaning up expired or orphaned
// data from the DB.
var CleanupDBJob = cli.Command{
	Name:   "cleanup-db",
	Usage:  "Perform cleanup of expired and orphaned data",
	Action: runCleanupDB,
}

func runCleanupDB(*cli.Context) error {
	tlsConfig := GetTLSConfig(tls.NoClientCert)

	log.Info("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	schemaMgr := schema.NewRethinkDBManager(dbSession)

	log.Info("Deleting expired identity token signing keys...")
	if err := schemaMgr.DeleteExpiredSigningKeys(); err != nil {
		log.Errorf("unable to delete expired identity token signing keys: %s", err)
	}

	log.Info("Deleting expired cached service token signing keys...")
	if err := schemaMgr.DeleteExpiredServiceKeys(); err != nil {
		log.Errorf("unable to delete expired cached service signing keys: %s", err)
	}

	log.Info("Deleting expired service authorization codes...")
	if err := schemaMgr.DeleteExpiredServiceAuthCodes(); err != nil {
		log.Errorf("unable to delete expired service authorization codes: %s", err)
	}

	log.Info("Deleting expired user login sessions...")
	if err := schemaMgr.DeleteExpiredSessions(); err != nil {
		log.Errorf("unable to delete expired login sessions: %s", err)
	}

	// TODO: Cleanup orphaned values from the following tables:
	// - orgMembership:
	//   organization or member may have been deleted
	// - teams:
	//   organization may have been deleted
	// - teamMembership:
	//   team or member may have been deleted
	// - services:
	//   owner may have been deleted
	// - serviceSessions:
	//   root sessions may have been deleted

	return nil
}
