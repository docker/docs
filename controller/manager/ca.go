package manager

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/pkg/pki"
	"github.com/docker/orca/utils"
)

var (
	clusterCAKvKeyPrefix = "clusterca"
	clientCAKvKeyPrefix  = "clientca"
	DummyPKI             = &pki.DummyClient{}
)

// Common
type CAConfiguration struct {
	Addr      string        `json:"addr"`
	PKIClient pki.PKIClient `json:"-"`
}

type CAConfigSubsystem struct {
	ksKey string
	m     *DefaultManager
	cfg   *CAConfiguration
}

//caType := strings.Split(key, "_")[0]
//if strings.HasPrefix(caType, "cluster")

func NewCAConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	if !strings.HasPrefix(key, clusterCAKvKeyPrefix) && !strings.HasPrefix(key, clientCAKvKeyPrefix) {
		return nil, fmt.Errorf("Malformed CA key: %s", key)
	}
	log.Debugf("Registering new CA %s", key)

	// TODO - might want to parse the instance information too and validate...

	s := CAConfigSubsystem{
		m:     m,
		ksKey: path.Join(KsConfigDir, key),
		cfg:   &CAConfiguration{PKIClient: DummyPKI},
	}

	cfgInt, err := s.ValidateConfig(jsonConfig, false)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgInt.(CAConfiguration)
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration type")
	}
	*s.cfg = cfg
	m.configSubsystems[key] = s
	return s, nil
}

func setupCA(m *DefaultManager) {
	// Not a singleton... only per-instance
	m.RegisterConfigSubsystem(clusterCAKvKeyPrefix, NewCAConfigSubsystem)
	m.RegisterConfigSubsystem(clientCAKvKeyPrefix, NewCAConfigSubsystem)

	// TODO - is this where we should migrate legacy CA configuration from pre 1.1?
}

func (s CAConfigSubsystem) GetKvKey() string {
	return s.ksKey
}

func (s CAConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(s.cfg)
	return string(data), err
}

// Attempt to initialize the PKI.  If something goes wrong stuff in the dummy
func (cfg *CAConfiguration) InitPKI(tlsConfig *tls.Config) {
	client, err := pki.NewDefaultClient(cfg.Addr, tlsConfig)
	if err != nil {
		log.Warnf("Failed to create PKI: (CA may be down, we will retry later) %s: %s", cfg.Addr, err)
		cfg.PKIClient = DummyPKI
		return
	}
	// Verify it's a good CA
	caCert, err := client.GetRootCertificate()
	if err != nil {
		log.Warnf("Failed to create PKI: (CA may be down, we will retry later) %s: %s", cfg.Addr, err)
		cfg.PKIClient = DummyPKI
		return
	}
	// We're only looking at the first PEM block
	der, _ := pem.Decode([]byte(caCert))
	if der == nil {
		log.Warnf("Failed to create PKI: (Unable to parse cert, we will retry later) %s", cfg.Addr)
		cfg.PKIClient = DummyPKI
		return
	}

	cert, err := x509.ParseCertificate(der.Bytes)
	if err != nil {
		log.Warnf("Failed to create PKI: (Unable to parse cert, we will retry later) %s: %s", cfg.Addr, err)
		cfg.PKIClient = DummyPKI
		return
	}

	// Now we can verify it's actually a CA
	if cert.BasicConstraintsValid && !cert.IsCA {
		log.Infof("Configured CA %s - %s is not a valid CA.  Will retry later.", cfg.Addr, cert.Subject.CommonName)
		cfg.PKIClient = DummyPKI
		return
	}
	cfg.PKIClient = client
}

func (s CAConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg CAConfiguration
	if jsonConfig == "" {
		return nil, fmt.Errorf("CA configuration must include a host address and can not be empty")
	}
	if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
		return nil, fmt.Errorf("Malformed CA configuration: %s", err)
	}

	// TODO - consider parsing the address and verifying it looks reasonable...

	(&cfg).InitPKI(s.m.swarmTLSConfig)

	return cfg, nil
}

func (s CAConfigSubsystem) UpdateConfig(cfgInt interface{}) error {
	cfg, ok := cfgInt.(CAConfiguration)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	*s.cfg = cfg
	return nil
}

func (m DefaultManager) ClusterSignCSR(csr *pki.CertificateSigningRequest) (*pki.CertificateResponse, error) {
	return m.DoSignCSR(clusterCAKvKeyPrefix, csr)
}

func (m DefaultManager) ClientSignCSR(csr *pki.CertificateSigningRequest) (*pki.CertificateResponse, error) {
	return m.DoSignCSR(clientCAKvKeyPrefix, csr)
}

func (m DefaultManager) DoSignCSR(prefix string, csr *pki.CertificateSigningRequest) (*pki.CertificateResponse, error) {

	available := []CAConfigSubsystem{}
	for key, subsystem := range m.configSubsystems {
		if strings.HasPrefix(key, prefix+"_") {
			ca := subsystem.(CAConfigSubsystem)
			available = append(available, ca)
		}
	}
	if len(available) == 0 {
		return nil, fmt.Errorf("No registered CAs detected")
	}

	// TODO - make this algorithm a little smarter
	// Prefer local
	// verify it's actually a root
	// Maybe keep track of which one worked last time to short-circuit
	var resp *pki.CertificateResponse
	var err error
	for _, ca := range available {
		resp, err = ca.cfg.PKIClient.SignCSR(csr)
		if err == nil {
			return resp, err
		}
		// TODO
		// Detect if it is down, and nil out the PKIClient
		// Detect if it's not a valid root, and nil out the PKIClient
	}

	// Try to revive a backup one
	for _, ca := range available {
		if ca.cfg.PKIClient != DummyPKI {
			continue
		}
		ca.cfg.InitPKI(m.swarmTLSConfig)
		resp, err = ca.cfg.PKIClient.SignCSR(csr)
		if err == nil {
			return resp, err
		}
		// TODO
		// Detect if it's not a valid root, and nil out the PKIClient
	}
	return nil, fmt.Errorf("None of the registered CAs are available.  Please try again later")
}

// Migration logic from pre 1.1 to 1.1 or later
func (m *DefaultManager) migrateLegacyCAConfig() error {
	kv := m.Datastore()

	// Check to see if we have at least one CA with the new style registration
	// TODO - drop these to debug once things are sorted out
	log.Info("Checking for new style CA config")
	configKeys, err := kv.List(KsConfigDir)
	for _, kvPair := range configKeys {
		subsystemKey := path.Base(kvPair.Key)
		if strings.HasPrefix(subsystemKey, clusterCAKvKeyPrefix+"_") {
			log.Info("At least one new-style CA detected, not migrating existing config")
			return nil
		}
	}

	log.Info("Migrating legacy CA configuration to current format")
	type PluginRegistration struct {
		//Name      string           `json:"Name"`
		Addr string `json:"Addr"`
		//TLSConfig *PluginTLSConfig `json:"TLSConfig,omitempty"`
	}

	// Look in the KV store for our CA configuration
	kvpair, err := kv.Get(path.Join(datastoreVersion, "plugins", "certificate_authority", "swarm"))
	if err == nil {
		var swarmPlugin PluginRegistration
		if err := json.Unmarshal(kvpair.Value, &swarmPlugin); err != nil {
			log.Debug("Malformed CA registration for swarm")
			return err
		}
		addr, err := url.Parse(swarmPlugin.Addr)
		if err != nil {
			return err
		}
		// Write out the config, and the event will wire it up
		kvKey := path.Join(KsConfigDir, fmt.Sprintf("%s_%s", clusterCAKvKeyPrefix, addr.Host))
		newCfgString := fmt.Sprintf(`{"addr":"%s"}`, swarmPlugin.Addr)
		if err := kv.Put(kvKey, []byte(newCfgString), nil); err != nil {
			err = utils.MaybeWrapEtcdClusterErr(err)
			log.Warnf("Unable to update %s config in kv store: %s", kvKey, err)
		}
	}
	kvpair, err = kv.Get(path.Join(datastoreVersion, "plugins", "certificate_authority", "orca"))
	if err == nil {
		var orcaPlugin PluginRegistration
		if err := json.Unmarshal(kvpair.Value, &orcaPlugin); err != nil {
			log.Debug("Malformed CA registration for orca")
			return err
		}
		addr, err := url.Parse(orcaPlugin.Addr)
		if err != nil {
			return err
		}
		kvKey := path.Join(KsConfigDir, fmt.Sprintf("%s_%s", clientCAKvKeyPrefix, addr.Host))
		newCfgString := fmt.Sprintf(`{"addr":"%s"}`, orcaPlugin.Addr)
		if err := kv.Put(kvKey, []byte(newCfgString), nil); err != nil {
			err = utils.MaybeWrapEtcdClusterErr(err)
			log.Warnf("Unable to update %s config in kv store: %s", kvKey, err)
		}
	}
	return nil
}

func (m *DefaultManager) getCABanners() []Banner {
	// Just check the cluster and assume user is consistent
	prefix := clusterCAKvKeyPrefix
	available := []CAConfigSubsystem{}
	for key, subsystem := range m.configSubsystems {
		if strings.HasPrefix(key, prefix+"_") {
			ca := subsystem.(CAConfigSubsystem)
			available = append(available, ca)
		}
	}
	placeholders := []string{}

	// NOTE: This algorithm doesn't actually query each CA to make sure they're alive right now.
	//       We just verify it was configured. The operating assumption is we'll pick up
	//       controller health if the whole controller goes away.  This means partial controller
	//       failure wont be immediately detected by this.

	// TODO in the future once we better parallelize banner generation, this can spend
	//      some more time probing the CA for health.
	goodCount := 0
	for _, ca := range available {
		// We assume that once a CA is good, it wont revert to being a dummy
		if ca.cfg.PKIClient != DummyPKI {
			goodCount++
			continue
		}

		// The dummy might have been updated since we last checked.
		// TODO - this is a bit expensive, we might want to have a background task to
		//        do this periodically instead of doing it within the banner logic
		ca.cfg.InitPKI(m.swarmTLSConfig)
		// if it's no longer a dummy, add it to the good list
		if ca.cfg.PKIClient != DummyPKI {
			goodCount++
		} else {
			caURL, err := url.Parse(ca.cfg.Addr)
			if err == nil {
				host, _, err := net.SplitHostPort(caURL.Host)
				if err == nil {
					placeholders = append(placeholders, host)
				} else {
					log.Debugf("Malformed CA registration: %s: %s", ca.cfg.Addr, err)
					placeholders = append(placeholders, ca.cfg.Addr) // Better than nothing...
				}
			} else {
				log.Debugf("Malformed CA registration: %s: %s", ca.cfg.Addr, err)
				placeholders = append(placeholders, ca.cfg.Addr) // Better than nothing...
			}
		}
	}
	switch goodCount {
	case 0:
		return []Banner{
			{
				Level:   BannerCRIT,
				Message: "No valid CAs detected.  Users will not be able to download certificate bundles, and new members can not join the cluster.  For more information visit https://success.docker.com/?cid=UCP3",
			},
		}
	case 1:
		if len(m.GetManagers()) == 1 { // TODO - redundant call while doing banners -- see getHABanners
			// For a single controller node system, don't nag
			return []Banner{}
		}
		return []Banner{
			{
				Level:   BannerWARN,
				Message: fmt.Sprintf("Only one controller has Root CA key material.  If that controller fails you will be unable to generate new user certificate bundles or add new nodes to the cluster.  Uninitialized placeholders detected on %s. For more information visit https://success.docker.com/?cid=UCP4", strings.Join(placeholders, ",")),
			},
		}
	}
	return []Banner{}
}
