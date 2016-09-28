package dtrutil

import (
	"net/http"
	"testing"
)

// TestHTTPClientGoogleSecure ensures that an http client with no CA and no
// insecure can GET https://www.google.com
func TestHTTPClientGoogleSecure(t *testing.T) {
	c, err := HTTPClient(false)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.Get("https://www.google.com")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

// TestHTTPClientGoogleInsecure ensures that an insecure http client with no CA
// can GET https://www.google.com
func TestHTTPClientGoogleInsecure(t *testing.T) {
	c, err := HTTPClient(true)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.Get("https://www.google.com")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}
