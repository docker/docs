package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

// Configuration is a versioned Garant token server configuration, intended to
// be provided by a yaml file, and optionally modified by environment vars.
type Configuration struct {
	// Version is the version which defines the format of the rest of the configuration
	Version Version `yaml:"version"`

	// Logging supports setting various parameters related to the logging
	// subsystem.
	Logging Logging `yaml:"logging"`

	// SigningKey is the path to the token signing key in JWK format. The
	// default signing key location is `/etc/garant/signing_key.json`.
	SigningKey string `yaml:"signingkey"`

	// Issuer is the issuer identity string that Garant will use in the issuer
	// ("iss") claim of the JWTs that it produces. Required.
	Issuer string `yaml:"issuer"`

	// Auth allows configuration of the authentication and authorization
	// backend to use.
	Auth Auth `yaml:"auth"`

	// Reporting is the configuration for error reporting
	Reporting Reporting `yaml:"reporting"`

	// HTTP contains configuration parameters for the server's http
	// interface.
	HTTP HTTP `yaml:"http"`
}

// v0_1Configuration is a Version 0.1 Configuration struct
// This is currently aliased to Configuration, as it is the current version
type v0_1Configuration Configuration

// Logging defines the configuration for the logging subsystem.
type Logging struct {
	// Level is the granularity at which registry operations are logged.
	Level string `yaml:"level"`

	// Formatter overrides the default formatter with another. Options
	// include "text", "json" and "logstash".
	Formatter string `yaml:"formatter"`

	// Fields allows users to specify static string fields to include in
	// the logger context.
	Fields map[string]string `yaml:"fields"`
}

// Parameters defines a key-value parameters mapping
type Parameters map[string]interface{}

// Auth defines the configuration for an authentication and authorization
// driver.
type Auth struct {
	// BackendName is the auth driver name
	BackendName string
	// Parameters is a map of extra configuration parameters
	Parameters Parameters
}

// UnmarshalYAML implements the yaml.Unmarshaler interface. Unmarshals a single
// item map into a Auth or a string into a Auth type with no parameters.
func (auth *Auth) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var (
		val         map[string]Parameters
		backendName string
		parameters  Parameters
	)

	err := unmarshal(&val)
	if err == nil {
		if len(val) != 1 {
			types := make([]string, 0, len(val))
			for k := range val {
				types = append(types, k)
			}
			return fmt.Errorf("must provide exactly one type. Provided: %v", types)

		}
		for backendName, parameters = range val {
			// Get the one pair from the map.
		}
	} else if err = unmarshal(&backendName); err == nil {
		parameters = Parameters{}
	}

	*auth = Auth{
		BackendName: backendName,
		Parameters:  parameters,
	}

	return nil
}

// MarshalYAML implements the yaml.Marshaler interface
func (auth Auth) MarshalYAML() (interface{}, error) {
	if auth.Parameters == nil {
		return auth.BackendName, nil
	}
	return map[string]Parameters{auth.BackendName: auth.Parameters}, nil
}

// Reporting defines error reporting methods.
type Reporting struct {
	// Bugsnag configures error reporting for Bugsnag (bugsnag.com).
	Bugsnag BugsnagReporting `yaml:"bugsnag,omitempty"`
	// NewRelic configures error reporting for NewRelic (newrelic.com)
	NewRelic NewRelicReporting `yaml:"newrelic,omitempty"`
}

// BugsnagReporting configures error reporting for Bugsnag (bugsnag.com).
type BugsnagReporting struct {
	// APIKey is the Bugsnag api key.
	APIKey string `yaml:"apikey,omitempty"`
	// ReleaseStage tracks where the Garant is deployed.
	// Examples: production, staging, development
	ReleaseStage string `yaml:"releasestage,omitempty"`
	// Endpoint is used for specifying an enterprise Bugsnag endpoint.
	Endpoint string `yaml:"endpoint,omitempty"`
	// The list of ReleaseStages to notify in. By default Bugsnag will notify
	// you in all release stages, but you can use this to silence development
	// errors.
	NotifyReleaseStages []string `yaml:"notifyreleasestages,omitempty"`
	// If you use a versioning scheme for deploys of your app, Bugsnag can use
	// the AppVersion to only re-open errors if they occur in later version of
	// the app.
	AppVersion string `yaml:"appversion,omitempty"`
	// In order to determine where a crash happens Bugsnag needs to know which
	// packages you consider to be part of your app (as opposed to a library).
	// By default this is set to []string{"main*"}. Strings are matched to
	// package names using filepath.Match.
	ProjectPackages []string `yaml:"projectpackages,omitempty"`
}

// NewRelicReporting configures error reporting for NewRelic (newrelic.com)
type NewRelicReporting struct {
	// LicenseKey is the NewRelic user license key
	LicenseKey string `yaml:"licensekey,omitempty"`
	// Name is the component name of the Garant token server in NewRelic
	Name string `yaml:"name,omitempty"`
}

// HTTP defines HTTP Server config.
type HTTP struct {
	// Addr specifies the bind address for the registry instance.
	// The default is `localhost:8080`.
	Addr string `yaml:"addr,omitempty"`
	// Prefix specified a URL path prefix to use.
	Prefix string `yaml:"prefix,omitempty"`
	// TLS instructs the http server to listen with a TLS configuration.
	// This only support simple tls configuration with a cert and key.
	// Mostly, this is useful for testing situations or simple deployments
	// that require tls. If more complex configurations are required, use
	// a proxy or make a proposal to add support here.
	TLS *TLS `yaml:"tls,omitempty"`
}

// TLS defines TLS certificate and key files.
type TLS struct {
	// Certificate specifies the path to an x509 certificate file to
	// be used for TLS.
	Certificate string `yaml:"certificate"`
	// Key specifies the path to the x509 key file, which should
	// contain the private portion for the file specified in
	// Certificate.
	Key string `yaml:"key"`
}

// Parse parses an input configuration yaml document into a Configuration
// struct. This should generally be capable of handling old configuration
// format versions.
//
// Environment variables may be used to override configuration parameters other
// than version, following the scheme below:
// Configuration.Abc may be replaced by the value of GARANT_ABC,
// Configuration.Abc.Xyz may be replaced by the value of GARANT_ABC_XYZ, and so
// forth.
func Parse(rd io.Reader) (*Configuration, error) {
	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	p := NewParser("garant", []VersionedParseInfo{
		{
			Version: MajorMinorVersion(0, 1),
			ParseAs: reflect.TypeOf(v0_1Configuration{}),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				if v0_1, ok := c.(*v0_1Configuration); ok {
					if v0_1.Auth.BackendName == "" {
						return nil, fmt.Errorf("no auth backend configuration provided")
					}
					if v0_1.SigningKey == "" {
						// use default signing key location.
						v0_1.SigningKey = "/etc/garant/signing_key.json"
					}
					if v0_1.Issuer == "" {
						return nil, fmt.Errorf("no issuer configuration provided")
					}
					return (*Configuration)(v0_1), nil
				}
				return nil, fmt.Errorf("expected *v0_1Configuration, received %#v", c)
			},
		},
	})

	config := new(Configuration)
	err = p.Parse(in, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
