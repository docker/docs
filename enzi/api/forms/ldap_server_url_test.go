package forms

import (
	"testing"
)

func TestValidateLDAPServerURL(t *testing.T) {
	type ldapServerURLTestCase struct {
		serverURL string
		startTLS  bool
		isValid   bool
	}

	testCases := []ldapServerURLTestCase{
		// Empty string.
		{serverURL: "", isValid: false},
		// No scheme, has hostname or IP, no startTLS.
		{serverURL: "192.168.1.189", startTLS: false, isValid: false},
		{serverURL: "54.128.12.17", startTLS: false, isValid: false},
		{serverURL: "ldap.example.com", startTLS: false, isValid: false},
		{serverURL: "ad.corpo.com", startTLS: false, isValid: false},
		// No scheme, has hostname or IP, with startTLS.
		{serverURL: "192.168.1.189", startTLS: true, isValid: false},
		{serverURL: "54.128.12.17", startTLS: true, isValid: false},
		{serverURL: "ldap.example.com", startTLS: true, isValid: false},
		{serverURL: "ad.corpo.com", startTLS: true, isValid: false},
		// No scheme, has hostname or IP with port, no startTLS.
		{serverURL: "192.168.1.189:389", startTLS: false, isValid: false},
		{serverURL: "54.128.12.17:3389", startTLS: false, isValid: false},
		{serverURL: "ldap.example.com:636", startTLS: false, isValid: false},
		{serverURL: "ad.corpo.com:3636", startTLS: false, isValid: false},
		// No scheme, has hostname or IP with port, with startTLS.
		{serverURL: "192.168.1.189:389", startTLS: true, isValid: false},
		{serverURL: "54.128.12.17:3389", startTLS: true, isValid: false},
		{serverURL: "ldap.example.com:636", startTLS: true, isValid: false},
		{serverURL: "ad.corpo.com:3636", startTLS: true, isValid: false},
		// 'ldap' scheme, has hostname or IP, no startTLS.
		{serverURL: "ldap://192.168.1.189", startTLS: false, isValid: true},
		{serverURL: "ldap://54.128.12.17", startTLS: false, isValid: true},
		{serverURL: "ldap://ldap.example.com", startTLS: false, isValid: true},
		{serverURL: "ldap://ad.corpo.com", startTLS: false, isValid: true},
		// 'ldap' scheme, has hostname or IP, with startTLS.
		{serverURL: "ldap://192.168.1.189", startTLS: true, isValid: true},
		{serverURL: "ldap://54.128.12.17", startTLS: true, isValid: true},
		{serverURL: "ldap://ldap.example.com", startTLS: true, isValid: true},
		{serverURL: "ldap://ad.corpo.com", startTLS: true, isValid: true},
		// 'ldap' scheme, has hostname or IP with port, no startTLS.
		{serverURL: "ldap://192.168.1.189:389", startTLS: false, isValid: true},
		{serverURL: "ldap://54.128.12.17:3389", startTLS: false, isValid: true},
		{serverURL: "ldap://ldap.example.com:636", startTLS: false, isValid: true},
		{serverURL: "ldap://ad.corpo.com:3636", startTLS: false, isValid: true},
		// 'ldap' scheme, has hostname or IP with port, with startTLS.
		{serverURL: "ldap://192.168.1.189:389", startTLS: true, isValid: true},
		{serverURL: "ldap://54.128.12.17:3389", startTLS: true, isValid: true},
		{serverURL: "ldap://ldap.example.com:636", startTLS: true, isValid: true},
		{serverURL: "ldap://ad.corpo.com:3636", startTLS: true, isValid: true},

		// 'ldaps' scheme, has hostname or IP, no startTLS.
		{serverURL: "ldaps://192.168.1.189", startTLS: false, isValid: true},
		{serverURL: "ldaps://54.128.12.17", startTLS: false, isValid: true},
		{serverURL: "ldaps://ldap.example.com", startTLS: false, isValid: true},
		{serverURL: "ldaps://ad.corpo.com", startTLS: false, isValid: true},
		// 'ldaps' scheme, has hostname or IP, with startTLS.
		{serverURL: "ldaps://192.168.1.189", startTLS: true, isValid: false},
		{serverURL: "ldaps://54.128.12.17", startTLS: true, isValid: false},
		{serverURL: "ldaps://ldap.example.com", startTLS: true, isValid: false},
		{serverURL: "ldaps://ad.corpo.com", startTLS: true, isValid: false},
		// 'ldaps' scheme, has hostname or IP with port, no startTLS.
		{serverURL: "ldaps://192.168.1.189:389", startTLS: false, isValid: true},
		{serverURL: "ldaps://54.128.12.17:3389", startTLS: false, isValid: true},
		{serverURL: "ldaps://ldap.example.com:636", startTLS: false, isValid: true},
		{serverURL: "ldaps://ad.corpo.com:3636", startTLS: false, isValid: true},
		// 'ldaps' scheme, has hostname or IP with port, with startTLS.
		{serverURL: "ldaps://192.168.1.189:389", startTLS: true, isValid: false},
		{serverURL: "ldaps://54.128.12.17:3389", startTLS: true, isValid: false},
		{serverURL: "ldaps://ldap.example.com:636", startTLS: true, isValid: false},
		{serverURL: "ldaps://ad.corpo.com:3636", startTLS: true, isValid: false},
	}

	for _, testCase := range testCases {
		err := ValidateLDAPServerURL("serverURL", testCase.serverURL, testCase.startTLS)
		if testCase.isValid && err != nil {
			t.Errorf("test case invalid %#v: %s", testCase, err)
		} else if !testCase.isValid && err == nil {
			t.Errorf("test case should be invalid %#v: %s", testCase, err)
		}
	}
}
