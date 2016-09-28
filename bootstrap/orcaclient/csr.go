package orcaclient

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/orca"
	"github.com/docker/orca/bootstrap/certs"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/utils"
	"golang.org/x/crypto/openpgp"
)

func sendNodeRequest(nodeReq *orca.NodeRequest, orcaUser, token string, httpClient *http.Client, orcaURL *url.URL) (*orca.NodeConfiguration, error) {
	reqJson, err := json.Marshal(*nodeReq)
	if err != nil {
		log.Debug("Failed to generate json csr")
		return nil, err
	}

	orcaURL.Path = "/api/nodes/authorize"

	// TODO - figure out what the common failure modes are and set up better messages
	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(reqJson))
	if err != nil {
		log.Debug("Failed to build request")
		return nil, err
	}
	req.Header.Set(GetTokenHeader(token, orcaUser))
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Debug("Failed to send request")
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Debugf("Response code: %d", resp.StatusCode)
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Errorf("Server response: %s", string(body))
			return nil, fmt.Errorf("Failed to add host to UCP: %s", string(body))
		}
		return nil, errors.New("Failed to add host to Orca")
	}

	var nodeConfig orca.NodeConfiguration
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		if err := json.Unmarshal(body, &nodeConfig); err != nil {
			return nil, fmt.Errorf("Failed to add host to UCP: %s", err)
		}
	} else {
		return nil, fmt.Errorf("Failed to add host to UCP: %s", err)
	}
	return &nodeConfig, nil
}

func DoJoin(nodeID, orcaUser, token, passphrase string, replica, byoCerts bool, httpClient *http.Client, orcaURL *url.URL) (*orca.NodeConfiguration, error) {

	// mount point, commonName
	// config.SwarmNodeCertVolumeMount, "swarm node",
	// config.SwarmKvCertVolumeMount, "kv"
	// config.SwarmControllerCertVolumeMount, "controller", true
	// config.OrcaServerCertVolumeMount, "controller" false (non-byo scenario)

	var nodeConfig *orca.NodeConfiguration

	// Set up the common CSRs
	type certTarget struct {
		caMount   string
		certMount string
		cn        string
		ou        string
	}
	targets := []certTarget{
		{
			caMount:   config.SwarmRootCAVolumeMount,
			certMount: config.SwarmNodeCertVolumeMount,
			cn:        nodeID,
			ou:        "swarm",
		},
	}

	// Append the replica items
	if replica {
		log.Debug("Joining as a replica")
		// Check to see if the user provided a backup with CA material
		if _, err := os.Stat(config.BackupFile); err == nil {
			file, err := os.Open(config.BackupFile)
			if err != nil {
				return nil, fmt.Errorf("Failed to read %s: %s", config.BackupFile, err)
			}
			log.Debug("Injecting user provided root CA cert/key pair")
			var in io.Reader
			if passphrase != "" {
				log.Debug("Decrypting with user provided passphrase")
				md, err := openpgp.ReadMessage(
					file,
					nil,
					func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
						return []byte(passphrase), nil
					},
					nil,
				)
				if err != nil {
					return nil, fmt.Errorf("Failed to setup decryption: %s", err)
				}
				in = md.UnverifiedBody
			} else {
				in = file
			}
			tr := tar.NewReader(in)

			log.Debug("Injecting user supplied replica CA certs/keys")
			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					log.Debug("Detected end of tar stream")
					// end of tar archive
					break
				}
				if err != nil {
					return nil, err
				}
				// Only restore files for the two CA cert paths
				if !(strings.HasPrefix(hdr.Name, "./"+config.OrcaRootCAVolumeName) ||
					strings.HasPrefix(hdr.Name, "./"+config.SwarmRootCAVolumeName)) {
					log.Debugf("Skipping non-CA file %s", hdr.Name)
					continue
				}
				filename := filepath.Join(config.Phase2VolMountDir, hdr.Name)
				if hdr.Typeflag == tar.TypeDir {
					if _, err := os.Stat(filename); os.IsNotExist(err) {
						err := os.Mkdir(filename, os.FileMode(hdr.Mode))
						if err != nil {
							return nil, fmt.Errorf("Failed to recreate directory %s: %s", filename, err)
						}
					}
				} else {
					log.Debugf("%s - %d", hdr.Name, hdr.Size)
					fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.FileMode(hdr.Mode))
					if err != nil {
						return nil, fmt.Errorf("Failed to restore file %s: %s", filename, err)
					}
					if _, err := io.Copy(fp, tr); err != nil {
						fp.Close()
						return nil, fmt.Errorf("Failed to restore file %s: %s", filename, err)
					}
					fp.Close()
				}
			}
		} else {
			log.Debug("Generating placeholder root CAs")
			targets = append(targets, certTarget{
				caMount:   config.OrcaRootCAVolumeMount,
				certMount: config.OrcaRootCAVolumeMount,
				cn:        "Placeholder invalid client root CA - please replace me",
			})
			targets = append(targets, certTarget{
				caMount:   config.SwarmRootCAVolumeMount,
				certMount: config.SwarmRootCAVolumeMount,
				cn:        "Placeholder invalid cluster root CA - please replace me",
			})
		}
		if !byoCerts {
			targets = append(targets, certTarget{
				caMount:   config.OrcaRootCAVolumeMount,
				certMount: config.OrcaServerCertVolumeMount,
				cn:        "ucp",
			})

		}
		targets = append(targets,
			certTarget{
				caMount:   config.SwarmRootCAVolumeMount,
				certMount: config.SwarmKvCertVolumeMount,
				cn:        nodeID,
				ou:        "kv",
			},
			certTarget{
				caMount:   config.SwarmRootCAVolumeMount,
				certMount: config.SwarmControllerCertVolumeMount,
				cn:        nodeID,
				ou:        "ucp",
			},
			certTarget{
				caMount:   config.SwarmRootCAVolumeMount,
				certMount: config.AuthStoreCertsVolumeMount,
				cn:        nodeID,
				ou:        "auth-store",
			},
			certTarget{
				caMount:   config.SwarmRootCAVolumeMount,
				certMount: config.AuthAPICertsVolumeMount,
				cn:        nodeID,
				ou:        "auth-api",
			},
			certTarget{
				caMount:   config.SwarmRootCAVolumeMount,
				certMount: config.AuthWorkerCertsVolumeMount,
				cn:        nodeID,
				ou:        "auth-worker",
			},
		)
	} else {
		log.Debug("Joining as a non-replica")
	}

	// Now generate the request
	req := orca.NodeRequest{
		ClusterCertificateRequests: map[string]string{},
		UserCertificateRequests:    map[string]string{},
		Replica:                    false,
	}
	for _, tgt := range targets {
		csr, key, err := certs.GenerateCSR(tgt.cn, tgt.ou, config.OrcaHostnames)
		if err != nil {
			log.Debug("Failed to generate CSR")
			return nil, err
		}
		// Write out the keys to the proper cert locations
		if err := ioutil.WriteFile(filepath.Join(tgt.certMount, config.KeyFilename), key, 0600); err != nil {
			return nil, err
		}
		// Get the CSR signed by the intended CA
		if tgt.caMount == config.SwarmRootCAVolumeMount {
			// Use the mount point so we can write out to the location specified
			req.ClusterCertificateRequests[tgt.certMount] = string(csr)
		} else {
			req.UserCertificateRequests[tgt.certMount] = string(csr)
		}
	}
	var err error
	nodeConfig, err = sendNodeRequest(&req, orcaUser, token, httpClient, orcaURL)
	if err != nil {
		return nil, err
	}
	// Write out the results
	for dir, cert := range nodeConfig.UserCertificates {
		// Note: if we supported intermediates in the future, we most likely need to join like this
		//cert := utils.JoinCerts(cert, nodeConfig.UserCertificateChain)
		if err := ioutil.WriteFile(filepath.Join(dir, config.CertFilename), []byte(cert), 0644); err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filepath.Join(dir, config.CAFilename), []byte(nodeConfig.UserCertificateChain), 0644); err != nil {
			return nil, err
		}
	}
	for dir, cert := range nodeConfig.ClusterCertificates {
		// Note: if we supported intermediates in the future, we most likely need to join like this
		// cert := utils.JoinCerts(cert, nodeConfig.ClusterCertificateChain)
		if err := ioutil.WriteFile(filepath.Join(dir, config.CertFilename), []byte(cert), 0644); err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filepath.Join(dir, config.CAFilename), []byte(nodeConfig.ClusterCertificateChain), 0644); err != nil {
			return nil, err
		}
	}
	if byoCerts {
		log.Debug("Appending user provided chain of trust with root CA certs")
		caPath := filepath.Join(config.OrcaServerCertVolumeMount, "ca.pem")
		caData, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, fmt.Errorf("Failed to process provided CA: %s", err)
		}
		caChain := utils.JoinCerts(nodeConfig.ClusterCertificateChain, nodeConfig.UserCertificateChain, string(caData))
		clusterCaPath := filepath.Join(config.SwarmControllerCertVolumeMount, "ca.pem")
		err = ioutil.WriteFile(clusterCaPath, []byte(caChain), 0655)
		if err != nil {
			return nil, fmt.Errorf("Failed to write CA: %s", err)
		}
		// Note: we do *NOT* glue the chain of trust onto the server cert.pem, as it was signed externally
	}
	if nodeConfig.OrcaID != "" {
		log.Debugf("Joining UCP ID: %s", nodeConfig.OrcaID)
		config.OrcaInstanceID = nodeConfig.OrcaID
	} // XXX is there an else case here we should worry about?

	return nodeConfig, nil
}
