package main

import (
	"bytes"
	"os"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/registry/middleware/migration"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/distribution/configuration"

	"github.com/docker/distribution/context"

	log "github.com/Sirupsen/logrus"
)

func main() {
	// 1. connect to KV store to get regsitry config
	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to initialize key value storage")
	}
	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))

	config, err := getRegistryConfig(settingsStore)
	if err != nil {
		log.WithField("error", err).Fatal("unable to retrieve registry config")
	}

	ctx := context.Background()

	reg, err := dtrutil.NewRegistry(ctx, config)
	if err != nil {
		log.WithField("error", err).Fatal("unable to construct registry")
	}

	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
	session, err := dtrutil.GetRethinkSession(replicaID)
	if err != nil {
		log.WithField("error", err).Fatal("unable to connect to rethink")
	}
	store := schema.NewMetadataManager(session)

	// 2. Is this being resumed?
	migrated, err := kvStore.Get(deploy.HasMigratedToTagstore)
	if err != nil {
		log.WithField("error", err).Fatal("unable to retrieve migration info")
	}
	if bytes.Equal(migrated, []byte("1")) {
		// already migrated; we can quit
		return
	}

	m := migration.NewMigration(reg, store)

	// The migration is resumable; the previous repository that was migrated is stored
	// in migrationRepo. This is where we should continue from.
	rawRepo, err := kvStore.Get(deploy.MigrationRepo)
	if err != nil {
		log.WithField("error", err).Fatal("unable to retrieve previous migration repository")
	}
	repoStr := string(rawRepo)
	if repoStr != "" {
		m.Resume(repoStr)
	}

	if repo, err := m.Migrate(ctx); err != nil {
		log.WithField("error", err).Error("error during migration")
		if err = kvStore.Put(deploy.MigrationRepo, []byte(repo)); err != nil {
			log.WithField("error", err).Error("error saving migration progress to kv store")
		}
		return
	}

	kvStore.Put(deploy.HasMigratedToTagstore, []byte("1"))
	kvStore.Delete(deploy.MigrationRepo)
}

func getRegistryConfig(settingsStore hubconfig.SettingsStore) (*configuration.Configuration, error) {
	config, err := settingsStore.RegistryConfig()
	if err != nil {
		return config, err
	}

	if config.Storage.Type() == "filesystem" {
		params := config.Storage["filesystem"]
		params["rootdirectory"] = "/storage"
		config.Storage["filesystem"] = params
	}

	return config, nil
}
