package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	sampleConfig := "{\"server\": {\"addr\":\"testAddr\",\"tls_cert_file\":\"testCertFile\",\"tls_key_file\":\"testKeyFile\"}}"
	sampleConfigStruct := &Configuration{
		Server: ServerConf{
			Addr:        "testAddr",
			TLSCertFile: "testCertFile",
			TLSKeyFile:  "testKeyFile",
		},
	}
	conf, err := Load(strings.NewReader(sampleConfig))
	if err != nil {
		t.Fatalf("Error parsing config: %s", err.Error())
	}

	if !reflect.DeepEqual(conf, sampleConfigStruct) {
		t.Fatalf("Parsed config did not match expected.")
	}
}
