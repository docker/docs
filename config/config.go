package config

import (
	"encoding/json"
	"os"
)

// Configuration is the top level object that
// all other configuration is namespaced under
type Configuration struct {
	Server ServerConf `json:"server,omitempty"`
}

// ServerConf specifically addresses configuration related to
// the http server.
type ServerConf struct {
	Addr        string `json:"addr"`
	TLSCertFile string `json:"tls_cert_file"`
	TLSKeyFile  string `json:"tls_key_file"`
}

// Load takes a filename (relative path from pwd) and attempts
// to parse the file as a JSON obejct into the Configuration
// struct
func Load(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	conf := Configuration{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
