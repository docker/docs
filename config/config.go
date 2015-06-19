package config

import (
	"encoding/json"
	"io"

	"github.com/Sirupsen/logrus"
)

// Configuration is the top level object that
// all other configuration is namespaced under
type Configuration struct {
	Server       ServerConf       `json:"server,omitempty"`
	TrustService TrustServiceConf `json:"trust_service,omitempty"`
	Logging      LoggingConf      `json:"logging,omitempty"`
}

// ServerConf specifically addresses configuration related to
// the http server.
type ServerConf struct {
	Addr        string `json:"addr"`
	TLSCertFile string `json:"tls_cert_file"`
	TLSKeyFile  string `json:"tls_key_file"`
}

// TrustServiceConf specificies the service to use for signing.
// `Type` will be `local` for library based signing implementations,
// `remote` will be used for
type TrustServiceConf struct {
	Type      string `json:"type"`
	Hostname  string `json:"hostname,omitempty"`
	Port      string `json:"port,omitempty"`
	TLSCAFile string `json:"tls_ca_file,omitempty"`
}

type LoggingConf struct {
	Level uint8 `json:"level,omitempty"`
}

// Load takes a filename (relative path from pwd) and attempts
// to parse the file as a JSON obejct into the Configuration
// struct
func Load(data io.Reader) (*Configuration, error) {
	conf := Configuration{}
	decoder := json.NewDecoder(data)
	err := decoder.Decode(&conf)
	if err != nil {
		logrus.Error("[Notary Server] : Failed to parse configuration: ", err.Error())
		return nil, err
	}
	return &conf, nil
}
