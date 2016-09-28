package licensing

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/licensing/util"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/dhe-license-server/requests"
	"github.com/docker/dhe-license-server/tiers"
	"github.com/docker/dhe-license-server/verify"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libtrust"
)

var (
	// this is "true", but can be overridden by running make with NOLICENSE=true
	enabled  = "true"
	enforced = !(enabled == "false")
)

type checker struct {
	settingsStore hubconfig.SettingsStore
	privateKey    string
	keyID         string
	licenseLock   sync.Mutex
	autoRefresh   bool
	isValid       bool

	publicKey         libtrust.PublicKey
	signature         *libtrust.JSONSignature
	licenseExpiration time.Time
	licenseToken      string
	licenseTier       string
	hourlyVerifier    HourlyVerifier
}

func NewChecker(settingsStore hubconfig.SettingsStore) hubconfig.LicenseChecker {
	publicKey, _ := certToPubKey(`-----BEGIN CERTIFICATE-----
MIICSzCCAbSgAwIBAgIETJC2dzANBgkqhkiG9w0BAQUFADBqMQswCQYDVQQGEwJV
UzETMBEGA1UECBMKV2FzaGluZ3RvbjEQMA4GA1UEBxMHU2VhdHRsZTEYMBYGA1UE
ChMPQW1hem9uLmNvbSBJbmMuMRowGAYDVQQDExFlYzIuYW1hem9uYXdzLmNvbTAe
Fw0xMDA5MTUxMjA1MTFaFw0xMDEyMTQxMjA1MTFaMGoxCzAJBgNVBAYTAlVTMRMw
EQYDVQQIEwpXYXNoaW5ndG9uMRAwDgYDVQQHEwdTZWF0dGxlMRgwFgYDVQQKEw9B
bWF6b24uY29tIEluYy4xGjAYBgNVBAMTEWVjMi5hbWF6b25hd3MuY29tMIGfMA0G
CSqGSIb3DQEBAQUAA4GNADCBiQKBgQCHvRjf/0kStpJ248khtIaN8qkDN3tkw4Vj
vA9nvPl2anJO+eIBUqPfQG09kZlwpWpmyO8bGB2RWqWxCwuB/dcnIob6w420k9WY
5C0IIGtDRNauN3kuvGXkw3HEnF0EjYr0pcyWUvByWY4KswZV42X7Y7XSS13hOIcL
6NLA+H94/QIDAQABMA0GCSqGSIb3DQEBBQUAA4GBADrBEztYdJwz3bsEwqTsgMRC
IN6EHK95v5/x1DDmzlHV7SH+ok9By/zDlCXRyUctKgxhVXnR0ZIF0tssl7TpAHGs
VFBW4j8M1aZGC0AtTtnhUvc3u6Pu2cf1re4o+I6MmoqJRKcNYh1ejzrYdfx9ebY7
M/apafynee2Beib6aqo6
-----END CERTIFICATE-----`)
	return &checker{
		settingsStore: settingsStore,
		hourlyVerifier: HourlyVerifier{
			configFile:     "/usr/local/etc/dtr/.hourly",
			awsPublicKey:   publicKey,
			awsMetadataURL: "http://169.254.169.254",
			// AWS forced us to disable this check because they are worried about us invalidating the credentials used for it
			awsConfirmProductInstance: false,
			// NO TRAILING SLASH! Otherwise it gives authentication errors.
			awsAPIURL: "https://ec2.amazonaws.com",
			// This AWS account has only one permission: given an instance and product code we own, is this instance legit?
			awsAPIAccount:  "",
			awsAPISecret:   "",
			awsProductCode: "7dpyto4eocr6th6o6hngpa4f3",
		},
	}
}

func (ch *checker) cacheLicense(checkLicenseResponse *requests.CheckLicenseResponse, signature *libtrust.JSONSignature) {
	ch.isValid = true
	ch.signature = signature
	ch.licenseTier = checkLicenseResponse.Tier
	ch.licenseToken = checkLicenseResponse.Token
	ch.licenseExpiration = checkLicenseResponse.Expiration
	_ = ch.settingsStore.SetLicenseStatus(ch.locklessIsValid() || !enforced)
}

func (ch *checker) Initialize() error {
	publicKey, err := util.GetPublicKey()
	if err != nil {
		return err
	} else if publicKey == nil {
		return errors.New("Invalid public key")
	}

	ch.publicKey = publicKey

	return ch.loadLicenseFromEtcd()
}

func (ch *checker) locklessIsValid() bool {
	if ch.LicenseType() == tiers.Hourly {
		return ch.isValid
	} else {
		return ch.isValid && !ch.locklessIsExpired()
	}
}
func (ch *checker) IsValid() bool {
	if ch.LicenseType() == tiers.Hourly {
		return ch.isValid
	} else {
		return ch.isValid && !ch.IsExpired()
	}
}

func (ch *checker) LicensingEnforced() bool {
	return enforced
}

func (ch *checker) BeginLicenseSyncing() {
	go func() {
		for range time.Tick(24 * time.Hour) {
			ch.makeCheck(false)
		}
	}()
}

func (ch *checker) makeCheck(override bool) {
	// sync local state with etcd
	ch.loadLicenseFromEtcd()

	ch.licenseLock.Lock()
	defer ch.licenseLock.Unlock()
	if override || ch.autoRefresh {
		log.WithField("keyId", ch.keyID).Info("Attempting to refresh license")
		// we will get the latest values of the ch struct's keyID and privateKey values
		jsonSignature, checkResponse, err := ch.makeCheckAndVerify(ch.keyID, ch.privateKey, ch.publicKey)
		if err == nil {
			ch.cacheLicense(checkResponse, jsonSignature)
			if err := ch.storeLicenseConfig(); err != nil {
				log.WithField("error", err).Error("Failed to store the license after a periodic license refresh")
			} else {
				log.WithField("KeyID", ch.keyID).Info("Periodic license refresh succeeded")
			}
		} else {
			log.WithField("error", err).Warn("Failed license check")
		}
	}
}

func (ch *checker) IsExpired() bool {
	if ch.LicenseType() == tiers.Hourly {
		return false
	}

	if !ch.isValid {
		return false
	}
	ch.licenseLock.Lock()
	defer ch.licenseLock.Unlock()
	return time.Now().After(ch.licenseExpiration)
}
func (ch *checker) locklessIsExpired() bool {
	if ch.LicenseType() == tiers.Hourly {
		return false
	}

	if !ch.isValid {
		return false
	}
	return time.Now().After(ch.licenseExpiration)
}

func (ch *checker) Expiration() time.Time {
	return ch.licenseExpiration
}

func (ch *checker) GetLicenseID() string {
	return ch.keyID
}

func (ch *checker) LicenseType() string {
	typ, _ := tiers.GetTypeFromTier(ch.licenseTier)
	return typ
}

func (ch *checker) LicenseTier() string {
	tier := ch.licenseTier
	if tier == "" {
		// All licenses that are old enough to not have a tier should be prod
		tier = "Production"
	}
	return tier
}

func (ch *checker) ToggleAutoRefresh(autoRefresh bool) error {
	// Refresh instantly on toggling to true
	if autoRefresh {
		ch.makeCheck(autoRefresh)
	}

	ch.autoRefresh = autoRefresh
	ch.licenseLock.Lock()
	defer ch.licenseLock.Unlock()
	return ch.storeLicenseConfig()
}

func (ch *checker) ChangeLicenseFromId(keyID, privateKey string) error {
	log.WithField("keyID", keyID).Info("Attempting to change License")

	ch.licenseLock.Lock()
	defer ch.licenseLock.Unlock()
	jsonSignature, checkResponse, err := ch.makeCheckAndVerify(keyID, privateKey, ch.publicKey)
	if err != nil {
		return err
	}

	ch.privateKey = privateKey
	ch.keyID = keyID
	ch.cacheLicense(checkResponse, jsonSignature)

	return ch.storeLicenseConfig()
}

func (ch *checker) loadLicenseFromEtcd() error {
	config, err := ch.settingsStore.LicenseConfig()
	if err != nil {
		return err
	} else if config == nil {
		return errors.New("Invalid license config")
	}

	return ch.LoadLicenseFromConfig(config, false)
}

func (ch *checker) LoadLicenseFromConfig(config *hubconfig.LicenseConfig, newLicense bool) error {
	authorizationBytes, err := base64.StdEncoding.DecodeString(config.Authorization)
	if err != nil {
		return err
	}

	jsonSignature, err := libtrust.ParseJWS(authorizationBytes)
	if err != nil {
		return err
	}

	checkResponse, err := util.VerifyJsonSignature(config.PrivateKey, jsonSignature, ch.publicKey)
	if err != nil {
		return err
	}

	ch.licenseLock.Lock()
	defer ch.licenseLock.Unlock()
	ch.privateKey = config.PrivateKey
	ch.keyID = config.KeyID

	if newLicense && checkResponse.Type == tiers.Online {
		ch.autoRefresh = true
	} else {
		ch.autoRefresh = config.AutoRefresh
	}

	if newLicense && ch.autoRefresh {
		// try to update the license if it's new and we have autorefresh on just so that we have the latest version
		log.WithField("keyId", config.KeyID).Info("Attempting to refresh license after config change")
		onlineSignature, onlineResponse, err := ch.makeCheckAndVerify(config.KeyID, config.PrivateKey, ch.publicKey)
		if err == nil {
			jsonSignature = onlineSignature
			checkResponse = onlineResponse
			log.Info("Online license refresh succeeded")
		} else {
			log.WithField("error", err).Warn("Failed to update license")
		}
	}

	ch.cacheLicense(checkResponse, jsonSignature)
	if newLicense {
		return ch.storeLicenseConfig()
	}

	return nil
}

func (ch *checker) makeCheckAndVerify(keyID, privateKey string, publicKey libtrust.PublicKey) (*libtrust.JSONSignature, *requests.CheckLicenseResponse, error) {
	response, err := ch.makeCheckLicenseRequest(keyID, privateKey)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()
	jsonSignatureBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	jsonSignature, err := libtrust.ParseJWS(jsonSignatureBytes)
	if err != nil {
		return nil, nil, err
	}

	checkResponse, err := util.VerifyJsonSignature(privateKey, jsonSignature, publicKey)
	if err != nil {
		return nil, nil, err
	}

	if time.Now().After(checkResponse.Expiration) {
		return nil, nil, errors.New("License is expired")
	}

	return jsonSignature, checkResponse, nil
}

func (ch *checker) makeCheckLicenseRequest(keyID, privateKey string) (*http.Response, error) {
	now := time.Now()

	token, err := verify.GenerateToken(now.Format(time.RFC3339), privateKey)
	if err != nil {
		return nil, err
	}

	checkLicenseRequest := &requests.CheckLicenseRequest{KeyID: keyID, Timestamp: now, Token: token}
	jsonedCheckLicenseRequest, err := json.Marshal(checkLicenseRequest)
	if err != nil {
		return nil, err
	}

	hubConfig, err := ch.settingsStore.UserHubConfig()
	if err != nil {
		return nil, err
	}

	releaseChannel := deploy.DefaultReleaseChannel
	if hubConfig != nil && hubConfig.ReleaseChannel != "" {
		releaseChannel = hubConfig.ReleaseChannel
	}

	licenseServerURL := "https://license.enterprise.docker.com"
	if deploy.ParseReleaseChannel(releaseChannel).Namespace != "docker" {
		licenseServerURL = "https://license-staging.enterprise.docker.com"
	}

	req, err := http.NewRequest("POST", licenseServerURL+"/v1/check", bytes.NewReader(jsonedCheckLicenseRequest))
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var response *http.Response
	err = dtrutil.Poll(time.Millisecond*100, 3, func() error {
		response, err = dtrutil.DoRequest(req)
		if err != nil {
			log.Warn(err)
			return err
		} else if response.StatusCode != http.StatusOK {
			response.Body.Close()
			log.Warn(errors.New(response.Status))
			return errors.New(response.Status)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

// You should only call this if you already have a ch.Lock()
func (ch *checker) storeLicenseConfig() error {
	authorizationBytes, err := ch.signature.JWS()
	if err != nil {
		return err
	}

	authorization := base64.StdEncoding.EncodeToString(authorizationBytes)

	config := hubconfig.LicenseConfig{
		KeyID:         ch.keyID,
		PrivateKey:    ch.privateKey,
		AutoRefresh:   ch.autoRefresh,
		Authorization: authorization,
	}

	return ch.settingsStore.SetLicenseConfig(&config)
}
