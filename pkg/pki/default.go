package pki

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

type DefaultClient struct {
	addr      string
	rootCert  string
	tlsConfig *tls.Config
}

type InfoReq struct {
	Label   string `json:"label"`
	Profile string `json:"profile"`
}

type InfoResp struct {
	Certificate  string   `json:"certificate"`
	Usage        []string `json:"usages"`
	ExpiryString string   `json:"expiry"`
}

func NewDefaultClient(addr string, tlsConfig *tls.Config) (*DefaultClient, error) {
	// Load up the root CA for this server at the beginning
	c := &DefaultClient{
		addr:      addr,
		tlsConfig: tlsConfig,
	}

	info := InfoReq{
		Profile: "node", // TODO Should probably be a const someplace
	}

	buf, err := json.Marshal(info)
	if err != nil {
		log.Debug("Failed to marshal info")
		return nil, err
	}

	resp, err := c.doRequest("POST", "/api/v1/cfssl/info", nil, bytes.NewBuffer(buf))
	if err != nil {
		log.Debugf("Post to remote server failed: %s", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		msg := "unknown error"
		if resp.Body != nil {
			d, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			msg = string(d)
		}
		return nil, fmt.Errorf("error from api: %s", msg)
	}

	d, _ := ioutil.ReadAll(resp.Body)
	var r struct {
		// TODO - hardening based on these flags
		//Success  bool      `json:"success"`
		Result *InfoResp `json:"result"`
		//Errors   []string  `json:"errors"`
		//Messages []string  `json:"messages"`
	}

	if err := json.Unmarshal(d, &r); err != nil {
		return nil, err
	}

	c.rootCert = r.Result.Certificate
	log.Debug("PKI client succesfully configured")
	return c, nil
}

func (c *DefaultClient) GetRootCertificate() (string, error) {
	return c.rootCert, nil
}

func (c *DefaultClient) Address() string {
	return c.addr
}

func (c *DefaultClient) getURL(p string) string {
	return c.addr + p
}

func (c *DefaultClient) doRequest(method string, p string, params *url.Values, body io.Reader) (*http.Response, error) {
	urlParams := ""
	if params != nil {
		urlParams = params.Encode()
	}

	u := fmt.Sprintf("%s?%s", c.getURL(p), urlParams)

	log.Debugf("request: method=%s url=%s", method, u)

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	tr := &http.Transport{}
	if c.tlsConfig != nil {
		tr.TLSClientConfig = c.tlsConfig
	}
	client := &http.Client{Transport: tr}
	return client.Do(req)
}

// SignCSR sends a certificate signing request to the API for signing
// and returns the signed certificate as a string if successful
func (c *DefaultClient) SignCSR(csr *CertificateSigningRequest) (*CertificateResponse, error) {
	buf, err := json.Marshal(csr)
	if err != nil {
		log.Debug("Failed to unmarshal csr")
		return nil, err
	}

	resp, err := c.doRequest("POST", "/api/v1/cfssl/sign", nil, bytes.NewBuffer(buf))
	if err != nil {
		log.Debug("Post to remote server failed")
		return nil, err
	}

	if resp.StatusCode != 200 {
		msg := "unknown error"
		if resp.Body != nil {
			d, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			msg = string(d)
		}
		return nil, fmt.Errorf("error from api: %s", msg)
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response from api: %s", err)
	}

	var r struct {
		// TODO - hardening based on these flags
		// Success  bool                 `json:"success"`
		Result *CertificateResponse `json:"result"`
		// Errors   []string             `json:"errors"`
		// Messages []string             `json:"messages"`
	}

	if err := json.Unmarshal(d, &r); err != nil || r.Result == nil {
		return nil, err
	}

	// Glue in the known root CA
	r.Result.CertificateChain = c.rootCert

	return r.Result, nil
}
