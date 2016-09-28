package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/strslice"

	"github.com/docker/orca/bootstrap/config"
)

// CleanupVolumes removes any state left in volume mount locations before
// performing an install or join.
func CleanupVolumes(c *cli.Context) error {
	// Start with data volumes.
	cleanup := []string{
		config.OrcaKVVolumeMount,
		config.AuthStoreDataVolumeMount,
		config.AuthWorkerDataVolumeMount,
	}

	// Leave the all the certs alone if preserve-certs.
	if !c.Bool("preserve-certs") {
		// Purge all the swarm related certs
		cleanup = append(cleanup,
			config.OrcaRootCAVolumeMount,
			config.SwarmRootCAVolumeMount,
			config.SwarmNodeCertVolumeMount,
			config.SwarmKvCertVolumeMount,
			config.SwarmControllerCertVolumeMount,
			config.AuthStoreCertsVolumeMount,
			config.AuthAPICertsVolumeMount,
			config.AuthWorkerCertsVolumeMount,
		)

		// Only purge the controller server cert if it was not brought
		// by the user from an external CA.
		if !(c.Bool("external-ucp-ca") || c.Bool("external-server-cert")) {
			cleanup = append(cleanup, config.OrcaServerCertVolumeMount)
		}
	}

	for _, mount := range cleanup {
		log.Debugf("Cleaning stale data from %s", mount)
		if err := config.CleanupState(mount); err != nil {
			return err
		}
	}

	return nil
}

func VerifyPermissions(doCerts bool) error {
	nobodyPaths := []string{
		config.OrcaKVVolumeMount,
		config.AuthStoreDataVolumeMount,
		config.AuthWorkerDataVolumeMount,
		// TODO - more?
	}
	if doCerts {
		nobodyPaths = append(nobodyPaths,
			config.SwarmNodeCertVolumeMount,
			config.SwarmKvCertVolumeMount,
			config.SwarmControllerCertVolumeMount,
			config.AuthStoreCertsVolumeMount,
			config.AuthAPICertsVolumeMount,
			config.AuthWorkerCertsVolumeMount,
			config.OrcaServerCertVolumeMount,
		)
	}
	var lasterr error
	for _, mount := range nobodyPaths {
		err := os.Chown(mount, 65534, 65534)
		if err != nil {
			log.Warnf("Failed to update permissions on mount: %s %s", mount, err)
			lasterr = err
		}
		err = filepath.Walk(mount, func(path string, info os.FileInfo, err error) error {
			err = os.Chown(path, 65534, 65534)
			if err != nil {
				log.Warnf("Failed to update permissions: %s %s", path, err)
				lasterr = err
			}
			return nil
		})
		if err != nil {
			log.Warnf("Failed to update permissions on walk: %s %s", mount, err)
			lasterr = err
		}
	}
	return lasterr
}

// DecodeContainerInfo generates a ContainerJSON object from a []byte representing its json encoding
func DecodeContainerInfo(jsonData []byte) (types.ContainerJSON, error) {
	var res types.ContainerJSON
	err := json.Unmarshal(jsonData, &res)
	return res, err
}

// GetEtcdAdvertiseHostFromCmd extracts the ip used in the --advertise-client-urls argument of a container
func GetEtcdAdvertiseHostFromCmd(cmd strslice.StrSlice) (string, error) {
	for idx, arg := range cmd {
		if arg == "--advertise-client-urls" && idx < len(cmd)-1 {
			urlObj, err := url.Parse(cmd[idx+1])
			if err != nil {
				return "", fmt.Errorf("unable to parse the etcd advertise url: %s", err)
			}
			return urlObj.Host, nil
		}
	}

	return "", fmt.Errorf("unable to locate a value for the --advertise-client-urls argument in the etcd configuration of the backup file")
}
