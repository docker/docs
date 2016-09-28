package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/docker/orca/auth"
	"github.com/go-ldap/ldap"
)

func resolveURL(hostURLString string) (*url.URL, error) {
	if !strings.Contains(hostURLString, "//") {
		hostURLString = "//" + hostURLString
	}

	hostURL, err := url.Parse(hostURLString)
	if err != nil {
		return nil, err
	}

	hostURL.Scheme = strings.ToLower(hostURL.Scheme)
	if hostURL.Scheme == "" {
		hostURL.Scheme = "ldap"
	}

	if !strings.Contains(hostURL.Host, ":") {
		switch hostURL.Scheme {
		case "ldap":
			hostURL.Host += ":389"
		case "ldaps":
			hostURL.Host += ":636"
		}
	}

	hostURL.Path = ""

	return hostURL, nil
}

func getTLSConfig(settings auth.LDAPSettings, serverName string) (*tls.Config, error) {
	config := &tls.Config{
		InsecureSkipVerify: settings.TLSSkipVerify,
		ServerName:         serverName,
	}

	if settings.RootCerts != "" {
		rootCertPool := x509.NewCertPool()
		if !rootCertPool.AppendCertsFromPEM([]byte(settings.RootCerts)) {
			return nil, errors.New("unable to parse one or more root certificates")
		}
		config.RootCAs = rootCertPool
	}

	return config, nil
}

func GetConn(hostURLString string, settings auth.LDAPSettings) (ldapConn *ldap.Conn, err error) {
	hostURL, err := resolveURL(hostURLString)
	if err != nil {
		return nil, err
	}

	hostname, port, err := net.SplitHostPort(hostURL.Host)
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}

	for i := range addrs {
		addrs[i] += ":" + port
	}

	var conn net.Conn

	defer func() {
		if err == nil {
			return
		}

		// Always close the connection(s) in case of an error.
		if ldapConn != nil {
			ldapConn.Close()
		} else if conn != nil {
			conn.Close()
		}
	}()

	switch hostURL.Scheme {
	case "ldap":
		conn, err = dialMulti("tcp", addrs, new(net.Dialer).Dial)
		if err != nil {
			return nil, err
		}

		ldapConn = ldap.NewConn(conn, false)
		ldapConn.Start()

		if settings.StartTLS {
			tlsConfig, err := getTLSConfig(settings, hostname)
			if err != nil {
				return nil, err
			}

			if err := ldapConn.StartTLS(tlsConfig); err != nil {
				return nil, err
			}
		}

		return ldapConn, nil
	case "ldaps":
		if settings.StartTLS {
			return nil, fmt.Errorf("cannot use StartTLS with ldaps://")
		}

		tlsConfig, err := getTLSConfig(settings, hostname)
		if err != nil {
			return nil, err
		}

		conn, err = dialMulti("tcp", addrs, func(network, addr string) (net.Conn, error) {
			return tls.Dial(network, addr, tlsConfig)
		})
		if err != nil {
			return nil, err
		}

		ldapConn = ldap.NewConn(conn, true)
		ldapConn.Start()

		return ldapConn, nil
	default:
		return nil, fmt.Errorf("%s is an invalid scheme for a LDAP url", hostURL.Scheme)
	}
}

// dialMulti attempts to establish connections to each destination of
// the list of addresses. It will return the first established
// connection and close the other connections. Otherwise it returns
// error on the last attempt.
// This method is based off of the go "net" package internal dialMulti
func dialMulti(network string, addrs []string, dialFunc func(network, addr string) (net.Conn, error)) (net.Conn, error) {
	type racer struct {
		net.Conn
		error
	}
	// sig controls the flow of dial results on lane. It passes a
	// token to the next racer and also indicates the end of flow
	// by using closed channel.
	sig := make(chan bool, 1)
	lane := make(chan racer, 1)
	for _, addr := range addrs {
		go func(addr string) {
			c, err := dialFunc(network, addr)
			if _, ok := <-sig; ok {
				lane <- racer{c, err}
			} else if err == nil {
				// We have to return the resources
				// that belong to the other
				// connections here for avoiding
				// unnecessary resource starvation.
				c.Close()
			}
		}(addr)
	}
	defer close(sig)
	lastErr := errTimeout
	nracers := len(addrs)
	for nracers > 0 {
		sig <- true
		racer := <-lane
		if racer.error == nil {
			return racer.Conn, nil
		}
		lastErr = racer.error
		nracers--
	}
	return nil, lastErr
}

// timeoutError is taken from the go "net" package internals
type timeoutError struct{}

func (e *timeoutError) Error() string   { return "i/o timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }

var errTimeout error = new(timeoutError)
