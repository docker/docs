package config

import (
	"time"

	"github.com/cloudflare/cfssl/config"
)

// API Path constants.
const (
	APIPathPrefix = "/api/v1/cfssl/"
	APIInfoPath   = APIPathPrefix + "info"
	APISignPath   = APIPathPrefix + "sign"
)

// Signing returns the signing policy config used by our CA services.
func Signing() *config.Signing {
	return &config.Signing{
		Default: config.DefaultConfig(),
		Profiles: map[string]*config.SigningProfile{
			"client": {
				Usage: []string{
					"signing",
					"key encipherment",
					"client auth",
				},
				Expiry:       8760 * time.Hour,
				ExpiryString: "8760h",
			},
		},
	}
}
