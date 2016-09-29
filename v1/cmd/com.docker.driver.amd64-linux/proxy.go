package main

import (
	"fmt"
	"net/url"
	"strings"
)

func parseNoProxy(p string) string {
	if p == "" {
		return p
	}
	elements := strings.Split(p, ",")
	for i, e := range elements {
		// trim whitespace
		elements[i] = strings.TrimSpace(e)
	}
	result := strings.Join(elements, ",")
	return result
}

func parseHTTPProxy(p string) (string, error) {
	if p == "" {
		return p, nil
	}
	u, err := url.Parse(p)
	if err != nil {
		return "", err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		// Default to HTTP
		u.Host = p
		u.Scheme = "http"
		u.Path = ""
		u.Opaque = ""
	}
	result := u.String()
	if strings.Contains(result, "localhost") || strings.Contains(result, "127.0.0.1") {
		return "", fmt.Errorf("%s is not reachable from the VM. Skipping", result)
	}
	return result, nil
}
