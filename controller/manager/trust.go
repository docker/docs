package manager

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

const SingletonTrustKvKey = "trust"

type TrustConfiguration struct {
	RequireContentTrustForDTR bool `json:"require_content_trust_for_dtr"`
	RequireContentTrustForHub bool `json:"require_content_trust_for_hub"`
}

type TrustConfigSubsystem struct {
	ksKey string
	m     *DefaultManager
	cfg   *TrustConfiguration
}

func NewTrustConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	// trust doesn't support instances, just a single one
	if key != SingletonTrustKvKey {
		log.Debugf("Malformed trust config key: %s", key)
		return nil, fmt.Errorf("Only one trust configuration supported")
	}

	t := TrustConfigSubsystem{
		m:     m,
		ksKey: path.Join(KsConfigDir, SingletonTrustKvKey),
		cfg: &TrustConfiguration{
			RequireContentTrustForDTR: false,
			RequireContentTrustForHub: false,
		},
	}

	cfgInt, err := t.ValidateConfig(jsonConfig, false)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgInt.(TrustConfiguration)
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration type")
	}
	*t.cfg = cfg
	m.configSubsystems[filepath.Base(t.ksKey)] = t
	return t, nil
}

func setupTrust(m *DefaultManager) {
	setupSingletonConfigSubsystem(m, SingletonTrustKvKey, "", NewTrustConfigSubsystem)
}

func (t TrustConfigSubsystem) GetKvKey() string {
	return t.ksKey
}

func (t TrustConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg TrustConfiguration

	if jsonConfig == "" {
		// use default settings
		cfg = TrustConfiguration{
			RequireContentTrustForDTR: false,
			RequireContentTrustForHub: false,
		}
		return cfg, nil
	}

	if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
		return nil, fmt.Errorf("Malformed trust configuration: %s", err)
	}
	return cfg, nil
}

func (t TrustConfigSubsystem) UpdateConfig(cfgInt interface{}) error {
	cfg, ok := cfgInt.(TrustConfiguration)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	if cfg == *t.cfg {
		log.Debug("trust configuration was unchanged")
		return nil
	}

	*t.cfg = cfg
	return nil
}

func (t TrustConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(t.cfg)
	return string(data), err
}
