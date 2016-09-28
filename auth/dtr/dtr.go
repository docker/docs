package dtr

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/tlsconfig"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/auth/builtin"
)

var (
	defaultHTTPTimeout = 30 * time.Second
)

type (
	DtrAuthenticator struct {
		Url             string
		TLSClientConfig *tls.Config
		builtin         *builtin.BuiltinAuthenticator
	}

	DtrAuthToken struct {
		SessionSecret string
		CsrfSecret    string
	}

	DtrConfig struct {
		Url       string `json:"url"`
		AdminUser string `json:"admin"`
		Insecure  bool   `json:"insecure"`
	}
)

var _ auth.Authenticator = (*DtrAuthenticator)(nil)

func NewAuthenticator(store *kvstore.Store, orcaID string, certPEM, keyPEM []byte, url string, allowInsecure bool) *DtrAuthenticator {
	builtinAuthenticator := builtin.NewAuthenticator(store, orcaID, certPEM, keyPEM)

	tlsClientConfig, err := tlsconfig.Client(tlsconfig.Options{InsecureSkipVerify: allowInsecure})
	if err != nil {
		log.Fatalf("Couldn't set up TLS Configuration to Docker Trusted Registry: %s", err)
	}

	return &DtrAuthenticator{
		Url:             url,
		TLSClientConfig: tlsClientConfig,
		builtin:         builtinAuthenticator,
	}
}

func (d DtrAuthenticator) doRequest(method string, path string, body []byte, headers map[string]string) ([]byte, []*http.Cookie, error) {
	b := bytes.NewBuffer(body)

	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: d.TLSClientConfig},
		Timeout:   defaultHTTPTimeout,
	}

	req, err := http.NewRequest(method, path, b)
	if err != nil {
		log.Errorf("couldn't create request: %s", err)
		return nil, nil, err
	}

	if headers != nil {
		for header, value := range headers {
			req.Header.Add(header, value)
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		if err == http.ErrHandlerTimeout {
			log.Error("Login timed out to Docker Trusted Registry")
			return nil, nil, err
		}
		log.Errorf("There was an error while authenticating: %s", err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp != nil && resp.StatusCode == 401 {
		// Unauthorized
		return nil, nil, auth.ErrUnauthorized
	} else if resp.StatusCode >= 400 {
		log.Errorf("Docker Trusted Registry returned an unexpected status code while authenticating: %s", resp.Status)
		return nil, nil, auth.ErrUnknown
	}

	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("couldn't read body: %s", err)
		return nil, nil, err
	}

	return rBody, resp.Cookies(), nil
}

func (d DtrAuthenticator) Login(username, password string) (*DtrAuthToken, error) {
	url := fmt.Sprintf("%s/%s", d.Url, "admin/login")
	data := []byte(fmt.Sprintf("username=%s&password=%s", neturl.QueryEscape(username), neturl.QueryEscape(password)))
	headers := map[string]string{
		"Content-type": "application/x-www-form-urlencoded",
	}

	_, cookies, err := d.doRequest("POST", url, data, headers)
	if err != nil {
		return nil, err
	}

	var sessSecret, csrfSecret string
	sessRegex := regexp.MustCompile("^session=(.*?);.*")
	csrfRegex := regexp.MustCompile("^csrftoken=(.*?);.*")

	for _, cookie := range cookies {
		if strings.HasPrefix(cookie.String(), "session") {
			sessSecret = sessRegex.FindStringSubmatch(cookie.String())[1]
		} else if strings.HasPrefix(cookie.String(), "csrftoken") {
			csrfSecret = csrfRegex.FindStringSubmatch(cookie.String())[1]
		}
	}

	// Warn if we don't have the CSRF secret.  Older version of DTR?
	if csrfSecret == "" {
		log.Warn("couldn't obtain csrf secret - no POST operations permitted to DTR")
	}
	if sessSecret == "" {
		// DTR didn't throw an error but it didn't return any credentials
		log.Errorf("Docker Trusted Registry didn't return credentials for user: %s", username)
		return nil, auth.ErrUnauthorized
	}

	authToken := &DtrAuthToken{
		SessionSecret: sessSecret,
		CsrfSecret:    csrfSecret,
	}

	return authToken, nil
}

func (d DtrAuthenticator) Name() string {
	return "dtr"
}

func (d DtrAuthenticator) AuthenticateUsernamePassword(username, password string) (*auth.Context, error) {
	dtrToken, err := d.Login(username, password)

	if err != nil {
		return nil, err
	}

	acct, err := d.GetUser(nil, username)
	if err == auth.ErrAccountDoesNotExist {
		acct = &auth.Account{
			Username: username,
		}

		_, err := d.SaveUser(nil, acct)
		if err != nil {
			log.Errorf("error saving account: %s", err)
			return nil, err
		}
	}

	// check account
	extraClaims := map[string]string{
		"sess": dtrToken.SessionSecret,
		"csrf": dtrToken.CsrfSecret,
	}

	tokenStr, err := d.builtin.GenerateToken(username, extraClaims)
	if err != nil {
		return nil, err
	}

	ctx := &auth.Context{
		User:         acct,
		SessionToken: tokenStr,
	}

	return ctx, nil
}

func (d DtrAuthenticator) AuthenticateSessionToken(tokenStr string) (*auth.Context, error) {
	return d.builtin.AuthenticateSessionToken(tokenStr)
}

func (d DtrAuthenticator) AuthenticatePublicKey(pubKey crypto.PublicKey) (*auth.Context, error) {
	return d.builtin.AuthenticatePublicKey(pubKey)
}

func (d DtrAuthenticator) AddUserPublicKey(user *auth.Account, label string, publicKey crypto.PublicKey) error {
	return d.builtin.AddUserPublicKey(user, label, publicKey)
}

func (d DtrAuthenticator) Logout(ctx *auth.Context) error {
	return d.builtin.Logout(ctx)
}

func (d DtrAuthenticator) GetUser(ctx *auth.Context, username string) (*auth.Account, error) {
	// TODO: We could pass back DTR account info here, instead of our cached account
	return d.builtin.GetUser(ctx, username)
}

func (d DtrAuthenticator) ListUsers(ctx *auth.Context) ([]*auth.Account, error) {
	// TODO: We could pass back DTR account info here, instead of our cached account
	return d.builtin.ListUsers(ctx)
}

func (d DtrAuthenticator) ListTeamMembers(ctx *auth.Context, teamID string) ([]*auth.Account, error) {
	return nil, auth.ErrUnsupported
}

func (d DtrAuthenticator) DeleteUser(ctx *auth.Context, account *auth.Account) error {
	return d.builtin.DeleteUser(ctx, account)
}

func (d DtrAuthenticator) SaveUser(ctx *auth.Context, account *auth.Account) (string, error) {
	return d.builtin.SaveUser(ctx, account)
}

func (d DtrAuthenticator) CanChangePassword(ctx *auth.Context) bool {
	return false
}

func (d DtrAuthenticator) ChangePassword(ctx *auth.Context, username, oldPassword, newPassword string) error {
	return auth.ErrUnsupported
}

func (d DtrAuthenticator) GetTeam(ctx *auth.Context, id string) (*auth.Team, error) {
	return nil, auth.ErrUnsupported
}

func (d DtrAuthenticator) ListTeams(ctx *auth.Context) ([]*auth.Team, error) {
	return nil, auth.ErrUnsupported
}

func (d DtrAuthenticator) ListUserTeams(ctx *auth.Context, username string) ([]*auth.Team, error) {
	return nil, auth.ErrUnsupported
}

func (d DtrAuthenticator) SaveTeam(ctx *auth.Context, team *auth.Team) (string, error) {
	return "", auth.ErrUnsupported
}

func (d DtrAuthenticator) DeleteTeam(ctx *auth.Context, team *auth.Team) error {
	return auth.ErrUnsupported
}

func (d DtrAuthenticator) AddTeamMember(ctx *auth.Context, teamID, username string) error {
	return auth.ErrUnsupported
}

func (d DtrAuthenticator) DeleteTeamMember(ctx *auth.Context, teamID, username string) error {
	return auth.ErrUnsupported
}
func (d DtrAuthenticator) Sync(ctx *auth.Context, force, onlyAdmin bool) error {
	// nothing to see here, move along
	return nil
}
func (d DtrAuthenticator) LastSyncStatus(ctx *auth.Context) string {
	// nothing to see here, move along
	return ""
}
