package server

import (
	"testing"

	"github.com/docker/dhe-deploy/hubconfig"
)

const exampleNeither = "http://example.com"
const exampleUsername = "http://abc@example.com"
const exampleNeitherNoHTTP = "example.com"
const exampleUsernameAndPassword = "http://abc:def@example.com"
const exampleUsernameAndPasswordNoHTTP = "abc:def@example.com"
const exampleUsername2AndPassword2 = "http://abcd:defg@example.com"
const exampleMasked = "http://abc:password_hidden@example.com"
const exampleAttack1 = "password_hidden"
const exampleAttack2 = ":password_hidden@"
const exampleAttack3 = "http://:password_hidden@"
const exampleAttack4 = "http://abc:password_hidden@example2.com"
const exampleAttack4_fail = "http://abc:def@example2.com"

func TestPasswordsMasking(t *testing.T) {
	for _, test := range []struct {
		url    string
		masked string
	}{
		{
			url:    exampleNeither,
			masked: exampleNeither,
		},
		{
			url:    exampleUsername,
			masked: exampleUsername,
		},
		{
			url:    exampleUsernameAndPassword,
			masked: exampleMasked,
		},
		{
			url:    exampleNeitherNoHTTP,
			masked: exampleNeitherNoHTTP,
		},
		{
			url:    exampleUsernameAndPasswordNoHTTP,
			masked: exampleUsernameAndPasswordNoHTTP,
		},
	} {
		haConfig := new(hubconfig.HAConfig)
		haConfig.HTTPProxy = test.url
		haConfig.HTTPSProxy = test.url
		haConfig = maskPasswordsInHAConfig(haConfig)

		if haConfig.HTTPProxy != test.masked {
			t.Fatalf("Masked %s into %s", test.url, haConfig.HTTPProxy)
		}
		if haConfig.HTTPSProxy != test.masked {
			t.Fatalf("Masked %s into %s", test.url, haConfig.HTTPSProxy)
		}
	}

	// make sure maskPasswordsInHAConfig doesn't modify the original config
	haConfig := new(hubconfig.HAConfig)
	haConfig.HTTPProxy = exampleUsernameAndPassword
	haConfig.HTTPSProxy = exampleUsernameAndPassword
	maskPasswordsInHAConfig(haConfig)

	if haConfig.HTTPProxy != exampleUsernameAndPassword {
		t.Fatalf("Modified %s into %s in place", exampleUsernameAndPassword, haConfig.HTTPProxy)
	}
	if haConfig.HTTPSProxy != exampleUsernameAndPassword {
		t.Fatalf("Modified %s into %s in place", exampleUsernameAndPassword, haConfig.HTTPSProxy)
	}
}

func TestPasswordsUnmasking(t *testing.T) {
	for _, test := range []struct {
		oldURL    string
		inputURL  string
		resultURL string
	}{
		{
			oldURL:    exampleNeither,
			inputURL:  exampleNeither,
			resultURL: exampleNeither,
		},
		{
			oldURL:    exampleNeither,
			inputURL:  exampleNeitherNoHTTP,
			resultURL: exampleNeitherNoHTTP,
		},
		{
			oldURL:    exampleUsername,
			inputURL:  exampleNeither,
			resultURL: exampleNeither,
		},
		{
			oldURL:    exampleUsernameAndPassword,
			inputURL:  exampleNeither,
			resultURL: exampleNeither,
		},
		{
			oldURL:    exampleNeither,
			inputURL:  exampleUsername,
			resultURL: exampleUsername,
		},
		{
			oldURL:    exampleUsername,
			inputURL:  exampleUsername,
			resultURL: exampleUsername,
		},
		{
			oldURL:    exampleUsernameAndPassword,
			inputURL:  exampleUsername,
			resultURL: exampleUsername,
		},
		{
			oldURL:    exampleNeither,
			inputURL:  exampleUsernameAndPassword,
			resultURL: exampleUsernameAndPassword,
		},
		{
			oldURL:    exampleUsername,
			inputURL:  exampleUsernameAndPassword,
			resultURL: exampleUsernameAndPassword,
		},
		{
			oldURL:    exampleUsernameAndPassword,
			inputURL:  exampleUsernameAndPassword,
			resultURL: exampleUsernameAndPassword,
		},
		{
			oldURL:    exampleNeither,
			inputURL:  exampleMasked,
			resultURL: exampleMasked,
		},
		{
			oldURL:    exampleUsername,
			inputURL:  exampleMasked,
			resultURL: exampleMasked,
		},
		{
			oldURL:    exampleUsernameAndPassword,
			inputURL:  exampleMasked,
			resultURL: exampleUsernameAndPassword,
		},
		{
			oldURL:    exampleUsername2AndPassword2,
			inputURL:  exampleUsernameAndPassword,
			resultURL: exampleUsernameAndPassword,
		},
	} {
		unmasked := unmaskURLWithPassword(test.inputURL, test.oldURL)
		if unmasked != test.resultURL {
			t.Fatalf("Unmasked %s with oldURL %s into %s", test.inputURL, test.oldURL, unmasked)
		}
	}
}

func TestPasswordsUnmaskingAttacks(t *testing.T) {
	// the first 3 attacks check if the UI will show the password to the user
	for _, attackStr := range []string{
		exampleAttack1,
		exampleAttack2,
		exampleAttack3,
	} {
		remasked := maskURLWithPassword(unmaskURLWithPassword(attackStr, exampleUsernameAndPassword))
		if remasked != attackStr {
			t.Fatalf("Re-masking turned %s into %s", attackStr, remasked)
		}
	}

	// this last attack checks if the user can change the proxy server to their own to trick
	// the system into logging into the wrong server with the credentials
	unmasked := unmaskURLWithPassword(exampleAttack4, exampleUsernameAndPassword)
	if unmasked == exampleAttack4_fail {
		t.Fatalf("Changing the url turned %s into %s", exampleAttack4, unmasked)
	}
}
