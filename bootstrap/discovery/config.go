package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
)

const (
	ClusterAdvertise = "cluster-advertise"
	ClusterStore     = "cluster-store"
	ClusterStoreOpts = "cluster-store-opts"
)

var (
	DiscoveryCertsFolderName = "discovery_certs"
	DiscoveryCertsDir        = filepath.Join(config.EngineLibDir, DiscoveryCertsFolderName)
)

func LoadCurrentConfiguration() (map[string]interface{}, error) {
	// Process the configuration file generically so we don't drop any fields we don't recognize
	var cfg map[string]interface{}
	data, err := ioutil.ReadFile(config.EngineConfigFile)
	if err != nil && !os.IsNotExist(err) { // Unexpected error
		return nil, fmt.Errorf("Failed to open configuration file: %s", err)
	} else if err == nil { // Update existing file
		var f interface{}
		d := json.NewDecoder(bytes.NewReader(data))
		d.UseNumber()
		err := d.Decode(&f)
		if err != nil {
			return nil, fmt.Errorf("Malformed daemon configuration file %s - %s", config.EngineConfigFile, err)
		}
		cfg = f.(map[string]interface{})
	} else { // Config doesn't exist yet
		cfg = map[string]interface{}{}
	}
	return cfg, nil
}

func PrettyConfig(cfg map[string]interface{}) (string, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		// Shouldn't happen...
		return "", err
	}
	var prettyConfig bytes.Buffer
	json.Indent(&prettyConfig, data, "", "\t")
	return prettyConfig.String(), nil
}

func GenerateUpdatedConfig(cfg map[string]interface{}, advertise string, controllers []string) (map[string]interface{}, bool) {
	changed := false
	expectedCA := filepath.Join(DiscoveryCertsDir, "ca.pem")
	expectedCert := filepath.Join(DiscoveryCertsDir, "cert.pem")
	expectedKey := filepath.Join(DiscoveryCertsDir, "key.pem")
	// Synthesize the controllers URL
	clusterStore := config.KvType + "://" + strings.Join(controllers, ",")

	oldClusterAdvertise := cfg[ClusterAdvertise]
	if oldClusterAdvertise == nil || reflect.TypeOf(oldClusterAdvertise).Kind() != reflect.String || oldClusterAdvertise.(string) != advertise {
		changed = true
		cfg[ClusterAdvertise] = advertise
	}

	oldClusterStore := cfg[ClusterStore]
	if oldClusterStore == nil || reflect.TypeOf(oldClusterStore).Kind() != reflect.String || oldClusterStore != clusterStore {
		changed = true
		cfg[ClusterStore] = clusterStore
	}

	oldClusterOpts := cfg[ClusterStoreOpts]
	if oldClusterOpts != nil {
		switch oldVal := oldClusterOpts.(type) {
		case map[string]interface{}:
			// Verify all 3 entries, update if needed
			opts := oldClusterOpts.(map[string]interface{})
			oldCA := opts["kv.cacertfile"]
			if oldCA == nil || reflect.TypeOf(oldCA).Kind() != reflect.String || oldCA != expectedCA {
				changed = true
				opts["kv.cacertfile"] = expectedCA
			}
			oldCert := opts["kv.certfile"]
			if oldCert == nil || reflect.TypeOf(oldCert).Kind() != reflect.String || oldCert != expectedCert {
				changed = true
				opts["kv.certfile"] = expectedCert
			}
			oldKey := opts["kv.keyfile"]
			if oldKey == nil || reflect.TypeOf(oldKey).Kind() != reflect.String || oldKey != expectedKey {
				changed = true
				opts["kv.keyfile"] = expectedKey
			}
		default:
			log.Debugf("Unexpected type for %s %V (updating)", ClusterStoreOpts, oldVal)
			changed = true
			cfg[ClusterStoreOpts] = map[string]string{
				"kv.cacertfile": expectedCA,
				"kv.certfile":   expectedCert,
				"kv.keyfile":    expectedKey,
			}
		}
	} else {
		changed = true
		cfg[ClusterStoreOpts] = map[string]string{
			"kv.cacertfile": expectedCA,
			"kv.certfile":   expectedCert,
			"kv.keyfile":    expectedKey,
		}
	}
	return cfg, changed
}

func GenerateUninstalledConfig(cfg map[string]interface{}) map[string]interface{} {
	delete(cfg, ClusterAdvertise)
	delete(cfg, ClusterStore)
	delete(cfg, ClusterStoreOpts)

	return cfg
}

func ApplyConfigChanges(ec *client.EngineClient, cfg map[string]interface{}, controllerCount int, isControllerReplica bool) error {
	initialConfig, err := isInitialConfig()
	if err != nil {
		return err
	}

	// Since humans may edit, make it purty
	prettyConfig, err := PrettyConfig(cfg)
	if err != nil {
		return err
	}
	// Show updated configuration at debug level
	log.Debugf("New configuration as follows:\n%s", prettyConfig)
	// Write out the updated configuration file
	if err := ioutil.WriteFile(config.EngineConfigFile, []byte(prettyConfig), 0644); err != nil {
		return fmt.Errorf("Unable to update %s - %s", config.EngineConfigFile, err)
	}

	if initialConfig && ec.EngineSupportsSignal() {
		err := signalEngine()
		if err != nil {
			log.Debugf("Signaling engine failed: %s", err)
		} else {
			// Success so don't print the warnings below
			return nil
		}
	}

	// Warn about some long restart time gotchas
	log.Warn("Configuration updated. You will have to manually restart the docker daemon for the changes to take effect.")
	if isControllerReplica {
		if controllerCount == 2 {
			log.Warn("Your cluster only has two controllers. Adding a third node prior to restarting the daemon is strongly recommended to prevent long restart times.")
		} else if controllerCount > 2 {
			// Since we asked them to not restart earlier, they might restart
			// multiple controllers at once, so warn against that too.
			log.Warn("If you have to restart daemons on multiple controllers, restart them one by one to prevent long restart times.")
		}
	}
	return nil
}

func isInitialConfig() (bool, error) {
	_, err := os.Stat(config.EngineConfigFile)
	if err != nil && os.IsNotExist(err) {
		return true, nil
	}
	return false, err
}

func signalEngine() error {
	// Signalling the engine is sometimes useful for us, depending on version:
	// - explodes on <=1.9 (unsupported by UCP)
	// - does not load libnetwork config on 1.10 OSS
	// - hits a bad bug in libnetwork on 1.10 CS (libnetwork issue 1051)
	// - will load initial config on 1.11 OSS and 1.10 CS
	// - will not reload existing config on any version

	log.Info("New configuration established.  Signalling the daemon to load it...")
	// Make sure the message gets delivered in case things go poof
	time.Sleep(1 * time.Second)

	data, err := ioutil.ReadFile(filepath.Join(config.EnginePidDir, "docker.pid")) // XXX Is this going to work across all supported distros?
	if err != nil {
		return fmt.Errorf("Unable to locate the docker daemon PID file.  You will have to manually restart your docker daemon for the changes to take effect.")
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		log.Debugf("Malformed pid file: %s", string(data))
		return fmt.Errorf("Unable to process the docker daemon PID file.  You will have to manually restart your docker daemon for the changes to take effect")
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("Unable to locate the docker daemon with the PID %d.  You will have to manually restart your docker daemon for the changes to take effect.", pid)
	}
	log.Debugf("Sending signal %d to pid %d", config.EngineConfigReloadSignal, pid)
	err = proc.Signal(config.EngineConfigReloadSignal)
	if err != nil {
		log.Warnf("Failed to signal docker daemon (%s) You will have to manually restart your docker daemon for the changes to take effect.", err)
	} else {
		log.Info("Successfully delivered signal to daemon")
	}
	return nil
}
