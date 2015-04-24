package token

import (
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/docker/libtrust"

	"github.com/docker/vetinari/auth"
	"github.com/docker/vetinari/utils"
)

// ConfigSection is the name used to identify the tokenAuthorizer config
// in the config json
const ConfigSection string = "token_auth"

type tokenConfig struct {
	Realm          string `json:"realm"`
	Issuer         string `json:"issuer"`
	Service        string `json:"service"`
	RootCertBundle string `json:"root_cert_bundle"`
}

// authChallenge implements the auth.Challenge interface.
type authChallenge struct {
	err     error
	realm   string
	service string
	scopes  []auth.Scope
}

// Error returns the internal error string for this authChallenge.
func (ac authChallenge) Error() string {
	return ac.err.Error()
}

// Status returns the HTTP Response Status Code for this authChallenge.
func (ac *authChallenge) Status() int {
	return http.StatusUnauthorized
}

// challengeParams constructs the value to be used in
// the WWW-Authenticate response challenge header.
// See https://tools.ietf.org/html/rfc6750#section-3
func (ac *authChallenge) challengeParams() string {
	str := fmt.Sprintf("Bearer realm=%q,service=%q", ac.realm, ac.service)

	scope := make([]string, 0, len(ac.scopes))
	for _, s := range ac.scopes {
		scope = append(scope, s.ID())
	}
	if len(scope) > 0 {
		scopeStr := strings.Join(scope, " ")
		str = fmt.Sprintf("%s,scope=%q", str, scopeStr)
	}

	if ac.err == ErrInvalidToken || ac.err == ErrMalformedToken {
		str = fmt.Sprintf("%s,error=%q", str, "invalid_token")
	} else if ac.err == ErrInsufficientScope {
		str = fmt.Sprintf("%s,error=%q", str, "insufficient_scope")
	}

	return str
}

// SetHeader sets the WWW-Authenticate value for the given header.
func (ac *authChallenge) SetHeader(header http.Header) {
	header.Add("WWW-Authenticate", ac.challengeParams())
}

// ServeHttp handles writing the challenge response
// by setting the challenge header and status code.
func (ac *authChallenge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ac.SetHeader(w.Header())
	w.WriteHeader(ac.Status())
}

// accessController implements the auth.AccessController interface.
type tokenAuthorizer struct {
	realm       string
	issuer      string
	service     string
	rootCerts   *x509.CertPool
	trustedKeys map[string]libtrust.PublicKey
}

// NewTokenAuthorizer creates an Authorizer that operates with JWTs.
func NewTokenAuthorizer(conf []byte) (auth.Authorizer, error) {
	tokenConf := new(tokenConfig)
	err := json.Unmarshal(conf, tokenConf)
	if err != nil {
		return nil, fmt.Errorf("unable to parse TokenAuthorizer configuration: %s", err)
	}

	fp, err := os.Open(tokenConf.RootCertBundle)
	if err != nil {
		return nil, fmt.Errorf("unable to open token auth root certificate bundle file %q: %s", tokenConf.RootCertBundle, err)
	}
	defer fp.Close()

	rawCertBundle, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read token auth root certificate bundle file %q: %s", tokenConf.RootCertBundle, err)
	}

	var rootCerts []*x509.Certificate
	pemBlock, rawCertBundle := pem.Decode(rawCertBundle)
	for pemBlock != nil {
		cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse token auth root certificate: %s", err)
		}

		rootCerts = append(rootCerts, cert)

		pemBlock, rawCertBundle = pem.Decode(rawCertBundle)
	}

	if len(rootCerts) == 0 {
		return nil, errors.New("token auth requires at least one token signing root certificate")
	}

	rootPool := x509.NewCertPool()
	trustedKeys := make(map[string]libtrust.PublicKey, len(rootCerts))
	for _, rootCert := range rootCerts {
		rootPool.AddCert(rootCert)
		pubKey, err := libtrust.FromCryptoPublicKey(crypto.PublicKey(rootCert.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("unable to get public key from token auth root certificate: %s", err)
		}
		trustedKeys[pubKey.KeyID()] = pubKey
	}

	return &tokenAuthorizer{
		realm:       tokenConf.Realm,
		issuer:      tokenConf.Issuer,
		service:     tokenConf.Service,
		rootCerts:   rootPool,
		trustedKeys: trustedKeys,
	}, nil
}

// Authorize handles checking whether the given request is authorized
// for actions on resources described by the given Scopes.
func (ac *tokenAuthorizer) Authorize(r *http.Request, scopes ...auth.Scope) (*auth.User, error) {
	challenge := &authChallenge{
		realm:   ac.realm,
		service: ac.service,
		scopes:  scopes,
	}

	token, err := parseToken(r)
	if err != nil {
		challenge.err = err
		return nil, challenge
	}

	resource := auth.Resource{"repo", utils.ResourceName(r)}

	return ac.authorize(token, resource, scopes...)
}

// authorize separates out the code that needs to know how to handle a request from the rest of the
// authorization code making it easier to test this part by injecting test values
func (ac *tokenAuthorizer) authorize(token *Token, resource auth.Resource, scopes ...auth.Scope) (*auth.User, error) {
	challenge := &authChallenge{
		realm:   ac.realm,
		service: ac.service,
		scopes:  scopes,
	}

	verifyOpts := VerifyOptions{
		TrustedIssuers:    []string{ac.issuer},
		AcceptedAudiences: []string{ac.service},
		Roots:             ac.rootCerts,
		TrustedKeys:       ac.trustedKeys,
	}

	if err := token.Verify(verifyOpts); err != nil {
		challenge.err = err
		return nil, challenge
	}

	tokenScopes := token.scopes(resource)
	for _, scope := range scopes {
		match := false
		for _, tokenScope := range tokenScopes {
			if scope.Compare(tokenScope) {
				match = true
				break
			}
		}
		if !match {
			challenge.err = ErrInsufficientScope
			return nil, challenge
		}
	}

	return &auth.User{token.Claims.Subject}, nil

}

func parseToken(r *http.Request) (*Token, error) {
	parts := strings.Split(r.Header.Get("Authorization"), " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, ErrTokenRequired
	}

	rawToken := parts[1]

	token, err := NewToken(rawToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}
