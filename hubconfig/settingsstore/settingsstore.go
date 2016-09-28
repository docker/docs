package settingsstore

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/libtrust"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/util"
	licenseutil "github.com/docker/dhe-deploy/licensing/util"
	"github.com/docker/dhe-deploy/shared"
	"github.com/docker/dhe-deploy/shared/containers"
	distributionconfig "github.com/docker/distribution/configuration"
	garantconfig "github.com/docker/garant/config"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/samalba/dockerclient"
	"gopkg.in/yaml.v2"
)

func New(store hubconfig.KeyValueStore) hubconfig.SettingsStore {
	return &settingsStore{
		store: store,
	}
}

type settingsStore struct {
	rootDirectory string
	store         hubconfig.KeyValueStore
}

// TODO: maybe store every key separately instead of one big yaml?
func (s *settingsStore) SetUserHubConfig(userHubConfig *hubconfig.UserHubConfig) error {
	// before saving the config, make sure that the tls certs are in order
	err := util.HubConfigTLSDomainConsistent(userHubConfig)
	if err != nil {
		return err
	}

	if configBytes, err := yaml.Marshal(userHubConfig); err != nil {
		return err
	} else {
		if err := s.store.Put(deploy.UserHubConfigFilename, configBytes); err != nil {
			return err
		}
		// this also needs to generate the nginx config because the nginx config is based on userhubconfig
		if err := s.refreshNginxConfig(); err != nil {
			return err
		}
		// this also needs to generate the notary server config because the notary config needs to know
		// garant information
		return s.refreshNotaryGarantConfig(userHubConfig)
	}
}

func (s *settingsStore) refreshNotaryGarantConfig(userHubConfig *hubconfig.UserHubConfig) error {
	domainName := userHubConfig.DTRHost
	if domainName == "" {
		return nil
	}

	serverConfig := &defaultconfigs.DefaultNotaryServerConfig
	serverConfig.Auth.Options = util.GetRegistryAuthConfig(domainName)

	if err := s.setNotaryServerConfig(serverConfig); err != nil {
		return err
	}

	signerConfig := &defaultconfigs.DefaultNotarySignerConfig
	return s.setNotarySignerConfig(signerConfig)
}

func (s *settingsStore) UserHubConfig() (*hubconfig.UserHubConfig, error) {
	var userHubConfig hubconfig.UserHubConfig
	if configBytes, err := s.store.Get(deploy.UserHubConfigFilename); err != nil || len(configBytes) == 0 {
		return nil, err
	} else if err := yaml.Unmarshal(configBytes, &userHubConfig); err != nil {
		return nil, err
	} else {
		return &userHubConfig, nil
	}
}

func (s *settingsStore) AuthConfig() (*garantconfig.Configuration, error) {
	var authConfig garantconfig.Configuration
	if configBytes, err := s.store.Get(deploy.AuthConfigFilename); err != nil || len(configBytes) == 0 {
		return nil, err
	} else if err := yaml.Unmarshal(configBytes, &authConfig); err != nil {
		return nil, err
	} else {
		return &authConfig, nil
	}
}

func (s *settingsStore) SetAuthConfig(authConfig *garantconfig.Configuration) error {
	if configBytes, err := yaml.Marshal(authConfig); err != nil {
		return err
	} else {
		return s.store.Put(deploy.AuthConfigFilename, configBytes)
	}
}

func (s *settingsStore) LicenseConfig() (*hubconfig.LicenseConfig, error) {
	var licensingConfig hubconfig.LicenseConfig
	if configBytes, err := s.store.Get(deploy.LicenseConfigFilename); err != nil || len(configBytes) == 0 {
		return nil, err
	} else if err := json.Unmarshal(configBytes, &licensingConfig); err != nil {
		return nil, err
	} else {
		return &licensingConfig, nil
	}
}

func (s *settingsStore) SetLicenseStatus(isValid bool) error {
	log.WithField("isValid", isValid).Info("refreshing nginx config because of license change")
	err := s.refreshNginxConfig()
	if err != nil {
		log.WithField("error", err).Warn("failed to update nginx config while changing license status")
	}
	return err
}

func (s *settingsStore) SetLicenseConfig(licensingConfig *hubconfig.LicenseConfig) error {
	if configBytes, err := json.Marshal(licensingConfig); err != nil {
		return err
	} else {
		if err := s.store.Put(deploy.LicenseConfigFilename, configBytes); err != nil {
			return err
		}
		return s.refreshNginxConfig()
	}
}

func (s *settingsStore) HAConfig() (*hubconfig.HAConfig, error) {
	var haConfig hubconfig.HAConfig
	if configBytes, err := s.store.Get(deploy.HAConfigFilename); err != nil || len(configBytes) == 0 {
		return nil, err
	} else if err := json.Unmarshal(configBytes, &haConfig); err != nil {
		return nil, err
	} else {
		return &haConfig, nil
	}
}

func (s *settingsStore) SetHAConfig(haConfig *hubconfig.HAConfig) error {
	if configBytes, err := json.Marshal(haConfig); err != nil {
		return err
	} else {
		if err := s.store.Put(deploy.HAConfigFilename, configBytes); err != nil {
			return err
		}
		if err := s.refreshNginxConfig(); err != nil {
			return err
		}
		return nil
	}
}

func (s *settingsStore) writeTLSCertificate(certificate *tls.Certificate) error {
	cert, err := util.CertificateToPEM(certificate)
	if err != nil {
		return err
	}
	return s.store.Put(deploy.TLSPEMFilename, cert)
}

func (s *settingsStore) RegistryConfig() (*distributionconfig.Configuration, error) {
	var registryConfig distributionconfig.Configuration
	if configBytes, err := s.store.Get(deploy.RegistryConfigFilename); err != nil || len(configBytes) == 0 {
		return nil, err
	} else if err := yaml.Unmarshal(configBytes, &registryConfig); err != nil {
		return nil, err
	} else {
		// fix goyaml issue where serialization/deserialization converts a map[string]interface{} to a map[interface{}]interface{}
		for _, props := range registryConfig.Storage {
			for key, val := range props {
				if imap, ok := val.(map[interface{}]interface{}); ok {
					smap := make(map[string]interface{})
					for k, v := range imap {
						// only keep keys which are strings, others are invalid
						if kstr, ok := k.(string); ok {
							smap[kstr] = v
						}
					}
					props[key] = smap
				}
			}
		}
		if registryConfig.Storage.Type() == "filesystem" {
			params := registryConfig.Storage["filesystem"]
			params["rootdirectory"] = fmt.Sprintf("/var/lib/docker/volumes/%s/_data", containers.RegistryVolume.ReplicaName(os.Getenv(deploy.ReplicaIDEnvVar)))
			registryConfig.Storage["filesystem"] = params
		}

		return &registryConfig, nil
	}
}

func (s *settingsStore) SetRegistryConfig(registryConfig *distributionconfig.Configuration) error {
	if configBytes, err := yaml.Marshal(registryConfig); err != nil {
		return err
	} else {
		return s.store.Put(deploy.RegistryConfigFilename, configBytes)
	}
}

func (s *settingsStore) setNotaryServerConfig(serverConfig *hubconfig.NotaryServerConfig) error {
	if configBytes, err := json.Marshal(serverConfig); err != nil {
		return err
	} else {
		return s.store.Put(deploy.NotaryServerConfigFilename, configBytes)
	}
}

func (s *settingsStore) setNotarySignerConfig(signerConfig *hubconfig.NotarySignerConfig) error {
	if configBytes, err := json.Marshal(signerConfig); err != nil {
		return err
	} else {
		return s.store.Put(deploy.NotarySignerConfigFilename, configBytes)
	}
}

func (s *settingsStore) HubCredentials() (*dockerclient.AuthConfig, error) {
	var dockerCfg map[string]interface{}
	if hubCredentialBytes, err := s.store.Get(deploy.HubCredentialsFilename); err != nil || len(hubCredentialBytes) == 0 {
		return nil, err
	} else if err := json.Unmarshal(hubCredentialBytes, &dockerCfg); err != nil {
		return nil, err
	} else if authEntry, ok := (dockerCfg[deploy.DockerIndexURL]).(map[string]interface{}); !ok {
		return nil, fmt.Errorf("cannot read json object")
	} else {
		var authBytes []byte
		var email string
		if authString, ok := authEntry["auth"].(string); !ok {
			return nil, fmt.Errorf("Failed to read auth from second level of json object")
		} else {
			authBytes = []byte(authString)
		}
		if email, ok = authEntry["email"].(string); !ok {
			return nil, fmt.Errorf("Failed to read email from second level of json object")
		}
		decLen := base64.StdEncoding.DecodedLen(len(authBytes))
		decoded := make([]byte, decLen)
		n, err := base64.StdEncoding.Decode(decoded, authBytes)
		if err != nil {
			return nil, err
		}
		if n > decLen {
			return nil, fmt.Errorf("Something went wrong decoding auth config")
		}
		arr := strings.SplitN(string(decoded), ":", 2)
		if len(arr) != 2 {
			return nil, fmt.Errorf("Invalid auth configuration file")
		}
		username := arr[0]
		password := strings.Trim(arr[1], "\x00")
		return &dockerclient.AuthConfig{
			Username: username,
			Password: password,
			Email:    email,
		}, nil
	}
}

func (s *settingsStore) GarantRootCert() (string, error) {
	dest := path.Join(deploy.GeneratedConfigsDir, deploy.GarantRootCertFilename)
	data, err := s.store.Get(dest)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *settingsStore) SetGarantRootCert(cert *x509.Certificate) error {
	encoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	return s.store.Put(path.Join(deploy.GeneratedConfigsDir, deploy.GarantRootCertFilename), encoded)
}

func (s *settingsStore) GarantSigningKey() (libtrust.PrivateKey, error) {
	dest := path.Join(deploy.GeneratedConfigsDir, deploy.GarantSigningKeyFilename)
	data, err := s.store.Get(dest)
	if err != nil {
		return nil, err
	}
	return libtrust.UnmarshalPrivateKeyJWK(data)
}

func (s *settingsStore) SetGarantSigningKey(key libtrust.PrivateKey) error {
	dest := path.Join(deploy.GeneratedConfigsDir, deploy.GarantSigningKeyFilename)
	if key == nil {
		return s.store.Put(dest, []byte{})
	}
	jsonKey, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("unable to encode root key JWK: %s", err)
	}
	return s.store.Put(dest, jsonKey)
}

func (s settingsStore) RawEnziSigningKey() (string, error) {
	dest := path.Join(deploy.GeneratedConfigsDir, deploy.EnziSigningKeyFilename)
	data, err := s.store.Get(dest)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s settingsStore) EnziSigningKey() (libtrust.PrivateKey, error) {
	dest := path.Join(deploy.GeneratedConfigsDir, deploy.EnziSigningKeyFilename)
	data, err := s.store.Get(dest)
	if err != nil {
		return nil, err
	}
	return libtrust.UnmarshalPrivateKeyJWK(data)
}

func (s settingsStore) SetEnziSigningKey(key libtrust.PrivateKey) error {
	dest := path.Join(deploy.GeneratedConfigsDir, deploy.EnziSigningKeyFilename)
	if key == nil {
		return s.store.Put(dest, []byte{})
	}
	jsonKey, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("unable to encode root key JWK: %s", err)
	}
	return s.store.Put(dest, jsonKey)
}

func (s *settingsStore) EnziService() (*enziresponses.Service, error) {
	var enziService enziresponses.Service
	if configBytes, err := s.store.Get(deploy.EnziServiceFilename); err != nil || len(configBytes) == 0 {
		return nil, err
	} else if err := json.Unmarshal(configBytes, &enziService); err != nil {
		return nil, err
	} else {
		return &enziService, nil
	}
}

func (s *settingsStore) SetEnziService(service *enziresponses.Service) error {
	if configBytes, err := json.Marshal(service); err != nil {
		return err
	} else {
		if err := s.store.Put(deploy.EnziServiceFilename, configBytes); err != nil {
			return err
		}
		return nil
	}
}

func (s *settingsStore) SetHubCredentials(hubCredentials *dockerclient.AuthConfig) error {
	authStr := hubCredentials.Username + ":" + hubCredentials.Password
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	dockerCfg := map[string]interface{}{
		deploy.DockerIndexURL: map[string]string{
			"auth":  string(encoded),
			"email": hubCredentials.Email,
		},
	}
	if dockerCfgBytes, err := json.Marshal(dockerCfg); err != nil {
		return err
	} else {
		return s.store.Put(deploy.HubCredentialsFilename, dockerCfgBytes)
	}
}

func (s *settingsStore) Ping() error {
	return s.store.Ping()
}

func (s *settingsStore) refreshNginxConfig() error {
	// fetch the things we need in order to create the nginx config
	userHubConfig, err := s.UserHubConfig()
	if err != nil {
		return fmt.Errorf("Failed to get user hub config: %s", err)
	}
	haConfig, err := s.HAConfig()
	if err != nil {
		return fmt.Errorf("Failed to get ha config: %s", err)
	}
	// don't set the nginx config yet if all the data is not available yet
	if userHubConfig == nil || haConfig == nil {
		return nil
	}

	parts := strings.Split(userHubConfig.DTRHost, ":")
	colonHTTPSPort := ""
	if len(parts) > 1 {
		colonHTTPSPort = fmt.Sprintf(":%s", parts[1])
	}
	ucpVerifyCert := haConfig.UCPVerifyCert

	nginxConfigTemplate := template.Must(template.New("nginx.conf").Parse(shared.NginxTemplate))
	configBuffer := new(bytes.Buffer)

	authBypassOU := ""
	if userHubConfig.AuthBypassCA != "" {
		if userHubConfig.AuthBypassOU != "" {
			authBypassOU = regexp.QuoteMeta(userHubConfig.AuthBypassOU)
		} else {
			authBypassOU = deploy.DefaultAuthBypassOU
		}
	}

	licenseConfig, err := s.LicenseConfig()
	if err != nil {
		return err
	}

	licenseIsValid, err := licenseutil.IsValidFromLicenseConfig(licenseConfig)
	if err != nil {
		return err
	}

	// the replica id is found and replaced by confd with the one from its local env var
	replicaID := "XXXREPLICA_IDXXX"

	templateArgs := struct {
		HasLicense               bool
		StorageServerName        string
		StorageServerPort        uint16
		AdminServerFullName      string
		RethinkFullName          string
		RethinkPort              uint16
		EtcdFullName             string
		EtcdPort1                uint16
		EtcdPort2                uint16
		EnziHost                 string
		AdminSubroute            string
		AdminPort                uint16
		GarantSubroute           string
		GarantPort               uint16
		NotaryHost               string
		NotaryPort               uint16
		NotaryCACert             string
		NotaryClientCert         string
		NotaryClientKey          string
		UCPVerifyCert            bool
		EnziVerifyCert           bool
		AuthBypassOU             string
		RegistryEventsHeaderName string
		ColonHTTPSPort           string
	}{
		HasLicense:               licenseIsValid,
		StorageServerName:        containers.Registry.BridgeName(replicaID),
		StorageServerPort:        deploy.StorageContainerPort,
		AdminServerFullName:      containers.APIServer.BridgeName(replicaID),
		RethinkFullName:          containers.Rethinkdb.BridgeName(replicaID),
		RethinkPort:              deploy.RethinkdbPort,
		EtcdFullName:             containers.Etcd.BridgeName(replicaID),
		EtcdPort1:                containers.EtcdClientPort1,
		EtcdPort2:                containers.EtcdClientPort2,
		AdminSubroute:            deploy.AdminSubroute,
		AdminPort:                deploy.AdminPort,
		GarantSubroute:           deploy.GarantSubroute,
		GarantPort:               deploy.GarantPort,
		NotaryHost:               containers.NotaryServer.BridgeName(replicaID),
		NotaryPort:               deploy.NotaryServerHTTPPort,
		NotaryCACert:             containers.NotaryCACertStore.CertPath(),
		NotaryClientCert:         containers.NotaryCertStore.CertPath(),
		NotaryClientKey:          containers.NotaryCertStore.KeyPath(),
		UCPVerifyCert:            ucpVerifyCert,
		AuthBypassOU:             authBypassOU,
		RegistryEventsHeaderName: deploy.RegistryEventsHeaderName,
		ColonHTTPSPort:           colonHTTPSPort,
	}

	err = nginxConfigTemplate.Execute(configBuffer, templateArgs)
	if err != nil {
		return fmt.Errorf("Failed to execute template: %s", err)
	}

	// create nginx config
	err = s.store.Put(deploy.NginxConfigFilename, configBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to save nginx config: %s", err)
	}

	// We store the auth bypass CA separately so nginx can read it
	// Note that we need to set it even if it's empty or confd in the nginx container will refuse to ever return success
	var authBypassCA []byte
	if userHubConfig != nil {
		authBypassCA = []byte(userHubConfig.AuthBypassCA)
	}
	err = s.store.Put(deploy.AuthBypassCAFilename, authBypassCA)
	if err != nil {
		return fmt.Errorf("Failed to write auth bypass CA cert: %s", err)
	}

	// now we write out all the certs from the hub config to their own etcd keys, yay :/

	// tls cert
	if userHubConfig != nil {
		// TODO: don't both parsing the cert at this point because it should be validated already
		cert, err := tls.X509KeyPair([]byte(userHubConfig.WebTLSCert), []byte(userHubConfig.WebTLSKey))
		if err != nil {
			return err
		}
		err = s.writeTLSCertificate(&cert)
		if err != nil {
			return fmt.Errorf("Failed to write public cert: %s", err)
		}
	}

	// we can't assume haConfig exists because we have a chicken and egg problem otherwise
	if haConfig != nil {
		err = s.store.Put(deploy.UCPCAPEMFilename, []byte(haConfig.UCPCA))
		if err != nil {
			return fmt.Errorf("Failed to write UCP CA cert: %s", err)
		}
	}

	return nil
}
