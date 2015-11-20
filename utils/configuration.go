// Common configuration elements that may be resused

package utils

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	bugsnag_hook "github.com/Sirupsen/logrus/hooks/bugsnag"
	"github.com/bugsnag/bugsnag-go"
	"github.com/spf13/viper"
)

// Storage is a configuration about what storage backend a server should use
type Storage struct {
	Backend string
	Source  string
}

// ParseServerTLS tries to parse out a valid ServerTLSOpts from a Viper:
// - If TLS is required, both the cert and key must be provided
// - If TLS is not requried, either both the cert and key must be provided or
//	 neither must be provided
func ParseServerTLS(configuration *viper.Viper, tlsRequired bool) (*ServerTLSOpts, error) {
	//  unmarshalling into objects does not seem to pick up env vars
	tlsOpts := ServerTLSOpts{
		ServerCertFile: configuration.GetString("server.tls_cert_file"),
		ServerKeyFile:  configuration.GetString("server.tls_key_file"),
		ClientCAFile:   configuration.GetString("server.client_ca_file"),
	}

	cert, key := tlsOpts.ServerCertFile, tlsOpts.ServerKeyFile
	if tlsRequired {
		if cert == "" || key == "" {
			return nil, fmt.Errorf("both the TLS certificate and key are mandatory")
		}
	} else {
		if (cert == "" && key != "") || (cert != "" && key == "") {
			return nil, fmt.Errorf(
				"either include both a cert and key file, or neither to disable TLS")
		}
	}

	if cert == "" && key == "" && tlsOpts.ClientCAFile == "" {
		return nil, nil
	}

	return &tlsOpts, nil
}

// ParseLogLevel tries to parse out a log level from a Viper.  If there is no
// configuration, defaults to the provided error level
func ParseLogLevel(configuration *viper.Viper, defaultLevel logrus.Level) (
	logrus.Level, error) {

	logStr := configuration.GetString("logging.level")
	if logStr == "" {
		return defaultLevel, nil
	}
	return logrus.ParseLevel(logStr)
}

// ParseStorage tries to parse out Storage from a Viper.  If backend and
// URL are not provided, returns a nil pointer.
func ParseStorage(configuration *viper.Viper, allowedBackeneds []string) (*Storage, error) {
	store := Storage{
		Backend: configuration.GetString("storage.backend"),
		Source:  configuration.GetString("storage.db_url"),
	}
	if store.Backend == "" && store.Source == "" {
		return nil, nil
	}

	if store.Source == "" {
		return nil, fmt.Errorf("must provide a non-empty database source")
	}
	store.Backend = strings.ToLower(store.Backend)
	for _, backend := range allowedBackeneds {
		if backend == store.Backend {
			return &store, nil
		}

	}
	return nil, fmt.Errorf(
		"must specify one of these supported backends: %s",
		strings.Join(allowedBackeneds, ", "))
}

// ParseBugsnag tries to parse out a Bugsnag Configuration from a Viper.
// If no values are provided, returns a nil pointer.
func ParseBugsnag(configuration *viper.Viper) (*bugsnag.Configuration, error) {
	// can't unmarshal because we can't add tags to the bugsnag.Configuration
	// struct
	bugconf := bugsnag.Configuration{
		APIKey:       configuration.GetString("reporting.bugsnag.api_key"),
		ReleaseStage: configuration.GetString("reporting.bugsnag.release_stage"),
		Endpoint:     configuration.GetString("reporting.bugsnag.endpoint"),
	}
	if bugconf.APIKey == "" && bugconf.ReleaseStage == "" && bugconf.Endpoint == "" {
		return nil, nil
	}
	if bugconf.APIKey == "" {
		return nil, fmt.Errorf("must provide an API key for bugsnag")
	}
	return &bugconf, nil
}

// utilities for setting up/acting on common configurations

// SetupViper sets up an instance of viper to also look at environment
// variables
func SetupViper(configuration *viper.Viper, envPrefix string) {
	configuration.SetEnvPrefix(envPrefix)
	configuration.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configuration.AutomaticEnv()
}

// SetUpBugsnag configures bugsnag and sets up a logrus hook
func SetUpBugsnag(config *bugsnag.Configuration) error {
	if config != nil {
		bugsnag.Configure(*config)
		hook, err := bugsnag_hook.NewBugsnagHook()
		if err != nil {
			return err
		}
		logrus.AddHook(hook)
		logrus.Debug("Adding logrus hook for Bugsnag")
	}
	return nil
}
