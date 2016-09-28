package cmd

import (
	"crypto/tls"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/distribution/context"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/authn/ldap/sync"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
)

// LDAPSyncJob is the worker subcommand for performing an LDAP sync.
var LDAPSyncJob = cli.Command{
	Name:   "ldap-sync",
	Usage:  "Sync Users and Teams with an LDAP Directory",
	Action: runLdapSync,
}

func runLdapSync(*cli.Context) error {
	tlsConfig := GetTLSConfig(tls.NoClientCert)

	log.Println("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	schemaMgr := schema.NewRethinkDBManager(dbSession)

	ctx := context.Background()

	authConfig, err := config.GetAuthConfig(schemaMgr)
	if err != nil {
		errCtx := context.WithValue(ctx, "error", err)
		context.GetLogger(errCtx, "error").Fatal("failed to get current auth configuration")
	}

	if authConfig.Backend != config.AuthBackendLDAP {
		context.GetLogger(ctx).Info("current auth configuration does not use LDAP integration")
		os.Exit(0)
	}

	ldapSettings, err := ldapconfig.GetLDAPConfig(schemaMgr)
	if err != nil {
		errCtx := context.WithValue(ctx, "error", err)
		context.GetLogger(errCtx, "error").Fatal("failed to get current LDAP configuration")
	}

	syncer := sync.NewLdapSyncer(ctx, schemaMgr, ldapSettings)

	if err := syncer.Run(); err != nil {
		context.GetLogger(ctx).Fatal(err)
	}

	return nil
}
