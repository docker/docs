package manager

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca"
	"github.com/docker/orca/registry/v2"
)

var (
	// TODO - fix this for multiple registries
	ksRegistry = KsConfigDir + "/registry" // single registry (for V1)
)

type (
	RegistryConfigSubsystem struct {
		ksKey string
		m     *DefaultManager
		cfg   *orca.RegistryConfig
	}
)

// TODO - this doesn't actually support multiple registries yet, but the new config infra will allow it...
// We'll add support for that in a subsequent change with added UI support, and integration tests...
func NewRegistryConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	if key != "registry" { // TODO - this is wrong! support multiples
		log.Debugf("Malformed config key: %s", key)
		return nil, fmt.Errorf("Only one configuration supported")
	}
	s := RegistryConfigSubsystem{
		m:     m,
		ksKey: path.Join(KsConfigDir, "registry"),
		cfg:   &orca.RegistryConfig{},
	}

	cfgInt, err := s.ValidateConfig(jsonConfig, false)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgInt.(orca.RegistryConfig)
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration type")
	}
	s.UpdateConfig(cfgInt)
	*s.cfg = cfg
	m.configSubsystems[filepath.Base(s.ksKey)] = s
	return s, nil
}

func setupRegistry(m *DefaultManager) {
	setupSingletonConfigSubsystem(m, "registry", "{}", NewRegistryConfigSubsystem)
}

func (s RegistryConfigSubsystem) GetKvKey() string {
	return s.ksKey
}

func (s RegistryConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg orca.RegistryConfig

	if jsonConfig != "" {
		if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
			return nil, fmt.Errorf("Malformed registry configuration: %s", err)
		}

		uri, err := url.Parse(cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("Malformed URL: %s", err)
		}
		// set the ID so we can check against it easily
		cfg.ID = uri.Host
	}

	// TODO - consider more validation...

	return cfg, nil
}
func (s RegistryConfigSubsystem) UpdateConfig(cfgInt interface{}) error {
	cfg, ok := cfgInt.(orca.RegistryConfig)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	if cfg == *s.cfg {
		log.Debug("registry was unchanged")
		return nil
	}

	tlsConfig := *s.m.swarmTLSConfig
	registry, err := v2.NewRegistry(&cfg, &tlsConfig)
	if err != nil {
		return err
	}

	s.m.logEvent("update-registry", fmt.Sprintf("name=%s endpoint=%s", cfg.ID, cfg.URL), []string{"registry"})

	// XXX - this should probably be an atomic operation, or the manager should get the registry stored in the
	//       subsystem instead of directly on the manager
	s.m.registry = registry
	*s.cfg = cfg
	return nil
}

func (s RegistryConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(s.cfg)
	return string(data), err
}

func getRegistryKeyHash(name string) string {
	return GetKeyHash(name)
}

func (m DefaultManager) Registry(host string) (orca.Registry, error) {
	// we may not (yet) have a registry configured
	if m.registry == nil {
		return nil, ErrRegistryDoesNotExist
	}

	// we only support single registries for now.  Normally this
	// would look the host up here
	cfg := m.registry.GetConfig()
	if cfg.ID != host {
		return nil, ErrRegistryDoesNotExist
	}

	return m.registry, nil
}
