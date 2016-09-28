package ucpclient

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/orca"
)

type APIClient struct {
	host               string
	apiClientPort      uint16
	apiClientUrlScheme string
	client             *http.Client
	username           string
	password           string
	jwt                string
}

func New(host string, httpClient *http.Client) *APIClient {
	httpClient.Jar, _ = cookiejar.New(nil)
	return &APIClient{
		host:               host,
		client:             httpClient,
		apiClientUrlScheme: "https",
		apiClientPort:      443,
	}
}

func (c *APIClient) JWT() string {
	return c.jwt
}

func (c *APIClient) makeRequest(method, route string, payload interface{}) (*http.Response, error) {
	var reader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	// url builder
	var url string
	url = c.apiClientUrlScheme + "://" + c.host // ex: https://example.com
	// only specify the port in non-standard cases
	if (c.apiClientUrlScheme == "https" && c.apiClientPort != 443) ||
		(c.apiClientUrlScheme == "http" && c.apiClientPort != 80) {
		url = url + ":" + strconv.Itoa(int(c.apiClientPort)) // ex: https://example.com:8080
	}

	url = url + route // ex: http://example.com:8080/resource

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if method != "HEAD" && method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.jwt))

	// always close the connection after sending the request so that we
	// don't get bitten by net/http's bugs with reusing connections
	req.Close = true

	return dtrutil.DoRequestWithClient(req, c.client)
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthToken string `json:"auth_token"`
}

// logs in as the given user
func (c *APIClient) Login(username, password string) error {
	c.username = username
	c.password = password
	res, err := c.makeRequest("POST", "/auth/login", LoginForm{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to authenticate. Status code: %d", res.StatusCode)
	}
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	parsed := LoginResponse{}
	err = decoder.Decode(&parsed)
	if err != nil {
		return err
	}
	c.jwt = parsed.AuthToken
	return nil
}

type Bundle struct {
	CaPEM   []byte
	CertPEM []byte
	CertPub []byte
	KeyPEM  []byte
}

func BundleToDisk(bundle *Bundle, dest string) error {
	os.MkdirAll(dest, 0600)
	err := ioutil.WriteFile(filepath.Join(dest, "ca.pem"), bundle.CaPEM, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dest, "cert.pem"), bundle.CertPEM, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dest, "cert.pub"), bundle.CertPub, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dest, "key.pem"), bundle.KeyPEM, 0600)
	if err != nil {
		return err
	}
	return nil
}

func zipFileToBytes(file *zip.File) ([]byte, error) {
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()
	buff := bytes.NewBuffer([]byte{})
	_, err = io.Copy(buff, fileReader)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func ZipToBundle(reader io.ReaderAt, length int64) (*Bundle, error) {
	zipReader, err := zip.NewReader(reader, length)
	if err != nil {
		return nil, err
	}

	bundle := Bundle{}
	for _, file := range zipReader.File {
		switch file.Name {
		case "ca.pem":
			bundle.CaPEM, err = zipFileToBytes(file)
		case "cert.pem":
			bundle.CertPEM, err = zipFileToBytes(file)
		case "cert.pub":
			bundle.CertPub, err = zipFileToBytes(file)
		case "key.pem":
			bundle.KeyPEM, err = zipFileToBytes(file)
		}
		if err != nil {
			return nil, err
		}
	}
	return &bundle, nil
}

func (c *APIClient) GetBundle() (*Bundle, error) {
	res, err := c.makeRequest("GET", "/api/clientbundle", nil)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dataReader := bytes.NewReader(data)
	bundle, err := ZipToBundle(dataReader, int64(len(data)))
	return bundle, err
}

// this struct is copied from ucp
type LicenseConfig struct {
	KeyID         string `yaml:"key_id" json:"key_id"`
	PrivateKey    string `yaml:"private_key" json:"private_key,omitempty"`
	Authorization string `yaml:"authorization" json:"authorization,omitempty"`
}

// this struct is copied from ucp
type LicenseSubsystemConfig struct {
	// Tune if we automatically attempt to refresh expired licenses (read-write)
	AutoRefresh bool `json:"auto_refresh"`
	// The underlying license returned from the license server (read-write)
	License LicenseConfig `json:"license_config"`
	// Details derived from License (read-only)
	//Details LicenseDetails `json:"details"`
	// Any error reported from the license server on the last attempt to auto refresh
	LastUpdateError string `json:"last_refresh_error"`
}

func (c *APIClient) Ping() error {
	res, err := c.makeRequest("GET", "/", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return fmt.Errorf("Bad status code from /: %d", res.StatusCode)
	}
	return nil
}

func (c *APIClient) GetNodes() ([]orca.Node, error) {
	var nodes []orca.Node
	var err error
	path := "/api/nodes"
	res, err := c.makeRequest("GET", path, nil)
	if err != nil {
		return nil, err
	} else if res.StatusCode < 200 || res.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad status code from %s: %d", path, res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&nodes)
	if err != nil {
		return nil, err
	} else {
		return nodes, nil
	}
}

func (c *APIClient) GetCA() (string, error) {
	res, err := c.makeRequest("GET", "/ca", nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Bad status code fetching CA: %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *APIClient) GetLicenseConfig() (*LicenseSubsystemConfig, error) {
	res, err := c.makeRequest("GET", "/api/config/license", nil)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	licenseConfig := LicenseSubsystemConfig{}
	err = decoder.Decode(&licenseConfig)
	if err != nil {
		return nil, err
	}
	return &licenseConfig, nil
}

type PublicKey struct {
	Label     string `json:"label"`
	PublicKey string `json:"public_key"`
}

type Account struct {
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Username   string      `json:"username"`
	Password   string      `json:"password"`
	Role       int         `json:"role"`
	Admin      bool        `json:"admin"`
	PublicKeys []PublicKey `json:"public_keys"`
}

func (c *APIClient) GetAccount() (*Account, error) {
	res, err := c.makeRequest("GET", fmt.Sprintf("/api/accounts/%s", c.username), nil)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	account := Account{}
	err = decoder.Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *APIClient) UpdateAccount(account *Account) error {
	res, err := c.makeRequest("POST", fmt.Sprintf("/api/accounts"), account)
	if err != nil {
		return err
	}
	if res.StatusCode != 204 {
		return fmt.Errorf("Failed to update account. Status code: %d", res.StatusCode)
	}
	return nil
}

func (c *APIClient) LabelOwnBundle(publicKey, label string) error {
	account, err := c.GetAccount()
	if err != nil {
		return err
	}
	success := false
	for i, pk := range account.PublicKeys {
		if pk.PublicKey == publicKey {
			pk.Label = label
			success = true
			account.PublicKeys[i] = pk
		}
	}
	if !success {
		return fmt.Errorf("Failed to find bundle for renaming with public key: %s", publicKey)
	}
	return c.UpdateAccount(account)
}

func (c *APIClient) DeleteOwnBundle(publicKey string) error {
	account, err := c.GetAccount()
	if err != nil {
		return err
	}
	// filter out the bundles with matching public key
	success := false
	newList := []PublicKey{}
	for _, pk := range account.PublicKeys {
		if pk.PublicKey == publicKey {
			success = true
		} else {
			newList = append(newList, pk)
		}
	}
	if !success {
		return fmt.Errorf("Failed to find bundle for deletion with public key: %s", publicKey)
	}
	account.PublicKeys = newList
	return c.UpdateAccount(account)
}

func (c *APIClient) Logout() error {
	res, err := c.makeRequest("POST", fmt.Sprintf("/auth/logout"), nil)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to log out. Status code: %d", res.StatusCode)
	}
	return nil
}
