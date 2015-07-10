package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/Sirupsen/logrus"
)

// HTTPStore manages pulling and pushing metadata from and to a remote
// service over HTTP. It assumes the URL structure of the remote service
// maps identically to the structure of the TUF repo:
// <baseURL>/<metaPrefix>/(root|targets|snapshot|timestamp).json
// <baseURL>/<targetsPrefix>/foo.sh
//
// If consistent snapshots are disabled, it is advised that caching is not
// enabled. Simple set a cachePath (and ensure it's writeable) to enable
// caching.
type HTTPStore struct {
	baseURL       url.URL
	metaPrefix    string
	metaExtension string
	targetsPrefix string
	keyExtension  string
	roundTrip     http.RoundTripper
}

func NewHTTPStore(baseURL, metaPrefix, metaExtension, targetsPrefix, keyExtension string, roundTrip http.RoundTripper) (*HTTPStore, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if !base.IsAbs() {
		return nil, errors.New("HTTPStore requires an absolute baseURL")
	}
	return &HTTPStore{
		baseURL:       *base,
		metaPrefix:    metaPrefix,
		metaExtension: metaExtension,
		targetsPrefix: targetsPrefix,
		keyExtension:  keyExtension,
		roundTrip:     roundTrip,
	}, nil
}

// GetMeta downloads the named meta file with the given size. A short body
// is acceptable because in the case of timestamp.json, the size is a cap,
// not an exact length.
func (s HTTPStore) GetMeta(name string, size int64) (json.RawMessage, error) {
	url, err := s.buildMetaURL(name)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.roundTrip.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, &ErrMetaNotFound{role: name}
	}
	b := io.LimitReader(resp.Body, int64(size))
	body, err := ioutil.ReadAll(b)

	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

func (s HTTPStore) SetMeta(name string, blob json.RawMessage) error {
	url, err := s.buildMetaURL(name)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(blob))
	if err != nil {
		return err
	}
	_, err = s.roundTrip.RoundTrip(req)
	return err
}

func (s HTTPStore) buildMetaURL(name string) (*url.URL, error) {
	filename := fmt.Sprintf("%s.%s", name, s.metaExtension)
	uri := path.Join(s.metaPrefix, filename)
	return s.buildURL(uri)
}

func (s HTTPStore) buildTargetsURL(name string) (*url.URL, error) {
	uri := path.Join(s.targetsPrefix, name)
	return s.buildURL(uri)
}

func (s HTTPStore) buildKeyURL(name string) (*url.URL, error) {
	filename := fmt.Sprintf("%s.%s", name, s.keyExtension)
	uri := path.Join(s.metaPrefix, filename)
	return s.buildURL(uri)
}

func (s HTTPStore) buildURL(uri string) (*url.URL, error) {
	sub, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	return s.baseURL.ResolveReference(sub), nil
}

// GetTarget returns a reader for the desired target or an error.
// N.B. The caller is responsible for closing the reader.
func (s HTTPStore) GetTarget(path string) (io.ReadCloser, error) {
	url, err := s.buildTargetsURL(path)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Attempting to download target: ", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.roundTrip.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (s HTTPStore) GetKey(role string) ([]byte, error) {
	url, err := s.buildKeyURL(role)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.roundTrip.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
