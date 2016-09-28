package gc

import (
	"bytes"
	"time"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/distribution/configuration"

	log "github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"
)

const RestartPollChecks = 300

func NewSetup() (s Setup, err error) {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	// TODO: make GC log level configurable
	log.SetLevel(log.InfoLevel)

	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		return s, stacktrace.Propagate(err, "failed to initialize key value storage")
	}

	// Ensure that the migration to tagstore has finished before returning a new setup instance.
	migrated, err := kvStore.Get(deploy.HasMigratedToTagstore)
	if err != nil {
		return s, stacktrace.Propagate(err, "unable to retrieve migration info")
	}
	if !bytes.Equal(migrated, []byte("1")) {
		// already migrated; we can quit
		return s, stacktrace.NewError("migration has not yet finished; failing to construct setup")
	}

	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))
	config, err := settingsStore.RegistryConfig()
	if err != nil {
		return s, stacktrace.Propagate(err, "failed to get config")
	}

	if config.Storage.Type() == "filesystem" {
		params := config.Storage["filesystem"]
		params["rootdirectory"] = "/storage"
		config.Storage["filesystem"] = params
	}

	return Setup{
		kv:     kvStore,
		ss:     settingsStore,
		config: config,
	}, nil

}

type Setup struct {
	kv     hubconfig.KeyValueStore
	ss     hubconfig.SettingsStore
	config *configuration.Configuration
}

func (s Setup) SetReadOnly() error {
	return s.setRegistryMode(true)
}

func (s Setup) SetReadWrite() error {
	return s.setRegistryMode(false)
}

func (s Setup) GetConfig() configuration.Configuration {
	return *s.config
}

func (s Setup) GetGCMode() (string, error) {
	hubConfig, err := s.ss.UserHubConfig()
	if err != nil {
		return "", stacktrace.Propagate(err, "unable to fetch gc mode")
	}
	if hubConfig.GCMode == "" {
		return ModeByTag, nil
	}
	return hubConfig.GCMode, nil
}

func (s Setup) setRegistryMode(ro bool) error {
	var (
		err         error
		numReplicas int
	)

	// Set the config to read only mode
	util.SetReadonlyMode(&s.config.Storage, ro)
	if err = s.ss.SetRegistryConfig(s.config); err != nil {
		return err
	}

	// we only need to wait for containers to switch modes if we're heading into
	// ro mode
	if ro {
		numReplicas, err = getNumReplicas(s.ss)
	}

	log.WithFields(log.Fields{
		"numReplicas": numReplicas,
		"isReadOnly":  ro,
	}).Infof("waiting for registries to switch modes")

	// wait for 10 seconds before polling all registries for RO status
	time.Sleep(deploy.ContainerRestartTimeout * time.Second)

	// poll to ensure that all containers are in RO mode
	return dtrutil.Poll(time.Second, RestartPollChecks, func() error {
		keys, err := s.kv.List(deploy.RegistryROStatePath)
		if err != nil {
			return stacktrace.Propagate(err, "failed to list read-only registries list")
		}
		if len(keys) != numReplicas {
			return stacktrace.NewError("registry containers in read-only mode: %d, needed: %d", len(keys), numReplicas)
		}
		return nil
	})
}

func getNumReplicas(settingsStore hubconfig.SettingsStore) (int, error) {
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return 0, err
	}
	return len(haConfig.ReplicaConfig), nil
}
