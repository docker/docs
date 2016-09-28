package server

import (
	"net/url"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/distribution/registry/api/errcode"
)

func CheckRegistryError(actual error, expected errcode.ErrorCode) bool {
	errors, ok := actual.(errcode.Errors)
	if !ok {
		return false
	}
	if len(errors) != 1 {
		return false
	}

	theErrorCode, ok := errors[0].(errcode.ErrorCode)
	if !ok {
		theError, ok := errors[0].(errcode.Error)
		if !ok {
			return false
		}
		theErrorCode = theError.Code
	}
	if theErrorCode != expected {
		return false
	}
	return true
}

func maskPasswordsInHAConfig(haConfig *hubconfig.HAConfig) *hubconfig.HAConfig {
	newHAConfig := *haConfig
	newHAConfig.HTTPProxy = maskURLWithPassword(haConfig.HTTPProxy)
	newHAConfig.HTTPSProxy = maskURLWithPassword(haConfig.HTTPSProxy)
	return &newHAConfig
}

// maskURLWithPassword replaces the password in a url string with `password_hidden` if necessary
func maskURLWithPassword(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	if parsedURL.User == nil {
		return rawURL
	}
	if _, hasPass := parsedURL.User.Password(); !hasPass {
		return rawURL
	}
	parsedURL.User = url.UserPassword(parsedURL.User.Username(), "password_hidden")
	return parsedURL.String()
}

// unmaskURLWithPassword replaces password_hidden in inputURL with the password from oldURL if necessary
func unmaskURLWithPassword(inputURL, oldURL string) string {
	// The only case where we are okay with keeping the same credentials from the old URL is when there have been no changes
	if inputURL == maskURLWithPassword(oldURL) {
		return oldURL
	}
	return inputURL
}
