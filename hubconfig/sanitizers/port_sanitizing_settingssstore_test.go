package sanitizers

import (
	"testing"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/memory"
)

func TestSanitizedPort(t *testing.T) {
	testcases := []struct {
		cfg    hubconfig.HAConfig
		errors bool
	}{
		{
			cfg: hubconfig.HAConfig{
				ReplicaConfig: map[string]hubconfig.ReplicaConfig{
					"deadbeef": {
						HTTPPort:  deploy.AdminPort,
						HTTPSPort: deploy.AdminTlsPort,
					},
				},
			},
			errors: false,
		},
		{
			cfg: hubconfig.HAConfig{
				ReplicaConfig: map[string]hubconfig.ReplicaConfig{
					"deadbeef": {
						HTTPPort:  deploy.AdminTlsPort,
						HTTPSPort: deploy.AdminTlsPort,
					},
				},
			},
			errors: true,
		},
		{
			cfg: hubconfig.HAConfig{
				ReplicaConfig: map[string]hubconfig.ReplicaConfig{
					"deadbeef": {
						HTTPPort:  1,
						HTTPSPort: deploy.AdminTlsPort,
					},
				},
			},
			errors: true,
		},
		{
			cfg: hubconfig.HAConfig{
				ReplicaConfig: map[string]hubconfig.ReplicaConfig{
					"deadbeef": {
						HTTPPort:  0,
						HTTPSPort: deploy.AdminTlsPort,
					},
				},
			},
			errors: false,
		},
		{
			cfg: hubconfig.HAConfig{
				ReplicaConfig: map[string]hubconfig.ReplicaConfig{
					"deadbeef": {
						HTTPPort:  6000,
						HTTPSPort: deploy.AdminTlsPort,
					},
				},
			},
			errors: true,
		},
		{
			cfg: hubconfig.HAConfig{
				ReplicaConfig: map[string]hubconfig.ReplicaConfig{
					"deadbeef": {
						HTTPPort:  9999,
						HTTPSPort: deploy.AdminTlsPort,
					},
				},
			},
			errors: false,
		},
		{
			cfg: hubconfig.HAConfig{
				ReplicaConfig: map[string]hubconfig.ReplicaConfig{
					"deadbeef": {
						HTTPPort:  deploy.AdminPort,
						HTTPSPort: 0,
					},
				},
			},
			errors: false,
		},
	}

	for _, testcase := range testcases {
		ss := PortSanitizingSettingsStore{memory.NewSettingsStore()}
		err := ss.SetHAConfig(&testcase.cfg)
		if (err != nil) != testcase.errors {
			t.Fatalf("Test config (%v) expected error? %v but found err = %v instead", testcase.cfg, testcase.errors, err)
		}
	}
}
