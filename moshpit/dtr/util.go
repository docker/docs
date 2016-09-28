package dtr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type tokenResponse struct {
	Token string `json:"token"`
}

func getToken(username, password, registry string, scopes []string, httpClient *http.Client) (token string, err error) {
	req, err := http.NewRequest("GET", "https://"+registry+"/auth/token", nil)
	if err != nil {
		return "", err
	}

	reqParams := req.URL.Query()

	reqParams.Add("service", registry)

	for _, scope := range scopes {
		reqParams.Add("scope", scope)
	}

	if username != "" {
		reqParams.Add("account", username)
		req.SetBasicAuth(username, password)
	}

	req.URL.RawQuery = reqParams.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token auth attempt: %s request failed with status: %d %s", req.URL, resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)

	tr := new(tokenResponse)
	if err = decoder.Decode(tr); err != nil {
		return "", fmt.Errorf("unable to decode token response: %s", err)
	}

	if tr.Token == "" {
		return "", errors.New("authorization server did not include a token in the response")
	}

	return tr.Token, nil
}
