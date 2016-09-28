package dtrutil

import (
	"archive/tar"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	version "github.com/docker/engine-api/types/versions"
	enzierrors "github.com/docker/orca/enzi/api/errors"
)

const maxHTTPRetries = 3

func CheckContainsEnziError(actual error, expectedCode string) bool {
	errors, ok := actual.(*enzierrors.APIErrors)
	if actual == nil || errors == nil || !ok {
		return false
	}

	for _, err := range errors.Errors {
		if err.Code == expectedCode {
			return true
		}
	}

	return false
}

func CombineErrors(errs ...error) error {
	var errSlice []error
	for _, err := range errs {
		if err == nil {
			// do nothing
		} else if m, ok := err.(multiError); ok {
			errSlice = append(errSlice, m.Errors()...)
		} else {
			errSlice = append(errSlice, err)
		}
	}
	if len(errSlice) == 0 {
		return nil
	}
	return &multiErr{errors: errSlice}
}

type multiError interface {
	Error() string
	Errors() []error
}

type multiErr struct {
	errors []error
}

func (e *multiErr) Error() string {
	switch len(e.errors) {
	case 0:
		return ""
	case 1:
		return e.errors[0].Error()
	default:
		errMsg := "errors:"
		for _, err := range e.errors {
			errMsg += "\n\t" + err.Error()
		}
		return errMsg
	}
}

func (e *multiErr) Errors() []error {
	return e.errors
}

func DoRequest(req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(req, http.DefaultClient)
}

func DoRequestWithClient(req *http.Request, client *http.Client) (resp *http.Response, err error) {
	resp, err = client.Do(req)
	switch req.Method {
	case "GET", "HEAD", "OPTIONS", "TRACE":
		for retryNumber := 0; retryNumber < maxHTTPRetries; retryNumber++ {
			if urlErr, ok := err.(*url.Error); !ok || urlErr.Err != io.EOF {
				break
			}
			req.Close = true
			log.WithFields(log.Fields{
				"error":       err,
				"retryNumber": retryNumber,
				"method":      req.Method,
				"url":         req.URL,
			}).Debug("retrying HTTP request")
			resp, err = client.Do(req)
		}
	}
	return
}

func Poll(interval time.Duration, retries int, run func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	for i := 0; i < retries; i++ {
		err = run()
		if err != nil {
			time.Sleep(interval)
		} else {
			return nil
		}
	}
	return fmt.Errorf("Polling failed with %d attempts %s apart: %s", retries, interval, err)
}

// Duration is used to facilitate automatic serialization/deserialization of a
// time.Duration in string form
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	d.Duration, err = time.ParseDuration(s)
	return err
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func DockerVersionCheck(client *client.Client, dockerRequired, ucpRequired string) (string, error) {
	verInfo, err := client.ServerVersion(context.Background())
	if err != nil {
		return "", fmt.Errorf("Failed to get docker version: %s", err)
	}

	// Drop the portion before the "/" since UCP adds it
	ucpSplit := strings.Split(verInfo.Version, "/")
	if len(ucpSplit) > 1 {
		if ucpSplit[0] != "ucp" {
			return verInfo.Version, fmt.Errorf("Unknown docker daemon type: %s", ucpSplit[0])
		}
		if version.LessThan(ucpSplit[1], ucpRequired) {
			return verInfo.Version, fmt.Errorf("Your engine version %v is too old. DTR requires at least engine version %v or UCP version %s.", verInfo.Version, dockerRequired, ucpRequired)
		}
	} else {
		if version.LessThan(verInfo.Version, dockerRequired) {
			return verInfo.Version, fmt.Errorf("Your engine version %v is too old. DTR requires at least engine version %v or UCP version %s.", verInfo.Version, dockerRequired, ucpRequired)
		}
	}

	return verInfo.Version, nil
}

type httpClientKey struct {
	insecure bool
	cas      string
}

var httpClientCache = map[httpClientKey]*http.Client{}

// see #1917 for why we can't just create a new http client every time we want one
func HTTPClient(insecure bool, cas ...string) (*http.Client, error) {
	casKey := ""
	for _, ca := range cas {
		casKey += ca
	}
	key := httpClientKey{
		insecure: insecure,
		cas:      casKey,
	}
	client, ok := httpClientCache[key]
	if !ok {
		client = new(http.Client)
		var err error
		client.Transport, err = HTTPTransport(insecure, cas, "", "")
		if err != nil {
			return nil, err
		}
		httpClientCache[key] = client
	}
	return client, nil
}

func HTTPTransport(insecure bool, cas []string, clientCertFile, clientKeyFile string) (*http.Transport, error) {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		//ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: insecure,
		},
	}
	if !insecure {
		caPool := systemRootsPool()
		for _, ca := range cas {
			if ca == "" {
				continue
			}
			ok := caPool.AppendCertsFromPEM([]byte(ca))
			if !ok {
				return nil, fmt.Errorf("TLS CA provided, but we could not parse a PEM encoded cert from it. CA provided: \n%s", ca)
			}
		}
		transport.TLSClientConfig.RootCAs = caPool

		if clientCertFile != "" || clientKeyFile != "" {
			// load client certs needed to connect to notary
			clientAuth, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
			if err != nil {
				return nil, fmt.Errorf("TLS client cert provided at (%s) and/or client key provided at (%s), but we could not load a keypair from.", clientCertFile, clientKeyFile)
			}
			transport.TLSClientConfig.Certificates = []tls.Certificate{clientAuth}
		}
	}
	return &transport, nil
}

func GetBytes(tableDocument interface{}) ([]byte, error) {
	if tableBytes, err := json.Marshal(tableDocument); err != nil {
		return nil, err
	} else {
		return tableBytes, nil
	}
}

func GetMap(tableDocument []byte) (interface{}, error) {
	var document map[string]interface{}
	if err := json.Unmarshal(tableDocument, &document); err != nil {
		return nil, err
	} else {
		return document, nil
	}
}

func AddBytesToTar(writer *tar.Writer, bytes []byte, name string) error {
	header := &tar.Header{
		Name:    name,
		Mode:    0600, // only user who creates tar file can read it, assuming it's DTR admin
		ModTime: time.Now(),
		Size:    int64(len(bytes)),
	}
	if err := writer.WriteHeader(header); err != nil {
		return err
	}

	_, err := writer.Write(bytes)
	return err
}
