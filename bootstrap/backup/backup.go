package backup

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
)

const BackupVersionPrefix = "dtr-backup-v"
const BackupVersion = BackupVersionPrefix + deploy.Version

func phase1(c *cli.Context) error {
	log.Debug("backup phase 1")

	err := bootstrap.RequireNoTTY()
	if err != nil {
		return err
	}

	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}

	replica, err := bs.ExistingReplicaFlagPicker("Choose a replica to back up from", true)
	if err != nil {
		return err
	}

	if err := replica.CheckEqualVersion(deploy.Version); err != nil {
		return err
	}

	bs.SetReplicaID(flags.ExistingReplicaID)

	cfg := *dropperOpts.MakeContainerConfig([]string{"backup"})
	cfg.Environment["affinity:container="] = containers.Etcd.ReplicaName(flags.ExistingReplicaID)
	// XXX: is this needed or does something else already set the replica id?
	bootstrap.SetEnvFromFlags(c, &cfg, flags.BackupFlags...)
	cfg.Tty = false
	cfg.Volumes = append(cfg.Volumes, containers.RethinkVolume)

	cResp, err := bs.ContainerCreateFromContainerConfig(cfg)
	if err != nil {
		return fmt.Errorf("Problem running container '%s' from image '%s': %s", cfg.Name, cfg.Image, err)
	}
	cId := cResp.ID

	aResp, err := bs.GetDockerClient().ContainerAttach(context.Background(), cId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  false,
		Stdout: true,
		Stderr: true,
	})

	if err != nil {
		return fmt.Errorf("Couldn't attach to dtr-phase2 container: %s", err)
	}

	if err := bs.ContainerStart(cId); err != nil {
		return fmt.Errorf("There was a problem running the dtr-phase2 container: %s", err)
	}

	go func() {
		defer aResp.Close()
		stdcopy.StdCopy(os.Stdout, os.Stderr, aResp.Reader)
	}()

	rc, err := bs.ContainerWait(cId)
	if err != nil {
		return fmt.Errorf("Error waiting for dtr-phase2 to finish: %s, code: %d", err, rc)
	}

	if err := bs.ContainerRemove(cId, types.ContainerRemoveOptions{}); err != nil {
		return fmt.Errorf("Error removing the dtr-phase2 container: %s", err)
	}
	log.Info("Backup complete.")

	return nil
}

func phase2(c *cli.Context) error {
	// make sure log output goes to stderr while we setup the environment
	// so the stdout stream can be used for the tar archive

	replicaID := flags.ExistingReplicaID
	if replicaID == "" {
		// Should never reach here...
		err := fmt.Errorf("You must specify DTR replica with --%s flag.", flags.ExistingReplicaIDFlag.Name)
		return err
	}
	configOnly := flags.ConfigOnly

	var err error
	var out io.WriteCloser
	out = os.Stdout

	bs, kvStore, err := SetupEnv(c, replicaID)
	if err != nil {
		err := fmt.Errorf("Can't set up environment to backup the configuration: %s", err)
		return err
	}

	tw := tar.NewWriter(out)
	defer tw.Close()

	versionHeader := &tar.Header{
		Name:     fmt.Sprintf("%s/", BackupVersion),
		Typeflag: tar.TypeDir,
		Mode:     0700,
	}
	if err := tw.WriteHeader(versionHeader); err != nil {
		log.Errorf("Couldn't write header to indicate backup version: %s", err)
		return err
	}

	if !configOnly {
		if err := backupRethink(tw, bs, replicaID); err != nil {
			return err
		}
	}

	if err := backupEtcd(c, bs, kvStore, replicaID, tw); err != nil {
		return err
	}

	log.Debug("The replica's containers are up and running again.")
	return nil
}

func backupRethink(tw *tar.Writer, bs bootstrap.Bootstrap, replicaID string) error {
	header := &tar.Header{
		Name:     fmt.Sprintf("%s/%s/", BackupVersion, deploy.RethinkdbDirectory),
		Typeflag: tar.TypeDir,
		Mode:     0700,
	}
	if err := tw.WriteHeader(header); err != nil {
		log.Errorf("Couldn't write header to indicate rethinkdb directory: %s", err)
		return err
	}

	schemaManager := bootstrap.CreateManager(bs.GetReplicaID())
	if err := schemaManager.Initialize(); err != nil {
		log.Errorf("Can't initialize schema manager: %s", err)
		return err
	}
	if err := schemaManager.DumpAllTables(tw, BackupVersion); err != nil {
		log.Errorf("Couldn't dump rethink data: %s", err)
		return err
	}
	return nil
}

func backupEtcd(c *cli.Context, bs bootstrap.Bootstrap, kvStore hubconfig.KeyValueStore, replicaID string, tw *tar.Writer) error {
	var errors bool
	// Only backup certain files as others hold keys and such which we re-generated
	// when we do a reinstall so backing those would mean
	// 1) we're backing up useless data (because it'll be stale)...
	// 2) DTR will break as it won't authenticate properly...
	configFiles := []string{
		deploy.UserHubConfigFilename,
		deploy.RegistryConfigFilename,
		deploy.AuthConfigFilename,
		deploy.LicenseConfigFilename,
		deploy.HAConfigFilename,
	}

	header := &tar.Header{
		Name:     fmt.Sprintf("%s/%s/", BackupVersion, deploy.EtcdDirectory),
		Typeflag: tar.TypeDir,
		Mode:     0700,
	}
	if err := tw.WriteHeader(header); err != nil {
		log.Errorf("Couldn't write header to indicate rethinkdb directory: %s", err)
		return err
	}

	for _, fileName := range configFiles {
		if configBytes, err := kvStore.Get(fileName); err == nil {
			log.Debugf("Backing up: %s", fileName)
			if fileName == deploy.HAConfigFilename {
				configBytes, err = updateHAConfig(configBytes, replicaID)
				if err != nil {
					return err
				}
			}
			fileName = fmt.Sprintf("%s/%s/%s/", BackupVersion, deploy.EtcdDirectory, fileName)
			if err := dtrutil.AddBytesToTar(tw, configBytes, fileName); err != nil {
				errors = true
				errMessage := strings.Join([]string{"Can't write ", fileName, "to tar: ", err.Error()}, " ")
				log.Error(errMessage)
			}
		} else {
			errors = true
			errMessage := strings.Join([]string{"Can't get data for file ", fileName, ": ", err.Error()}, " ")
			log.Error(errMessage)
		}
	}

	if errors {
		return fmt.Errorf("Backup may be incomplete")
	} else {
		return nil
	}
}

func updateHAConfig(configBytes []byte, replicaID string) (updatedConfigBytes []byte, err error) {
	var haConfig hubconfig.HAConfig
	if err := json.Unmarshal(configBytes, &haConfig); err != nil {
		log.Errorf("Can't back up specified replica HA config: %s", err)
		return nil, err
	}
	// back up all fields of HAConfig except ReplicaConfig, so set it null
	// we could just not look at the ReplicaConfig during restore and ignore nulling
	// it here but why waste space when this field of backup won't be used?
	haConfig.ReplicaConfig = map[string]hubconfig.ReplicaConfig{}

	if updateConfigBytes, err := json.Marshal(haConfig); err != nil {
		return nil, err
	} else {
		return updateConfigBytes, nil
	}
}

func SetupEnv(c *cli.Context, replicaID string) (bootstrap.Bootstrap, hubconfig.KeyValueStore, error) {
	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return nil, nil, err
	}
	bs.SetReplicaID(replicaID)

	err = bootstrap.SetupNode(c, bs)
	if err != nil {
		return nil, nil, err
	}

	netName := bootstrap.GetBridgeNetworkName(c, bs)
	err = bootstrap.Phase2NetworkConnect(bs, netName)
	if err != nil {
		return nil, nil, err
	}

	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		log.Debugf("Couldn't set up kvStore: %s", err)
		return nil, nil, err
	}

	return bs, kvStore, nil
}

func backup(c *cli.Context) error {
	if !bootstrap.IsPhase2() {
		return phase1(c)
	} else {
		return phase2(c)
	}
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	log.SetOutput(os.Stderr)
	if err := backup(c); err != nil {
		log.Fatal(err)
	}
}
