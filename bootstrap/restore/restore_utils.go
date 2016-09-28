package restore

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/openpgp"

	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/utils"
	orcaconfig "github.com/docker/orca/config"
)

func isStdinAttached(ec *client.EngineClient) (bool, error) {
	info, err := ec.InspectContainer(orcaconfig.BootstrapContainerName)
	if err != nil {
		return false, fmt.Errorf("unable to inspect the \"ucp\" container: %s", err)
	}
	return info.Config.AttachStdin, nil
}

func getTarReader(reader io.Reader, passphrase string) (*tar.Reader, error) {
	if passphrase == "" {
		return tar.NewReader(reader), nil
	}

	log.Debug("Decrypting stdin tar stream with user provided passphrase")
	md, err := openpgp.ReadMessage(
		reader,
		nil,
		func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
			return []byte(passphrase), nil
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup tar file decryption: %s", err)
	}
	return tar.NewReader(md.UnverifiedBody), nil
}

// parseInspectDumps returns the IP of this node as inspected from the backed up ucp-kv container
// parseInspectDumps iterates through files in a tar Reader until a non-json file is found.
// parseInspectDumps is meant to be run asynchronously and either fail or block if the backup file -
// appears to be corrupt.
func parseInspectDumps(ec *client.EngineClient, tr *tar.Reader, rootCAOnly bool) (string, error) {
	var ip string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			log.Debug("Detected end of tar stream")
			return ip, fmt.Errorf("Detected end of tar stream before reaching volume files")
		}
		if err != nil {
			return ip, err
		}

		filename := filepath.Join(config.Phase2VolMountDir, hdr.Name)
		if hdr.Typeflag == tar.TypeDir {
			// If a directory was found, the are no more container inspect dumps in the tar stream
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				err := os.Mkdir(filename, os.FileMode(hdr.Mode))
				if err != nil {
					return ip, fmt.Errorf("Failed to recreate directory %s: %s", filename, err)
				}
			}

			break
		}

		// Verify that the name in the header correctly matches an inspect dump
		if !strings.Contains(hdr.Name, ".json") ||
			len(strings.Split(hdr.Name, "/")) != 2 ||
			!strings.HasPrefix(hdr.Name, "./") {
			break
		}

		// If we are only interested in a Root CA restore, move on to the next file
		if rootCAOnly {
			continue
		}

		// Handle the top-level .json files containing container inspect JSON dumps
		// It is OK to return an error from in here as we haven't actually restored any data
		// Skip everything except ./ucp-kv.json in this version
		// TODO: this is where we would also handle upgrade flows during restore
		if hdr.Name != "./"+orcaconfig.OrcaKvContainerName+".json" {
			continue
		}

		// Process the ./ucp-kv.json file
		buf := new(bytes.Buffer)
		numBytes, err := buf.ReadFrom(tr)
		if err != nil {
			return ip, fmt.Errorf("Failed to restore file %s: %s", filename, err)
		}
		if numBytes != hdr.Size {
			log.Warnf("did not read the expected number of bytes from %s", filename)
		}
		cinfo, err := utils.DecodeContainerInfo(buf.Bytes())
		if err != nil {
			return ip, fmt.Errorf("Failed to decode container info from file %s: %s", filename, err)
		}
		// Get the IP that the backed up etcd instance used to advertise
		host, err := utils.GetEtcdAdvertiseHostFromCmd(cinfo.Config.Cmd)
		if err != nil {
			return ip, fmt.Errorf("Failed to obtain advertise ip: %s", err)
		}
		hostParts := strings.Split(host, ":")
		ip = hostParts[0]

		// Get the IP that was advertised on the stopped ucp-etcd container in this node
		cinfo, err = ec.InspectContainer(orcaconfig.OrcaKvContainerName)
		if err != nil {
			return ip, err
		}
		localHost, err := utils.GetEtcdAdvertiseHostFromCmd(cinfo.Config.Cmd)
		localIP, _, err := net.SplitHostPort(localHost)
		if err != nil {
			return ip, err
		}

		if localIP != ip {
			return ip, fmt.Errorf("Unable to perform restore operation: the IP of the controller stored in the backup (%s) is different than the IP of this host (%s)", ip, localIP)
		}
		log.Debug("Local etcd IP and backup etcd IP match")
	}

	// Verify the ip and cmd are non-empty if we care about them
	if !rootCAOnly && ip == "" {
		return ip, fmt.Errorf("recovered an empty controller IP from the backup: %s", ip)
	}
	return ip, nil
}

func purgeVolumes(matches []string) []error {
	var errors []error

	// Point of no return (without corruption, that is)
	for _, filename := range matches {
		log.Debugf("Purging %s", filename)
		err := os.RemoveAll(filename)
		if err != nil {
			errors = append(errors, fmt.Errorf("Failed to remove %s - %s", filename, err))
		}
	}
	return errors
}

func restoreVolumes(tr *tar.Reader) []error {
	var errors []error

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			log.Debug("Detected end of tar stream")
			break
		}
		if err != nil {
			return append(errors, err)
		}

		log.Debugf("%s - %d", hdr.Name, hdr.Size)

		// Construct the path to the target file for extraction
		var filename string
		filename = filepath.Join(config.Phase2VolMountDir, hdr.Name)
		if hdr.Typeflag == tar.TypeDir {
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				err := os.Mkdir(filename, os.FileMode(hdr.Mode))
				if err != nil {
					errors = append(errors, fmt.Errorf("Failed to recreate directory %s: %s", filename, err))
				}
			}
			if err := os.Chown(filename, hdr.Uid, hdr.Gid); err != nil {
				errors = append(errors, fmt.Errorf("Failed to set uid/gid for dir %s: %s", filename, err))
			}
			continue
		}

		fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.FileMode(hdr.Mode))
		if err != nil {
			errors = append(errors, fmt.Errorf("Failed to restore file %s: %s", filename, err))
		}
		if _, err := io.Copy(fp, tr); err != nil {
			errors = append(errors, fmt.Errorf("Failed to restore file %s: %s", filename, err))
		}
		if err := fp.Chown(hdr.Uid, hdr.Gid); err != nil {
			errors = append(errors, fmt.Errorf("Failed to set uid/gid for file %s: %s", filename, err))
		}
		fp.Close()

	}
	return errors
}

type member struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	PeerURLs   []string `json:"peerURLs,omitempty"`
	ClientURLs []string `json:"clientURLs,omitempty"`
}

type membersResponse struct {
	Members []member `json:"members"`
}

type peerURLRequest struct {
	PeerURLs []string `json:"peerURLs"`
}

func EtcdRestoreMembers(client *http.Client, localIP string) error {
	kvURL := fmt.Sprintf("https://%s:12379", localIP)
	resp, err := client.Get(kvURL + "/v2/members")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var membersResp membersResponse
	err = json.Unmarshal(body, &membersResp)
	if err != nil {
		return err
	}

	if len(membersResp.Members) == 0 {
		return fmt.Errorf("No members were found in the KV store")
	}

	log.Debug("Recovering KV cluster members")
	var localMemberID string
	for _, member := range membersResp.Members {
		if strings.HasSuffix(member.Name, localIP) {
			// Fix the Peer URLs for this member after deleting all other members
			localMemberID = member.ID
		} else {
			// Delete all KV members that are not orca-kv-<localIP>
			_, err := kvDelete(client, kvURL+"/v2/members")
			if err != nil {
				return err
			}
		}
	}

	if localMemberID == "" {
		return fmt.Errorf("Unable to find the member ID of the local host in etcd")
	}
	peerURLReq := peerURLRequest{
		PeerURLs: []string{"https://" + localIP + ":12380"},
	}
	reqData, err := json.Marshal(peerURLReq)
	if err != nil {
		return fmt.Errorf("Error while marshalling etcd Peer update request: %s", err)
	}

	log.Debug("Recovering Peer URLs for this host")
	req, err := http.NewRequest("PUT", kvURL+"/v2/members/"+localMemberID, bytes.NewBuffer(reqData))
	if err != nil {
		return fmt.Errorf("Error while updating etcd Peer urls: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while updating etcd Peer urls: %s", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debug(bytes.NewBuffer(body).String())

	return nil
}

func kvDelete(c *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while removing kv entry: %s", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while removing kv entry: %s", err)
	}
	return resp, nil
}
