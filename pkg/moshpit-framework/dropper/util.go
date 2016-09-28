package dropper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func getJWT(host string, httpClient *http.Client, username, password string) (string, error) {
	type LoginForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data, err := json.Marshal(LoginForm{username, password})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/auth/login", host), bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Failed to authenticate. Status code: %d", res.StatusCode)
	}
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	type LoginResponse struct {
		AuthToken string `json:"auth_token"`
	}
	parsed := LoginResponse{}
	err = decoder.Decode(&parsed)
	if err != nil {
		return "", err
	}
	return parsed.AuthToken, nil
}
