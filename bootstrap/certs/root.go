package certs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	cfssllog "github.com/cloudflare/cfssl/log"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types/swarm"

	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/utils"
)

// certPattern is a Go regular expression that matches an EC Private Key
var certPattern = "(?s)(-+BEGIN EC PRIVATE KEY-+)(.+?)(-+END EC PRIVATE KEY-+)"

// Initialize a CA at mount with name caName
func InitCA(caName, mount string) error {
	// Quiet down the cfssl logging
	cfssllog.Level = cfssllog.LevelWarning

	// Check for existing files, skip init if all present
	if existing, err := verifyExisting(mount, true); err != nil {
		// We could just clobber them and regenerate, but that might mask real failures
		return fmt.Errorf("Existing certs for %s appear to be corrupt: %s", caName, err)
	} else if existing {
		log.Debugf("Reusing existing certs for %s", caName)
		return nil
	} // else proceed and generate them

	log.Infof("Generating %s", caName)
	req := csr.CertificateRequest{
		CN: caName,
		KeyRequest: &csr.BasicKeyRequest{
			A: config.KeyAlgo,
			S: config.KeySize,
		},
	}
	caCert, _, key, err := initca.New(&req)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(mount, config.CertFilename), caCert, 0644); err != nil {
		return err
	}
	if err := WriteCert(filepath.Join(mount, config.KeyFilename), key, 0600); err != nil {
		return err
	}

	return nil
}

func SetupOrcaTrust() error {
	// Make sure the client CA is trusted by the server cert chain
	orcaCA, err := ioutil.ReadFile(filepath.Join(config.OrcaRootCAVolumeMount, "cert.pem"))
	if err != nil {
		return err
	}
	// If the user brought their own cert, make sure we trust that root too
	externalCA, err := ioutil.ReadFile(filepath.Join(config.OrcaServerCertVolumeMount, "ca.pem"))
	if err != nil {
		// Give a slightly more detailed error since this is a potential user error scenario
		return fmt.Errorf("Unable to open the user provided server ca.pem: %s", err)
	}
	filename := filepath.Join(config.SwarmControllerCertVolumeMount, "ca.pem")
	serverCA, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	chain := utils.JoinCerts(string(serverCA), string(orcaCA), string(externalCA))
	if err := WriteCert(filename, []byte(chain), 0644); err != nil {
		return err
	}
	return nil
}

func mostRecentSnapshot() (string, error) {
	dir := "/var/lib/docker/swarm/raft/snap/"
	minPath := ""
	minModTime := time.Time{}
	// Get the last file to be modified
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip subfolders
		if info.IsDir() {
			return nil
		}

		modTime := info.ModTime()
		if modTime.After(minModTime) {
			minModTime = modTime
			minPath = path
		}
		return nil
	})
	if err != nil {
		return minPath, err
	}
	return minPath, nil
}

func triggerSnapshot(dclient *client.Client) error {
	// Inspect the swarm
	info, err := dclient.SwarmInspect(context.TODO())
	if err != nil {
		return err
	}
	// Update the snapshot interval to 1
	prevInterval := info.Spec.Raft.SnapshotInterval
	info.Spec.Raft.SnapshotInterval = 1

	err = dclient.SwarmUpdate(context.TODO(), info.Meta.Version, info.Spec, swarm.UpdateFlags{})
	if err != nil {
		return err
	}

	// Get the newer swarm version
	info, err = dclient.SwarmInspect(context.TODO())
	if err != nil {
		return err
	}
	// Set the snapshot interval back to its original version
	info.Spec.Raft.SnapshotInterval = prevInterval

	err = dclient.SwarmUpdate(context.TODO(), info.Meta.Version, info.Spec, swarm.UpdateFlags{})
	if err != nil {
		return err
	}
	return nil
}

func getRootCASubstring(dclient *client.Client) (string, error) {
	var res string
	snapshotFile, err := mostRecentSnapshot()
	if err != nil {
		return res, err
	}

	// If there's no local snapshot, try to trigger a snapshot
	for snapshotFile == "" {
		err := triggerSnapshot(dclient)
		if err != nil {
			return res, err
		}

		time.Sleep(3 * time.Second)

		snapshotFile, err = mostRecentSnapshot()
		if err != nil {
			return res, err
		}
	}

	buffRes, err := ioutil.ReadFile(snapshotFile)
	if err != nil {
		return res, err
	}

	r, err := regexp.Compile(certPattern)
	if err != nil {
		return "", err
	}

	results := r.FindAllString(string(buffRes), -1)
	if len(results) != 1 {
		return "", fmt.Errorf("found %d private keys instead of one", len(results))
	}

	return results[0], nil
}

// CopySwarmRootCA obtains the Root Key and Cert from the present swarm manager
// and copies them to the required location in the ucp-cluster-root-ca volume
func CopySwarmRootCA(dclient *client.Client) error {
	log.Info("Establishing mutual Cluster Root CA with Swarm")
	rootKey, err := getRootCASubstring(dclient)
	if err != nil {
		return err
	}

	// Read the Root CA Cert
	sourceCrt, err := os.Open("/var/lib/docker/swarm/certificates/swarm-root-ca.crt")
	if err != nil {
		return err
	}

	// Write the Root CA Cert to the target volume
	destCrt, err := os.OpenFile(config.SwarmRootCAVolumeMount+"/cert.pem",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer destCrt.Close()

	_, err = io.Copy(destCrt, sourceCrt)
	if err != nil {
		return err
	}

	// Write the Root Key
	destKey, err := os.OpenFile(config.SwarmRootCAVolumeMount+"/key.pem",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer destKey.Close()
	_, err = destKey.WriteString(rootKey)
	if err != nil {
		return err
	}

	return nil
}
