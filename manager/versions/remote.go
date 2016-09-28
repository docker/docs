package versions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

const (
	REGISTRY_V2_URL = "https://registry-1.docker.io/v2"
)

var USER_AGENT = fmt.Sprintf("docker/1.6.2 docker-trusted-registry/%s", deploy.Version)

type tokenResponse struct {
	Token string `json:"token"`
}

type registryTagList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	Code string   `json:"code"`
	Msg  string   `json:"message"`
}

func getRegistryAuthToken(authConfig *dockerclient.AuthConfig, params map[string]string) (string, error) {
	if authConfig == nil {
		authConfig = new(dockerclient.AuthConfig)
	}

	log.WithField("user", authConfig.Username).Info("Retrieving registry auth token")
	realm, err := url.Parse(params["realm"])
	if err != nil {
		return "", err
	}

	query := make(url.Values)
	query.Set("service", params["service"])
	query.Set("scope", params["scope"])
	query.Set("account", authConfig.Username)
	realm.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", realm.String(), nil)
	if err != nil {
		return "", err
	}
	if authConfig.Username != "" {
		req.SetBasicAuth(authConfig.Username, authConfig.Password)
	}
	resp, err := dtrutil.DoRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token auth attempt for registry failed with status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)

	tr := new(tokenResponse)
	if err = decoder.Decode(tr); err != nil {
		return "", fmt.Errorf("unable to decode token response: %s", err)
	}

	if tr.Token == "" {
		return "", fmt.Errorf("authorization server did not include a token in the response")
	}

	return tr.Token, nil
}

func getRemoteTagList(authConfig *dockerclient.AuthConfig, managerRepoName, token string) ([]string, error) {
	tagsUrl := fmt.Sprintf("%s/%s/tags/list", REGISTRY_V2_URL, managerRepoName)
	tagsRequest, err := http.NewRequest("GET", tagsUrl, nil)
	if token != "" {
		tagsRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	tagsRequest.Header.Set("User-Agent", USER_AGENT)
	resp, err := dtrutil.DoRequest(tagsRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var repoTags registryTagList
	if resp.StatusCode == http.StatusUnauthorized {
		if token == "" {
			challengeHeader := resp.Header.Get("WWW-Authenticate")
			if challengeHeader == "" {
				return nil, fmt.Errorf("Status unauthorized but no Authenticate header provided.")
			}
			_, params := parseValueAndParams(challengeHeader)
			newToken, err := getRegistryAuthToken(authConfig, params)
			if err != nil {
				return nil, err
			}
			return getRemoteTagList(authConfig, managerRepoName, newToken)
		} else {
			log.WithFields(log.Fields{
				"username": authConfig.Username,
				"repoName": managerRepoName,
			}).Error("DTR's logged in user is not authorized to list tags")
			return nil, fmt.Errorf("DTR's logged in user does not have the appropriate credentials to list tags on " + managerRepoName)
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&repoTags); err != nil {
		return nil, fmt.Errorf("Unable to decode tag list (%v)", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"statusCode":   resp.StatusCode,
			"errorCode":    repoTags.Code,
			"errorMessage": repoTags.Msg,
		}).Error("failed to get tags from v2 registry")
		return nil, fmt.Errorf("Error while retrieving tags from v2 registry")
	}
	return repoTags.Tags, nil
}
