package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/shared/containers"
)

// This binary outputs the arguments that rethinkdb should be started with on stdout and uses stderr for logging and error reporting

func printRethinkArgs() error {
	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		return err
	}
	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return err
	}

	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
	args := "--bind all --no-update-check"
	args += fmt.Sprintf(" --directory %s", filepath.Join(containers.RethinkVolume.Location, "rethink"))
	args += fmt.Sprintf(" --driver-tls-key %s", containers.RethinkCertStore.KeyPath())
	args += fmt.Sprintf(" --driver-tls-cert %s", containers.RethinkCertStore.CertPath())
	args += fmt.Sprintf(" --driver-tls-ca %s", containers.RethinkCACertStore.CertPath())
	args += fmt.Sprintf(" --cluster-tls-key %s", containers.RethinkCertStore.KeyPath())
	args += fmt.Sprintf(" --cluster-tls-cert %s", containers.RethinkCertStore.CertPath())
	args += fmt.Sprintf(" --cluster-tls-ca %s", containers.RethinkCACertStore.CertPath())
	args += fmt.Sprintf(" --http-port %s", fmt.Sprintf("%d", deploy.RethinkAdminPort))
	args += fmt.Sprintf(" --server-tag %s", containers.Rethinkdb.RethinkServerTagName(replicaID))
	args += fmt.Sprintf(" --server-name %s", containers.Rethinkdb.RethinkServerTagName(replicaID))
	args += fmt.Sprintf(" --canonical-address %s", containers.Rethinkdb.OverlayName(replicaID))

	for replicaID := range haConfig.ReplicaConfig {
		args += fmt.Sprintf(" --join %s", containers.Rethinkdb.OverlayName(replicaID))
	}
	fmt.Println(args)
	return nil
}

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.Info("arggen starting")
	err := printRethinkArgs()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
