package authn

import (
	"fmt"
	"net/http"
)

func ShouldBypassAuth(r *http.Request) (bool, error) {
	// We trust nginx to tell us if the client cert is authorized to bypass auth
	if r.Header.Get("X-Client-Cert-Valid") == "SUCCESS" {
		if r.Header.Get("X-Client-OU-Valid") == "true" {
			return true, nil
		} else {
			return false, fmt.Errorf("failed to validate OU field")
		}
	}
	// if the client cert is bad, we need to alert UCP (the client) of this fact
	if r.Header.Get("X-Client-Cert-Valid") == "FAILED" {
		return false, fmt.Errorf("failed to validate client certificate")
	}
	return false, nil
}
