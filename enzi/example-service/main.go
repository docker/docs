package main

import (
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/docker/orca/enzi/api/client/openid"
	"github.com/docker/orca/enzi/api/client/openid/oautherrors"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/docker/orca/enzi/jose"
)

// The name of the CSRF Token header required for non-HEAD/GET requests.
const CSRFTokenHeaderName = "X-Csrf-Token"

func main() {
	var (
		providerHost       string
		serviceID          string
		serviceHost        string
		providerPathPrefix string
	)

	flag.StringVar(&providerHost, "provider-host", "api.enzi:4443", "the host address of the eNZi OpenID Connect Provider")
	flag.StringVar(&serviceID, "service-id", "", "this service's ID on the eNZi OpenID Connect Provider")
	flag.StringVar(&serviceHost, "service-host", "example-service.enzi:4043", "the host address of this service; used for OpenID Connect authorization redirect URL")
	flag.StringVar(&providerPathPrefix, "provider-path-prefix", "", "prefix all provider HTTP routes with this path")
	flag.Parse()

	log.Print("initializing service ...")

	rootCertsPEM, err := ioutil.ReadFile("/tls/ca.pem")
	if err != nil {
		log.Fatalf("unable to read root CA certificates: %s", err)
	}

	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(rootCertsPEM) {
		log.Fatalf("unable to parse root CA certificates")
	}

	tlsConfig := &tls.Config{RootCAs: rootCAs}

	exampleService, err := NewExampleService(tlsConfig, providerHost, serviceID, serviceHost, providerPathPrefix)
	if err != nil {
		log.Fatalf("unable to initialize service: %s", err)
	}

	server := &http.Server{
		Addr:    ":4043",
		Handler: exampleService,
	}

	log.Print("listening for connections ...")
	log.Fatal(server.ListenAndServeTLS("/tls/cert.pem", "/tls/key.pem"))
}

// ExampleService is a simple example of an OpenID Connect Client which uses
// eNZi for authentication, accounts, and session management.
type ExampleService struct {
	*http.ServeMux

	signingKey     *jose.PrivateKey
	providerClient *openid.Client
}

// NewExampleService returns a new ExampleService.
func NewExampleService(tlsConfig *tls.Config, providerHost, serviceID, serviceHost, providerPathPrefix string) (*ExampleService, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("unable to generate 2048-bit RSA key: %s", err)
	}

	signingKey, err := jose.NewPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to create JWT signing key: %s", err)
	}

	redirectURI := fmt.Sprintf("https://%s/openid_callback", serviceHost)

	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			TLSClientConfig:       tlsConfig,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			MaxIdleConnsPerHost:   5,
		},
	}

	authorizationPath := path.Join("/", providerPathPrefix, "authorize")
	tokenPath := path.Join("/", providerPathPrefix, "v0/id/token")

	service := &ExampleService{
		ServeMux:       http.NewServeMux(),
		signingKey:     signingKey,
		providerClient: openid.NewClient(httpClient, signingKey, serviceID, redirectURI, providerHost, authorizationPath, tokenPath),
	}

	service.registerHandlers()

	return service, nil
}

func (service *ExampleService) registerHandlers() {
	routes := map[string]http.HandlerFunc{
		"/":                service.handleHome,
		"/openid_begin":    service.handleOpenIDBegin,
		"/openid_callback": service.handleOpenIDCallback,
		"/openid_error":    service.handleOpenIDError,
		"/openid_keys":     service.handleOpenIDKeys,
	}

	for route, handler := range routes {
		service.HandleFunc(route, handler)
	}
}

// authenticateSession attempts to authenticate the given request using
// session cookies. If there is no session cookie OR the session is invalid,
// a nil account is returned. If there is any error, a non-nil error is
// returned with details.
func (service *ExampleService) authenticateSession(r *http.Request) *responses.Account {
	sessionCookie, _ := r.Cookie("session") // Error would be http.ErrNoCookie.
	if sessionCookie == nil || sessionCookie.Value == "" {
		// Handle like an unauthenticated request.
		return nil
	}

	tokenResponse, err := service.providerClient.GetTokenWithServiceSession(sessionCookie.Value)
	if err != nil {
		log.Printf("unable to get identity token from service session: %s", err)
		// Handle like an unauthenticated request.
		return nil
	}

	csrfToken := tokenResponse.SessionCSRFToken
	if !(r.Method == "GET" || r.Method == "HEAD" || csrfToken == r.Header.Get(CSRFTokenHeaderName) || csrfToken == r.FormValue("csrftoken")) {
		// Handle like an unauthenticated request.
		return nil
	}

	return tokenResponse.Account
}

func (service *ExampleService) handleHome(w http.ResponseWriter, r *http.Request) {
	if account := service.authenticateSession(r); account != nil {
		// Render a hello message to the user.
		w.Write([]byte(fmt.Sprintf("Hello, %s!", account.Name)))
		return
	}

	// Redirect to the authentication login process.
	redirectParams := url.Values{}
	redirectParams.Set("next", r.URL.String())

	redirectURL := fmt.Sprintf("/openid_begin?%s", redirectParams.Encode())

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (service *ExampleService) handleOpenIDBegin(w http.ResponseWriter, r *http.Request) {
	authorizeURL, redirectStateKey := service.providerClient.PrepareAuthorizationRequest(r.URL.Query().Get("next"))

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_redirect_state_key",
		Value:    redirectStateKey,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 5),
	})

	http.Redirect(w, r, authorizeURL, http.StatusFound)
}

func validateRedirectState(r *http.Request) *openid.RedirectState {
	encodedState := r.URL.Query().Get("state")
	if encodedState == "" {
		log.Print("no redirect state in OpenID Connect authorization callback")
		return nil
	}

	cookie, _ := r.Cookie("csrf_redirect_state_key")
	if cookie == nil || cookie.Value == "" {
		log.Print("OpenID Connect authorization callback csrf redirect state key cookie not set")
		return nil
	}

	redirectState, err := openid.DecodeRedirectState(encodedState, cookie.Value)
	if err != nil {
		log.Printf("unable to decode OpenID Connect authorization callback redirect state: %s", err)
		return nil
	}

	return redirectState
}

func (service *ExampleService) handleOpenIDCallback(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	state := validateRedirectState(r)
	if state == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if errCode := params.Get("error"); errCode != "" {
		errorParams := url.Values{}
		errorParams.Set("error", errCode)
		errorParams.Set("error_description", params.Get("error_description"))

		http.Redirect(w, r, fmt.Sprintf("/openid_error?%s", errorParams.Encode()), http.StatusFound)
		return
	}

	code := params.Get("code")
	if code == "" {
		log.Print("OpenID Connect Provider authorize endpoint returned no code")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tokenResponse, err := service.providerClient.GetTokenWithAuthorizationCode(code, state, true)
	if err != nil {
		errorParams := url.Values{}
		if oauthErr, ok := err.(*oautherrors.ErrorResponse); ok {
			errorParams.Set("error", oauthErr.Code)
			errorParams.Set("error_description", oauthErr.Description)
		} else {
			log.Printf("unable to get identity token: %s", err)
			errorParams.Set("error", "unable to get token")
			errorParams.Set("error_description", "see server logs for details")
		}

		http.Redirect(w, r, fmt.Sprintf("/openid_error?%s", errorParams.Encode()), http.StatusFound)
		return
	}

	// Set browser cookies to expire at some point that is very, very far
	// in the future. Sessions should be managed by users with the auth
	// provider, not their browser, unless they wish to delete cookies.
	expiration := time.Now().Add(time.Hour * 24 * 365 * 5) // 5 years.

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    tokenResponse.SessionSecret,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  expiration,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrftoken",
		Value:    tokenResponse.SessionCSRFToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: false, // Available to javascript.
		Expires:  expiration,
	})

	if state.RedirectNext == "" {
		state.RedirectNext = "/"
	}

	http.Redirect(w, r, state.RedirectNext, http.StatusFound)
}

func (service *ExampleService) handleOpenIDError(w http.ResponseWriter, r *http.Request) {
	// Render an error message to the user.
	w.Write([]byte(fmt.Sprintf("OpenID Connect Error\n\n%s\n\n%s", r.URL.Query().Get("error"), r.URL.Query().Get("error_description"))))
}

func (service *ExampleService) handleOpenIDKeys(w http.ResponseWriter, r *http.Request) {
	response := map[string][]jose.PublicKey{
		"keys": {service.signingKey.PublicKey},
	}

	// The provider should cache our keys for 1 hour.
	w.Header().Set("Cache-Control", "max-age=3600")

	json.NewEncoder(w).Encode(response)
}
