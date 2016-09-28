package reconfigure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/ucpclient"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/dhe-deploy/manager/versions"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/engine-api/types"
)

// this is a separate error type so it can be ignored by the caller
func NewAuthVerifyError(err error) AuthVerifyError {
	return AuthVerifyError{
		err: err,
	}
}

type AuthVerifyError struct {
	err error
}

func (a AuthVerifyError) Error() string {
	return a.err.Error()
}

// This is a helper function for Reconfigure which sets up each of the flags and returns an HAConfig
// and whether or not we need to reregister with enzi
func setReconfigureOptions(bs bootstrap.Bootstrap, settingsStore hubconfig.SettingsStore, c *cli.Context) (*hubconfig.HAConfig, bool, error) {
	var haConfig *hubconfig.HAConfig
	var err error

	replicaID := bs.GetReplicaID()

	haConfig, err = settingsStore.HAConfig()
	if err != nil {
		return haConfig, false, err
	}
	hubConfig, err := settingsStore.UserHubConfig()
	if err != nil {
		return haConfig, false, err
	}

	needToReRegisterWithEnzi := false
	if flags.IsSet(c, flags.EnziCAFlag) {
		haConfig.EnziCA = flags.EnziCA
	}
	if flags.UCPInsecureTLS {
		log.Debug("Fetching UCP CA")
		httpClient, err := dtrutil.HTTPClient(true)
		if err != nil {
			return haConfig, false, err
		}
		ucp := ucpclient.New(flags.UCPHost(), httpClient)
		ca, err := ucp.GetCA()
		if err != nil {
			return haConfig, false, err
		}
		log.Debug("UCP CA: %s", ca)
		haConfig.UCPCA = ca
		flags.UCPCA = ca
	}
	haConfig.UCPVerifyCert = true
	if flags.IsSet(c, flags.UCPCAFlag) {
		haConfig.UCPCA = flags.UCPCA
	}
	if flags.IsSet(c, flags.EnziInsecureTLSFlag) {
		haConfig.EnziVerifyCert = !flags.EnziInsecureTLS
	}
	if flags.IsSet(c, flags.UCPHostFlag) {
		needToReRegisterWithEnzi = true
		haConfig.UCPHost = flags.UCPHost()
	}
	if flags.IsSet(c, flags.EnziHostFlag) {
		needToReRegisterWithEnzi = true
		haConfig.EnziHost = flags.EnziHost()
	}
	if flags.IsSet(c, flags.HTTPProxyFlag) {
		haConfig.HTTPProxy = flags.HTTPProxy
	}
	if flags.IsSet(c, flags.HTTPSProxyFlag) {
		haConfig.HTTPSProxy = flags.HTTPSProxy
	}
	if flags.IsSet(c, flags.NoProxyFlag) {
		haConfig.NoProxy = flags.NoProxy
	}
	// make sure the map is set
	if haConfig.ReplicaConfig == nil {
		haConfig.ReplicaConfig = map[string]hubconfig.ReplicaConfig{}
	}
	// ports for the replica
	replicaConfig := haConfig.ReplicaConfig[replicaID]
	if flags.ReplicaHTTPPort != 0 {
		replicaConfig.HTTPPort = uint16(flags.ReplicaHTTPPort)
	}
	if flags.ReplicaHTTPSPort != 0 {
		replicaConfig.HTTPSPort = uint16(flags.ReplicaHTTPSPort)
	}
	node := flags.UCPNode
	// if we don't know what node we're on, figure it out
	if !flags.NoUCP && node == "" {
		node, err = bs.ContainerNode(deploy.BootstrapPhase2ContainerName)
		if err != nil {
			return haConfig, false, fmt.Errorf("Couldn't figure out node for the second phase: %s", err)
		}
	}
	log.Debugf("Node name is '%s'", node)
	if node != "" {
		replicaConfig.Node = node
	}
	haConfig.ReplicaConfig[replicaID] = replicaConfig
	// ports for the load balancer
	if flags.DTRHost() != "" && flags.DTRHost() != hubConfig.DTRHost {
		needToReRegisterWithEnzi = true
		hubConfig.DTRHost = flags.DTRHost()
		// TODO: deduplicate cert regeneration
		domainName := strings.Split(hubConfig.DTRHost, ":")[0]
		// if you change the domain name, we may need a new cert
		err := util.HubConfigTLSDomainConsistent(hubConfig)
		// if it's no longer consistent for whatever reason, generate a new cert
		if err != nil {
			cert, err := util.GenTLSCert(domainName)
			if err != nil {
				return haConfig, false, err
			}

			err = util.SetTLSCertificateInHubConfig(hubConfig, cert, cert)
			if err != nil {
				return haConfig, false, err
			}
		}
	}

	if flags.IsSet(c, flags.EnablePProfFlag) {
		haConfig.EnablePProf = flags.EnablePProf
	}

	if flags.LogProtocol != "" {
		haConfig.LogProtocol = flags.LogProtocol
	}
	if flags.LogHost() != "" {
		haConfig.LogHost = flags.LogHost()
	}
	if flags.LogLevel != "" {
		haConfig.LogLevel = flags.LogLevel
	}
	// Perhaps flags parser or something is stripping the newline from the
	// certs so we do this quick hack to add back the \n for now
	if flags.LogProtocol == "tcp+tls" {
		if flags.IsSet(c, flags.LogTLSCACertFlag) {
			haConfig.LogTLSCACert = flags.LogTLSCACert + "\n"
		}
		if flags.IsSet(c, flags.LogTLSCertFlag) {
			haConfig.LogTLSCert = flags.LogTLSCert + "\n"
		}
		if flags.IsSet(c, flags.LogTLSKeyFlag) {
			haConfig.LogTLSKey = flags.LogTLSKey + "\n"
		}
		if flags.IsSet(c, flags.LogTLSSkipVerifyFlag) {
			haConfig.LogTLSSkipVerify = flags.LogTLSSkipVerify
		}
	}

	if flags.IsSet(c, flags.EtcdHeartbeatIntervalFlag) {
		haConfig.EtcdHeartbeatInterval = flags.EtcdHeartbeatInterval
	}
	if flags.IsSet(c, flags.EtcdElectionTimeoutFlag) {
		haConfig.EtcdElectionTimeout = flags.EtcdElectionTimeout
	}
	if flags.IsSet(c, flags.EtcdSnapshotCountFlag) {
		haConfig.EtcdSnapshotCount = flags.EtcdSnapshotCount
	}

	// validate log config
	if err := haConfig.LogTest(); err != nil {
		return haConfig, false, err
	}
	enziConfig := util.GetEnziConfig(haConfig)
	// validate enzi config
	httpClient, err := dtrutil.HTTPClient(!enziConfig.VerifyCert, enziConfig.CA)
	if err != nil {
		log.Fatalf(err.Error())
	}
	response, err := httpClient.Get(fmt.Sprintf("https://%s", enziConfig.Host))
	if err != nil {
		return haConfig, false, fmt.Errorf("Failed to connect to enzi: %s", err)
	}
	response.Body.Close()

	err = settingsStore.SetHAConfig(haConfig)
	if err != nil {
		return haConfig, false, err
	}
	err = settingsStore.SetUserHubConfig(hubConfig)
	if err != nil {
		return haConfig, false, err
	}

	return haConfig, needToReRegisterWithEnzi, nil
}

// TODO: refactor the arguments to this func into a struct and/or break it up
func Reconfigure(bs bootstrap.Bootstrap, kvStore hubconfig.KeyValueStore, settingsStore hubconfig.SettingsStore, c *cli.Context, dtrContainers []containers.DTRContainer, installTime bool, expectHealthy bool, upgrade bool, dontConfigureAuth bool) error {
	if len(dtrContainers) == 0 {
		dtrContainers = containers.AllContainers
	}
	log.Infof("Getting container configuration and starting containers...")
	replicaID := bs.GetReplicaID()
	if replicaID == "" {
		return fmt.Errorf("Internal error: replica ID not passed to second phase of reconfigure correctly")
	}
	schemaManager := bootstrap.CreateManager(replicaID)
	reregister := false

	var haConfig *hubconfig.HAConfig
	var err error

	if !upgrade {
		// this sets each of the params for the reconfigure
		haConfig, reregister, err = setReconfigureOptions(bs, settingsStore, c)
		if err != nil {
			return err
		}
	} else {
		haConfig, err = settingsStore.HAConfig()
		if err != nil {
			return err
		}
	}

	os.MkdirAll(deploy.LogsCertPathInContainer, 0600)
	if err = ioutil.WriteFile(filepath.Join(deploy.LogsCertPathInContainer, "ca.pem"), []byte(haConfig.LogTLSCACert), 0600); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(deploy.LogsCertPathInContainer, "cert.pem"), []byte(haConfig.LogTLSCert), 0600); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(deploy.LogsCertPathInContainer, "key.pem"), []byte(haConfig.LogTLSKey), 0600); err != nil {
		return err
	}

	allReplicaIDs := []string{}
	for id := range haConfig.ReplicaConfig {
		allReplicaIDs = append(allReplicaIDs, id)
	}
	// sort the replicas so that we always reconfigure in the same order
	// this mitigates potential loss of quorum when restarting containers if the
	// reconfigure is aborted in the middle
	sort.Strings(allReplicaIDs)
	log.Debug("Replicas to upgrade: %q", allReplicaIDs)

	// now that we've changed the configs in etcd, let's make them take effect
	// WARNING: replicaID inside this loop is not the same as replicaID outside this loop
	for _, replicaID := range allReplicaIDs {
		replicaConfig := haConfig.ReplicaConfig[replicaID]
		// TODO: wait for each replica to come up after restarting it
		bs.SetReplicaID(replicaID)
		bs.SetNodeName(replicaConfig.Node)
		// set the version of each replica if we're installing or upgrading
		if installTime || upgrade {
			log.Debugf("Setting version to: %s", deploy.ShortVersion)
			replicaConfig.Version = deploy.ShortVersion
			haConfig.ReplicaConfig[replicaID] = replicaConfig
		}
		if expectHealthy {
			log.Infof("Waiting for database to stabilize for up to %d seconds before attempting to reconfigure replica %s", schema.WaitTime, replicaID)
			err := schemaManager.WaitForReady()
			if err != nil {
				return fmt.Errorf("Failed to wait for rethink to be ready %s", err)
			}
		}
		for _, container := range containers.ContainerConfigs(c.Command.Name, replicaID, allReplicaIDs, replicaConfig.Node, haConfig, dtrContainers) {
			diff := false
			if !upgrade {
				diff, err = bs.ContainerDiff(container)
				if err != nil {
					return fmt.Errorf("Failed to diff container %s: %s", container.Name, err)
				}
			}
			if diff || upgrade {
				// first we check what image was used to create the container because we don't want to upgrade unexpectedly
				info, err := bs.ContainerInspect(container.Name)
				if err != nil {
					log.Debugf("Failed to find outdated container %s: %s", container.Name, err)
					container.Environment["DTR_VERSION"] = deploy.ShortVersion
				} else {
					if !upgrade {
						container.Image = info.Config.Image
						container.Environment["DTR_VERSION"] = getDTRVersion(info.Config.Env)
					} else {
						ver := getDTRVersion(info.Config.Env)
						upgradeVer, err := versions.TagToSemver(deploy.AllowUpgradesStartingFrom)
						if err != nil {
							return err
						}
						currVer, err := versions.TagToSemver(ver)
						if err != nil {
							return err
						}
						bsVer, err := versions.TagToSemver(deploy.ShortVersion)
						if err != nil {
							return err
						}
						if currVer.LT(upgradeVer) || currVer.GE(bsVer) {
							log.Infof("Container %s is at version %s which cannot be upgraded by this bootstrapper; skipping", container.Name, ver)
							continue
						}
						container.Environment["DTR_VERSION"] = deploy.ShortVersion
					}
				}

				// TODO: kill gracefully and handle recovery in case of creation error instead of just rm -f
				log.Infof("Recreating %s...", container.Name)
				err = bs.ContainerRemove(container.Name, types.ContainerRemoveOptions{
					Force: true,
				})
				if err != nil {
					log.Debugf("Failed to remove outdated container %s: %s", container.Name, err)
					if !bootstrap.IsNoSuchImageErr(err.Error()) {
						// if there was an error we don't know if the container was deleted or not,
						// so we try to create it bug don't get too upset if we fail
						cResp, err := bs.ContainerCreateFromContainerConfig(container)
						if err != nil {
							log.Warnf("Failed to create new container %s: %s", container.Name, err)
						} else {
							cId := cResp.ID
							if err = bs.ContainerStart(cId); err != nil {
								log.Warnf("Failed to start container %s: %s", container.Name, err)
							}
						}
						continue
					}
				}

				cResp, err := bs.ContainerCreateFromContainerConfig(container)
				if err != nil {
					return fmt.Errorf("Failed to create new container %s: %s", container.Name, err)
				}
				cId := cResp.ID
				if err = bs.ContainerStart(cId); err != nil {
					return fmt.Errorf("Failed to start container %s: %s", container.Name, err)
				}
			}
		}
	}
	// we set this back just in case something else needs it in the future :/
	bs.SetReplicaID(replicaID)
	bs.SetNodeName(haConfig.ReplicaConfig[replicaID].Node)

	// save the version if we're installing or upgrading
	if installTime || upgrade {
		err = settingsStore.SetHAConfig(haConfig)
		if err != nil {
			return err
		}
	}

	err = kvStore.Ping()
	if err != nil {
		log.Errorf("KV store did not come back up after reconfigure")
		return err
	}

	if (installTime || reregister) && !dontConfigureAuth {
		hubConfig, err := settingsStore.UserHubConfig()
		if err != nil {
			return NewAuthVerifyError(fmt.Errorf("Failed to get hubconfig after installation is over: %s", err))
		}

		enziSession, err := util.GetAuthAPISession(flags.Username, flags.Password, settingsStore)
		if err != nil {
			return err
		}
		if err = util.RegisterAuth(enziSession, settingsStore); err != nil {
			return fmt.Errorf("Couldn't register with enzi: %#v", err)
		}

		log.Info("Verifying auth settings...")

		err = kvStore.Ping()
		if err != nil {
			return NewAuthVerifyError(fmt.Errorf("KV store did not come back up after reconfigure: %s", err))
		}

		// even though we are talking basically to the same machine, we use the real certs dtr generated by reading them out of etcd
		httpClient, err := dtrutil.HTTPClient(false, hubConfig.WebTLSCA)
		if err != nil {
			return NewAuthVerifyError(fmt.Errorf("Failed to create http client after installation is over: %s", err))
		}
		// redirect handling hack
		var RedirectAttemptedError = errors.New("redirect")
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return RedirectAttemptedError
		}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://%s", hubConfig.DTRHost), nil)
		if err != nil {
			return NewAuthVerifyError(fmt.Errorf("error making request to authorize endpoint %s", err))
		}
		req.SetBasicAuth(flags.Username, flags.Password)

		// wait for nginx to reload itself
		err = dtrutil.Poll(time.Second*5, 30, func() error {
			log.Infof("Waiting for DTR to start...")
			response, err := dtrutil.DoRequestWithClient(req, httpClient)
			if err != nil {
				// warn here so we know that it's polling
				log.Debugf("Failed to connect to DTR: %s", err)
				return fmt.Errorf("Failed to connect to DTR: %s", err)
			}

			response.Body.Close()
			if response.StatusCode != 200 {
				// warn here so we know that it's polling
				log.Debugf("Unexpected status code when trying to log in: %d. Make sure DTR's host is accessible from UCP.", response.StatusCode)
				return fmt.Errorf("Unexpected status code when trying to log in: %d. Make sure DTR's host is accessible from UCP.", response.StatusCode)
			}
			return nil
		})
		if err != nil {
			return NewAuthVerifyError(fmt.Errorf("Failed to wait for dtr to come back up: %s", err))
		}
		log.Info("Authentication test passed.")
	}

	return nil
}

func getDTRVersion(envs []string) string {
	for _, env := range envs {
		s := strings.Split(env, "=")
		if s[0] == "DTR_VERSION" {
			return s[1]
		}
	}
	// If we don't find a version, assume 2.0.0
	log.Debugf("No version associated with container.  Assuming 2.0.0")
	return "2.0.0"
}

func Phase1(c *cli.Context) error {
	log.Info("Starting phase1 reconfigure")
	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}
	if flags.IsSet(c, flags.EnziHostFlag) || flags.UCPHost() != "" {
		if flags.Username == "" {
			return fmt.Errorf("username required for changing these settings because we need to re-configure auth")
		}
		if flags.Password == "" {
			return fmt.Errorf("password required for changing these settings because we need to re-configure auth")
		}
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}

	replica, err := bs.ExistingReplicaFlagPicker("Choose a replica to reconfigure", true)
	if err != nil {
		return err
	}

	if err := replica.CheckEqualVersion(deploy.Version); err != nil {
		return err
	}

	container := dropperOpts.MakeContainerConfig([]string{"reconfigure"})
	bootstrap.SetEnvFromFlags(c, container, flags.ReconfigureFlags...)

	bs.SetReplicaID(flags.ExistingReplicaID)
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
		return err
	}
	return nil
}

func phase2(c *cli.Context) error {
	log.Info("Starting phase2 reconfigure")

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

	err = Reconfigure(bs, kvStore, settingsStore, c, []containers.DTRContainer{}, false, true, false, false)

	return err
}

func reconfigure(c *cli.Context) error {
	if bootstrap.IsPhase2() {
		return phase2(c)
	}
	err := Phase1(c)
	if err != nil {
		log.Error(bootstrap.IdempotentMsg("Reconfigure"))
	}
	return err
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	if err := reconfigure(c); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
