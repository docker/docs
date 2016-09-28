package manager

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libtrust"
	"github.com/docker/orca/utils"
)

var (
	singletonLicenceKvKey = "license"
	UnlicensedID          = "unlicensed"
)

// Downloaded from HUB or the license server (common with DTR)
type LicenseConfig struct {
	KeyID         string `yaml:"key_id" json:"key_id"`
	PrivateKey    string `yaml:"private_key" json:"private_key,omitempty"`
	Authorization string `yaml:"authorization" json:"authorization,omitempty"`
}

// This is the top-level configuration used by the subsystem
type LicenseSubsystemConfig struct {
	// Tune if we automatically attempt to refresh expired licenses (read-write)
	AutoRefresh bool `json:"auto_refresh"`
	// The underlying license returned from the license server (read-write)
	License LicenseConfig `json:"license_config"`
	// Details derived from License (read-only)
	Details LicenseDetails `json:"details"`
	// Any error reported from the license server on the last attempt to auto refresh
	LastUpdateError string `json:"last_refresh_error"`
}

type LicenseConfigSubsystem struct {
	ksKey             string
	m                 *DefaultManager
	cfg               *LicenseSubsystemConfig
	publicKey         libtrust.PublicKey
	recentEngineCount *int
}

// Used for read-only display purposes in the UI
type LicenseDetails struct {
	MaxEngines int       `json:"max_engines"`
	Expiration time.Time `json:"expiration"`
	Tier       string
	// Other fields from the License response are ignored currently
	// but we may add them later (tiers, type, etc.)
}

const (
	// This is the public key from the production license server
	publicRSAKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0Ka2lkOiBKN0xEOjY3VlI6TDVIWjpVN0JBOjJPNEc6NEFMMzpPRjJOOkpIR0I6RUZUSDo1Q1ZROk1GRU86QUVJVAoKTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUF5ZEl5K2xVN283UGNlWSs0K3MrQwpRNU9FZ0N5RjhDeEljUUlXdUs4NHBJaVpjaVk2NzMweUNZbndMU0tUbHcrVTZVQy9RUmVXUmlvTU5ORTVEczVUCllFWGJHRzZvbG0ycWRXYkJ3Y0NnKzJVVUgvT2NCOVd1UDZnUlBIcE1GTXN4RHpXd3ZheThKVXVIZ1lVTFVwbTEKSXYrbXE3bHA1blEvUnhyVDBLWlJBUVRZTEVNRWZHd20zaE1PL2dlTFBTK2hnS1B0SUhsa2c2L1djb3hUR29LUAo3OWQvd2FIWXhHTmw3V2hTbmVpQlN4YnBiUUFLazIxbGc3OThYYjd2WnlFQVRETXJSUjlNZUU2QWRqNUhKcFkzCkNveVJBUENtYUtHUkNLNHVvWlNvSXUwaEZWbEtVUHliYncwMDBHTyt3YTJLTjhVd2dJSW0waTVJMXVXOUdrcTQKempCeTV6aGdxdVVYYkc5YldQQU9ZcnE1UWE4MUR4R2NCbEp5SFlBcCtERFBFOVRHZzR6WW1YakpueFpxSEVkdQpHcWRldlo4WE1JMHVrZmtHSUkxNHdVT2lNSUlJclhsRWNCZi80Nkk4Z1FXRHp4eWNaZS9KR1grTEF1YXlYcnlyClVGZWhWTlVkWlVsOXdYTmFKQitrYUNxejVRd2FSOTNzR3crUVNmdEQwTnZMZTdDeU9IK0U2dmc2U3QvTmVUdmcKdjhZbmhDaVhJbFo4SE9mSXdOZTd0RUYvVWN6NU9iUHlrbTN0eWxyTlVqdDBWeUFtdHRhY1ZJMmlHaWhjVVBybQprNGxWSVo3VkQvTFNXK2k3eW9TdXJ0cHNQWGNlMnBLRElvMzBsSkdoTy8zS1VtbDJTVVpDcXpKMXlFbUtweXNICjVIRFc5Y3NJRkNBM2RlQWpmWlV2TjdVQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo="
	// How close to expiration should we start trying to refresh
	licenseRefreshDelta = 12 * time.Hour
	licenseServerURL    = "https://license.enterprise.docker.com/v1/check"
)

func NewLicenseConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	// Licensing doesn't support instances, just a single one
	recentEngineCount := 0
	if key != singletonLicenceKvKey {
		log.Debugf("Malformed license config key: %s", key)
		return nil, fmt.Errorf("Only one license configuration supported")
	}
	pemBytes, err := base64.StdEncoding.DecodeString(publicRSAKey)
	if err != nil {
		// Should never happen
		log.Fatalf("Failed to decode the embedded license public key during init: %s", err)
	}

	publicKey, err := libtrust.UnmarshalPublicKeyPEM(pemBytes)
	if err != nil {
		// Should never happen
		log.Fatalf("Failed to parse the embedded license public key during init: %s", err)
	}
	s := LicenseConfigSubsystem{
		m:                 m,
		ksKey:             path.Join(KsConfigDir, singletonLicenceKvKey),
		publicKey:         publicKey,
		recentEngineCount: &recentEngineCount,
		cfg:               &LicenseSubsystemConfig{},
	}

	cfgInt, err := s.ValidateConfig(jsonConfig, false)
	if err != nil {
		return nil, err
	}
	cfg, ok := cfgInt.(LicenseSubsystemConfig)
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration type")
	}
	*s.cfg = cfg
	m.configSubsystems[filepath.Base(s.ksKey)] = s
	return s, nil
}

// Wire up a callback handler for logging configuration
func setupLicensing(m *DefaultManager) {
	setupSingletonConfigSubsystem(m, singletonLicenceKvKey, "", NewLicenseConfigSubsystem)
}

func (s LicenseConfigSubsystem) UpdateConfig(cfgInt interface{}) error {
	cfg, ok := cfgInt.(LicenseSubsystemConfig)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	// We don't take any actions on update
	*s.cfg = cfg
	return nil
}
func (s LicenseConfigSubsystem) GetKvKey() string {
	return s.ksKey
}

func (s LicenseConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(s.cfg)
	return string(data), err
}

func (s LicenseConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg LicenseSubsystemConfig
	if jsonConfig == "" {
		// If no config, set up unlicensed default setup
		cfg = LicenseSubsystemConfig{
			AutoRefresh: true,
			License: LicenseConfig{
				KeyID: UnlicensedID,
			},
			Details: LicenseDetails{
				Expiration: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
	} else {
		if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
			return nil, fmt.Errorf("Malformed license configuration: %s", err)
		}
	}

	if cfg.License.KeyID == UnlicensedID {
		// No further processing for unlicensed
		return cfg, nil
	}

	authorizationBytes, err := base64.StdEncoding.DecodeString(cfg.License.Authorization)
	if err != nil {
		return nil, fmt.Errorf("Malformed license key - Failed to decode key authorization: %s", err)
	}

	jsonSignature, err := libtrust.ParseJWS(authorizationBytes)
	if err != nil {
		return nil, fmt.Errorf("Malformed license key - Failed to parse JWT: %s", err)
	}

	checkResponse, err := verifyJsonSignature(cfg.License.PrivateKey, jsonSignature, s.publicKey)
	if err != nil {
		return nil, fmt.Errorf("Malformed license key - Signature verification failure: %s", err)
	}

	// always clobber the details from the actual parsed license
	cfg.Details.MaxEngines = checkResponse.MaxEngines
	cfg.Details.Expiration = checkResponse.Expiration
	cfg.Details.Tier = checkResponse.Tier

	return cfg, nil
}

func (m DefaultManager) GetLicenseKeyID() string {
	s := m.configSubsystems[singletonLicenceKvKey].(LicenseConfigSubsystem)
	return s.cfg.License.KeyID
}

func (m DefaultManager) GetLicenseTier() string {
	license := m.GetLicense()
	return license.Details.Tier
}

func (m DefaultManager) getLicenseBanners() []Banner {
	s := m.configSubsystems[singletonLicenceKvKey].(LicenseConfigSubsystem)
	if s.cfg.License.KeyID == UnlicensedID {
		return []Banner{{
			Level:   BannerWARN,
			Message: "Your system is unlicensed.  Please visit https://store.docker.com/bundles/docker-datacenter to download a free trial license or purchase an enterprise license.",
		}}
	}

	now := time.Now().UTC()
	if now.After(s.cfg.Details.Expiration) {
		msg := "Your license has expired.  "
		if s.cfg.AutoRefresh {
			msg += "We've been unable to auto renew.  " + s.cfg.LastUpdateError
		} else {
			msg += "Your deployment is out of compliance and we may be unable to provide support as a result. Please visit https://hub.docker.com/account/licenses/ to download a renewed license. "
		}
		return []Banner{{
			Level:   BannerWARN,
			Message: msg,
		}}
	}

	// TODO - decide how close to expiring we want to nag them and catch that scenario here

	if *s.recentEngineCount > s.cfg.Details.MaxEngines {
		return []Banner{{
			Level:   BannerWARN,
			Message: fmt.Sprintf("You are exceeding your licensed engine count of %d by %d engines", s.cfg.Details.MaxEngines, *s.recentEngineCount-s.cfg.Details.MaxEngines),
		}}
	}
	return []Banner{}
}

func (m DefaultManager) GetLicense() LicenseSubsystemConfig {
	s := m.configSubsystems[singletonLicenceKvKey].(LicenseConfigSubsystem)
	return *s.cfg
}

// Called periodically to see if we should do any license renewals
func (m DefaultManager) periodicLicenseCheck() {
	s := m.configSubsystems[singletonLicenceKvKey].(LicenseConfigSubsystem)

	// Gather node count for licensing purposes
	nodes, err := m.Nodes()
	if err == nil {
		*s.recentEngineCount = len(nodes)
	}

	// If we're unlicensed, don't bother trying to update
	if s.cfg.License.KeyID == UnlicensedID {
		return
	}

	// Check for almost expiring
	now := time.Now().UTC()
	if now.After(s.cfg.Details.Expiration.Add(-licenseRefreshDelta)) && s.cfg.AutoRefresh {
		cfg := *s.cfg
		defer func() {
			*s.cfg = cfg
			data, err := json.Marshal(s.cfg)
			if err == nil {
				kv := m.Datastore()
				if err := kv.Put(s.GetKvKey(), data, nil); err != nil {
					err = utils.MaybeWrapEtcdClusterErr(err)
					log.Warnf("Unable to update license config in kv store: %s", err)
				}
			}
		}()

		// TODO - better support for HA cluster (less wasted redundancy)
		//        license refresh is largely idempotent, so having multiple nodes
		//        refresh concurrently shouldn't be a problem.
		log.Debug("License is close to expiring, attempting refresh")

		token, err := GenerateToken(now.Format(time.RFC3339), cfg.License.PrivateKey)
		if err != nil {
			cfg.LastUpdateError = fmt.Sprintf("Failed to generate renewal token: %s", err)
			return
		}
		checkLicenseRequest := &CheckLicenseRequest{KeyID: cfg.License.KeyID, Timestamp: now, Token: token}
		jsonedCheckLicenseRequest, err := json.Marshal(checkLicenseRequest)
		if err != nil {
			cfg.LastUpdateError = fmt.Sprintf("Failed to marshal renewal token: %s", err)
			return
		}
		req, err := http.NewRequest("POST", licenseServerURL, bytes.NewReader(jsonedCheckLicenseRequest))
		if err != nil {
			cfg.LastUpdateError = fmt.Sprintf("Failed to make renewal request: %s", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			cfg.LastUpdateError = fmt.Sprintf("Failed to make request: %s", err)
			return
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			// Shoudn't happen
			cfg.LastUpdateError = fmt.Sprintf("Failed to read license server response: %s", err)
			return
		}
		if response.StatusCode != http.StatusOK {
			// TODO - it looks like the license server is sending json docs back - probably want to parse...
			// e.g., "{\"error\": \"This License is invalid\"}"
			cfg.LastUpdateError = string(body)
			return
		}

		// Process the new license we just got
		jsonSignature, err := libtrust.ParseJWS(body)
		if err != nil {
			cfg.LastUpdateError = fmt.Sprintf("Malformed renewal from license server: %s", err)
			return
		}
		checkResponse, err := verifyJsonSignature(cfg.License.PrivateKey, jsonSignature, s.publicKey)
		if err != nil {
			cfg.LastUpdateError = fmt.Sprintf("Malformed license key - Signature verification failure: %s", err)
			return
		}

		// If we get this far, then the license looks good
		cfg.LastUpdateError = ""
		// Base64 encode the content which is the "authorization"
		cfg.License.Authorization = base64.StdEncoding.EncodeToString(body)
		// always clobber the details from the actual parsed license
		cfg.Details.MaxEngines = checkResponse.MaxEngines
		cfg.Details.Expiration = checkResponse.Expiration
	}
}

// Copied ~verbatim from DTR so we don't have to vendor in a mountain of unnecessary dependencies
func verifyJsonSignature(privateKey string, jsonSignature *libtrust.JSONSignature, publicKey libtrust.PublicKey) (*CheckLicenseResponse, error) {
	keys, err := jsonSignature.Verify()
	if err != nil {
		return nil, err
	} else if len(keys) != 1 || keys[0].KeyID() != publicKey.KeyID() {
		return nil, errors.New("Bad signature")
	}

	payload, err := jsonSignature.Payload()
	if err != nil {
		return nil, err
	}

	checkLicenseResponse := new(CheckLicenseResponse)
	if err := json.NewDecoder(bytes.NewReader(payload)).Decode(checkLicenseResponse); err != nil {
		return nil, err
	}

	ok, err := CheckToken(checkLicenseResponse.Expiration.Format(time.RFC3339), checkLicenseResponse.Token, privateKey)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("Invalid token")
	}

	return checkLicenseResponse, nil
}

// Copied from dhe-license-server/requests
type CheckLicenseResponse struct {
	Expiration time.Time `json:"expiration"` // In RFC3339 time format
	Token      string    `json:"token"`
	MaxEngines int       `json:"maxEngines"`
	Type       string    `json:"licenseType"`
	Tier       string    `json:"tier"`
}
type CheckLicenseRequest struct {
	KeyID     string    `json:"keyId"`
	Timestamp time.Time `json:"timestamp"` // In RFC3339 time format
	Token     string    `json:"token"`
}

// Copied from dhe-license-server/verify
func CheckToken(message, token, privateKey string) (bool, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false, err
	}

	generatedToken, err := GenerateToken(message, privateKey)
	if err != nil {
		return false, err
	}

	generatedBytes, err := base64.URLEncoding.DecodeString(generatedToken)
	if err != nil {
		return false, err
	}

	return hmac.Equal(tokenBytes, generatedBytes), nil
}

// Copied from dhe-license-server/verify
func GenerateToken(message, privateKey string) (string, error) {
	key, err := base64.URLEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}
