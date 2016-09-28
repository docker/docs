package auth

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/controller/ctx"
	"github.com/docker/orca/controller/manager"
)

const (
	ErrInvalidAuthTokenStr = "invalid authentication token provided"
)

func defaultDeniedHostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

type AuthRequired struct {
	deniedHostHandler http.Handler
	manager           manager.Manager
	whitelistCIDRs    []string
}

func NewAuthRequired(m manager.Manager, whitelistCIDRs []string) *AuthRequired {
	return &AuthRequired{
		deniedHostHandler: http.HandlerFunc(defaultDeniedHostHandler),
		manager:           m,
		whitelistCIDRs:    whitelistCIDRs,
	}
}

func (a *AuthRequired) isWhitelisted(addr string) (bool, error) {
	parts := strings.Split(addr, ":")
	src := parts[0]

	srcIp := net.ParseIP(src)

	// check each whitelisted ip
	for _, c := range a.whitelistCIDRs {
		_, ipNet, err := net.ParseCIDR(c)
		if err != nil {
			return false, err
		}

		// global from anywhere
		if ipNet.String() == "0.0.0.0/0" {
			return true, nil
		}

		if ipNet.Contains(srcIp) {
			return true, nil
		}
	}

	return false, nil
}

func (a *AuthRequired) handleRequest(w http.ResponseWriter, r *http.Request) (*ctx.OrcaRequestContext, error) {
	requestContext := ctx.OrcaRequestContext{}

	whitelisted, err := a.isWhitelisted(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	if whitelisted {
		return &requestContext, nil
	}

	// set the remote IP from x-forwarded-for if present
	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}

	var authCtx *auth.Context
	// check for peer certificates; check client cert auth
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		clientCert := r.TLS.PeerCertificates[0]
		authCtx, err = a.manager.AuthenticatePublicKey(clientCert.PublicKey, remoteAddr)
	} else {
		authCtx, err = a.manager.AuthenticateSessionToken(r.Header.Get("Authorization"))
	}

	if err != nil {
		return nil, err
	}

	// Admin bool to Role conversion - should eventually be deprecated at the authenticator level
	if authCtx.User.Admin {
		authCtx.User.Role = auth.Admin
	}

	// Populate the Request Context
	requestContext.Auth = authCtx
	requestContext.Request = r

	// Extract variables from the request
	err = requestContext.ParseVars()
	if err != nil {
		return nil, err
	}

	return &requestContext, nil
}

func (a *AuthRequired) Initializer(w http.ResponseWriter, r *http.Request) (*ctx.OrcaRequestContext, error) {
	rc, err := a.handleRequest(w, r)
	if err == nil {
		return rc, nil
	}

	// Handle the unauthorized flow
	a.manager.SaveEvent(&orca.Event{
		Type:       "unauthorized",
		Time:       time.Now(),
		RemoteAddr: r.RemoteAddr,
		Tags:       []string{"api", r.URL.Path},
		Message:    err.Error(),
	})
	a.deniedHostHandler.ServeHTTP(w, r)
	return nil, err
}
