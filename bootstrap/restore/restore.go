package restore

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/backup"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/install"
	"github.com/docker/dhe-deploy/bootstrap/reconfigure"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/manager/versions"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
	version "github.com/docker/engine-api/types/versions"
)

func phase1(c *cli.Context) error {
	log.Debug("restore phase 1")

	var in io.Reader
	in = os.Stdin
	tr := tar.NewReader(in)

	hdr, err := tr.Next()
	if err != nil {
		return err
	}

	restoreVersion := findBackupVersion(hdr.Name)
	if err := checkRestoreVersion(restoreVersion); err != nil {
		return err
	}

	// first install DTR
	replicaID, err := install.Phase1(c)
	if err != nil {
		return err
	}

	// from now on, pretend like the replicaID flag was set whether or not it was
	flags.ReplicaID = replicaID

	log.Info("Install done")

	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}
	bs.SetReplicaID(flags.ReplicaID)

	cfg := *dropperOpts.MakeContainerConfig([]string{"restore"})
	cfg.AttachStdin = true
	cfg.StdinOnce = true
	cfg.OpenStdin = true
	cfg.Tty = false
	cfg.Environment["affinity:container="] = containers.Etcd.ReplicaName(flags.ReplicaID)
	bootstrap.SetEnvFromFlags(c, &cfg, flags.RestoreFlags...)

	cResp, err := bs.ContainerCreateFromContainerConfig(cfg)
	if err != nil {
		return fmt.Errorf("Problem running container '%s' from image '%s': %s", cfg.Name, cfg.Image, err)
	}
	cId := cResp.ID

	aResp, err := bs.GetDockerClient().ContainerAttach(context.Background(), cId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: false,
	})

	if err != nil {
		return fmt.Errorf("Couldn't attach to phase2 container: %s", err)
	}

	if err := bs.ContainerStart(cId); err != nil {
		return fmt.Errorf("There was a problem running the dtr-phase2 container: %s", err)
	}

	doneChan := make(chan int)
	var exitCode int
	go func() {
		rc, err := bs.ContainerWait(cId)
		if err != nil {
			log.Errorf("Failed to wait for phase 2 container: %s", err)
		}
		doneChan <- rc
	}()

	rd := bufio.NewReader(aResp.Reader)
	go func() {
		if _, err := io.Copy(aResp.Conn, os.Stdin); err != nil {
			log.Warnf("Stdin copy interruped: %s", err)
		}
	}()

	go func() {
		oldLevel := log.GetLevel()
		log.SetLevel(log.InfoLevel)
		defer log.SetLevel(oldLevel)
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, rd); err != nil {
			log.Errorf("Failed to stream logs: %s\n", err)
			return
		}
	}()

	for {
		exitCode = <-doneChan
		if exitCode < 0 {
			continue
		}
		if err := bs.ContainerRemove(cId, types.ContainerRemoveOptions{}); err != nil {
			log.Errorf("Problem removing the phase 2 container: %s", err)
			return err
		}
		log.Debug("Removed container")
		break
	}

	flags.ExistingReplicaID = replicaID
	err = reconfigure.Phase1(c)
	if err != nil {
		log.Errorf("Couldn't reconfigure: %s", err)
		return err
	}
	log.Info("Restore complete.")

	return nil
}

func phase2(c *cli.Context) error {
	replicaID := flags.ReplicaID
	if replicaID == "" {
		// Should never reach here...
		err := fmt.Errorf("Must specify DTR replica with --%s setting.", flags.ReplicaIDFlag.Name)
		return err
	}
	configOnly := flags.ConfigOnly

	var schemaManager schema.Manager
	var in io.Reader
	in = os.Stdin

	bs, kvStore, err := backup.SetupEnv(c, replicaID)
	if err != nil {
		errMessage := strings.Join([]string{"Can't set up environment to do a restore:", err.Error()}, " ")
		log.Error(errMessage)
		return err
	}

	configBuffer := new(bytes.Buffer)
	bufWriter := bufio.NewWriter(configBuffer)
	tr := tar.NewReader(in)

	log.Info("Reading input from stdin...")
	if !configOnly {
		schemaManager = bootstrap.CreateManager(bs.GetReplicaID())
		if err := schemaManager.Initialize(); err != nil {
			log.Errorf("Can't initialize schema manager: %s", err)
		}
	}

	versionIdentified := false
	var restoreVersion string

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			log.Info("Detected end of tar stream")
			break
		}
		if err != nil {
			return err
		}

		// The root folder will have the version number in its name
		// so we can immediately detect the version of the backup
		if !versionIdentified {
			restoreVersion = findBackupVersion(hdr.Name)
			versionIdentified = true
			continue
		}

		// Check if the following files are versioned properly according
		// to the version identified above, if not skip them
		rootDirName := strings.Split(hdr.Name, "/")[0]                                                         // no clean library function to get root dir in path :(
		if !strings.Contains(rootDirName, string(restoreVersion)) && !version.Equal(restoreVersion, "2.0.0") { // since 2.0.0 isn't versioned in the file name
			continue
		}

		if version.GreaterThanOrEqualTo(restoreVersion, "2.0.1") && version.LessThanOrEqualTo(restoreVersion, deploy.Version) {
			// trim the directory name so we can detect prefix in the next stage
			hdr.Name = strings.TrimPrefix(hdr.Name, fmt.Sprintf("dtr-backup-v%s/", restoreVersion))
		} else if version.Equal(restoreVersion, "2.0.0") {
			if hdr.Name == deploy.EtcdDirectory || hdr.Name == deploy.RethinkdbDirectory {
				// reached the etcd/rethinkdb section of the backup tar archive
				// continue to the individual files within each folder
				continue
			}
		} else {
			log.Warn("The backup version you're trying to restore isn't supported")
			continue
		}

		if strings.HasPrefix(hdr.Name, deploy.EtcdDirectory) {
			if hdr.Typeflag == tar.TypeDir {
				continue
			}
			if err := restoreEtcdData(tr, hdr, bufWriter, configBuffer, bs, kvStore, replicaID, restoreVersion); err != nil {
				log.Errorf("Couldn't restore etcd data: %s", err.Error())
				return err
			}
		} else if strings.HasPrefix(hdr.Name, deploy.RethinkdbDirectory) && !configOnly {
			if hdr.Typeflag == tar.TypeDir {
				continue
			}
			if err := restoreRethinkData(tr, hdr, bufWriter, configBuffer, bs, replicaID, restoreVersion, schemaManager); err != nil {
				log.Errorf("Couldn't restore rethink data: %s", err.Error())
				return err
			}
		} else {
			// Should never reach here while we're only backing up etcd and rethinkdb stuff...
			log.Infof("%s", hdr.Name)
			log.Warn("Unexpected header entry: not rethinkdb or etcd content")
		}

	}

	log.Infof("Restore successfully completed for replica %s", replicaID)
	return nil
}

func checkRestoreVersion(restoreVersion string) error {
	versionSemver, err := versions.TagToSemver(deploy.Version)
	if err != nil {
		return err
	}
	restoreSemver, err := versions.TagToSemver(restoreVersion)
	if err != nil {
		return err
	}

	// Only allow patch version differences
	if versionSemver.Major != restoreSemver.Major || versionSemver.Minor != restoreSemver.Minor {
		return fmt.Errorf("Couldn't restore from backup version %s on %s. Supported versions to restore from: %d.%d.x", deploy.Version, restoreVersion, restoreSemver.Major, restoreSemver.Minor)
	}

	return nil
}

func findBackupVersion(fileName string) string {
	// The first backup doesn't have a version so do an explicit check
	// to determine if the given backup tar is versioned or not
	if strings.Contains(fileName, backup.BackupVersionPrefix) {
		backupVersion := strings.TrimPrefix(fileName, backup.BackupVersionPrefix)
		// the root entry is a directory and we only want the
		// version within the directory name
		return strings.TrimSuffix(backupVersion, "/")
	} else {
		// default to the first version
		return "2.0.0"
	}
}

func restoreRethinkData(tr *tar.Reader, hdr *tar.Header, bufWriter *bufio.Writer, configBuffer *bytes.Buffer, bs bootstrap.Bootstrap, replicaID string, restoreVersion string, schemaManager schema.Manager) error {
	var tableName string
	var document []byte
	var err error

	if _, err := io.Copy(bufWriter, tr); err != nil {
		return fmt.Errorf("Error copying from tar archive to buffer writer: %s", err)
	}

	// if a version is greater than deploy.Version we won't even get to here
	// we'll error out upstream
	if version.GreaterThanOrEqualTo(restoreVersion, "2.0.1") {
		// a file for rethink is named like
		// dtr_backup_v2.0.1/rethink/repositories/1
		// so we're interested in the last directory as it's the table name
		dir := strings.Split(hdr.Name, "/")
		tableName = dir[len(dir)-2]
		document = []byte(configBuffer.String())
	} else {
		tableName = strings.TrimPrefix(hdr.Name, deploy.RethinkdbDirectory+"/")
		document, err = base64.StdEncoding.DecodeString(configBuffer.String())
		if err != nil {
			return err
		}
	}
	configBuffer.Reset()

	if err := schemaManager.RestoreTableDocument(tableName, document); err != nil {
		return fmt.Errorf("Failed to backup a rethink table %s's row: %s", tableName, err)
	}
	return nil
}

func restoreEtcdData(tr *tar.Reader, hdr *tar.Header, bufWriter *bufio.Writer, configBuffer *bytes.Buffer, bs bootstrap.Bootstrap, kvStore hubconfig.KeyValueStore, replicaID string, restoreVersion string) error {
	var err error
	var entryName string
	var configBytes []byte
	if _, err := io.Copy(bufWriter, tr); err != nil {
		log.Errorf("Error copying tar data to buffer: %s", err.Error())
		return err
	}

	// if a version is greater than deploy.Version we won't even get to here
	// we'll error out upstream
	if version.GreaterThanOrEqualTo(restoreVersion, "2.0.1") {
		// similar to rethink's restore, an etcd backup file structure is like
		// dtr_backup_v2.0.1/etcd/hub.yml
		// but for etcd we're interested in the file name
		_, entryName = filepath.Split(hdr.Name)
		configBytes = []byte(configBuffer.String())
	} else {
		entryName = strings.TrimPrefix(hdr.Name, deploy.EtcdDirectory+"/")
		configBytes, err = base64.StdEncoding.DecodeString(configBuffer.String())
		if err != nil {
			return err
		}
	}
	configBuffer.Reset()

	if entryName == deploy.HAConfigFilename {
		configBytes, err = updateHAConfig(kvStore, configBytes)
		if err != nil {
			log.Errorf("Couldn't properly restore HA config: %s", err)
			return err
		}
	}

	if err := kvStore.Put(entryName, configBytes); err != nil {
		log.Errorf("Can't restore state for %s: %s", hdr.Name, err.Error())
		return err
	}

	return nil
}

func updateHAConfig(kvStore hubconfig.KeyValueStore, backedUpConfigBytes []byte) (updatedConfigBytes []byte, err error) {
	var existingHAConfig, backedupHAConfig hubconfig.HAConfig
	if existingConfigBytes, err := kvStore.Get(deploy.HAConfigFilename); err != nil {
		log.Errorf("Couldn't get existing config to update: %s", err)
		return nil, err
	} else {
		if err := json.Unmarshal(existingConfigBytes, &existingHAConfig); err != nil {
			log.Errorf("Couldn't unmarshal HA config data: %s", err)
			return nil, err
		}
	}

	if err := json.Unmarshal(backedUpConfigBytes, &backedupHAConfig); err != nil {
		log.Errorf("Couldn't unmarshal HA config data: %s", err)
		return nil, err
	}
	// Restoring everything except the ReplicaConfig field or the ucp config
	// The ucp config doesn't need to be preserved, but we do it in case the restore
	// fails for some reason
	newHAConfig := backedupHAConfig
	newHAConfig.ReplicaConfig = existingHAConfig.ReplicaConfig
	newHAConfig.UCPHost = existingHAConfig.UCPHost
	newHAConfig.UCPCA = existingHAConfig.UCPCA
	newHAConfig.UCPVerifyCert = existingHAConfig.UCPVerifyCert

	if updatedConfigBytes, err := json.Marshal(newHAConfig); err != nil {
		log.Errorf("Couldn't marshal HA config data: %s", err)
		return nil, err
	} else {
		return updatedConfigBytes, nil
	}
}

func restore(c *cli.Context) error {
	if !bootstrap.IsPhase2() {
		return phase1(c)
	} else {
		return phase2(c)
	}
	return nil
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	if err := restore(c); err != nil {
		log.Fatal(err)
	}
}
