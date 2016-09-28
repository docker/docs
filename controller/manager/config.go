package manager

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libkv/store"
	"github.com/docker/orca/utils"
)

var (
	// All orca controller configuration lives in this dir in the KV store
	KsConfigDir = path.Join(datastoreVersion, "config")

	// deprecatedConfigSubsystemKeys contains keys for config subsystems
	// which are no longer used. The config dir watcher will ignore these
	// keys and not try to instantiate config subsystem instances when they
	// are still in the kv store or if controllers running old code alter
	// state.
	deprecatedConfigSubsystemKeys = map[string]struct{}{
		"auth": {},
	}
)

// Each config subsystem implements this interface
type ConfigSubsystem interface {
	// Get the current configuration as serialized json
	GetConfiguration() (string, error)

	// Get the KV store location where the config is stored for this subsystem
	GetKvKey() string

	// Validate the given configuration - if valid, return a type specific config object
	ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error)

	// Replace the current configuration with the provided config (already validated)
	UpdateConfig(cfg interface{}) error
}

func (m DefaultManager) ListConfigSubsystems() []string {
	subsystems := []string{}
	for subsystem := range m.configSubsystems {
		subsystems = append(subsystems, subsystem)
	}
	return subsystems
}

// Common constructor function type for all subsystems
// keyOrInstance may be type, or type_$instance
// e.g., "logging" or "dtr_1.2.3.4"
// jsonConfig may be empty, which indicates a default config
type NewConfigSubsystem func(keyOrInstance, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error)

func (m DefaultManager) RegisterConfigSubsystem(key string, fp NewConfigSubsystem) {
	m.configSubsystemCtors[key] = fp
}

func (m DefaultManager) NewConfigSubsystemInstance(key, jsonConfig string) error {
	subsystemType := strings.Split(key, "_")[0]
	if fp, ok := m.configSubsystemCtors[subsystemType]; ok {
		_, err := fp(key, jsonConfig, &m)
		return err
	} else {
		return fmt.Errorf("Unrecognized type %s for key %s", subsystemType, key)
	}

}

func (m DefaultManager) GetSubsystemConfig(subsystemName string) (string, error) {
	subsystem := m.configSubsystems[subsystemName]
	if subsystem == nil {
		return "", fmt.Errorf("Invalid subsystem %s", subsystemName)
	}
	return subsystem.GetConfiguration()
}

func (m DefaultManager) UserConfigUpdate(subsystemName, jsonConfig string) error {
	// TODO
	// split the subsystemName on "_" and lookup
	subsystem := m.configSubsystems[subsystemName]
	if subsystem == nil {
		// Try to create a new instance
		return m.NewConfigSubsystemInstance(subsystemName, jsonConfig)
	}

	cfgInt, err := subsystem.ValidateConfig(jsonConfig, true)
	if err != nil {
		return err
	}
	oldCfgString, err := subsystem.GetConfiguration()
	if err != nil {
		// Should not happen
		return err
	}
	err = subsystem.UpdateConfig(cfgInt)
	if err != nil {
		return err
	}
	newCfgString, err := subsystem.GetConfiguration()
	if err != nil {
		// Should not happen
		return err
	}
	if oldCfgString != newCfgString {
		// Config actually changed, so write it out
		kv := m.Datastore()
		kvKey := subsystem.GetKvKey()
		if err := kv.Put(kvKey, []byte(newCfgString), nil); err != nil {
			err = utils.MaybeWrapEtcdClusterErr(err)
			log.Warnf("Unable to update %s config in kv store: %s", kvKey, err)
			return err
		}
		// Note: at this point we'll get an update event, but the config will match, so it'll be a no-op
	}
	return nil
}

func (m DefaultManager) setupConfigWatcher() {
	kv := m.Datastore()
	m.configStopCh = make(chan struct{}) // TODO - implement shutdown logic someday...

	// We can't watch a non-existent directory, so make sure it exists
	settings, err := kv.List(KsConfigDir)
	if len(settings) == 0 {
		kv.Put(KsConfigDir, nil, &store.WriteOptions{IsDir: true})
	}
	ch, err := kv.WatchTree(KsConfigDir, m.configStopCh)
	if err != nil {
		log.Error("Unable to watch the configuration - dynamic updates will not work")
		return
	}
	go m.configWatcher(ch)
}

// Follow updates to the controller configuration
func (m DefaultManager) configWatcher(ch <-chan []*store.KVPair) {
	log.Debug("Config update watcher activated")

	for list := range ch {
		// Every update iteration will contain the full set of config
		// To detect deletions (removal of configuration) we look for the absence
		// of expected keys
		processed := make(map[string]interface{})
		if list == nil {
			continue
		}
		for _, kvpair := range list {
			// XXX - this should handle directories as well as documents
			key := filepath.Base(kvpair.Key)

			if _, ok := deprecatedConfigSubsystemKeys[key]; ok {
				// Skip this now-unused subsystem.
				continue
			}

			processed[key] = struct{}{}
			if m.configSubsystems[key] != nil {
				if cfg, err := m.configSubsystems[key].ValidateConfig(string(kvpair.Value), false); err != nil {
					log.Errorf("Failed to validate %s configuration %s", key, err)
				} else {
					err := m.configSubsystems[key].UpdateConfig(cfg)
					if err != nil {
						log.Errorf("Failed to update %s configuration %s", key, err)
					}
				}
			} else {
				log.Infof("Detected new config subsystem instance %s", key)
				err := m.NewConfigSubsystemInstance(key, string(kvpair.Value))
				if err != nil {
					log.Error(err)
				}
			}
		}
		// Now look for deletions
		for key := range m.configSubsystems {
			if processed[key] == nil {
				// Send the empty string to signify no config stored in kv store
				if cfg, err := m.configSubsystems[key].ValidateConfig("", false); err != nil {
					log.Errorf("Failed to reset %s configuration to default %s", key, err)
				} else {
					err := m.configSubsystems[key].UpdateConfig(cfg)
					if err != nil {
						log.Errorf("Failed to update %s with default configuration %s", key, err)
					}
				}
			}
		}
	}
	// If we ever wire up the stop channel this should be reduced in severity and handled
	log.Warn("Config update watcher exiting")
}

// Helper routine to set up a singleton config subsystem
func setupSingletonConfigSubsystem(m *DefaultManager, key, defaultConfig string, ctor NewConfigSubsystem) {
	// Check for an existing config in the kv store, if not found, wire up the default
	m.RegisterConfigSubsystem(key, ctor)
	var s ConfigSubsystem
	existingCfgString, err := m.getKvValue(path.Join(KsConfigDir, key))
	if existingCfgString != "" && err == nil {
		s, err = ctor(key, existingCfgString, m)
		if err != nil {
			log.Warnf("Existing %s config appears malformed: %s", key, err)
		}
	}
	if s == nil {
		s, err = ctor(key, defaultConfig, m)
		if err != nil {
			log.Fatalf("Internal error setting up %s: %s", key, err)
		}
		// Write out the default config
		kv := m.Datastore()
		kvKey := s.GetKvKey()
		newCfgString, err := s.GetConfiguration()
		if err != nil {
			log.Warnf("Failed to serialize default config: %s", err)
		}
		if err := kv.Put(kvKey, []byte(newCfgString), nil); err != nil {
			err = utils.MaybeWrapEtcdClusterErr(err)
			log.Warnf("Unable to update %s config in kv store: %s", key, err)
		}
		// Note: at this point we'll get an update event, but the config will match, so it'll be a no-op
	}
}
