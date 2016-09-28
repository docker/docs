package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/docker/dhe-deploy/hubconfig"

	"github.com/docker/dhe-license-server/requests"
	"github.com/docker/dhe-license-server/verify"
	"github.com/docker/libtrust"
)

const (
	publicRSAKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0Ka2lkOiBKN0xEOjY3VlI6TDVIWjpVN0JBOjJPNEc6NEFMMzpPRjJOOkpIR0I6RUZUSDo1Q1ZROk1GRU86QUVJVAoKTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUF5ZEl5K2xVN283UGNlWSs0K3MrQwpRNU9FZ0N5RjhDeEljUUlXdUs4NHBJaVpjaVk2NzMweUNZbndMU0tUbHcrVTZVQy9RUmVXUmlvTU5ORTVEczVUCllFWGJHRzZvbG0ycWRXYkJ3Y0NnKzJVVUgvT2NCOVd1UDZnUlBIcE1GTXN4RHpXd3ZheThKVXVIZ1lVTFVwbTEKSXYrbXE3bHA1blEvUnhyVDBLWlJBUVRZTEVNRWZHd20zaE1PL2dlTFBTK2hnS1B0SUhsa2c2L1djb3hUR29LUAo3OWQvd2FIWXhHTmw3V2hTbmVpQlN4YnBiUUFLazIxbGc3OThYYjd2WnlFQVRETXJSUjlNZUU2QWRqNUhKcFkzCkNveVJBUENtYUtHUkNLNHVvWlNvSXUwaEZWbEtVUHliYncwMDBHTyt3YTJLTjhVd2dJSW0waTVJMXVXOUdrcTQKempCeTV6aGdxdVVYYkc5YldQQU9ZcnE1UWE4MUR4R2NCbEp5SFlBcCtERFBFOVRHZzR6WW1YakpueFpxSEVkdQpHcWRldlo4WE1JMHVrZmtHSUkxNHdVT2lNSUlJclhsRWNCZi80Nkk4Z1FXRHp4eWNaZS9KR1grTEF1YXlYcnlyClVGZWhWTlVkWlVsOXdYTmFKQitrYUNxejVRd2FSOTNzR3crUVNmdEQwTnZMZTdDeU9IK0U2dmc2U3QvTmVUdmcKdjhZbmhDaVhJbFo4SE9mSXdOZTd0RUYvVWN6NU9iUHlrbTN0eWxyTlVqdDBWeUFtdHRhY1ZJMmlHaWhjVVBybQprNGxWSVo3VkQvTFNXK2k3eW9TdXJ0cHNQWGNlMnBLRElvMzBsSkdoTy8zS1VtbDJTVVpDcXpKMXlFbUtweXNICjVIRFc5Y3NJRkNBM2RlQWpmWlV2TjdVQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo="
)

func IsValidFromLicenseConfig(config *hubconfig.LicenseConfig) (bool, error) {
	if config == nil {
		return false, nil
	}

	// TODO Lock the license lock in the kv here!!!!
	authorizationBytes, err := base64.StdEncoding.DecodeString(config.Authorization)
	if err != nil {
		return false, err
	}

	jsonSignature, err := libtrust.ParseJWS(authorizationBytes)
	if err != nil {
		return false, err
	}

	publicKey, err := GetPublicKey()
	if err != nil {
		return false, err
	} else if publicKey == nil {
		return false, errors.New("Invalid public key")
	}

	checkResponse, err := VerifyJsonSignature(config.PrivateKey, jsonSignature, publicKey)
	if err != nil {
		return false, err
	}

	licenseExpiration := checkResponse.Expiration
	return time.Now().Before(licenseExpiration), nil
}

func GetPublicKey() (libtrust.PublicKey, error) {
	pemBytes, err := base64.StdEncoding.DecodeString(publicRSAKey)
	if err != nil {
		return nil, err
	}

	return libtrust.UnmarshalPublicKeyPEM(pemBytes)
}

func VerifyJsonSignature(privateKey string, jsonSignature *libtrust.JSONSignature, publicKey libtrust.PublicKey) (*requests.CheckLicenseResponse, error) {
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

	checkLicenseResponse := new(requests.CheckLicenseResponse)
	if err := json.NewDecoder(bytes.NewReader(payload)).Decode(checkLicenseResponse); err != nil {
		return nil, err
	}

	ok, err := verify.CheckToken(checkLicenseResponse.Expiration.Format(time.RFC3339), checkLicenseResponse.Token, privateKey)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("Invalid token")
	}

	return checkLicenseResponse, nil
}
