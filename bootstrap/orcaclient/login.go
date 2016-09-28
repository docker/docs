package orcaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"

	log "github.com/Sirupsen/logrus"
)

// Lifted from orca - might want to import these...
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}
type AuthToken struct {
	Token     string `json:"auth_token,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}

func Login(client *http.Client, url, username, password string) (string, error) {
	orcaURL, err := neturl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("Failed to parse connection URL: %s", err)
	}
	orcaURL.Path = "/auth/login"
	creds := Credentials{
		Username: username,
		Password: password,
	}
	reqJson, err := json.Marshal(creds)
	if err != nil {
		return "", err
	}
	resp, err := client.Post(orcaURL.String(), "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var authToken AuthToken
			if err := json.Unmarshal(body, &authToken); err != nil {
				return "", err
			} else {
				return authToken.Token, nil
			}
		} else {
			return "", err
		}
	} else if resp.StatusCode == 401 {
		return "", fmt.Errorf("Bad username or password. Failed to login to UCP.")

	} else {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Warn(string(body))
		}
		return "", fmt.Errorf("Unexpected error logging in to UCP: %d", resp.StatusCode)
	}
}

func GetTokenHeader(token, username string) (string, string) {
	return "Authorization", "Bearer " + token

}
