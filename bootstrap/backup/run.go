// Package backup will stream a tar out containing all the local nodes state
package backup

import (
	"archive/tar"
	"bufio"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"

	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

func pathWalker(tw *tar.Writer, striplen int) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		new_path := "." + path[striplen:]
		if len(new_path) == 1 {
			return nil
		}
		// Make sure we can open it first, skip if we can't
		fd, err := os.Open(path)
		if err != nil {
			// Fail hard if the backup will be corrupted. We don't like corrupt backups
			return fmt.Errorf("Unable to open %s - backup will be incomplete: %s", path, err)
		}
		defer fd.Close()
		hdr, err := tar.FileInfoHeader(info, new_path)
		if err != nil {
			// Should never happen
			return err
		}
		// FileInfoHeader doesn't get the name right
		if info.Mode().IsDir() {
			hdr.Name = new_path + "/"
		} else {
			hdr.Name = new_path
		}
		log.Debugf("%s - %d", hdr.Name, hdr.Size)
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if !info.Mode().IsDir() {
			if _, err := io.Copy(tw, fd); err != nil {
				return err
			} else {
				tw.Flush()
			}
		} // else no body for directories
		return nil
	}
}

func backup(c *cli.Context) (int, error) {
	// Make sure all logging goes to stderr so stdout can be our tar file
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}
	if ec.IsTty() {
		return 1, fmt.Errorf("Backups will be corrupted with TTY mode enabled.  Please re-run without the '-t' flag.")
	}

	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)
	if !config.InPhase2 {
		nodeType := ec.DetectNodeType()
		if nodeType != client.Controller {
			return 1, fmt.Errorf("This node is not a controller.  In the future if you encounter problems on this node, re-run the 'join' command to reconnect to the cluster.")
		}

		// Figure out the instance ID of the local instance

		// TODO - does this work for already stopped containers?
		// TODO - should we remember the initial state so we don't restart if they were stopped?
		ids := client.GetInstanceIDs(containers)
		if len(ids) == 0 {
			return 1, fmt.Errorf("No running UCP instances detected on this engine")
		} else if len(ids) > 1 {
			log.Warnf("Multiple UCP instances detected: %v", ids)
		}
		id := ids[0]

		flagID := c.String("id")
		if flagID != "" {
			if flagID != id {
				return 1, fmt.Errorf("The provided UCP instance ID argument %s is not equal to the ID of the local UCP instance: %s", flagID, id)
			}
			config.OrcaInstanceID = id
		} else {
			// Empty --id flag flow
			if c.Bool("interactive") {
				if c.Bool("root-ca-only") {
					log.Infof("We're about to temporarily stop the local UCP CA containers for UCP ID: %s", id)
				} else {
					log.Infof("We're about to temporarily stop all local UCP containers for UCP ID: %s", id)
				}
				fmt.Fprintf(os.Stderr, "Do you want proceed with the backup? (y/n): ")

				reader := bufio.NewReader(os.Stdin)
				value, err := reader.ReadString('\n')
				if err != nil {
					log.Debugf("Failed to read input: %s", err)
					return 1, err
				}
				value = strings.TrimSpace(strings.ToLower(value))
				if value == "y" || value == "yes" {
					config.OrcaInstanceID = id
				} else {
					return 1, fmt.Errorf("Aborting backup by user request")
				}

			} else {
				log.Infof("We detected local components of UCP instance %s", id)
				return 1, fmt.Errorf(`Re-run the command with "--id %s" or --interactive to confirm you want to backup this UCP instance.`, id)
			}
		}

		containerIDs, err := ec.FindContainerIDsByOrcaInstanceID(config.OrcaInstanceID)
		if err != nil {
			log.Debug("Failed to find specified UCP instances")
			return 1, err
		}

		// Figure out all the local volumes that we need to mount
		volumes, err := ec.ListExistingOrcaVolumes()
		if err != nil {
			return 1, fmt.Errorf("Failed to lookup volumes: %s", err)
		}

		// Verify this node has both Root CA volumes
		caVolumes := 0
		for _, volume := range volumes {
			if volume == config.OrcaRootCAVolumeName || volume == config.SwarmRootCAVolumeName {
				caVolumes++
			}
		}
		if caVolumes != 2 {
			return 1, fmt.Errorf("Unable to detect both the UCP and Swarm Root CA volumes. The backup operation can only be performed on UCP controllers.")
		}

		// Verify this node has both Root CA containers
		caContainers := []string{}
		for _, id := range containerIDs {
			info, err := ec.InspectContainer(id)
			if err != nil {
				return 1, fmt.Errorf("Failed to lookup container %s name", id)
			}
			if strings.Contains(info.Name, orcaconfig.OrcaCAContainerName) || strings.Contains(info.Name, orcaconfig.OrcaSwarmCAContainerName) {
				caContainers = append(caContainers, id)
			}
		}
		if len(caContainers) != 2 {
			return 1, fmt.Errorf("This system does not appear to be running both Root CA containers")
		}

		if c.Bool("root-ca-only") {
			// Reduce the scope of volumes and containers to backup
			volumes = []string{config.OrcaRootCAVolumeName, config.SwarmRootCAVolumeName}
			containerIDs = caContainers
		} else {
			// If backing up everything, make sure no existing KV backup volume is present in the node
			// This is to prevent etcdtl from ever touching an existing backup that was not completed properly
			if ec.VolumeExists(config.OrcaKVBackupVolumeName) {
				return 1, fmt.Errorf("A %s volume is already present in this host. Backup operation aborted.", config.OrcaKVBackupVolumeName)
			}
		}

		for _, volume := range volumes {
			sourceVolume := volume
			targetVolume := volume
			if targetVolume == config.OrcaKVVolumeName {
				// Instead of the actual ucp-kv volume, mount ucp-kv-backup in the ucp-kv directory
				// Assuming the backup at the end of Phase 1 completes successfully
				sourceVolume = config.OrcaKVBackupVolumeName
			}
			config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
				fmt.Sprintf("%s:%s", sourceVolume,
					path.Join(config.Phase2VolMountDir, targetVolume)))
		}

		// Only stop the root CA
		if c.Bool("root-ca-only") {
			log.Infof("Temporarily stopping the local CA containers to ensure a consistent backup")
		} else {
			log.Infof("Temporarily stopping all local UCP containers to ensure a consistent backup")
		}

		for _, id := range containerIDs {
			log.Debugf("Stopping container %s", id)
			if err := ec.StopContainer(id); err != nil {
				log.Warningf("Failed to stop container %s (%s) Your backup may contain inconsistent data", id, err)
				// TODO - Do we want to give the user a flag to control if we bail on any stop failures?
			}
		}
		// Wire up the restart for after we're done
		defer func() {
			log.Infof("Resuming stopped UCP containers")
			for _, id := range containerIDs {
				log.Debugf("Restarting container %s", id)
				if err := ec.ContainerRestart(id); err != nil {
					// TODO - might want to give them the friendly name not id here, but inspect could fail too...
					log.Warningf("Failed to restart container %s (%s) ", id, err)
				}
			}
		}()

		if !c.Bool("root-ca-only") {
			// Backup etcd in Phase 1 after the containers are stopped
			log.Info("Backing up internal KV store")
			defer ec.RemoveVolume(config.OrcaKVBackupVolumeName)
			err = ec.TakeKVBackup()
			if err != nil {
				return 1, err
			}
			// The backup volume should be cleaned up when the backup operation is finished
			log.Debug("etcd backup completed successfully")
		}

		return ec.StartPhase2(os.Args[1:], false)

	} else { // This is phase 2
		log.Debug("Entering Phase 2")

		// Fix up kv backup permissions
		err := os.Chown(filepath.Join(config.Phase2VolMountDir, config.OrcaKVVolumeName), 65534, 65534)
		if err != nil {
			return 1, err
		}
		err = filepath.Walk(filepath.Join(config.Phase2VolMountDir, config.OrcaKVVolumeName), func(path string, info os.FileInfo, err error) error {
			err = os.Chown(path, 65534, 65534)
			if err != nil {
				log.Warn("Failed to update permissions: %s %s", path, err)
				return err
			}
			return nil
		})
		if err != nil {
			return 1, err
		}

		// First thing in phase 2 is to verify that the Root CAs are not signed with a Placeholder CommonName
		log.Debug("Validating Root CAs")

		validated := 0
		// Iterate over all the directories mounted in the phase2 container
		err = filepath.Walk(config.Phase2VolMountDir, func(path string, info os.FileInfo, err error) error {
			_, volumeName := filepath.Split(path)
			if volumeName == config.OrcaRootCAVolumeName || volumeName == config.SwarmRootCAVolumeName {
				log.Debug(volumeName)
				file, err := os.Open(filepath.Join(path, config.CertFilename))
				if err != nil {
					log.Errorf("Unable to find certificate file in the %s volume", volumeName)
					return err
				}
				data, err := ioutil.ReadAll(file)
				if err != nil {
					return err
				}

				der, _ := pem.Decode(data)
				if der == nil {
					return err
				}
				cert, err := x509.ParseCertificate(der.Bytes)
				if err != nil {
					return err
				}
				if strings.Contains(cert.Subject.CommonName, "Placeholder invalid") ||
					strings.Contains(cert.Issuer.CommonName, "Placeholder invalid") {
					return fmt.Errorf("This controller replica does not contain valid Root CA material. Please run the backup operation on a controller with valid Root CAs.")
				}
				validated++
			}
			return nil
		})
		if err != nil {
			return 1, err
		}

		if validated != 2 {
			return 1, fmt.Errorf("Failed to verify the validity of the UCP and Swarm Root CAs on this node. Aborting backup procedure")
		}

		// Create a tar file bound to stdout
		passphrase := c.String("passphrase")
		var out io.WriteCloser
		if passphrase != "" {
			log.Debug("Encrypting tar file with user provided passphrase")
			out, err = openpgp.SymmetricallyEncrypt(
				os.Stdout,
				[]byte(passphrase),
				&openpgp.FileHints{
					IsBinary: true,
					FileName: "_CONSOLE",
				},
				&packet.Config{},
			)
			if err != nil {
				return 1, fmt.Errorf("Failed to setup encryption: %s", err)
			}
			defer out.Close()
		} else {
			out = os.Stdout
		}
		log.Info("Beginning backup")
		tw := tar.NewWriter(out)
		defer tw.Close()

		// Dump a json file for the output of inspect for all UCP containers
		// Get all container IDs of UCP containers on the node
		log.Debug("Beginning tar dump of UCP container inspect information")

		for _, cinfo := range containers {
			// Strip any leading prefixes to the container name (typically just a single slash)
			cnameParts := strings.Split(cinfo.Name, "/")
			cname := cnameParts[len(cnameParts)-1]
			filename := cname + ".json"

			data, err := json.Marshal(cinfo)
			if err != nil {
				return 1, err
			}

			hdr := &tar.Header{
				Name: "./" + filename,
				Mode: 0644,
				Size: int64(len(data)),
			}

			err = tw.WriteHeader(hdr)
			if err != nil {
				return 1, err
			}
			log.Debugf("%s - %d", hdr.Name, hdr.Size)

			// Copy the file data to the tar stream
			_, err = tw.Write(data)
			if err != nil {
				return 1, err
			}
		}

		log.Debugf("Beginning tar dump of UCP volumes mounted in %s", config.Phase2VolMountDir)
		// Walk the tree rooted at
		if err := filepath.Walk(config.Phase2VolMountDir, pathWalker(tw, len(config.Phase2VolMountDir))); err != nil {
			log.Debug("Tar dump failed")
			return 1, err
		}

		log.Info("Backup completed successfully")
	}
	return 0, nil
}

// Run the backup flow
func Run(c *cli.Context) {
	if code, err := backup(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
