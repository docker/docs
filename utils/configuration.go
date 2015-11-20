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

// Server is a configuration about what addresses a server should listen on
type Server struct {
	*ServerTLSOpts
	HTTPAddr string
	GRPCAddr string
}

// Storage is a configuration about what storage backend a server should use
type Storage struct {
	Backend string `mapstructure:"backend"`
	URL     string `mapstructure:"db_url"`
}

// ParseServer tries to parse out a valid Server from a Viper:
// - Either or both of HTTP and GRPC address must be provided
// - If TLS is required, both the cert and key must be provided
// - If TLS is not requried, either both the cert and key must be provided or
//	 neither must be provided
func ParseServer(configuration *viper.Viper, tlsRequired bool) (*Server, error) {
	// mapstructure does not support unmarshalling into a pointer
	var tlsOpts ServerTLSOpts
	err := configuration.UnmarshalKey("server", &tlsOpts)
	if err != nil {
		return nil, err
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

	server := Server{
		HTTPAddr:      configuration.GetString("server.http_addr"),
		GRPCAddr:      configuration.GetString("server.grpc_addr"),
		ServerTLSOpts: &tlsOpts,
	}
	if cert == "" && key == "" && tlsOpts.ClientCAFile == "" {
		server.ServerTLSOpts = nil
	}

	if server.HTTPAddr == "" && server.GRPCAddr == "" {
		return nil, fmt.Errorf("server must have an HTTP and/or GRPC address")
	}

	return &server, nil
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
func ParseStorage(configuration *viper.Viper) (*Storage, error) {
	var store Storage
	err := configuration.UnmarshalKey("storage", &store)
	if err != nil {
		return nil, err
	}
	if store.Backend == "" && store.URL == "" {
		return nil, nil
	}
	store.Backend = strings.ToLower(store.Backend)
	if store.Backend != "mysql" {
		return nil, fmt.Errorf(
			"must specify one of these supported backends: mysql")
	}
	if store.URL == "" {
		return nil, fmt.Errorf("must provide a non-empty database URL")
	}
	return &store, nil
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

// utilities for handling common configurations

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
