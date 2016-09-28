package manager

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	singletonTrackingKvKey = "tracking"
)

type TrackingConfiguration struct {
	DisableUsageInfo  bool `json:"disable_usageinfo"`
	DisableTracking   bool `json:"disable_tracking"`
	AnonymizeTracking bool `json:"anonymize_tracking"`
}

type TrackingConfigSubsystem struct {
	ksKey string
	m     *DefaultManager
	cfg   *TrackingConfiguration
}

func NewTrackingConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	// tracking doesn't support instances, just a single one
	if key != singletonTrackingKvKey {
		log.Debugf("Malformed tracking config key: %s", key)
		return nil, fmt.Errorf("Only one tracking configuration supported")
	}
	s := TrackingConfigSubsystem{
		m:     m,
		ksKey: path.Join(KsConfigDir, singletonTrackingKvKey),
		cfg:   &TrackingConfiguration{},
	}

	cfgInt, err := s.ValidateConfig(jsonConfig, false)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgInt.(TrackingConfiguration)
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration type")
	}
	*s.cfg = cfg
	m.configSubsystems[filepath.Base(s.ksKey)] = s
	return s, nil
}

func setupTracking(m *DefaultManager) {
	m.disableTracking = new(bool)
	m.disableUsageInfo = new(bool)
	m.anonymizeTracking = new(bool)
	m.disableUsageInfoCh = usageInfo(m)
	setupSingletonConfigSubsystem(m, singletonTrackingKvKey, "{}", NewTrackingConfigSubsystem)
	// Report usage on setup
	if !*m.disableUsageInfo {
		m.reportUsage()
	}
}

func usageInfo(m *DefaultManager) chan<- bool {
	usageInfoDisabledCh := make(chan bool)
	var u bool
	ticker := time.NewTicker(time.Minute * 60)

	go func() {
		for {
			select {
			case <-ticker.C:
				if !u {
					m.reportUsage()
				}
			case u = <-usageInfoDisabledCh:
			}
		}
	}()
	return usageInfoDisabledCh
}

func (s TrackingConfigSubsystem) GetKvKey() string {
	return s.ksKey
}

func (s TrackingConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg TrackingConfiguration

	if jsonConfig == "" {
		// use default settings
		cfg = TrackingConfiguration{
			DisableUsageInfo:  false,
			DisableTracking:   false,
			AnonymizeTracking: false,
		}
	} else {
		if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
			return nil, fmt.Errorf("Malformed tracking configuration: %s", err)
		}
	}
	return cfg, nil

}
func (s TrackingConfigSubsystem) UpdateConfig(cfgInt interface{}) error {
	cfg, ok := cfgInt.(TrackingConfiguration)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	if cfg == *s.cfg {
		log.Debug("tracking was unchanged")
		return nil
	}

	*s.m.disableUsageInfo = cfg.DisableUsageInfo
	*s.m.disableTracking = cfg.DisableTracking
	*s.m.anonymizeTracking = cfg.AnonymizeTracking
	s.m.disableUsageInfoCh <- cfg.DisableUsageInfo

	*s.cfg = cfg
	return nil
}
func (s TrackingConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(s.cfg)
	return string(data), err
}

func (m DefaultManager) GetTrackingDisabled() bool {
	return *m.disableTracking
}

func (m DefaultManager) AnonymizeTracking() bool {
	return *m.anonymizeTracking
}

func (m DefaultManager) GetUsageInfoDisabled() bool {
	return *m.disableUsageInfo
}
