package manager

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/docker/libtrust"
)

func getBlankLCS(t *testing.T) LicenseConfigSubsystem {
	pemBytes, err := base64.StdEncoding.DecodeString(publicRSAKey)
	if err != nil {
		// should not happen
		t.Fatalf("Failed to parse the embedded license public key during init: %s", err)
	}
	publicKey, err := libtrust.UnmarshalPublicKeyPEM(pemBytes)
	if err != nil {
		// should not happen
		t.Fatalf("Failed to unmarshal the embedded license public key during init: %s", err)
	}
	count := 0
	return LicenseConfigSubsystem{
		m:                 nil, // Can't use mock_test, import cycles :-(
		cfg:               &LicenseSubsystemConfig{},
		publicKey:         publicKey,
		recentEngineCount: &count,
	}
}

func TestLicensingValidateConfigEmpty(t *testing.T) {
	s := getBlankLCS(t)

	cfg, err := s.ValidateConfig("", false)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if cfg.(LicenseSubsystemConfig).License.KeyID != "unlicensed" {
		t.Fatalf("Expected %s found %#v", "unlicensed", cfg)
	}

}

func TestLicensingValidateConfigBadJson(t *testing.T) {
	s := getBlankLCS(t)

	expected := "Malformed license configuration"
	_, err := s.ValidateConfig("not json", false)
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("Expected %s found %s", expected, err.Error())
	}

}

func TestLicensingValidateConfigUnlicensed(t *testing.T) {
	s := getBlankLCS(t)

	_, err := s.ValidateConfig(`{"auto_refresh":true,"license_config":{"key_id":"unlicensed"}}`, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLicensingValidateConfigMalformedAuthField(t *testing.T) {
	s := getBlankLCS(t)

	expected := "Failed to decode key authorization:"
	_, err := s.ValidateConfig(`{"auto_refresh":true,"license_config":{"key_id":"S3E_XoqN7B3JI0nNfHAfYL44avPUv8KiULVLQNvN499t", "authorization":"not base64 data"}}`, false)
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("Expected %s found %s", expected, err.Error())
	}
}

func TestLicensingValidateConfigMalformedAuthFieldNotJWS(t *testing.T) {
	s := getBlankLCS(t)

	expected := "Failed to parse JWT:"
	_, err := s.ValidateConfig(`{"auto_refresh":true,"license_config":{"key_id":"S3E_XoqN7B3JI0nNfHAfYL44avPUv8KiULVLQNvN499t", "authorization":"SGVsbG8gd29ybGQK"}}`, false) // "Hello world"
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("Expected %s found %s", expected, err.Error())
	}
}

func TestLicensingValidateConfigFailVerification(t *testing.T) {
	s := getBlankLCS(t)

	expected := "Signature verification failure:"
	// bad key ID, but valid jwt
	_, err := s.ValidateConfig(`{"auto_refresh":true,"license_config":{"key_id":"SGVsbG8gd29ybGQK", "authorization":"ewogICAicGF5bG9hZCI6ICJleUpsZUhCcGNtRjBhVzl1SWpvaU1qQXhOUzB4TWkweE1sUXlNRG95TlRveU1DNHhPRE14TXpZeU9Wb2lMQ0owYjJ0bGJpSTZJbkpMWnpoc2MwcFpRVVpRVHpaS1kyaDFlWFF0V1VwaVdrVnJabXRsVUhGcmRWWmtZMmRzV0hOcGJFVTlJaXdpYldGNFJXNW5hVzVsY3lJNk1Td2liR2xqWlc1elpWUjVjR1VpT2lKUGJteHBibVVpTENKMGFXVnlJam9pVkhKcFlXd2lmUSIsCiAgICJzaWduYXR1cmVzIjogWwogICAgICB7CiAgICAgICAgICJoZWFkZXIiOiB7CiAgICAgICAgICAgICJqd2siOiB7CiAgICAgICAgICAgICAgICJlIjogIkFRQUIiLAogICAgICAgICAgICAgICAia2V5SUQiOiAiSjdMRDo2N1ZSOkw1SFo6VTdCQToyTzRHOjRBTDM6T0YyTjpKSEdCOkVGVEg6NUNWUTpNRkVPOkFFSVQiLAogICAgICAgICAgICAgICAia2lkIjogIko3TEQ6NjdWUjpMNUhaOlU3QkE6Mk80Rzo0QUwzOk9GMk46SkhHQjpFRlRIOjVDVlE6TUZFTzpBRUlUIiwKICAgICAgICAgICAgICAgImt0eSI6ICJSU0EiLAogICAgICAgICAgICAgICAibiI6ICJ5ZEl5LWxVN283UGNlWS00LXMtQ1E1T0VnQ3lGOEN4SWNRSVd1Szg0cElpWmNpWTY3MzB5Q1lud0xTS1Rsdy1VNlVDX1FSZVdSaW9NTk5FNURzNVRZRVhiR0c2b2xtMnFkV2JCd2NDZy0yVVVIX09jQjlXdVA2Z1JQSHBNRk1zeER6V3d2YXk4SlV1SGdZVUxVcG0xSXYtbXE3bHA1blFfUnhyVDBLWlJBUVRZTEVNRWZHd20zaE1PX2dlTFBTLWhnS1B0SUhsa2c2X1djb3hUR29LUDc5ZF93YUhZeEdObDdXaFNuZWlCU3hicGJRQUtrMjFsZzc5OFhiN3ZaeUVBVERNclJSOU1lRTZBZGo1SEpwWTNDb3lSQVBDbWFLR1JDSzR1b1pTb0l1MGhGVmxLVVB5YmJ3MDAwR08td2EyS044VXdnSUltMGk1STF1VzlHa3E0empCeTV6aGdxdVVYYkc5YldQQU9ZcnE1UWE4MUR4R2NCbEp5SFlBcC1ERFBFOVRHZzR6WW1YakpueFpxSEVkdUdxZGV2WjhYTUkwdWtma0dJSTE0d1VPaU1JSUlyWGxFY0JmXzQ2SThnUVdEenh5Y1plX0pHWC1MQXVheVhyeXJVRmVoVk5VZFpVbDl3WE5hSkIta2FDcXo1UXdhUjkzc0d3LVFTZnREME52TGU3Q3lPSC1FNnZnNlN0X05lVHZndjhZbmhDaVhJbFo4SE9mSXdOZTd0RUZfVWN6NU9iUHlrbTN0eWxyTlVqdDBWeUFtdHRhY1ZJMmlHaWhjVVBybWs0bFZJWjdWRF9MU1ctaTd5b1N1cnRwc1BYY2UycEtESW8zMGxKR2hPXzNLVW1sMlNVWkNxekoxeUVtS3B5c0g1SERXOWNzSUZDQTNkZUFqZlpVdk43VSIKICAgICAgICAgICAgfSwKICAgICAgICAgICAgImFsZyI6ICJSUzI1NiIKICAgICAgICAgfSwKICAgICAgICAgInNpZ25hdHVyZSI6ICJpdmNHS0dPZVhxb2U3cGZoM2hiVk5qZ2phTW1mc3hqOWJSbEVNTGptSjdwTW1STVhRdTBxOHhDNi13cC1tZGRLSi1IX2xyejUyMk1BLV9ueUk3Q005dlNRZ2UtV3hZV3hDN0QyMXFFckpwNUNFYnRGT3E3cGpGWFQzR2EyWHZLUUJPdjV2TUh4SzVXVDNTZHAyY2VYbmlLMTZLT3ByWGpxbVQyOWRjWkk4Z2tFZUJTY19CUjJLQWZHTmpRU1BMVi1LVFJHQ1djd29CeDR5bHZxei1kdnRTaUlvYm1YLU5fajNQZGxXRGxUajlwRGJTV3JQT2VsMnZUZ3JwMnJfVTBoYUNKMDdaZS1oV3BOTmpER0RxV1hzS0kwdjJDWDdNZFVWMFpTck5ZOUJNNzdhejRjbUhlMW5Qal9fRDJ3b1I0b0sxWTA0WVF5N2R2bGFmYlp5dW1ib1dLU2l0QnY0TGtOUm5Sdzl0dlhxSHEyRTBQNGFldjQ2dm5jU3MxdXFqNTc0ZUtPci1HWWVVWEdNY1M5ZXVQN2ctb25HOXIzRkYxWmc3UjFwYzgxX09pSVVkc0NaX0Z5YVZCQ05CWEZCYXRUSjNiQmFoQzNlWWJrYWY1RHJtR1NBUF9mYTExQkpBNGlGZVBka0VxR3pzYUlKZzlXNXcxZjJxeGNzT2Q1WU5jOEIzU3E0cGpMX3NiaDlxT1VXS2RWRlJwY0pGSk1YdWJoTktWOXlQUllZYTF0ZkthOW51bUhRd3E1SHlRUjJRSzg0RThpNjFXNmRfTGg1TGhUbUcxSTYydmVwd2FnbXV0TUZFZ0F3ODR4ZVFZZWtYWnRVTUlrQVI0T3Q5dkJWczRKUFdZc2t3Qm96dk1TdWVEcFpiSWp3OVg4N0UtNlpNdS1xOHc2eFpSdmliQSIsCiAgICAgICAgICJwcm90ZWN0ZWQiOiAiZXlKbWIzSnRZWFJNWlc1bmRHZ2lPakUxTXl3aVptOXliV0YwVkdGcGJDSTZJbVpSSWl3aWRHbHRaU0k2SWpJd01UVXRNVEl0TVRGVU1qQTZNalU2TWpCYUluMCIKICAgICAgfQogICBdCn0="}}`, false)
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("Expected %s found %s", expected, err.Error())
	}
}

func TestLicensingValidateConfigGood(t *testing.T) {
	s := getBlankLCS(t)

	cfg, err := s.ValidateConfig(`{"auto_refresh":true,"license_config":{"key_id": "FFUFtz4bBM76ds9vA7rJHmrZz8FOUEUBZVwWCGK-qKtU", "private_key":"S3E_XoqN7B3JI0nNfHAfYL44avPUv8KiULVLQNvN499t", "authorization":"ewogICAicGF5bG9hZCI6ICJleUpsZUhCcGNtRjBhVzl1SWpvaU1qQXhOUzB4TWkweE1sUXlNRG95TlRveU1DNHhPRE14TXpZeU9Wb2lMQ0owYjJ0bGJpSTZJbkpMWnpoc2MwcFpRVVpRVHpaS1kyaDFlWFF0V1VwaVdrVnJabXRsVUhGcmRWWmtZMmRzV0hOcGJFVTlJaXdpYldGNFJXNW5hVzVsY3lJNk1Td2liR2xqWlc1elpWUjVjR1VpT2lKUGJteHBibVVpTENKMGFXVnlJam9pVkhKcFlXd2lmUSIsCiAgICJzaWduYXR1cmVzIjogWwogICAgICB7CiAgICAgICAgICJoZWFkZXIiOiB7CiAgICAgICAgICAgICJqd2siOiB7CiAgICAgICAgICAgICAgICJlIjogIkFRQUIiLAogICAgICAgICAgICAgICAia2V5SUQiOiAiSjdMRDo2N1ZSOkw1SFo6VTdCQToyTzRHOjRBTDM6T0YyTjpKSEdCOkVGVEg6NUNWUTpNRkVPOkFFSVQiLAogICAgICAgICAgICAgICAia2lkIjogIko3TEQ6NjdWUjpMNUhaOlU3QkE6Mk80Rzo0QUwzOk9GMk46SkhHQjpFRlRIOjVDVlE6TUZFTzpBRUlUIiwKICAgICAgICAgICAgICAgImt0eSI6ICJSU0EiLAogICAgICAgICAgICAgICAibiI6ICJ5ZEl5LWxVN283UGNlWS00LXMtQ1E1T0VnQ3lGOEN4SWNRSVd1Szg0cElpWmNpWTY3MzB5Q1lud0xTS1Rsdy1VNlVDX1FSZVdSaW9NTk5FNURzNVRZRVhiR0c2b2xtMnFkV2JCd2NDZy0yVVVIX09jQjlXdVA2Z1JQSHBNRk1zeER6V3d2YXk4SlV1SGdZVUxVcG0xSXYtbXE3bHA1blFfUnhyVDBLWlJBUVRZTEVNRWZHd20zaE1PX2dlTFBTLWhnS1B0SUhsa2c2X1djb3hUR29LUDc5ZF93YUhZeEdObDdXaFNuZWlCU3hicGJRQUtrMjFsZzc5OFhiN3ZaeUVBVERNclJSOU1lRTZBZGo1SEpwWTNDb3lSQVBDbWFLR1JDSzR1b1pTb0l1MGhGVmxLVVB5YmJ3MDAwR08td2EyS044VXdnSUltMGk1STF1VzlHa3E0empCeTV6aGdxdVVYYkc5YldQQU9ZcnE1UWE4MUR4R2NCbEp5SFlBcC1ERFBFOVRHZzR6WW1YakpueFpxSEVkdUdxZGV2WjhYTUkwdWtma0dJSTE0d1VPaU1JSUlyWGxFY0JmXzQ2SThnUVdEenh5Y1plX0pHWC1MQXVheVhyeXJVRmVoVk5VZFpVbDl3WE5hSkIta2FDcXo1UXdhUjkzc0d3LVFTZnREME52TGU3Q3lPSC1FNnZnNlN0X05lVHZndjhZbmhDaVhJbFo4SE9mSXdOZTd0RUZfVWN6NU9iUHlrbTN0eWxyTlVqdDBWeUFtdHRhY1ZJMmlHaWhjVVBybWs0bFZJWjdWRF9MU1ctaTd5b1N1cnRwc1BYY2UycEtESW8zMGxKR2hPXzNLVW1sMlNVWkNxekoxeUVtS3B5c0g1SERXOWNzSUZDQTNkZUFqZlpVdk43VSIKICAgICAgICAgICAgfSwKICAgICAgICAgICAgImFsZyI6ICJSUzI1NiIKICAgICAgICAgfSwKICAgICAgICAgInNpZ25hdHVyZSI6ICJpdmNHS0dPZVhxb2U3cGZoM2hiVk5qZ2phTW1mc3hqOWJSbEVNTGptSjdwTW1STVhRdTBxOHhDNi13cC1tZGRLSi1IX2xyejUyMk1BLV9ueUk3Q005dlNRZ2UtV3hZV3hDN0QyMXFFckpwNUNFYnRGT3E3cGpGWFQzR2EyWHZLUUJPdjV2TUh4SzVXVDNTZHAyY2VYbmlLMTZLT3ByWGpxbVQyOWRjWkk4Z2tFZUJTY19CUjJLQWZHTmpRU1BMVi1LVFJHQ1djd29CeDR5bHZxei1kdnRTaUlvYm1YLU5fajNQZGxXRGxUajlwRGJTV3JQT2VsMnZUZ3JwMnJfVTBoYUNKMDdaZS1oV3BOTmpER0RxV1hzS0kwdjJDWDdNZFVWMFpTck5ZOUJNNzdhejRjbUhlMW5Qal9fRDJ3b1I0b0sxWTA0WVF5N2R2bGFmYlp5dW1ib1dLU2l0QnY0TGtOUm5Sdzl0dlhxSHEyRTBQNGFldjQ2dm5jU3MxdXFqNTc0ZUtPci1HWWVVWEdNY1M5ZXVQN2ctb25HOXIzRkYxWmc3UjFwYzgxX09pSVVkc0NaX0Z5YVZCQ05CWEZCYXRUSjNiQmFoQzNlWWJrYWY1RHJtR1NBUF9mYTExQkpBNGlGZVBka0VxR3pzYUlKZzlXNXcxZjJxeGNzT2Q1WU5jOEIzU3E0cGpMX3NiaDlxT1VXS2RWRlJwY0pGSk1YdWJoTktWOXlQUllZYTF0ZkthOW51bUhRd3E1SHlRUjJRSzg0RThpNjFXNmRfTGg1TGhUbUcxSTYydmVwd2FnbXV0TUZFZ0F3ODR4ZVFZZWtYWnRVTUlrQVI0T3Q5dkJWczRKUFdZc2t3Qm96dk1TdWVEcFpiSWp3OVg4N0UtNlpNdS1xOHc2eFpSdmliQSIsCiAgICAgICAgICJwcm90ZWN0ZWQiOiAiZXlKbWIzSnRZWFJNWlc1bmRHZ2lPakUxTXl3aVptOXliV0YwVkdGcGJDSTZJbVpSSWl3aWRHbHRaU0k2SWpJd01UVXRNVEl0TVRGVU1qQTZNalU2TWpCYUluMCIKICAgICAgfQogICBdCn0="}}`, false)
	if err != nil {
		t.Fatal(err)
	}
	err = s.UpdateConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Spot check the expiration/engines for that trial license
	expectedMaxEngines := 1
	if s.cfg.Details.MaxEngines != expectedMaxEngines {
		t.Fatalf("Expected %d max engines, found %d", expectedMaxEngines, s.cfg.Details.MaxEngines)
	}
	expectedExpiration := time.Date(2015, 12, 12, 20, 25, 20, 183136290, time.UTC)
	if s.cfg.Details.Expiration != expectedExpiration {
		t.Fatalf("Expected %v Expiration, found %v", expectedExpiration, s.cfg.Details.Expiration)
	}
}

func TestLicensingGetLicenseKeyID(t *testing.T) {
	expected := "1234"
	m := &DefaultManager{
		configSubsystems: map[string]ConfigSubsystem{
			singletonLicenceKvKey: LicenseConfigSubsystem{
				cfg: &LicenseSubsystemConfig{
					AutoRefresh: true,
					License: LicenseConfig{
						KeyID: expected,
					},
					Details: LicenseDetails{
						Expiration: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	res := m.GetLicenseKeyID()
	if res != expected {
		t.Fatalf("Expected %s got %s", expected, res)
	}
}

func TestLicensingGetBannerUnlicensed(t *testing.T) {
	s := getBlankLCS(t)
	s.cfg.License.KeyID = UnlicensedID
	m := &DefaultManager{
		configSubsystems: make(map[string]ConfigSubsystem),
	}
	m.configSubsystems[singletonLicenceKvKey] = s

	res := m.getLicenseBanners()
	if len(res) != 1 {
		t.Fatalf("Expected 1 banner, got %d", len(res))
	}
	expected := "Your system is unlicensed"
	if !strings.Contains(res[0].Message, expected) {
		t.Fatalf("Expected %s got %s", expected, res[0].Message)
	}
}

func TestLicensingGetBannerExpiredAutoRefresh(t *testing.T) {
	s := getBlankLCS(t)
	s.cfg.License.KeyID = UnlicensedID
	m := &DefaultManager{
		configSubsystems: make(map[string]ConfigSubsystem),
	}
	m.configSubsystems[singletonLicenceKvKey] = s

	s.cfg.License.KeyID = "somelicensekey"
	s.cfg.AutoRefresh = true
	s.cfg.Details.Expiration = time.Now().UTC().Add(-30 * time.Second)

	res := m.getLicenseBanners()
	if len(res) != 1 {
		t.Fatalf("Expected 1 banner, got %d", len(res))
	}
	expected := "been unable to auto renew"
	if !strings.Contains(res[0].Message, expected) {
		t.Fatalf("Expected %s got %s", expected, res[0].Message)
	}
}

func TestLicensingGetBannerExpiredManualRefresh(t *testing.T) {
	s := getBlankLCS(t)
	s.cfg.License.KeyID = UnlicensedID
	m := &DefaultManager{
		configSubsystems: make(map[string]ConfigSubsystem),
	}
	m.configSubsystems[singletonLicenceKvKey] = s

	s.cfg.License.KeyID = "somelicensekey"
	s.cfg.Details.Expiration = time.Now().UTC().Add(-30 * time.Second)

	res := m.getLicenseBanners()
	if len(res) != 1 {
		t.Fatalf("Expected 1 banner, got %d", len(res))
	}
	expected := "to download a renewed license"
	if !strings.Contains(res[0].Message, expected) {
		t.Fatalf("Expected %s got %s", expected, res[0].Message)
	}
}

func TestLicensingGetBannerTooManyEngines(t *testing.T) {
	s := getBlankLCS(t)
	s.cfg.License.KeyID = UnlicensedID
	m := &DefaultManager{
		configSubsystems: make(map[string]ConfigSubsystem),
	}
	m.configSubsystems[singletonLicenceKvKey] = s

	s.cfg.License.KeyID = "somelicensekey"
	s.cfg.Details.Expiration = time.Now().UTC().Add(30 * time.Second)
	*s.recentEngineCount = 100

	res := m.getLicenseBanners()
	if len(res) != 1 {
		t.Fatalf("Expected 1 banner, got %d", len(res))
	}
	expected := "exceeding your licensed engine count"
	if !strings.Contains(res[0].Message, expected) {
		t.Fatalf("Expected %s got %s", expected, res[0].Message)
	}
}

func TestLicensingGetBannerAllGood(t *testing.T) {
	s := getBlankLCS(t)
	s.cfg.License.KeyID = UnlicensedID
	m := &DefaultManager{
		configSubsystems: make(map[string]ConfigSubsystem),
	}
	m.configSubsystems[singletonLicenceKvKey] = s

	s.cfg.License.KeyID = "somelicensekey"
	s.cfg.Details.Expiration = time.Now().UTC().Add(30 * time.Second)
	s.cfg.Details.MaxEngines = 10
	*s.recentEngineCount = 1

	res := m.getLicenseBanners()
	if len(res) != 0 {
		t.Fatalf("Expected 0 banner, got %d - %s", len(res), res[0].Message)
	}
}

func TestLicensingGetLicense(t *testing.T) {
	s := getBlankLCS(t)
	s.cfg.License.KeyID = UnlicensedID
	m := &DefaultManager{
		configSubsystems: make(map[string]ConfigSubsystem),
	}
	m.configSubsystems[singletonLicenceKvKey] = s

	res := m.GetLicense()
	if res != *s.cfg {
		t.Fatalf("Expected %#v got %#v", *s.cfg, res)
	}
}

func TestLicensingGetLicenseCheckTier(t *testing.T) {
	s := getBlankLCS(t)
	cfg, err := s.ValidateConfig(`{"auto_refresh":true,"license_config":{"key_id": "FFUFtz4bBM76ds9vA7rJHmrZz8FOUEUBZVwWCGK-qKtU", "private_key":"S3E_XoqN7B3JI0nNfHAfYL44avPUv8KiULVLQNvN499t", "authorization":"ewogICAicGF5bG9hZCI6ICJleUpsZUhCcGNtRjBhVzl1SWpvaU1qQXhOUzB4TWkweE1sUXlNRG95TlRveU1DNHhPRE14TXpZeU9Wb2lMQ0owYjJ0bGJpSTZJbkpMWnpoc2MwcFpRVVpRVHpaS1kyaDFlWFF0V1VwaVdrVnJabXRsVUhGcmRWWmtZMmRzV0hOcGJFVTlJaXdpYldGNFJXNW5hVzVsY3lJNk1Td2liR2xqWlc1elpWUjVjR1VpT2lKUGJteHBibVVpTENKMGFXVnlJam9pVkhKcFlXd2lmUSIsCiAgICJzaWduYXR1cmVzIjogWwogICAgICB7CiAgICAgICAgICJoZWFkZXIiOiB7CiAgICAgICAgICAgICJqd2siOiB7CiAgICAgICAgICAgICAgICJlIjogIkFRQUIiLAogICAgICAgICAgICAgICAia2V5SUQiOiAiSjdMRDo2N1ZSOkw1SFo6VTdCQToyTzRHOjRBTDM6T0YyTjpKSEdCOkVGVEg6NUNWUTpNRkVPOkFFSVQiLAogICAgICAgICAgICAgICAia2lkIjogIko3TEQ6NjdWUjpMNUhaOlU3QkE6Mk80Rzo0QUwzOk9GMk46SkhHQjpFRlRIOjVDVlE6TUZFTzpBRUlUIiwKICAgICAgICAgICAgICAgImt0eSI6ICJSU0EiLAogICAgICAgICAgICAgICAibiI6ICJ5ZEl5LWxVN283UGNlWS00LXMtQ1E1T0VnQ3lGOEN4SWNRSVd1Szg0cElpWmNpWTY3MzB5Q1lud0xTS1Rsdy1VNlVDX1FSZVdSaW9NTk5FNURzNVRZRVhiR0c2b2xtMnFkV2JCd2NDZy0yVVVIX09jQjlXdVA2Z1JQSHBNRk1zeER6V3d2YXk4SlV1SGdZVUxVcG0xSXYtbXE3bHA1blFfUnhyVDBLWlJBUVRZTEVNRWZHd20zaE1PX2dlTFBTLWhnS1B0SUhsa2c2X1djb3hUR29LUDc5ZF93YUhZeEdObDdXaFNuZWlCU3hicGJRQUtrMjFsZzc5OFhiN3ZaeUVBVERNclJSOU1lRTZBZGo1SEpwWTNDb3lSQVBDbWFLR1JDSzR1b1pTb0l1MGhGVmxLVVB5YmJ3MDAwR08td2EyS044VXdnSUltMGk1STF1VzlHa3E0empCeTV6aGdxdVVYYkc5YldQQU9ZcnE1UWE4MUR4R2NCbEp5SFlBcC1ERFBFOVRHZzR6WW1YakpueFpxSEVkdUdxZGV2WjhYTUkwdWtma0dJSTE0d1VPaU1JSUlyWGxFY0JmXzQ2SThnUVdEenh5Y1plX0pHWC1MQXVheVhyeXJVRmVoVk5VZFpVbDl3WE5hSkIta2FDcXo1UXdhUjkzc0d3LVFTZnREME52TGU3Q3lPSC1FNnZnNlN0X05lVHZndjhZbmhDaVhJbFo4SE9mSXdOZTd0RUZfVWN6NU9iUHlrbTN0eWxyTlVqdDBWeUFtdHRhY1ZJMmlHaWhjVVBybWs0bFZJWjdWRF9MU1ctaTd5b1N1cnRwc1BYY2UycEtESW8zMGxKR2hPXzNLVW1sMlNVWkNxekoxeUVtS3B5c0g1SERXOWNzSUZDQTNkZUFqZlpVdk43VSIKICAgICAgICAgICAgfSwKICAgICAgICAgICAgImFsZyI6ICJSUzI1NiIKICAgICAgICAgfSwKICAgICAgICAgInNpZ25hdHVyZSI6ICJpdmNHS0dPZVhxb2U3cGZoM2hiVk5qZ2phTW1mc3hqOWJSbEVNTGptSjdwTW1STVhRdTBxOHhDNi13cC1tZGRLSi1IX2xyejUyMk1BLV9ueUk3Q005dlNRZ2UtV3hZV3hDN0QyMXFFckpwNUNFYnRGT3E3cGpGWFQzR2EyWHZLUUJPdjV2TUh4SzVXVDNTZHAyY2VYbmlLMTZLT3ByWGpxbVQyOWRjWkk4Z2tFZUJTY19CUjJLQWZHTmpRU1BMVi1LVFJHQ1djd29CeDR5bHZxei1kdnRTaUlvYm1YLU5fajNQZGxXRGxUajlwRGJTV3JQT2VsMnZUZ3JwMnJfVTBoYUNKMDdaZS1oV3BOTmpER0RxV1hzS0kwdjJDWDdNZFVWMFpTck5ZOUJNNzdhejRjbUhlMW5Qal9fRDJ3b1I0b0sxWTA0WVF5N2R2bGFmYlp5dW1ib1dLU2l0QnY0TGtOUm5Sdzl0dlhxSHEyRTBQNGFldjQ2dm5jU3MxdXFqNTc0ZUtPci1HWWVVWEdNY1M5ZXVQN2ctb25HOXIzRkYxWmc3UjFwYzgxX09pSVVkc0NaX0Z5YVZCQ05CWEZCYXRUSjNiQmFoQzNlWWJrYWY1RHJtR1NBUF9mYTExQkpBNGlGZVBka0VxR3pzYUlKZzlXNXcxZjJxeGNzT2Q1WU5jOEIzU3E0cGpMX3NiaDlxT1VXS2RWRlJwY0pGSk1YdWJoTktWOXlQUllZYTF0ZkthOW51bUhRd3E1SHlRUjJRSzg0RThpNjFXNmRfTGg1TGhUbUcxSTYydmVwd2FnbXV0TUZFZ0F3ODR4ZVFZZWtYWnRVTUlrQVI0T3Q5dkJWczRKUFdZc2t3Qm96dk1TdWVEcFpiSWp3OVg4N0UtNlpNdS1xOHc2eFpSdmliQSIsCiAgICAgICAgICJwcm90ZWN0ZWQiOiAiZXlKbWIzSnRZWFJNWlc1bmRHZ2lPakUxTXl3aVptOXliV0YwVkdGcGJDSTZJbVpSSWl3aWRHbHRaU0k2SWpJd01UVXRNVEl0TVRGVU1qQTZNalU2TWpCYUluMCIKICAgICAgfQogICBdCn0="}}`, false)
	if err != nil {
		t.Fatal(err)
	}
	err = s.UpdateConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if s.cfg.Details.Tier != "Trial" {
		t.Fatalf("expected trial license; received %s", s.cfg.Details.Tier)
	}
}
