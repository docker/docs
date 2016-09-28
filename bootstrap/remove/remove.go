package remove

import (
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	jschema "github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/engine-api/types/filters"
)

const RemoveReplicaIDEnvVar = "DTR_REMOVE_REPLICA_ID"

func destroyReplica(bs bootstrap.Bootstrap, removeReplicaID string, removeCA bool) error {
	// TODO: get rid of this ridiculous state thing
	bs.SetReplicaID(removeReplicaID)
	// let's remove everything about this replica first - containers and volumes (we don't remove networks for now until we switch individual bridge networks)
	log.Info("Stopping containers")
	if err := bs.StopDTRContainers(map[string]bool{}, removeReplicaID); err != nil {
		return err
	}
	log.Info("Removing containers")
	if err := bs.RemoveDTRContainers(map[string]bool{}, removeReplicaID); err != nil {
		return err
	}
	log.Info("Removing volumes")
	volumes, err := bs.GetDockerClient().VolumeList(context.Background(), filters.Args{})
	if err != nil {
		return err
	}
	for _, volume := range volumes.Volumes {
		// we do some hacks to chop off the node name for the purpose of comparing names
		volumeName := volume.Name
		parts := strings.SplitN(volumeName, "/", 2)
		if len(parts) > 1 {
			volumeName = parts[1]
		}

		for _, dtrvol := range containers.Volumes {
			if dtrvol == containers.CAVolume && !removeCA {
				continue
			}
			expectedVolumeName := dtrvol.ReplicaName(removeReplicaID)

			if expectedVolumeName == volumeName {
				log.Debugf("Removing volume '%s'", volume.Name)
				if err := bs.VolumeRemove(volume.Name); err != nil {
					log.Errorf("Couldn't remove volume '%s'", volume.Name)
					return err
				}
			}
		}
	}
	return nil
}

// phase 2 is invoked only when an existing replica id is given and we want to keep the cluster healthy after the removal
// and never when uninstalling aka force removing
func phase2(c *cli.Context) error {
	log.Debug("phase 2 starting...")
	var err error

	// this will not do any prompts if it's run correctly
	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}
	// this is not a cli flag, in phase 2 the replica id of the healthy node is flags.ReplicaID
	removeReplicaID := os.Getenv(RemoveReplicaIDEnvVar)

	log.Info("Checking your current configs")
	// let's also remove the etcd key for read-only registry just in case it's there
	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		return fmt.Errorf("Error creating kv store: %s", err)
	}
	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return fmt.Errorf("failed to get ha config %s", err)
	}

	// we need to remove the etcd container from the cluster before it's killed to maintain consensus
	log.Info("De-registering etcd")
	// now that the container is gone, let's re-register it from etcd
	// the replica we connect to here is on the healthy node
	clientURL := fmt.Sprintf("https://%s:%d", containers.Etcd.OverlayName(flags.ReplicaID), containers.EtcdClientPort1)
	// Oh, god, I'm sorry. This needs to be cleaned up later
	log.Debugf("connecting to etcd on %s", clientURL)
	err = RemoveEtcdNode(removeReplicaID, clientURL)
	if err != nil {
		return fmt.Errorf("Error removing node from etcd: %s", err)
	}
	// XXX: do we need to wait for etcd to settle here because we removed a node?

	// remove data about the replica from haConfig: number of replicas, port configs
	log.Info("Removing the replica from the configs")
	haConfig, err = settingsStore.HAConfig()
	if err != nil {
		return fmt.Errorf("failed to get ha config %s", err)
	}
	delete(haConfig.ReplicaConfig, removeReplicaID)

	err = settingsStore.SetHAConfig(haConfig)
	if err != nil {
		return fmt.Errorf("failed to set ha config %s", err)
	}

	log.Info("Removing left over config state from the replica")
	// let's also remove the etcd key for read-only registry just in case it's there
	roModePath := path.Join(deploy.RegistryROStatePath, removeReplicaID)
	err = kvStore.Delete(roModePath)
	if err != nil {
		log.Debugf("could not remove read-only mode key, this is probably normal: %s", err)
	}

	log.Info("Scaling down the number of database table replicas")
	// scale down rethinkdb number of copies for dtr tables
	schemaManager := bootstrap.CreateManager(flags.ReplicaID)

	err = schemaManager.DecommissionServer(containers.Rethinkdb.RethinkServerTagName(removeReplicaID))
	if err != nil {
		return fmt.Errorf("Failed to decomission server: %s", err)
	}

	// scale down rethinkdb number of copies for jobrunner tables
	dbSession, err := dtrutil.GetRethinkSession(flags.ReplicaID)
	if err != nil {
		return fmt.Errorf("Failed to create db session: %s", err)
	}
	schemaMgr := jschema.NewJobrunnerManager(deploy.JobrunnerDBName, dbSession)
	err = schemaMgr.ScaleDB(uint(len(haConfig.ReplicaConfig)), false)
	if err != nil {
		return fmt.Errorf("Failed to set up jobrunner tables: %s", err)
	}

	if err = schemaManager.SetReplication(len(haConfig.ReplicaConfig)); err != nil {
		return fmt.Errorf("Error changing rethink replication factor: %s", err)
	}

	// finally, we get rid of the replica
	err = destroyReplica(bs, removeReplicaID, true)
	if err != nil {
		return err
	}

	log.Infof("Done. Your cluster size is now %d", len(haConfig.ReplicaConfig))
	return nil
}

func phase1(c *cli.Context) error {
	log.Infof("Beginning Docker Trusted Registry replica remove")

	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}

	replicas, err := bs.ListReplicas(true)
	if err != nil {
		return err
	}
	replicaIDs := replicas.ReplicaIDs()

	// This replica ID flag is special.
	// Remove needs to be idempotent, so we trust the user if they provide an invalid replica ID and try to remove it anyway
	if flags.ReplicaID == "" {
		log.Infof("This cluster contains the replicas: %s", strings.Join(replicaIDs, " "))
		if len(replicaIDs) > 0 {
			flags.ReplicaID = bootstrap.PromptString(fmt.Sprintf("%s [%s]: ", "Choose a replica to remove", replicaIDs[0]), replicaIDs[0])
		} else {
			flags.ReplicaID = bootstrap.PromptString(fmt.Sprintf("%s: ", "Choose a replica to remove"), "")
		}
	}

	// we allow existing replica id to be anything if it's the same as flags.ReplicaID, so we can do idempotent uninstall
	replica, err := bs.ExistingReplicaFlagPicker("Choose any healthy replica", false)
	if err != nil {
		return err
	}

	if replica != nil {
		if err := replica.CheckEqualVersion(deploy.Version); err != nil {
			return err
		}
	}

	uninstall := flags.ReplicaID == flags.ExistingReplicaID
	forced := flags.ForceRemove

	if !uninstall && forced {
		return fmt.Errorf("The --%s flag does not apply to regular removals. Either don't supply the --%s flag or make --%s == --%s to remove the replica without notifying the rest of the cluster.",
			flags.ForceRemoveFlag.Name, flags.ForceRemoveFlag.Name, flags.ExistingReplicaIDFlag.Name, flags.ReplicaIDFlag.Name)
	}

	// uninstall is done from phase1 because we don't need to communicate with etcd or rethinkdb and we shouldn't mount volumes we are about to remove
	if uninstall {
		if !forced {
			return fmt.Errorf("Attempting to remove replica without notifying the rest of the cluster. You need to use --%s to do this to confirm that you understand that this may break your cluster.", flags.ForceRemoveFlag.Name)
		} else {
			log.Debug("Forced removal of %s", flags.ReplicaID)
			err := destroyReplica(bs, flags.ReplicaID, true)
			if err != nil {
				return err
			}
			log.Info("Replica removed.")
			return nil
		}
	}

	// for this command replicaID is the one you are removing and existingReplicaID is a healthy node
	// we drop onto the healthy node when we do the removal
	container := dropperOpts.MakeContainerConfig([]string{"remove"})
	bootstrap.SetEnvFromFlags(c, container, flags.RemoveFlags...)
	container.Environment[RemoveReplicaIDEnvVar] = flags.ReplicaID
	container.Networks = []containers.NetworkConfig{{Name: deploy.BridgeNetworkName}, {Name: deploy.OverlayNetworkName}}

	bs.SetReplicaID(flags.ExistingReplicaID)
	// we have to make sure the container is dropped in the right place to connect to the existing replica's etcd
	if !flags.NoUCP {
		if container.Node, err = bs.ContainerNode(containers.Etcd.ReplicaName(flags.ExistingReplicaID)); err != nil {
			log.Errorf("Couldn't find the node Etcd is running on: %s", err)
			return err
		}
		dropperOpts.Node = container.Node
	}

	err = bootstrap.Phase2Execute(container, bs)
	if err != nil {
		return err
	}

	return nil
}

func remove(c *cli.Context) error {
	if !bootstrap.IsPhase2() {
		err := phase1(c)
		if err != nil {
			log.Error(bootstrap.IdempotentMsg("Remove"))
		}
		return err
	} else {
		return phase2(c)
	}
	return nil
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	if err := remove(c); err != nil {
		log.Fatal(err)
	}
}

// TODO: deduplicate this with AddEtcdNode
func RemoveEtcdNode(replicaID, clientURL string) error {
	// XXX - the way the etcd client works here is different than the way it works in
	//       hubconfig/etcd/keyvaluestore.go.  We should refactor the other code to work
	//       similarly to this, plus also have a generic interface for getting a
	//       reusable etcd client.
	memberName := containers.Etcd.ReplicaName(replicaID)
	// TODO: make this a constant or function somewhere. It's too magic.

	log.Debugf("removing %s", memberName)
	etcdClient, err := bootstrap.GetEtcdConn(clientURL)
	if err != nil {
		log.Errorf("Couldn't get etcd connection: %s", err)
		return err
	}

	m := client.NewMembersAPI(etcdClient)
	log.Debugf("client = %q", m)

	ctx := context.Background()
	members, err := m.List(ctx)
	if err != nil {
		return err
	}

	memberID := ""
	for _, member := range members {
		if memberName == member.Name {
			memberID = member.ID
		}
	}
	// if we failed to find it, it doesn't exist already
	if memberID == "" {
		log.Debug("etcd member already removed")
		return nil
	}

	err = m.Remove(ctx, memberID)
	if err != nil {
		// TODO: oh, god, deduplicate this error handling with AddEtcdNode
		if err == context.Canceled {
			log.Error("canceled")
			return err
		} else if err == context.DeadlineExceeded {
			log.Error("deadline exceeded")
			return err
		} else if cerr, ok := err.(*client.ClusterError); ok {
			log.Errorf("cluster error:  %q", cerr.Errors)
			return err
		} else {
			log.Errorf("Couldn't remove from cluster: %s", err)
			return err
		}
		return err
	}

	return nil
}
