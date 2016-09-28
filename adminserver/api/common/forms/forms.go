package forms

import (
	"net/url"
	"strconv"

	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
)

type CreateRepo struct {
	Name             string `json:"name"`
	ShortDescription string `json:"shortDescription"`
	LongDescription  string `json:"longDescription"`
	Visibility       string `json:"visibility,omitempty" enum:"public|private"`
}

type Access struct {
	responses.WithAccessLevel
}

type UpdateRepo struct {
	ShortDescription *string `json:"shortDescription,omitempty"`
	LongDescription  *string `json:"longDescription,omitempty"`
	Visibility       *string `json:"visibility,omitempty" enum:"public|private"`
}

type Team struct {
	Name                     string `json:"name"`
	Description              string `json:"description"`
	Type                     string `json:"type" enum:"managed|ldap"`
	LdapDN                   string `json:"ldapDN,omitempty"`
	LdapGroupMemberAttribute string `json:"ldapGroupMemberAttribute,omitempty"`
}

type TeamUpdate struct {
	Name                     string  `json:"name,omitempty"`
	Description              *string `json:"description,omitempty"`
	Type                     string  `json:"type,omitempty"`
	LdapDN                   string  `json:"ldapDN,omitempty"`
	LdapGroupMemberAttribute string  `json:"ldapGroupMemberAttribute,omitempty"`
}

// accountForm is used only for documentation purposes. It combines userForm, organizationForm and ldap.Settings
type Account struct {
	Type      string `json:"type" enum:"user|organization" description:"user or organization" modelDescription:"A user or organization account"`
	Name      string `json:"name" description:"The user or org's namespace. It can contain lowercase letters, numbers, - and _. It must start with a letter or number."`
	LDAPLogin string `json:"ldapLogin,omitempty"`
	Password  string `json:"password,omitempty"`
}

type User struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	LDAPLogin string `json:"ldapLogin,omitempty"`
}

type Organization struct {
	Name string `json:"name"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword"`
}

type SearchOptions struct {
	IncludeRepositories bool
	IncludeAccounts     bool
	Namespace           string
	Query               string
	Limit               uint
}

type Settings struct {
	// TODO: document the enum fields
	DTRHost            *string `json:"dtrHost,omitempty"`
	AuthBypassCA       *string `json:"authBypassCA,omitempty"`
	AuthBypassOU       *string `json:"authBypassOU,omitempty"`
	DisableUpgrades    *bool   `json:"disableUpgrades,omitempty"`
	ReportAnalytics    *bool   `json:"reportAnalytics,omitempty"`
	AnonymizeAnalytics *bool   `json:"anonymizeAnalytics,omitempty"`
	ReleaseChannel     *string `json:"releaseChannel,omitempty"`
	WebTLSCert         *string `json:"webTLSCert,omitempty"`
	WebTLSKey          *string `json:"webTLSKey,omitempty"`
	WebTLSCA           *string `json:"webTLSCA,omitempty"`
	GCMode             *string `json:"gcMode,omitempty"`
}

func (opts SearchOptions) GetURLValues() url.Values {
	urlValues := url.Values{}
	if opts.IncludeRepositories {
		urlValues.Add("includeRepositories", "true")
	}
	if opts.IncludeAccounts {
		urlValues.Add("includeAccounts", "true")
	}
	urlValues.Add("namespace", opts.Namespace)
	urlValues.Add("query", opts.Query)
	urlValues.Add("limit", strconv.FormatUint(uint64(opts.Limit), 10))
	return urlValues
}
