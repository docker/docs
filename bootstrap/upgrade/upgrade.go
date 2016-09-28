package upgrade

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/install"
	"github.com/docker/dhe-deploy/bootstrap/reconfigure"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/shared/containers"
)

func upgrade(c *cli.Context) error {
	if !bootstrap.IsPhase2() {
		err := phase1(c)
		if err != nil {
			log.Error(bootstrap.IdempotentMsg("Upgrade"))
		}
		return err
	} else {
		return phase2(c)
	}

	return nil
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	if err := upgrade(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(0)
	}
}

func phase1(c *cli.Context) error {
	log.Infof("Beginning Docker Trusted Registry Upgrade")
	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}

	log.Infof("This command will upgrade all replicas in your cluster, however it needs a single replica to begin.")
	replica, err := bs.ExistingReplicaFlagPicker("Choose a replica to perform the upgrade", true)
	if err != nil {
		return err
	}

	if err := replica.CheckUpgradeableVersion(deploy.Version); err != nil {
		return err
	}

	container := dropperOpts.MakeContainerConfig([]string{"upgrade"})
	bootstrap.SetEnvFromFlags(c, container, flags.UpgradeFlags...)

	bs.SetReplicaID(flags.ExistingReplicaID)

	// Look to see if http/https ports are already in use, and exclude those nodes
	// XXX - NodePortConstraints are probably just setting this to run on the existing replica
	if !flags.NoUCP {
		if nodeName, err := bs.ContainerNode(containers.Etcd.ReplicaName(flags.ExistingReplicaID)); err != nil {
			log.Errorf("Couldn't find the node Etcd is running on: %s", err)
			return err
		} else {
			container.Node = nodeName
			bs.SetNodeName(nodeName)
		}
	}

	err = bootstrap.Phase2Execute(container, bs)
	if err != nil {
		bootstrap.CheckConstraintsError(err)
		return err
	}
	return nil
}

func phase2(c *cli.Context) error {
	log.Info("Starting phase2 upgrade")

	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}
	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}
	bs.SetReplicaID(flags.ExistingReplicaID)

	err = bootstrap.SetupNode(c, bs)
	if err != nil {
		return err
	}

	netName := bootstrap.GetBridgeNetworkName(c, bs)
	err = bootstrap.Phase2NetworkConnect(bs, netName)
	if err != nil {
		return err
	}

	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		return fmt.Errorf("Couldn't set up kvStore: %s", err)
	}
	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))

	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return fmt.Errorf("Couldn't get ha config: %s", err)
	}
	numReplicas := len(haConfig.ReplicaConfig)

	replicas, err := bs.ListReplicas(true)
	if err != nil {
		return err
	}

	replica := replicas.GetReplica(flags.ExistingReplicaID)
	if replica == nil {
		return fmt.Errorf("Couldn't find replica: %s", flags.ExistingReplicaID)
	}

	migrate, err := replica.RequiresMigration(deploy.Version)
	if err != nil {
		return err
	}

	if migrate {
		log.Infof("Migrating cluster to %s...", deploy.Version)
		if err := migrateReplica(bs, kvStore, settingsStore, c, numReplicas); err != nil {
			return err
		}
	}

	// Trigger regeneration of derived configs like notary configs and
	// derived constants like paths and directories
	if err := refreshSettingsStore(settingsStore); err != nil {
		return err
	}

	err = reconfigure.Reconfigure(bs, kvStore, settingsStore, c, nil, false, true, true, false)
	if err != nil {
		return err
	}

	if numReplicas < 3 {
		log.Info("Upgrade is complete, however it may take a few moments for Docker Trusted Registry to become available.")
	} else {
		log.Info("Upgrade is complete.")
	}

	return nil
}

// Migrations for 2.0.x to 2.1.x specifically
// Overwrite this function when we start writing 2.2.x
func migrateReplica(bs bootstrap.Bootstrap, kvStore hubconfig.KeyValueStore, settingsStore hubconfig.SettingsStore, c *cli.Context, numReplicas int) error {
	if err := install.SetupVolumes(c, bs); err != nil {
		return err
	}

	if err := install.CreateNotaryClientCerts(bs); err != nil {
		return err
	}

	var premigrateContainers []containers.DTRContainer
	premigrateContainers = append(premigrateContainers, containers.Etcd)
	premigrateContainers = append(premigrateContainers, containers.Rethinkdb)
	err := reconfigure.Reconfigure(bs, kvStore, settingsStore, c, premigrateContainers, false, true, true, false)
	if err != nil {
		return err
	}

	// Migrate databases
	_, err = bootstrap.MigrateDatabase(bs.GetReplicaID(), uint(numReplicas))
	if err != nil {
		log.Errorf("Couldn't migrate database: %s", err)
		return err
	}

	return nil
}

// Upon upgrading versions, paths for configs may change and various
// other changes made to the settingsStore will need a refresh of
// all the configs except for signing keys and certs
func refreshSettingsStore(settingsStore hubconfig.SettingsStore) error {
	userHubConfig, err := settingsStore.UserHubConfig()
	if err != nil {
		return err
	}
	if err := settingsStore.SetUserHubConfig(userHubConfig); err != nil {
		return err
	}

	authConfig, err := settingsStore.AuthConfig()
	if err != nil {
		return err
	}
	if err := settingsStore.SetAuthConfig(authConfig); err != nil {
		return err
	}

	licenseConfig, err := settingsStore.LicenseConfig()
	if err != nil {
		return err
	}
	if err := settingsStore.SetLicenseConfig(licenseConfig); err != nil {
		return err
	}

	registryConfig, err := settingsStore.RegistryConfig()
	if err != nil {
		return err
	}
	if err := settingsStore.SetRegistryConfig(registryConfig); err != nil {
		return err
	}

	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return err
	}
	if err := settingsStore.SetHAConfig(haConfig); err != nil {
		return err
	}

	return nil
}
