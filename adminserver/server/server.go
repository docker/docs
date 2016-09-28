package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/garant/authz"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/licensing"
	"github.com/docker/dhe-deploy/manager/versions"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	auditLogger := log.New()
	auditLogger.Formatter = new(log.JSONFormatter)
	var err error

	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to initialize key value storage")
	}
	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))
	versionChecker := versions.NewChecker(settingsStore)

	licenseChecker := licensing.NewChecker(settingsStore)
	err = licenseChecker.Initialize()
	if err != nil {
		log.Info("Initializing with an invalid license")
	} else {
		licenseChecker.BeginLicenseSyncing()
	}

	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
	session, err := dtrutil.GetRethinkSession(replicaID)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to connect to rethink")
	}

	authorizer := authz.NewAuthorizer(session, settingsStore)
	etcdClient := new(http.Client)
	etcdClient.Transport, err = bootstrap.GetEtcdTransport()
	if err != nil {
		log.WithField("error", err).Fatal("Failed to initialize server")
	}

	notaryClient := new(http.Client)
	notaryClient.Transport, err = bootstrap.GetNotaryTransport()
	if err != nil {
		log.WithField("error", err).Fatal("Failed to connect to notary server")
	}

	adminServer, err := adminserver.New(settingsStore, kvStore, etcdClient, notaryClient, licenseChecker, versionChecker, auditLogger, authorizer, session)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to initialize server")
	}

	adminServer.SubscribeToEvents()
	adminServer.ListenAndServe(fmt.Sprintf(":%d", deploy.AdminPort))
}
