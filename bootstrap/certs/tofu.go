package certs

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Attempt to get the CA for a server (assumes /ca endpoint works)
func getCA(config *tls.Config, addr string) (*x509.Certificate, error) {
	// Attempt to get the /ca first and favor that
	tr := &http.Transport{
		TLSClientConfig: config,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fmt.Sprintf("https://%s/ca", addr))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Server does not support /ca endpoint: %d, %s", resp.StatusCode, body)
	}
	der, _ := pem.Decode([]byte(body))
	if der == nil {
		log.Debug(body)
		return nil, fmt.Errorf("Server does not support /ca endpoint")
	}
	c, err := x509.ParseCertificate(der.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Server did not return valid cert at /ca endpoint: %s", err)
	}
	return c, nil
}

// Connect to an untrusted server, and get cert information about it
func Tofu(addr string) ([]string, error) {
	ret := []string{}
	if !strings.Contains(addr, ":") {
		addr = addr + ":443"
	}
	// First check to see if it's already trusted by the system
	config := &tls.Config{}
	conn, err := tls.Dial("tcp", addr, config)
	if err == nil {
		log.Debugf("Server %s is signed by a system-wide trusted CA", addr)
		return ret, nil
	}

	log.Debugf("Performing TOFU check on %s", addr)

	config.InsecureSkipVerify = true // By definition, we're trying to build the trust...
	conn, err = tls.Dial("tcp", addr, config)
	if err != nil {
		log.Infof("Unable to connect to %s: %s", addr, err)
		return nil, err
	} else {
		state := conn.ConnectionState()

		cert, err := getCA(config, addr)
		if err == nil {
			ret = append(ret, fmt.Sprintf("CA Subject: %s", cert.Subject.CommonName))
			ret = append(ret, fmt.Sprintf("Serial Number: %s", hex.EncodeToString(cert.SerialNumber.Bytes())))
			ret = append(ret, getFingerprint(cert.Raw, "SHA-256"))
		} else {

			for _, cert := range state.PeerCertificates {
				ret = append(ret, fmt.Sprintf("Subject: %s", cert.Subject.CommonName))
				ret = append(ret, fmt.Sprintf("Serial Number: %s", hex.EncodeToString(cert.SerialNumber.Bytes())))
				ret = append(ret, fmt.Sprintf("Issuer: %s", cert.Issuer.CommonName))
				ret = append(ret, getFingerprint(cert.Raw, "SHA-256"))
			}
		}
		return ret, nil
	}
}

// Given a trusted fingerprint, return an http client that will verify a match
func GetTofuClient(fingerprint, addr string) (*http.Client, error) {
	fingerprint = strings.ToUpper(fingerprint)
	config := &tls.Config{}
	if len(fingerprint) == 0 {
		// Assume a system-wide trusted cert
		if !strings.Contains(addr, ":") {
			addr = addr + ":443"
		}
		_, err := tls.Dial("tcp", addr, config)
		if err == nil {
			log.Debugf("Server %s is signed by a system-wide trusted CA", addr)
			return &http.Client{}, nil
		}
		return nil, fmt.Errorf("Server is not signed by an already trusted CA, you must specify a fingerprint")
	}

	// Now figure out if the fingerprint is a CA fingerprint, or server fingerprint
	config.InsecureSkipVerify = true
	cert, err := getCA(config, addr)
	if err == nil {
		sha1Fingerprint := strings.ToUpper(getFingerprint(cert.Raw, "SHA-1"))
		sha256Fingerprint := strings.ToUpper(getFingerprint(cert.Raw, "SHA-256"))
		if strings.Contains(sha1Fingerprint, fingerprint) || strings.Contains(sha256Fingerprint, fingerprint) {
			log.Debugf("Trusting all certs signed by CA: %s - serial number %s", cert.Subject.CommonName, hex.EncodeToString(cert.SerialNumber.Bytes()))
			caCertPool := x509.NewCertPool()
			caCertPool.AddCert(cert)
			return &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs: caCertPool,
					},
				},
			}, nil
		} // else fall through
	}

	log.Debugf("Trusting server certs with fingerprint: %s", fingerprint)
	// Figure out if we should use a CA trust, or server trust

	dial := func(network, addr string) (net.Conn, error) {
		config := &tls.Config{
			InsecureSkipVerify: true, // We'll verify ourselves
		}

		conn, err := tls.Dial(network, addr, config)
		if err != nil {
			log.Infof("Unable to connect to %s: %s", addr, err)
			return nil, err

		}
		state := conn.ConnectionState()
		now := time.Now()
		matched := false
		for _, cert := range state.PeerCertificates {
			// Allow substring and case insensitive matching to make it a little easier on the user
			sha1Fingerprint := strings.ToUpper(getFingerprint(cert.Raw, "SHA-1"))
			sha256Fingerprint := strings.ToUpper(getFingerprint(cert.Raw, "SHA-256"))
			if strings.Contains(sha1Fingerprint, fingerprint) {
				log.Debugf("%s ~= %s", sha1Fingerprint, fingerprint)
				matched = true
			} else if strings.Contains(sha256Fingerprint, fingerprint) {
				log.Debugf("%s ~= %s", sha256Fingerprint, fingerprint)
				matched = true
			} else {
				log.Debugf("%s != %s", sha256Fingerprint, fingerprint)
			}
			if now.Before(cert.NotBefore) {
				conn.Close()
				return nil, fmt.Errorf("Server certificate %s is not yet valid (too new)", cert.Subject.CommonName)
			}
			if now.After(cert.NotAfter) {
				conn.Close()
				return nil, fmt.Errorf("Server certificate %s has expired", cert.Subject.CommonName)
			}
		}
		if !matched {
			conn.Close()
			return nil, errors.New("Server certificate(s) didn't match your trusted fingerprint.  Re-run without --fingerprint to report server fingerprint")
		}
		// If we've gotten this far, then we can trust the server
		log.Debug("Server cert(s) passed TOFU tests")
		return conn, nil
	}
	return &http.Client{
		Transport: &http.Transport{
			DialTLS: dial,
		},
	}, nil

}
