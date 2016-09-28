package apiclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/docker/dhe-deploy/garant/authn/enzi"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	// "github.com/docker/dhe-deploy/garant/authn/enzi"

	enziclient "github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/forms"
	enziresponses "github.com/docker/orca/enzi/api/responses"
	"github.com/gorilla/websocket"
)

func (c *apiClient) makeRequest(method string, u url.URL, payload interface{}) (*http.Response, error) {
	req, err := c.createRequest(method, u, payload)
	if err != nil {
		return nil, err
	}

	c.AuthenticateRequest(req)

	return dtrutil.DoRequestWithClient(req, c.client)
}

func (c *apiClient) openUnauthenticatedWebsocket(route string) (*websocket.Conn, error) {
	return c.openWebsocket(route, nil)
}

func (c *apiClient) openAuthenticatedWebsocket(route string) (*websocket.Conn, error) {
	header := http.Header{}
	header.Add("Cookie", "session="+c.sessionSecret)

	return c.openWebsocket(route, header)
}

func (c *apiClient) openWebsocket(route string, header http.Header) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: c.host, Path: route}
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *apiClient) makeUnauthenticatedRequest(method string, u url.URL, payload interface{}) (*http.Response, error) {
	req, err := c.createRequest(method, u, payload)
	if err != nil {
		return nil, err
	}

	return dtrutil.DoRequestWithClient(req, c.client)
}

func (c *apiClient) createRequest(method string, u url.URL, payload interface{}) (*http.Request, error) {
	var reader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	if u.Host == "" {
		u.Host = c.host
	}
	if u.Scheme == "" {
		u.Scheme = c.apiClientUrlScheme
	}
	// only specify the port in non-standard cases
	if (c.apiClientUrlScheme == "https" && c.apiClientPort != 443) ||
		(c.apiClientUrlScheme == "http" && c.apiClientPort != 80) {
		u.Host = u.Host + ":" + strconv.Itoa(int(c.apiClientPort)) // ex: https://example.com:8080
	}

	req, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	if method != "HEAD" && method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	// always close the connection after sending the request so that we
	// don't get bitten by net/http's bugs with reusing connections
	req.Close = true

	return req, nil
}

// hack to make the go http client not follow redirects
var RedirectAttemptedError = errors.New("redirect")

func (c *apiClient) AuthenticateRequest(req *http.Request) {
	req.AddCookie(&http.Cookie{
		Name:     "session",
		Value:    c.sessionSecret,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 365 * 5),
	})

	if req.Method != "GET" {
		req.Header.Add(enzi.CSRFTokenHeaderName, c.csrfCookie)
	}
}

// XXX: nothing checks for login auth failures. We need to fix that.
func (c *apiClient) Login(username, password string) error {
	// ----- DTR Login (openid)

	// ----- OpenID Begin
	beginResp, err := c.makeUnauthenticatedRequest("GET", url.URL{Path: "/api/v0/openid/begin"}, nil)
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error making request to openid/begin %s", err)
	}
	defer beginResp.Body.Close()

	if beginResp.StatusCode != http.StatusFound {
		return fmt.Errorf("Unexpected status when hitting openid begin path %d", beginResp.StatusCode)
	}
	redirectTo := beginResp.Header.Get("Location")

	// This is the URL DTR redirects us to for auth (probably /login)
	redirectUrl, err := url.Parse(redirectTo)
	if err != nil {
		return err
	}

	c.ucpAsEnzi = false
	// figure out if we are using ucp or enzi for auth
	if redirectUrl.Path[:5] == "/enzi" {
		c.ucpAsEnzi = true
	}

	c.enziHost = redirectUrl.Host

	// at this point the browse is redirected to enzi to log in

	// ----- Enzi Login
	prefix := "/"
	if c.ucpAsEnzi {
		prefix = "/enzi"
	}
	session := enziclient.New(c.client, c.enziHost, prefix, nil)

	loginForm := forms.Login{
		Username: username,
		Password: password,
	}

	loginSessionResp, err := session.Login(loginForm)
	if err != nil {
		return fmt.Errorf("error making request to enzi login %s", err)
	}
	c.loginSession = loginSessionResp

	// ----- Enzi Authoriaztioin
	// we skip that part and overwrite the path to do just the submission of the authorization form
	if c.ucpAsEnzi {
		redirectUrl.Path = "/enzi/v0/id/authorize"
	} else {
		redirectUrl.Path = "/v0/id/authorize"
	}

	req, err := http.NewRequest("POST", redirectUrl.String(), strings.NewReader(redirectUrl.RawQuery))
	if err != nil {
		return fmt.Errorf("error making request to authorize endpoint %s", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", fmt.Sprintf("SessionToken %s", loginSessionResp.SessionToken))
	req.Close = true

	enziResp, err := dtrutil.DoRequestWithClient(req, c.client)
	if err != nil {
		return err
	}
	defer enziResp.Body.Close()

	type redirectResponse struct {
		Redirect string `json:"redirect"`
	}
	redirect := redirectResponse{}

	if err := json.NewDecoder(enziResp.Body).Decode(&redirect); err != nil {
		return fmt.Errorf("unable to decode API response: %s", err)
	}

	// ----- Hit the OpenID Callback on DTR's side to prove to DTR that we are authorized
	req2, err := http.NewRequest("GET", redirect.Redirect, nil)
	if err != nil {
		return fmt.Errorf("error making request object for callback %s", err)
	}
	req2.Close = true

	// Add the csrf cookie from the openIDBegin response
	// XXX: is this necessary?
	for _, cookie := range beginResp.Cookies() {
		req2.AddCookie(cookie)
	}

	callback, err := dtrutil.DoRequestWithClient(req2, c.client)
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error making request to callback %s", err)
	}
	defer callback.Body.Close()

	for _, cookie := range callback.Cookies() {
		if cookie.Name == "session" {
			c.sessionSecret = cookie.Value
		} else if cookie.Name == "csrftoken" {
			c.csrfCookie = cookie.Value
		}
	}

	return nil
}

func (c *apiClient) Logout() error {
	logoutResp, err := c.makeRequest("GET", url.URL{Path: "/logout"}, nil)
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("error making request to logout %s", err)
	}
	defer logoutResp.Body.Close()

	if logoutResp.StatusCode != http.StatusSeeOther {
		return fmt.Errorf("Unexpected status when hitting logout %d", logoutResp.StatusCode)
	}
	redirectTo := logoutResp.Header.Get("Location")

	redirectUrl, err := url.Parse(redirectTo)
	if err != nil {
		return err
	}

	redirectUrl.RawQuery = ""
	if c.ucpAsEnzi {
		redirectUrl.Path = "/enzi/v0/id/logout"
	} else {
		redirectUrl.Path = "/v0/id/logout"
	}

	req, err := http.NewRequest("POST", redirectUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating request to logout endpoint %s", err)
	}

	if c.loginSession != nil {
		req.Header.Set("Authorization", fmt.Sprintf("SessionToken %s", c.loginSession.SessionToken))
	} else {
		// not logged in
		return nil
	}
	req.Close = true

	_, err = dtrutil.DoRequestWithClient(req, c.client)
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		err = nil
	}
	c.loginSession = nil

	return err
}

func (c *apiClient) makeURLEncodedRequest(method, route string, form url.Values) (*http.Response, error) {
	reader := strings.NewReader(form.Encode())

	// url builder
	var url string
	url = c.apiClientUrlScheme + "://" + c.host // ex: https://example.com
	// only specify the port in non-standard cases
	if (c.apiClientUrlScheme == "https" && c.apiClientPort != 443) ||
		(c.apiClientUrlScheme == "http" && c.apiClientPort != 80) {
		url = url + ":" + strconv.Itoa(int(c.apiClientPort)) // ex: https://example.com:8080
	}

	url = url + route // ex: http://example.com:8080/resource

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	// always close the connection after sending the request so that we
	// don't get bitten by net/http's bugs with reusing connections
	req.Close = true

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	return dtrutil.DoRequestWithClient(req, c.client)
}

type EnziAuthenticator struct {
	loginSession *enziresponses.LoginSession
}

// AuthenticateRequest sets a basic authorization header on the given request.
func (a *EnziAuthenticator) AuthenticateRequest(req *http.Request) {
	if a.loginSession != nil {
		req.Header.Set("Authorization", fmt.Sprintf("SessionToken %s", a.loginSession.SessionToken))
	}
}

func (c *apiClient) EnziSession() *enziclient.Session {

	prefix := "/"
	if c.ucpAsEnzi {
		prefix = "/enzi"
	}
	return enziclient.New(c.client, c.enziHost, prefix, &EnziAuthenticator{c.loginSession})
}
