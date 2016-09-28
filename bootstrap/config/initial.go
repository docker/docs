package config

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/tlsconfig"
	"github.com/docker/libkv"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/orca/auth"
	enziauth "github.com/docker/orca/auth/enzi"
	orcaconfig "github.com/docker/orca/config"
	enziclient "github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/client/openid"
	enziforms "github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/jose"
	"github.com/docker/orca/types"
	"github.com/docker/orca/utils"
)

// WARNING: If we change the way Orca stores accounts, this might drift, so ultimately this should
//          probably be refactored to a common utility module that both the bootstrapper and orca can use
const OrcaPrefix = "orca/v1"

type Account struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Admin     bool   `json:"admin,omitempty"`
}

type TrackingConfiguration struct {
	DisableUsageInfo bool `json:"disable_usageinfo,omitempty"`
	DisableTracking  bool `json:"disable_tracking,omitempty"`
}

func init() {
	consul.Register()
	etcd.Register()
}

func GetKV(kvStoreURL *url.URL) (kvstore.Store, error) {
	// Connect to KV store and write out configuration settings
	kvType := strings.ToLower(kvStoreURL.Scheme)
	kvHost := kvStoreURL.Host
	var backend kvstore.Backend

	switch kvType {
	case "consul":
		backend = kvstore.CONSUL
	case "etcd":
		backend = kvstore.ETCD
	}

	tlsConfig, err := tlsconfig.Client(tlsconfig.Options{
		CAFile:   filepath.Join(SwarmKvCertVolumeMount, CAFilename),
		CertFile: filepath.Join(SwarmKvCertVolumeMount, CertFilename),
		KeyFile:  filepath.Join(SwarmKvCertVolumeMount, KeyFilename),
	})
	if err != nil {
		return nil, err
	}
	options := &kvstore.Config{
		ConnectionTimeout: time.Second * 10,
		TLS:               tlsConfig,
	}
	// TODO - We might need an extra "wait" in here to give the KV store time to come up...
	log.Debug("creating new KV object")
	return libkv.NewStore(
		backend,
		[]string{kvHost},
		options,
	)
}

func GetControllers(kvStoreURL *url.URL) ([]types.Controller, error) {
	controllers := []types.Controller{}
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return controllers, fmt.Errorf("Could not connect to KV store %s, maybe your host and UCP certificates are out of sync: %s", kvStoreURL, err)
	}

	path := path.Join(OrcaPrefix, "controllers")
	kvList, err := kv.List(path)
	if err != nil {
		return controllers, err
	}

	for _, kvPair := range kvList {
		var controller types.Controller
		if err := json.Unmarshal(kvPair.Value, &controller); err != nil {
			return []types.Controller{}, fmt.Errorf("Controller %s KV data is malformed: %s", kvPair.Key, err)
		}
		controllers = append(controllers, controller)
	}

	return controllers, nil
}

// GetController retrieves the controller registration info for the controller
// with the given label.
// HINT: the label is the host address (without port) used when the controller
// was installed or joined.
func GetController(kvStoreURL *url.URL, label string) (*types.Controller, error) {
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		log.Errorf("Could not connect to KV store %s: maybe your host and UCP certificates are out of sync?", kvStoreURL)
		return nil, err
	}

	key := path.Join(OrcaPrefix, "controllers", label)
	kvPair, err := kv.Get(key)
	if err != nil {
		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	var controller types.Controller
	if err := json.Unmarshal(kvPair.Value, &controller); err != nil {
		return nil, fmt.Errorf("Controller %s KV data is malformed: %s", key, err)
	}

	return &controller, nil
}

// TODO - This will need some refactoring once we figure out how we want to support multi-swarm
//        at present, it assumes each HA controller has a swarm manager, but that
//        assumption gets messy once we have multi-swarm

func AddSwarmManager(kvStoreURL *url.URL, controllerEndpoint, swarmEndpoint, proxyEndpoint string) error {
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return err
	}
	c := types.Controller{
		Label:               strings.Split(controllerEndpoint, ":")[0],
		Controller:          controllerEndpoint,
		SwarmClassicManager: swarmEndpoint,
		EngineProxy:         proxyEndpoint,
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := kv.Put(path.Join(OrcaPrefix, "controllers", c.Label), data, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}
	if err := kv.Put(path.Join(OrcaPrefix, "swarm", "0", "managers", swarmEndpoint),
		[]byte(fmt.Sprintf("tcp://%s", swarmEndpoint)), nil); err != nil {

		return utils.MaybeWrapEtcdClusterErr(err)
	}

	// Add pointers to our local CA
	if err := setKVCAConfig(kv); err != nil {
		return fmt.Errorf("unable to set CA config in KV store: %s", err)
	}

	return nil
}

// SetKVCAConfig adds the addresses of the cluster peer and client CA servers
// on this replica node into the KV store using the given kvStoreURL.
// OrcaHostAddress MUST be set.
func SetKVCAConfig(kvStoreURL *url.URL) error {
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return err
	}

	return setKVCAConfig(kv)
}

// setKVCAConfig is the unexported version of SetKVCAConfig which is used by
// several other methods in this package which have already resolved a kv store
// handle from the URL of the kv store.
func setKVCAConfig(kv kvstore.Store) error {
	kvKey := path.Join(OrcaPrefix, "config", fmt.Sprintf("clusterca_%s:%d", OrcaHostAddress, orcaconfig.SwarmCAPort))
	newCfgString := fmt.Sprintf(`{"addr":"https://%s:%d"}`, OrcaHostAddress, orcaconfig.SwarmCAPort)
	if err := kv.Put(kvKey, []byte(newCfgString), nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}
	kvKey = path.Join(OrcaPrefix, "config", fmt.Sprintf("clientca_%s:%d", OrcaHostAddress, orcaconfig.OrcaCAPort))
	newCfgString = fmt.Sprintf(`{"addr":"https://%s:%d"}`, OrcaHostAddress, orcaconfig.OrcaCAPort)
	if err := kv.Put(kvKey, []byte(newCfgString), nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	return nil
}

func RemoveSwarmManager(kvStoreURL *url.URL, swarmEndpointSubstring string) error {
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return err
	}

	controllerDataPaths := []string{
		path.Join("controllers", swarmEndpointSubstring),
		path.Join("config", fmt.Sprintf("clusterca_%s:%d", swarmEndpointSubstring, orcaconfig.SwarmCAPort)),
		path.Join("config", fmt.Sprintf("clientca_%s:%d", swarmEndpointSubstring, orcaconfig.OrcaCAPort)),
	}

	for _, dataPath := range controllerDataPaths {
		fullPath := path.Join(OrcaPrefix, dataPath)
		if err = kv.Delete(fullPath); err != nil {
			// TODO - might want to do a more exhaustive search...
			log.Warnf("Failed to delete KV data %q: %s", fullPath, err)
		}
	}
	kvList, err := kv.List(path.Join(OrcaPrefix, "swarm", "0", "managers"))
	if err != nil {
		return err
	}
	for _, kvpair := range kvList {
		if strings.Contains(kvpair.Key, swarmEndpointSubstring+":") || strings.HasSuffix(kvpair.Key, swarmEndpointSubstring) {
			log.Debugf("Detected swarm manager %s matching %s", kvpair.Key, swarmEndpointSubstring)
			if err := kv.Delete(kvpair.Key); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("Unable to located swarm manager endpoint registration in KV store")
}

func BootstrapOrcaConfig(kvStoreURL *url.URL, adminUsername, adminPassword, controllerEndpoint, swarmEndpoint, proxyEndpoint string, disableTracking bool, disableUsage bool) error {
	log.Debug("Writing out initial configuration to KV store")
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return err
	}

	// Create Default Org and Register UCP service.
	if err := ConfigureAuth(kv, adminUsername, adminPassword, controllerEndpoint); err != nil {
		return fmt.Errorf("unable to configure auth: %s", err)
	}

	// Orca ID information, PEM encoded - hidden using leading "_"
	if err := kv.Put(path.Join(OrcaPrefix, "_trust", "key"), []byte(OrcaInstanceKey), nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	// Controller Identity
	c := types.Controller{
		Label:               strings.Split(controllerEndpoint, ":")[0],
		Controller:          controllerEndpoint,
		SwarmClassicManager: swarmEndpoint,
		EngineProxy:         proxyEndpoint,
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := kv.Put(path.Join(OrcaPrefix, "controllers", c.Label), data, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	// Swarm configuration
	// Silly shim that should get cleaned up
	if err := kv.Put(path.Join(OrcaPrefix, "swarm", "0", "managers", swarmEndpoint),
		[]byte(fmt.Sprintf("tcp://%s", swarmEndpoint)), nil); err != nil {

		return utils.MaybeWrapEtcdClusterErr(err)
	}
	//  TODO - Might want to consider actually placing the cert bits in the KV instead of on the FS...
	if err := kv.Put(path.Join(OrcaPrefix, "swarm", "0", "ca-cert-path"),
		[]byte(path.Join(CertDir, "swarm", CAFilename)), nil); err != nil {

		return utils.MaybeWrapEtcdClusterErr(err)
	}
	if err := kv.Put(path.Join(OrcaPrefix, "swarm", "0", "cert-path"),
		[]byte(path.Join(CertDir, "swarm", CertFilename)), nil); err != nil {

		return utils.MaybeWrapEtcdClusterErr(err)
	}
	if err := kv.Put(path.Join(OrcaPrefix, "swarm", "0", "key-path"),
		[]byte(path.Join(CertDir, "swarm", KeyFilename)), nil); err != nil {

		return utils.MaybeWrapEtcdClusterErr(err)
	}

	// Write out the local CA information
	if err := setKVCAConfig(kv); err != nil {
		return fmt.Errorf("unable to set CA config in KV store: %s", err)
	}

	// Write out KV config parameters
	kvCfg := KVDeployCfg{
		Timeout: KVTimeout,
	}
	data, err = json.Marshal(kvCfg)
	if err := kv.Put(path.Join(OrcaPrefix, "config", "kv"), data, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	trackingCfg := TrackingConfiguration{
		DisableUsageInfo: disableUsage,
		DisableTracking:  disableTracking,
	}
	data, err = json.Marshal(trackingCfg)
	if err := kv.Put(path.Join(OrcaPrefix, "config", "tracking"), data, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	// Load a user supplied license if detected
	licenseData, err := ioutil.ReadFile(LicenseFile)
	if err == nil {
		// The file may contain a leading BOM, which will choke the json deserializer
		licenseData := bytes.Trim(licenseData, "\xef\xbb\xbf")

		// Definitions borrowed from orca/controller/manager/license.go
		type LicenseConfig struct {
			KeyID         string `yaml:"key_id" json:"key_id"`
			PrivateKey    string `yaml:"private_key" json:"private_key,omitempty"`
			Authorization string `yaml:"authorization" json:"authorization,omitempty"`
		}
		type LicenseSubsystemConfig struct {
			AutoRefresh bool          `json:"auto_refresh"`
			License     LicenseConfig `json:"license_config"`
			// Details and LastUpdateError omitted
		}

		var license LicenseConfig
		err := json.Unmarshal(licenseData, &license)
		if err != nil {
			log.Warnf("User supplied license appears corrupt - skipping injection. %s", err)
			log.Debug(string(licenseData))
		} else {
			cfg := LicenseSubsystemConfig{AutoRefresh: true, License: license}
			data, err = json.Marshal(cfg)
			if err := kv.Put(path.Join(OrcaPrefix, "config", "license"), data, nil); err != nil {
				return utils.MaybeWrapEtcdClusterErr(err)
			}
			log.Debug("User supplied license injected")
		}
	}

	// TODO - DTR configuration (if specified)

	return nil
}

// FIXME: We should share this kv store key with the manager package.
var authConfigKVStoreKey = path.Join(OrcaPrefix, "config", "auth2")
var ErrNoAuthConfig = errors.New("no auth config")

func GetAuthConfig(kvStoreURL *url.URL) (*auth.AuthenticatorConfiguration, error) {
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return nil, err
	}

	kvPair, err := kv.Get(authConfigKVStoreKey)
	if err != nil {
		if err == kvstore.ErrKeyNotFound {
			return nil, ErrNoAuthConfig
		}
		err = utils.MaybeWrapEtcdClusterErr(err)
		return nil, fmt.Errorf("unable to get auth config from KV store: %s", err)
	}

	var authConfig auth.AuthenticatorConfiguration
	if err := json.Unmarshal(kvPair.Value, &authConfig); err != nil {
		return nil, fmt.Errorf("unable to decode auth config from JSON: %s", err)
	}

	return &authConfig, nil
}

// ConfigureAuth handles configuring UCP to use the eNZi auth service:
//     1. Create a default "docker-datacenter" organization account and return
//        its account ID.
//     2. Register a privileged service to represent UCP and return its
//        unique service ID.
func ConfigureAuth(kv kvstore.Store, adminUsername, adminPassword, ucpControllerAddr string) error {
	authAPIServerAddr := net.JoinHostPort(OrcaHostAddress, fmt.Sprintf("%d", orcaconfig.AuthAPIPort))

	enziSession, err := getEnziUserSession(authAPIServerAddr, adminUsername, adminPassword)
	if err != nil {
		return err
	}

	// Create a default "Docker Datacenter" organization account which will
	// be used to organize all of the teams used with UCP. This org account
	// will "own" the cluster resources.
	createOrgForm := enziforms.CreateAccount{
		Name:     "docker-datacenter",
		FullName: "Docker Datacenter",
		IsOrg:    true,
		IsActive: true,
	}

	dockerDatacenterOrg, err := enziSession.CreateAccount(createOrgForm)
	if err != nil {
		log.Errorln("unable to create docker-datacenter organization")
		return err
	}

	// Next, register a privileged service to represent UCP, owned by
	// the docker-datacenter organization. UCP will authenticate to eNZi
	// with JWTs that include a cert that chains up to the swarm cluster
	// root CA. This CA cert is implicitly trusted by eNZi for all
	// privileged services (which UCP is an instance of).

	createServiceForm := enziforms.CreateService{
		Name:        "Docker Universal Control Plane",
		Description: "Docker Datacenter Container Orchestration",
		URL:         fmt.Sprintf("https://%s/", ucpControllerAddr),
		Privileged:  true, // VERY IMPORTANT.
		RedirectURIs: []string{
			// NOTE: It doesn't matter what we put here for now as
			// UCP wont be using the authorization code flow. It
			// might in the future.
			fmt.Sprintf("https://%s/openid_callback", ucpControllerAddr),
		},
		JWKsURIs: []string{
			// NOTE: This is the public API endpoint for listing
			// UCP's openid authentication keys. We currently opt
			// to use certificate chains in our JWTs instead but
			// we keep this here as a fallback that we could use
			// in the future.
			fmt.Sprintf("https://%s/openid_keys", ucpControllerAddr),
		},
		// The controller's authentication tokens must have an issuer
		// field that exactly equals one of these values.
		ProviderIdentities: []string{
			authAPIServerAddr,
		},
	}

	ucpService, err := enziSession.CreateService(dockerDatacenterOrg.Name, createServiceForm)
	if err != nil {
		log.Errorln("unable to register UCP with auth service")
		return err
	}

	if ucpService.ID == "" {
		return fmt.Errorf("service registration response did not have a service ID")
	}
	if ucpService.OwnerID == "" {
		return fmt.Errorf("service registration response did not have an owner ID")
	}

	authConfig := auth.AuthenticatorConfiguration{
		AuthenticatorType: auth.AuthenticatorEnzi,
		EnziConfig: auth.EnziConfig{
			ServiceID:     ucpService.ID,
			DefaultOrgID:  ucpService.OwnerID,
			ProviderAddrs: []string{authAPIServerAddr},
		},
	}

	authConfigJSON, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("unable to encode auth config to JSON: %s", err)
	}

	if err := kv.Put(authConfigKVStoreKey, authConfigJSON, nil); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return fmt.Errorf("unable to save auth config to KV Store: %s", err)
	}

	return nil
}

// UpdateUCPServiceAuthConfig handles adding an additional endpoint to the UCP
// service registered in eNZi.
func UpdateUCPServiceAuthConfig(kvStoreURL *url.URL) error {
	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return err
	}

	// TODO: Acquire a distributed lock (using the kv store) around this
	// code.

	// Get the service which represents UCP. Our Default Org ID and Service
	// ID are in the kv store.
	kvPair, err := kv.Get(authConfigKVStoreKey)
	if err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return fmt.Errorf("unable to get auth config from KV store: %s", err)
	}

	authConfigJSON := kvPair.Value

	var authConfig auth.AuthenticatorConfiguration
	if err := json.Unmarshal(authConfigJSON, &authConfig); err != nil {
		return fmt.Errorf("unable to decode auth config from JSON: %s", err)
	}

	if authConfig.EnziConfig.DefaultOrgID == "" {
		return fmt.Errorf("auth config is missing UCP default org ID")
	}
	if authConfig.EnziConfig.ServiceID == "" {
		return fmt.Errorf("auth config is missing UCP service ID")
	}
	if len(authConfig.EnziConfig.ProviderAddrs) == 0 {
		return fmt.Errorf("no existing auth provider addresses known")
	}

	// Connect to the initial auth provider addr to update the config.
	existingAuthAPIServerAddr := authConfig.EnziConfig.ProviderAddrs[0]
	enziSession, err := getUCPServiceOwnerEnziSession(existingAuthAPIServerAddr, authConfig.EnziConfig.DefaultOrgID, authConfig.EnziConfig.ServiceID)
	if err != nil {
		return err
	}

	ucpService, err := enziSession.GetService("id:"+authConfig.EnziConfig.DefaultOrgID, "id:"+authConfig.EnziConfig.ServiceID)
	if err != nil {
		log.Errorln("unable to get UCP service config from Auth Service")
		return err
	}

	newAuthAPIServerAddr := net.JoinHostPort(OrcaHostAddress, fmt.Sprintf("%d", orcaconfig.AuthAPIPort))
	providerAddrs := dedupeStrings(append(ucpService.ProviderIdentities, newAuthAPIServerAddr)...)

	// NOTE: We don't bother updating the JWKs URIs or Redirect URIs
	// because UCP does not use them.
	updateForm := enziforms.UpdateService{
		ProviderIdentities: &providerAddrs,
	}

	if _, err := enziSession.UpdateService("id:"+authConfig.EnziConfig.DefaultOrgID, "id:"+authConfig.EnziConfig.ServiceID, updateForm); err != nil {
		log.Errorln("unable to update UCP service config in Auth Service")
		return err
	}

	authConfig.EnziConfig.ProviderAddrs = providerAddrs

	authConfigJSON, err = json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("unable to encode auth config to JSON: %s", err)
	}

	if err := kv.Put(authConfigKVStoreKey, authConfigJSON, nil); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return fmt.Errorf("unable to save auth config to KV Store: %s", err)
	}

	return nil
}

func dedupeStrings(inVals ...string) []string {
	deduped := make(map[string]struct{}, len(inVals))
	for _, val := range inVals {
		deduped[val] = struct{}{}
	}

	outVals := make([]string, 0, len(deduped))
	for val := range deduped {
		outVals = append(outVals, val)
	}

	return outVals
}

func RemoveEnziProviderConfig(kvStoreURL *url.URL) error {
	// Infer what our local IP is from the local kvStore URL.
	localIP, _, err := net.SplitHostPort(kvStoreURL.Host)
	if err != nil {
		log.Errorf("unable to split host/port from kv store address: %s", kvStoreURL.Host)
		return err
	}

	providerAddr := net.JoinHostPort(localIP, fmt.Sprintf("%d", orcaconfig.AuthAPIPort))

	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return err
	}

	// TODO: Acquire a distributed lock (from the kv store) around this
	// code.

	kvPair, err := kv.Get(authConfigKVStoreKey)
	if err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return fmt.Errorf("unable to get auth config from KV store: %s", err)
	}

	authConfigJSON := kvPair.Value

	var authConfig auth.AuthenticatorConfiguration
	if err := json.Unmarshal(authConfigJSON, &authConfig); err != nil {
		return fmt.Errorf("unable to decode auth config from JSON: %s", err)
	}

	remainingProviderAddrs := make([]string, 0, len(authConfig.EnziConfig.ProviderAddrs))
	for _, existingAddr := range authConfig.EnziConfig.ProviderAddrs {
		if existingAddr != providerAddr {
			remainingProviderAddrs = append(remainingProviderAddrs, existingAddr)
		}
	}

	if len(remainingProviderAddrs) == len(authConfig.EnziConfig.ProviderAddrs) {
		// This indicates a configuration issue. It's okay since we are
		// removing it, but we should log a warning message.
		log.Warnf("no auth provider addrs were removed from current auth config: %s not found in %s", providerAddr, authConfig.EnziConfig.ProviderAddrs)
	}

	authConfig.EnziConfig.ProviderAddrs = remainingProviderAddrs

	authConfigJSON, err = json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("unable to encode auth config to JSON: %s", err)
	}

	// Put the updated config back in the KV Store.
	if err := kv.Put(authConfigKVStoreKey, authConfigJSON, nil); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return fmt.Errorf("unable to save auth config to KV Store: %s", err)
	}

	return nil
}

func getEnziUserSession(authAPIServerAddr, username, password string) (*enziclient.Session, error) {
	tlsConfig, err := getAuthAPIClientTLSConfig()
	if err != nil {
		log.Errorln(`Failed to create TLS config for auth API client`)
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	basicAuthenticator := &enziclient.BasicAuthenticator{
		Username: username,
		Password: password,
	}

	return enziclient.New(httpClient, authAPIServerAddr, "enzi", basicAuthenticator), nil
}

func getUCPServiceOwnerEnziSession(authAPIServerAddr, OwnerID, serviceID string) (*enziclient.Session, error) {
	tlsConfig, err := getAuthAPIClientTLSConfig()
	if err != nil {
		log.Errorln(`Failed to create TLS config for auth API client`)
		return nil, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	signingConfig, err := getUCPServiceTokenSigningKey()
	if err != nil {
		return nil, fmt.Errorf("unable to load UCP service token signing config: %v", err)
	}

	openidClient := openid.NewClient(
		httpClient,
		signingConfig.key,
		serviceID,
		"redirectURI", // NOT REQUIRED: We do not use the authorizaiton_code grant type.
		authAPIServerAddr,
		"authorizationPath", // NOT REQUIRED: We do not use the authorizaiton_code grant type.
		enziauth.TokenAPIEndpointPath,
		signingConfig.certChain...,
	)

	tokenResp, err := openidClient.GetTokenWithAccountID(OwnerID)
	if err != nil {
		log.Errorln("unable to get identity token for UCP service owner")
		return nil, err
	}

	return enziclient.New(httpClient, authAPIServerAddr, "enzi", tokenResp), nil
}

func getAuthAPIClientTLSConfig() (*tls.Config, error) {
	caFilename := filepath.Join(AuthAPICertsVolumeMount, "ca.pem")

	rootCertsPEM, err := ioutil.ReadFile(caFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to read root CA certificates: %s", err)
	}

	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(rootCertsPEM) {
		return nil, fmt.Errorf("unable to parse root CA certificates")
	}

	return &tls.Config{RootCAs: rootCAs}, nil
}

type serviceTokenSigningKey struct {
	key       *jose.PrivateKey
	certChain []string
}

func getUCPServiceTokenSigningKey() (*serviceTokenSigningKey, error) {
	keyPEM, err := ioutil.ReadFile(filepath.Join(SwarmControllerCertVolumeMount, KeyFilename))
	if err != nil {
		return nil, fmt.Errorf("unable to read swarm controller key: %v", err)
	}

	certPEM, err := ioutil.ReadFile(filepath.Join(SwarmControllerCertVolumeMount, CertFilename))
	if err != nil {
		return nil, fmt.Errorf("unable to read swarm controller cert: %v", err)
	}

	keyPair, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("unable to parse swarm controller key pair: %s", err)
	}

	signingKey, err := jose.NewPrivateKey(keyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to load JOSE private key: %s", err)
	}

	// There is guaranteed to be at least one in the chain or else the
	// call to tls.X509KeyPair() above would have failed.
	certChain := make([]string, len(keyPair.Certificate))
	for i, certBytes := range keyPair.Certificate {
		certChain[i] = base64.StdEncoding.EncodeToString(certBytes)
	}

	return &serviceTokenSigningKey{
		key:       signingKey,
		certChain: certChain,
	}, nil
}

type KVDeployCfg struct {
	Timeout int
}

func GetKVDeployConfig(kvStoreURL *url.URL) (*KVDeployCfg, error) {
	kvCfg := &KVDeployCfg{
		Timeout: KVTimeout,
	}

	kv, err := GetKV(kvStoreURL)
	if err != nil {
		return kvCfg, err
	}

	kvKey := path.Join(OrcaPrefix, "config", "kv")
	kvPair, err := kv.Get(kvKey)
	if err != nil {
		if err == kvstore.ErrKeyNotFound {
			return kvCfg, nil
		}
		err = utils.MaybeWrapEtcdClusterErr(err)
		return kvCfg, fmt.Errorf("unable to get KV config from KV store: %s", err)
	}

	if err := json.Unmarshal(kvPair.Value, &kvCfg); err != nil {
		return kvCfg, fmt.Errorf("unable to decode KV config from JSON: %s", err)
	}

	return kvCfg, nil
}
