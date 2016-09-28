package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

func HTTPClient(insecure bool, cas ...string) (*http.Client, error) {
	client := new(http.Client)
	client.Timeout = 10 * time.Minute
	var err error
	client.Transport, err = HTTPTransport(insecure, cas...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func HTTPTransport(insecure bool, cas ...string) (*http.Transport, error) {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Minute,
			KeepAlive: 10 * time.Minute,
		}).Dial,
		TLSHandshakeTimeout: 30 * time.Second,
		//ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
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
	}
	return &transport, nil
}
