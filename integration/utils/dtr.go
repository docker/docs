package utils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

// DTRTestsEnabled checks if the env variables for DTR are configured
func DTRTestsEnabled() bool {
	if os.Getenv("DTR_URL") == "" ||
		os.Getenv("DTR_USERNAME") == "" ||
		os.Getenv("DTR_PASSWORD") == "" ||
		os.Getenv("DTR_PRIVATE_REPO") == "" {
		return false
	}
	return true
}

// GetDTRHTTPClient returns an http.Client with a cookie jar containing a DTR session token
func GetDTRHTTPClient(dtrURL, username, password string) *http.Client {
	// Create an HTTP Client with a Cookie Jar
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 1,
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout:   time.Duration(10 * time.Second),
		Transport: tr,
		Jar:       jar,
	}

	// Create a login request
	path := dtrURL + "/admin/login"
	resp, err := client.PostForm(path, url.Values{
		"username": {username},
		"password": {password},
	})
	if err != nil {
		log.Debug(err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		log.Debug(resp.StatusCode)
		log.Errorf("Failed to login to DTR: %s %s", err, data)
		return nil
	}

	// Return the client, holding the session cookies in the Jar
	return client
}
