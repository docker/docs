package auth

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	enziclient "github.com/docker/orca/enzi/api/client"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidSessionToken = errors.New("invalid session token")
	ErrInvalidPublicKey    = errors.New("invalid public key")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNoUserInToken       = errors.New("no user sent in token")
	ErrAccountDisabled     = errors.New("this account has been disabled")
	ErrAccountExists       = errors.New("account already exists")
	ErrAccountDoesNotExist = errors.New("account does not exist")
	ErrTeamExists          = errors.New("team already exists")
	ErrTeamDoesNotExist    = errors.New("team does not exist")
	ErrUnsupported         = errors.New("feature unsupported")
	ErrUnknown             = errors.New("unknown error")
)

const (
	AuthenticatorBuiltin = iota
	AuthenticatorDTR
	AuthenticatorLDAP
	AuthenticatorEnzi
)

type (
	BuiltinAuthenticatorConfig struct{}

	DTRAuthenticatorConfig struct {
		AdminUser string `json:"admin_user"`
		URL       string `json:"url"`
		Insecure  bool   `json:"insecure,omitempty"`
	}

	LDAPSettings struct {
		AdminUsername string `json:"adminUsername"`
		AdminPassword string `json:"adminPassword"` // Used only during configuration change

		ServerURL     string `json:"serverURL"`
		StartTLS      bool   `json:"startTLS"`
		RootCerts     string `json:"rootCerts"`
		TLSSkipVerify bool   `json:"tlsSkipVerify"`

		ReaderDN       string `json:"readerDN"`
		ReaderPassword string `json:"readerPassword"`

		UserBaseDN        string `json:"userBaseDN"`
		UserLoginAttrName string `json:"userLoginAttrName"`
		UserSearchFilter  string `json:"userSearchFilter"`
		UserDefaultRole   Role   `json:"userDefaultRole"`

		SyncInterval int `json:"syncInterval"` // Units: Minutes
	}

	EnziConfig struct {
		ServiceID     string   `json:"serviceID"`
		DefaultOrgID  string   `json:"defaultOrgID"`
		ProviderAddrs []string `json:"providerAddrs"`

		UserDefaultRole Role `json:"userDefaultRole"`
	}

	AuthenticatorConfiguration struct {
		AuthenticatorType int                        `json:"auth_type"`
		BuiltinConfig     BuiltinAuthenticatorConfig `json:"builtin_auth"`
		DTRConfig         DTRAuthenticatorConfig     `json:"dtr_auth"`
		LDAPConfig        LDAPSettings               `json:"ldap_auth"`
		EnziConfig        EnziConfig                 `json:"enziConfig"`
	}

	AccountKey struct {
		Label     string `json:"label,omitempty"` // User editable label to identify this key
		KeyID     string `json:"keyID"`           // Hash of public key DER bytes.
		UserID    string `json:"userID"`          // ID of user which the key belongs to.
		PublicKey string `json:"public_key,omitempty"`
	}

	Account struct {
		ID              string       `json:"id"` // ID of user, only used with eNZi backend.
		FirstName       string       `json:"first_name,omitempty"`
		LastName        string       `json:"last_name,omitempty"`
		Username        string       `json:"username,omitempty"`
		Password        string       `json:"password,omitempty"`
		Admin           bool         `json:"admin"`
		PublicKeys      []AccountKey `json:"public_keys,omitempty"`
		Role            Role         `json:"role"`
		ManagedTeams    []string     `json:"managed_teams,omitempty"`
		DiscoveredTeams []string     `json:"discovered_teams,omitempty"` // Discovered teams are based on LDAP group membership
		LdapDN          string       `json:"ldap_dn,omitempty"`          // Maps to the LDAP user for this account
		Disabled        bool         `json:"disabled,omitempty"`         // Only used for Discovered LDAP accounts
	}

	AuthToken struct {
		ID        string `json:"token_id,omitempty"`
		Token     string `json:"auth_token,omitempty"`
		UserAgent string `json:"user_agent,omitempty"`
	}

	Context struct {
		// The authenticated client user.
		User *Account
		// Non-empty if a session token was used to authenticate.
		SessionToken string
		// Used by the eNZi backend for making requests to the
		// auth provider. Not used by other backends.
		ClientCreds enziclient.RequestAuthenticator
	}

	Authenticator interface {
		Name() string

		AuthenticateUsernamePassword(username, password string) (*Context, error)
		AuthenticateSessionToken(token string) (*Context, error)
		AuthenticatePublicKey(pubKey crypto.PublicKey) (*Context, error)
		Logout(ctx *Context) error

		CanChangePassword(ctx *Context) bool
		ChangePassword(ctx *Context, username, oldPassword, newPassword string) error

		AddUserPublicKey(user *Account, label string, publicKey crypto.PublicKey) error

		SaveUser(ctx *Context, account *Account) (string, error)
		GetUser(ctx *Context, username string) (*Account, error)
		ListUsers(ctx *Context) ([]*Account, error)
		DeleteUser(ctx *Context, account *Account) error

		SaveTeam(ctx *Context, team *Team) (string, error)
		GetTeam(ctx *Context, id string) (*Team, error)
		ListTeams(ctx *Context) ([]*Team, error)
		DeleteTeam(ctx *Context, team *Team) error

		AddTeamMember(ctx *Context, teamID, username string) error
		ListTeamMembers(ctx *Context, teamID string) ([]*Account, error)
		ListUserTeams(ctx *Context, username string) ([]*Team, error)
		DeleteTeamMember(ctx *Context, teamID, username string) error

		Sync(ctx *Context, force, adminOnly bool) error
		LastSyncStatus(ctx *Context) string
	}
)

func (c *AuthenticatorConfiguration) Equals(other *AuthenticatorConfiguration) bool {
	return ((c.AuthenticatorType == other.AuthenticatorType) &&
		(c.BuiltinConfig == other.BuiltinConfig) &&
		(c.DTRConfig == other.DTRConfig) &&
		(c.LDAPConfig == other.LDAPConfig) &&
		(c.EnziConfig.Equals(other.EnziConfig)))
}

func (c *EnziConfig) Equals(other EnziConfig) bool {
	if (c.DefaultOrgID != other.DefaultOrgID) || (c.ServiceID != other.ServiceID) {
		return false
	}

	if len(c.ProviderAddrs) != len(other.ProviderAddrs) {
		return false
	}

	if c.UserDefaultRole != other.UserDefaultRole {
		return false
	}

	for i := range c.ProviderAddrs {
		if c.ProviderAddrs[i] != other.ProviderAddrs[i] {
			return false
		}
	}

	return true
}

func Hash(data string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(h[:]), err
}

func GenerateToken() (string, error) {
	return Hash(time.Now().String())
}

func GenerateRandID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	r := base64.URLEncoding.EncodeToString(b)
	return r, nil
}
