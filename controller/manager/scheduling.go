package manager

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

var (
	singletonSchedulingKvKey = "scheduling"
)

type SchedulingConfiguration struct {
	EnableAdminUCPScheduling bool `json:"enable_admin_ucp_scheduling"`
	EnableUserUCPScheduling  bool `json:"enable_user_ucp_scheduling"`
}

type SchedulingConfigSubsystem struct {
	ksKey string
	m     *DefaultManager
	cfg   *SchedulingConfiguration
}

func NewSchedulingConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	// scheduling doesn't support instances, just a single one
	if key != singletonSchedulingKvKey {
		log.Debugf("Malformed scheduling config key: %s", key)
		return nil, fmt.Errorf("Only one scheduling configuration supported")
	}

	enableAdminUCPScheduling := true
	enableUserUCPScheduling := true

	s := SchedulingConfigSubsystem{
		m:     m,
		ksKey: path.Join(KsConfigDir, singletonSchedulingKvKey),
		cfg: &SchedulingConfiguration{
			EnableAdminUCPScheduling: enableAdminUCPScheduling,
			EnableUserUCPScheduling:  enableUserUCPScheduling,
		},
	}

	cfgInt, err := s.ValidateConfig(jsonConfig, false)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgInt.(SchedulingConfiguration)
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration type")
	}
	*s.cfg = cfg
	m.configSubsystems[filepath.Base(s.ksKey)] = s
	return s, nil
}

func setupScheduling(m *DefaultManager) {
	setupSingletonConfigSubsystem(m, singletonSchedulingKvKey, "", NewSchedulingConfigSubsystem)
}

func (s SchedulingConfigSubsystem) GetKvKey() string {
	return s.ksKey
}

func (s SchedulingConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg SchedulingConfiguration

	if jsonConfig == "" {
		// use default settings
		cfg = SchedulingConfiguration{
			EnableAdminUCPScheduling: true,
			EnableUserUCPScheduling:  true,
		}
	} else {
		if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
			return nil, fmt.Errorf("Malformed scheduling configuration: %s", err)
		}
	}

	return cfg, nil
}

func (s SchedulingConfigSubsystem) UpdateConfig(cfgInt interface{}) error {
	cfg, ok := cfgInt.(SchedulingConfiguration)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	if cfg == *s.cfg {
		log.Debug("scheduling was unchanged")
		return nil
	}

	*s.cfg = cfg
	return nil
}

func (s SchedulingConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(s.cfg)
	return string(data), err
}

func (m DefaultManager) EnableAdminUCPScheduling() bool {
	s := m.configSubsystems[singletonSchedulingKvKey].(SchedulingConfigSubsystem)
	return s.cfg.EnableAdminUCPScheduling
}

func (m DefaultManager) EnableUserUCPScheduling() bool {
	s := m.configSubsystems[singletonSchedulingKvKey].(SchedulingConfigSubsystem)
	return s.cfg.EnableUserUCPScheduling
}
