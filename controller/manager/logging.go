package manager

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/docker/orca/utils"
)

var (
	singletonLoggingKvKey = "logging"
)

type LogConfiguration struct {
	Protocol string   `json:"protocol,omitempty"`
	Host     string   `json:"host,omitempty"`
	Level    string   `json:"level,omitempty"`
	hook     log.Hook `json:"-"`
}

type LogConfigSubsystem struct {
	ksKey string
	m     *DefaultManager
	cfg   *LogConfiguration
}

func NewLogConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	// Logging doesn't support instances, just a single one
	if key != singletonLoggingKvKey {
		log.Debugf("Malformed log config key: %s", key)
		return nil, fmt.Errorf("Only one logging configuration supported")
	}
	s := LogConfigSubsystem{
		m:     m,
		ksKey: path.Join(KsConfigDir, singletonLoggingKvKey),
		cfg:   &LogConfiguration{},
	}

	cfgInt, err := s.ValidateConfig(jsonConfig, false)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgInt.(LogConfiguration)
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration type")
	}
	*s.cfg = cfg
	m.configSubsystems[filepath.Base(s.ksKey)] = s
	return s, nil
}

func setupLogging(m *DefaultManager) {
	setupSingletonConfigSubsystem(m, singletonLoggingKvKey, "{}", NewLogConfigSubsystem)
}

func (s LogConfigSubsystem) GetKvKey() string {
	return s.ksKey
}

func (s LogConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(s.cfg)
	return string(data), err
}

func (s LogConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg LogConfiguration
	if jsonConfig == "" {
		log.Debug("Detected unconfiguration of logging - reverting to defaults")
	} else {
		if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
			return nil, fmt.Errorf("Malformed log configuration: %s", err)
		}
	}

	// TODO - do some more sanity validation of the settings here (to catch bad user input)

	if cfg.Level == "" {
		if log.GetLevel() == log.DebugLevel {
			cfg.Level = "DEBUG"
		} else {
			cfg.Level = "INFO"
		}
	}

	if cfg.Protocol != "" {
		if cfg.Protocol != "tcp" && cfg.Protocol != "udp" {
			return nil, fmt.Errorf("Invalid logging protocol: %s - must be tcp or udp", cfg.Protocol)
		}
	}
	return cfg, nil
}

func (s LogConfigSubsystem) UpdateConfig(cfgInt interface{}) error {
	var err error
	cfg, ok := cfgInt.(LogConfiguration)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	// See if the log config actually changed
	cfg.hook = s.cfg.hook
	if cfg == *s.cfg {
		log.Debug("Logging config unchanged")
		return nil
	}
	cfg.hook = nil

	// Interpret the log level
	level := strings.ToUpper(cfg.Level)
	syslogLevel := syslog.LOG_INFO
	logrusLevel := log.InfoLevel
	if strings.Contains(level, "EMERG") {
		syslogLevel = syslog.LOG_EMERG
		logrusLevel = log.FatalLevel
	} else if strings.Contains(level, "ALERT") {
		syslogLevel = syslog.LOG_ALERT
		logrusLevel = log.FatalLevel
	} else if strings.Contains(level, "CRIT") {
		syslogLevel = syslog.LOG_CRIT
		logrusLevel = log.FatalLevel
	} else if strings.Contains(level, "ERR") {
		syslogLevel = syslog.LOG_ERR
		logrusLevel = log.ErrorLevel
	} else if strings.Contains(level, "WARN") {
		syslogLevel = syslog.LOG_WARNING
		logrusLevel = log.WarnLevel
	} else if strings.Contains(level, "NOTICE") {
		syslogLevel = syslog.LOG_NOTICE
		logrusLevel = log.InfoLevel
	} else if strings.Contains(level, "INFO") {
		syslogLevel = syslog.LOG_INFO
		logrusLevel = log.InfoLevel
	} else if strings.Contains(level, "DEBUG") {
		syslogLevel = syslog.LOG_DEBUG
		logrusLevel = log.DebugLevel
	} else {
		log.Infof("Unrecognised log level %s - you must use standard syslog levels", level)
		// Let the default levels set above apply
	}

	unhook := func() {
		hooksMap := log.StandardLogger().Hooks
		for hookLevel, hooksList := range hooksMap {
			// Assumes the hook wont be duplicated in the list
			for i, h := range hooksList {
				if h == s.cfg.hook {
					hooksMap[hookLevel] = append(hooksMap[hookLevel][:i], hooksMap[hookLevel][i+1:]...)
					break
				}
			}
		}
	}

	// TODO - take a lock protecting the global config here

	log.SetLevel(logrusLevel)
	if cfg.Host != "" {
		if cfg.Protocol == "" {
			cfg.Protocol = "tcp"
		}

		if s.cfg.hook != nil {
			// Would be nice to see if it *actually* changed and fix it up,
			// but odds are it did, so just rebuild it to keep the algo simpler
			log.Debug("Unhooking old syslog hook")
			unhook()

		}
		// Note: logrus immediately sets up the logging, so if the cfg.Host is bad, we'll find out now
		cfg.hook, err = logrus_syslog.NewSyslogHook(cfg.Protocol, cfg.Host, syslogLevel, "") // tag unused
		if err != nil {
			return fmt.Errorf("Unable to connect to remote syslog daemon: %s", err)
		}
		log.AddHook(cfg.hook)
		log.Infof("Remote syslog activated for %s:%s", cfg.Protocol, cfg.Host)
	} else {
		if s.cfg.hook != nil {
			// Send one final message before we turn it off
			log.Info("Deactivating remote syslog")
			unhook()
		}
	}
	// Finally replace the current config with the new config
	*s.cfg = cfg
	return nil
}

// TODO - nuke this
func (s LogConfigSubsystem) updateOldConfig() {
	kv := s.m.Datastore()
	log.Info("Checking for old log config that needs schema updates")
	type (
		Version7 struct {
			Protocol string
			Host     string
			Level    string
		}
		Version8 struct {
			Version7
			Protocol string `json:"protocol,omitempty"`
			Host     string `json:"host,omitempty"`
			Level    string `json:"level,omitempty"`
		}
	)

	kvPair, _ := kv.Get(s.ksKey)
	if kvPair == nil {
		return
	}

	var cfg Version8
	if err := json.Unmarshal(kvPair.Value, &cfg); err != nil {
		log.Warnf("Failed to process old config for schema update: %s - %s", string(kvPair.Key), err)
		return
	}
	if cfg.Protocol == "" && cfg.Host == "" && cfg.Level == "" {
		cfg.Protocol = cfg.Version7.Protocol
		cfg.Host = cfg.Version7.Host
		cfg.Level = cfg.Version7.Level

		log.Infof("Detected old logging configuration, updating schema")
		data, err := json.Marshal(cfg)
		if err != nil {
			log.Warnf("Failed to unmarshal updated logging config: %s", err)
			return
		}
		if err := kv.Put(s.ksKey, data, nil); err != nil {
			err = utils.MaybeWrapEtcdClusterErr(err)
			log.Warnf("Failed to save updated logging config: %s", err)
		}
	}
}
