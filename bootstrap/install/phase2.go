package install

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/certificates"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/reconfigure"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/shared/containers"
	garantconfig "github.com/docker/garant/config"
)

func phase2(c *cli.Context) error {
	log.Debug("Starting phase2")
	var err error

	// this will not do any prompts if it's run correctly
	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return err
	}

	dtrHost := flags.DTRHost()
	if dtrHost == "" {
		return fmt.Errorf("DTR host parameter is required")
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return err
	}
	err = bootstrap.SetupReplicaAndNode(c, bs)
	if err != nil {
		return err
	}
	err = bs.ConfirmNodeIsOkay()
	if err != nil {
		return err
	}

	err = SetupNetworks(c, bs)
	if err != nil {
		return err
	}

	if !flags.NoUCP {
		err := bootstrap.UCPConnTest()
		if err != nil {
			return err
		}
	}

	err = SetupVolumes(c, bs)
	if err != nil {
		return err
	}

	if err = CreateCACerts(bs); err != nil {
		log.Debugf("Couldn't create client certs: %s", err)
		return err
	}

	if err = CreateClientCerts(bs); err != nil {
		log.Debugf("Couldn't create client certs: %s", err)
		return err
	}

	// Start etcd first, since it has all of our config files
	_, err = StartEtcd(bs, c.Command.Name, nil)
	if err != nil {
		return err
	}

	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		log.Errorf("Couldn't set up kvStore: %s", err)
		return err
	}
	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))
	if err = setupConfigs(dtrHost, settingsStore, c); err != nil {
		log.Errorf("Couldn't setup configs: %s", err)
		return err
	}

	if !flags.NoUCP {
		ucp := bs.GetUCPClient()
		ucpConfig, err := ucp.GetLicenseConfig()
		if err != nil {
			log.Errorf("Couldn't get license config from UCP: %s", err)
			return err
		}
		ucpLicense := ucpConfig.License
		if ucpLicense.KeyID != "" && ucpLicense.PrivateKey != "" && ucpLicense.Authorization != "" {
			dtrLicense := hubconfig.LicenseConfig{
				KeyID:         ucpLicense.KeyID,
				PrivateKey:    ucpLicense.PrivateKey,
				Authorization: ucpLicense.Authorization,
				AutoRefresh:   ucpConfig.AutoRefresh,
			}
			err = settingsStore.SetLicenseConfig(&dtrLicense)
			if err != nil {
				log.Errorf("Failed to set license from UCP: %s", err)
				return err
			}
			log.Info("License config copied from UCP.")
		} else {
			log.Info("License config not copied from UCP because UCP has no valid license.")
		}
	}

	log.Info("Starting rethinkdb...")
	// Configure Enzi host, ports, etc. This will also start the database.
	err = reconfigure.Reconfigure(bs, kvStore, settingsStore, c, []containers.DTRContainer{containers.Rethinkdb}, false, false, false, true)
	if err != nil {
		return err
	}

	// Init databases
	_, err = bootstrap.MigrateDatabase(bs.GetReplicaID(), 1)
	if err != nil {
		log.Errorf("Couldn't migrate database: %s", err)
		return err
	}

	// This will also start the rest of the containers.
	log.Info("Starting all containers...")
	err = reconfigure.Reconfigure(bs, kvStore, settingsStore, c, nil, true, false, false, false)
	if err != nil {
		if err2, ok := err.(reconfigure.AuthVerifyError); ok {
			log.Warnf("Couldn't confirm authentication works, but still completing installation: %s", err2)
		} else {
			log.Errorf("Couldn't reconfigure: %s", err)
			return err
		}
	}

	log.Debug("configuration done")

	log.Info("Installation is complete")
	log.Infof("Replica ID is set to: %s", bs.GetReplicaID())
	log.Infof("You can use flag '--%s %s' when joining other replicas to your Docker Trusted Registry Cluster", flags.ExistingReplicaIDFlag.Name, bs.GetReplicaID())

	return nil
}

func StartEtcd(bs bootstrap.Bootstrap, installType string, peers []string) (string, error) {
	// TODO: do the config bootstrap dance to make logging configs work for etcd
	containerConfig := containers.Etcd.ContainerConfig(installType, bs.GetReplicaID(), []string{}, bs.GetNodeName(), &defaultconfigs.DefaultHAConfig)

	// XXX - It would be nice if we could add this into the ContainerConfig() call
	// add self to list of peers
	peers = append(peers,
		fmt.Sprintf("%s=https://%s:%d", containers.Etcd.ReplicaName(bs.GetReplicaID()), containers.Etcd.OverlayName(bs.GetReplicaID()), containers.EtcdPeerPort1),
	)

	initialCluster := strings.Join(peers, ",")
	containerConfig.Environment["ETCD_INITIAL_CLUSTER"] = initialCluster

	// Constrain the container to run on the same container as the phase2 bootstrap
	etcdId, err := bs.ContainerRunFromContainerConfig(containerConfig)
	if err != nil {
		log.Errorf("Problem running container '%s': %s", containerConfig.Name, err)
		return "", err
	}
	return etcdId, nil
}

func SetupNetworks(c *cli.Context, bs bootstrap.Bootstrap) error {
	var err error

	overlayType := func() string {
		if flags.NoUCP {
			return "bridge"
		}
		return "overlay"
	}

	type netType struct {
		Name string
		Type string
	}

	networks := []netType{
		{
			Name: bootstrap.GetBridgeNetworkName(c, bs),
			Type: "bridge",
		},
		{
			Name: deploy.OverlayNetworkName,
			Type: overlayType(),
		},
	}

	// Create the network if it doesn't exist, and attach the phase2 container to it
	for _, net := range networks {
		if !bs.NetworkExists(net.Name) {
			log.Infof("Creating network: %s", net.Name)
			_, err = bs.CreateNetwork(net.Name, net.Type)
			if err != nil {
				err = fmt.Errorf("Can't create %s network '%s': %s", net.Type, net.Name, err)
				return err
			}
		} else {
			log.Debugf("Network %s already exists", net.Name)
		}

		err = bootstrap.Phase2NetworkConnect(bs, net.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetupVolumes(c *cli.Context, bs bootstrap.Bootstrap) error {
	// Create persistent storage volumes
	log.Info("Setting up replica volumes...")
	for _, vol := range containers.Volumes {
		volumeName := vol.ReplicaName(bs.GetReplicaID())
		if !flags.NoUCP {
			volumeName = fmt.Sprintf("%s/%s", bs.GetNodeName(), volumeName)
		}
		if !bs.VolumeExists(volumeName) {
			log.Debugf("Creating volume '%s'", volumeName)
			if err := bs.VolumeCreate(volumeName); err != nil {
				log.Errorf("Couldn't create volume '%s'", volumeName)
				return err
			}
		}
	}
	return nil
}

func setupConfigs(dtrHost string, settingsStore hubconfig.SettingsStore, c *cli.Context) error {
	// User Config
	log.Debug("Using default user config settings")
	defaultHubConfig, err := util.DefaultUserHubConfig(dtrHost)
	if err != nil {
		log.Errorf("Couldn't generate user config settings: %s", err)
		return err
	}

	if err := settingsStore.SetUserHubConfig(defaultHubConfig); err != nil {
		log.Errorf("Couldn't set user config settings: %s", err)
		return err
	}

	// HA config
	log.Debug("Using default HA settings")
	if err := settingsStore.SetHAConfig(&defaultconfigs.DefaultHAConfig); err != nil {
		log.Errorf("Couldn't set HA config settings: %s", err)
		return err
	}

	// Registry Config
	log.Debug("Using default registry settings")
	if err := settingsStore.SetRegistryConfig(&defaultconfigs.DefaultRegistryConfig); err != nil {
		log.Errorf("Couldn't set registry config settings: %s", err)
		return err
	}

	// Auth Config
	log.Debug("Using default auth config settings")
	authCfg := new(garantconfig.Configuration)
	if err := settingsStore.SetAuthConfig(authCfg); err != nil {
		log.Errorf("Couldn't set auth config settings: %s", err)
		return err
	}

	// Certificate Generation
	certGen := &certificates.FSCertificateGenerator{
		SettingsStore: settingsStore,
	}
	certGen.Generate()
	return nil
}

func CreateCACerts(bs bootstrap.Bootstrap) error {
	// Bootstrap the two CAs, use the existing CA if there is already one there
	log.Info("Creating initial CA certificates")
	err := bootstrapCA(containers.EtcdCACertStore, containers.Etcd.OverlayName(bs.GetReplicaID()))
	if err != nil {
		log.Errorf("Couldn't bootstrap etcd CA: %s", err)
		return err
	}
	err = bootstrapCA(containers.RethinkCACertStore, containers.Rethinkdb.OverlayName(bs.GetReplicaID()))
	if err != nil {
		log.Errorf("Couldn't bootstrap rethink CA: %s", err)
		return err
	}
	return nil
}

func CreateClientCerts(bs bootstrap.Bootstrap) error {
	var err error
	var altNames []string

	altNames = []string{containers.Etcd.BridgeName(bs.GetReplicaID()), containers.Etcd.OverlayName(bs.GetReplicaID())}
	containers.EtcdCertStore.Mkdirp()
	if err = MakeSignedClientCert(containers.EtcdCACertStore, containers.EtcdCertStore, altNames); err != nil {
		log.Errorf("Couldn't create new signed etcd certificate: %s", err)
		return err
	}

	altNames = []string{containers.Rethinkdb.BridgeName(bs.GetReplicaID()), containers.Rethinkdb.OverlayName(bs.GetReplicaID())}
	containers.RethinkCertStore.Mkdirp()
	if err = MakeSignedClientCert(containers.RethinkCACertStore, containers.RethinkCertStore, altNames); err != nil {
		log.Errorf("Couldn't create new signed rethink certificate: %s", err)
		return err
	}

	if err = CreateNotaryClientCerts(bs); err != nil {
		return err
	}

	return nil
}

func CreateNotaryClientCerts(bs bootstrap.Bootstrap) error {
	var altNames []string

	if err := bootstrapCA(containers.NotaryCACertStore, containers.NotaryServer.BridgeName(bs.GetReplicaID())); err != nil {
		log.Errorf("Couldn't bootstrap throwaway notary CA: %s", err)
		return err
	}
	// this is a throwaway CA, so once the certs are created, remove the key
	defer os.Remove(containers.NotaryCACertStore.KeyPath())

	// create the notary server and signer certs, server client certs
	altNames = []string{containers.NotaryServer.BridgeName(bs.GetReplicaID())}
	containers.NotaryServerStore.Mkdirp()
	if err := MakeSignedClientCert(containers.NotaryCACertStore, containers.NotaryServerStore, altNames); err != nil {
		log.Errorf("Couldn't create new signed notaryserver certificate: %s", err)
		return err
	}

	containers.NotaryCertStore.Mkdirp()
	if err := MakeSignedClientCert(containers.NotaryCACertStore, containers.NotaryCertStore, altNames); err != nil {
		log.Errorf("Couldn't create new signed notaryserver client certificate: %s", err)
		return err
	}

	altNames = []string{containers.NotarySigner.BridgeName(bs.GetReplicaID())}
	containers.NotarySignerStore.Mkdirp()
	if err := MakeSignedClientCert(containers.NotaryCACertStore, containers.NotarySignerStore, altNames); err != nil {
		log.Errorf("Couldn't create new signed notarysigner certificate: %s", err)
		return err
	}

	return nil
}
