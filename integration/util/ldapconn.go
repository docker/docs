package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/docker/garant/auth/common"

	"github.com/go-ldap/ldap"
)

func resolveURL(hostURLString string) (*url.URL, error) {
	if !strings.Contains(hostURLString, "//") {
		hostURLString = "//" + hostURLString
	}
	hostURL, err := url.Parse(hostURLString)
	if err != nil {
		return nil, common.WithStackTrace(err)
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

func GetLDAPConn(hostURLString string, startTLS bool) (*ldap.Conn, error) {
	hostURL, err := resolveURL(hostURLString)
	if err != nil {
		return nil, err
	}

	hostname, port, err := net.SplitHostPort(hostURL.Host)
	if err != nil {
		return nil, common.WithStackTrace(err)
	}
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, common.WithStackTrace(err)
	}

	for i := range addrs {
		addrs[i] += ":" + port
	}

	var conn net.Conn
	switch hostURL.Scheme {
	case "ldap":
		conn, err = dialMulti("tcp", addrs, new(net.Dialer).Dial)
		if err != nil {
			return nil, common.WithStackTrace(err)
		}
		ldapConn := ldap.NewConn(conn, false)
		ldapConn.Start()
		if startTLS {
			if err := ldapConn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
				ldapConn.Close()
				return nil, common.WithStackTrace(err)
			}
		}
		return ldapConn, nil
	case "ldaps":
		if startTLS {
			return nil, common.WithStackTrace(fmt.Errorf("cannot use StartTLS with ldaps://"))
		}
		conn, err = dialMulti("tcp", addrs, func(network, addr string) (net.Conn, error) {
			return tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: true})
		})
		if err != nil {
			return nil, common.WithStackTrace(err)
		}
		ldapConn := ldap.NewConn(conn, true)
		ldapConn.Start()
		return ldapConn, nil
	default:
		return nil, common.WithStackTrace(fmt.Errorf("%s is an invalid scheme for a LDAP url", hostURL.Scheme))
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
