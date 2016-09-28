package forms

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/passwords"
	"github.com/docker/orca/enzi/worker"
	"github.com/robfig/cron"
	"golang.org/x/text/unicode/norm"
)

var (
	// validAccountNamePattern restricts account names to contain only
	// characters which are in the unicode categories "Letter, lowercase",
	// "Letter, other", arabic numerals 0 through 9, `@` signs, periods,
	// underscores, and hiphens. This pattern also ensures that the account
	// name is at least 1 character in length.
	validAccountNamePattern = regexp.MustCompile(`^[\p{Ll}\p{Lo}0-9$@._-]+$`)
	maxNameLength           = 100 // Some people have really long names.
	minPasswordLength       = 1
	maxDescriptonLength     = 140 // Should be succinct enough to be a Tweet.
)

type validator interface {
	Validate() []*errors.APIError
}

func validateJSONForm(r io.Reader, form validator) []*errors.APIError {
	if err := json.NewDecoder(r).Decode(form); err != nil {
		return []*errors.APIError{errors.InvalidJSON(err)}
	}

	return form.Validate()
}

// NormalizeAccountName normalizes the given account name string to be
// lowercased and use the Unicode "Composing Normal Form".
func NormalizeAccountName(name string) string {
	return norm.NFC.String(strings.ToLower(name))
}

// NormalizeFullName normalizes the given full name using the Unicode
// "Composing Normal Form" and removes leading and trailing space characters.
func NormalizeFullName(name string) string {
	return strings.TrimSpace(norm.NFC.String(name))
}

// ValidateAccountName validates that the given name can be used as a valid
// account name. The given name is first normalized using the
// NormalizeAccountName function (this is why it is passed in as a pointer).
// If the normalized name is not valid, returns a non-nil API error describing
// the validation error.
func ValidateAccountName(name *string, fieldName string) *errors.APIError {
	*name = NormalizeAccountName(*name)
	if len(*name) > maxNameLength || !validAccountNamePattern.MatchString(*name) {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("must match the pattern %s and be no longer than %d characters", validAccountNamePattern, maxNameLength))
	}

	return nil
}

// ValidatePassword validates that the given password meets our password
// requirements. If the password is not valid, returns a non-nil API error
// describing the validation error.
func ValidatePassword(password, fieldName string) *errors.APIError {
	if len(password) < minPasswordLength {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("must be at least %d characters", minPasswordLength))
	}

	return nil
}

// ValidatePasswordHash validates that the given password hash is a valid
// password hash. If the password hash is not valid, returns a non-nil API error
// describing the validation error.
func ValidatePasswordHash(passwordHash, fieldName string) *errors.APIError {
	if err := passwords.CheckPasswordHash(passwordHash); err != nil {
		return errors.InvalidFormField(fieldName, err.Error())
	}
	return nil
}

// ValidateFullName validates that the given fullName meets our fullName
// requirements. If the fullName is not valid, returns a non-nil API error
// describing the validation error.
func ValidateFullName(fullName, fieldName string) *errors.APIError {
	if len(fullName) < 1 || len(fullName) > maxNameLength {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("must be no longer than %d characters", maxNameLength))
	}

	return nil
}

// ValidateCronSpec validates that the given schedule in CRON format is valid
// and ignores seconds and the special '@every <duration>' spec. The value of
// the schedule may be trimmed of whitespace and altered to meet the full
// length cron spec if necessary. If no schedule is specified, it is set to
// "@hourly" by default.
func ValidateCronSpec(schedule *string, fieldName string) *errors.APIError {
	*schedule = strings.TrimSpace(*schedule)

	if *schedule == "" {
		*schedule = "@hourly"
		return nil
	}

	if strings.HasPrefix(*schedule, "@every") {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("unrecognized descriptor: %s", *schedule))
	}

	// If the schedule is not one of the @hourly, @daily, etc., specs ...
	if len(*schedule) > 0 && (*schedule)[0] != '@' {
		// Split on whitespace.  We require exactly 5 fields. We don't
		// accept a seconds field. Crons can only be run with by-minute
		// granularity:
		// (minute) (hour) (day of month) (month) (day of week)
		fields := strings.Fields(*schedule)
		if fields[0] != "0" {
			return errors.InvalidFormField(fieldName, fmt.Sprintf("seconds field of cronspec must be '0', found %q", fields[0]))
		}
	}

	// Ensure that the cronspec can be parsed.
	if _, err := cron.Parse(*schedule); err != nil {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("unable to parse cron spec %q: %s", *schedule, err))
	}

	return nil
}

// ImportAccount is an entry in the form submitted to create many new users
// or organizations through the ImportAccounts endpoint.
type ImportAccount struct {
	Name     string `json:"name"               description:"Name of account"`
	FullName string `json:"fullName"           description:"Full name of account"`
	IsOrg    bool   `json:"isOrg,omitempty"    description:"Whether the account is an organization"`

	// Fields for users only.
	IsAdmin      bool   `json:"isAdmin,omitempty"  description:"Whether the user is an admin (users only)"`
	IsActive     bool   `json:"isActive,omitempty" description:"Whether the user is active and can login (users only)"`
	PasswordHash string `json:"passwordHash,omitempty" description:"The password hash for the user (users only) - either using PBKDF2 SHA256 or bcrypt"`
}

type ImportAccounts struct {
	Accounts []ImportAccount `json:"accounts"`
}

func (form *ImportAccount) validate(baseFieldName string) (apiErrs []*errors.APIError) {
	if apiErr := ValidateAccountName(&form.Name, fqFieldName(baseFieldName, "name")); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	if form.FullName != "" {
		form.FullName = NormalizeFullName(form.FullName)
		if apiErr := ValidateFullName(form.FullName, fqFieldName(baseFieldName, "fullName")); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if !form.IsOrg {
		if apiErr := ValidatePasswordHash(form.PasswordHash, fqFieldName(baseFieldName, "passwordHash")); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	return apiErrs
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *ImportAccounts) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this CreateAccount form. Returns any validation errors
// for their corresponding form fields.
func (form *ImportAccounts) Validate() (apiErrs []*errors.APIError) {
	for i := range form.Accounts {
		fieldPrefix := fmt.Sprintf("[%d].", i)
		if subFormErrs := form.Accounts[i].validate(fieldPrefix); subFormErrs != nil {
			apiErrs = append(apiErrs, subFormErrs...)
		}
	}

	return apiErrs
}

// CreateAccount is a form submitted to create a new user or organization.
type CreateAccount struct {
	Name     string `json:"name"               description:"Name of account"`
	FullName string `json:"fullName"           description:"Full name of account"`
	IsOrg    bool   `json:"isOrg,omitempty"    description:"Whether the account is an organization"`

	// Fields for users only.
	IsAdmin  bool   `json:"isAdmin,omitempty"  description:"Whether the user is an admin (users only)"`
	IsActive bool   `json:"isActive,omitempty" description:"Whether the user is active and can login (users only)"`
	Password string `json:"password,omitempty" description:"Password for the user (users only)"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *CreateAccount) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this CreateAccount form. Returns any validation errors
// for their corresponding form fields.
func (form *CreateAccount) Validate() (apiErrs []*errors.APIError) {
	if apiErr := ValidateAccountName(&form.Name, "name"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	if form.FullName != "" {
		form.FullName = NormalizeFullName(form.FullName)
		if apiErr := ValidateFullName(form.FullName, "fullName"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if !form.IsOrg {
		if apiErr := ValidatePassword(form.Password, "password"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	return apiErrs
}

// MemberSyncOpts is a form submitted to update options for syncing members of
// an organization or team.
type MemberSyncOpts struct {
	EnableSync         bool   `json:"enableSync"         description:"Whether to enable LDAP syncing. If false, all other fields are ignored"`
	SelectGroupMembers bool   `json:"selectGroupMembers" description:"Whether to sync using a group DN and member attribute selection or to use a search filter (if false)"`
	GroupDN            string `json:"groupDN"            description:"The distinguished name of the LDAP group. Required if selectGroupMembers is true, ignored otherwise"`
	GroupMemberAttr    string `json:"groupMemberAttr"    description:"The name of the LDAP group entry attribute which corresponds to distinguished names of members. Required if selectGroupMembers is true, ignored otherwise"`
	SearchBaseDN       string `json:"searchBaseDN"       description:"The distinguished name of the element from which the LDAP server will search for users. Required if selectGroupMembers is false, ignored otherwise"`
	SearchScopeSubtree bool   `json:"searchScopeSubtree" description:"Whether to search for users in the entire subtree of the base DN or to only search one level under the base DN (if false). Required if selectGroupMembers is false, ignored otherwise"`
	SearchFilter       string `json:"searchFilter"       description:"The LDAP search filter used to select users if selectGroupMembers is false, may be left blank"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *MemberSyncOpts) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this MemberSyncOpts form. Returns any validation errors
// for their corresponding form fields.
func (form *MemberSyncOpts) Validate() (apiErrs []*errors.APIError) {
	return form.validate("")
}

func (form *MemberSyncOpts) validate(baseFieldName string) (apiErrs []*errors.APIError) {
	if !form.EnableSync {
		return nil
	}

	if form.SelectGroupMembers {
		// Ensure these fields are not blank.
		if form.GroupDN == "" {
			apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "groupDN"), "this field can not be left blank"))
		}
		if form.GroupMemberAttr == "" {
			apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "groupMemberAttr"), "this field can not be left blank"))
		}
	} else if form.SearchBaseDN == "" {
		// Note: it's okay if the filter is left blank.
		apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "searchBaseDN"), "this field can not be left blank"))
	}

	return apiErrs
}

// UpdateAccount is a form submitted to update a user or organization.
type UpdateAccount struct {
	FullName *string `json:"fullName,omitempty" description:"Full name of account, unchanged if null or omitted"`

	// Fields for users only.
	IsAdmin  *bool `json:"isAdmin,omitempty"    description:"Whether the user is an admin (users only), unchanged if null or omitted"`
	IsActive *bool `json:"isActive,omitempty"   description:"Whether the user is active and can login (users only), unchanged if null or omitted"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *UpdateAccount) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this UpdateAccount form. Returns any validation errors
// for their corresponding form fields.
func (form *UpdateAccount) Validate() (apiErrs []*errors.APIError) {
	if form.FullName != nil && *form.FullName != "" {
		*form.FullName = NormalizeFullName(*form.FullName)
		if apiErr := ValidateFullName(*form.FullName, "fullName"); apiErr != nil {
			return []*errors.APIError{apiErr}
		}
	}

	return nil
}

// ChangePassword is a form submitted to change a user's password.
type ChangePassword struct {
	OldPassword string `json:"oldPassword" description:"User's current password. Required if the client is changing their own password. May be omitted if an admin is changing another user's password"`
	NewPassword string `json:"newPassword" description:"User's new password"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *ChangePassword) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this ChangePassword form. Returns any validation errors
// for their corresponding form fields.
func (form *ChangePassword) Validate() (apiErrs []*errors.APIError) {
	// NOTE: Don't bother checking the old password. It will be verified
	// by the API handler.
	if apiErr := ValidatePassword(form.NewPassword, "newPassword"); apiErr != nil {
		return []*errors.APIError{apiErr}
	}

	return nil
}

// ValidateTeamDescription validates that the given description meets our
// description requirements. If the description is not valid, returns a non-nil
// API error describing the validation error.
func ValidateTeamDescription(description, fieldName string) *errors.APIError {
	if len(description) > maxDescriptonLength {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("must be no longer than %d characters", maxDescriptonLength))
	}

	return nil
}

// CreateTeam is a form submitted to create a new team in some org.
type CreateTeam struct {
	Name        string `json:"name"                  description:"Name of the team"`
	Description string `json:"description,omitempty" description:"Description of the team"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *CreateTeam) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this CreateTeam form. Returns any validation errors
// for their corresponding form fields.
func (form *CreateTeam) Validate() (apiErrs []*errors.APIError) {
	form.Name = NormalizeFullName(form.Name)
	if apiErr := ValidateFullName(form.Name, "name"); apiErr != nil {
		return []*errors.APIError{apiErr}
	}

	if apiErr := ValidateTeamDescription(form.Description, "description"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	return apiErrs
}

// UpdateTeam is a form submitted to update a team.
type UpdateTeam struct {
	Name        *string `json:"name,omitempty"        description:"Name of the team, unchanged if nil or omitted"`
	Description *string `json:"description,omitempty" description:"Description of the team, unchanged if nil or omitted"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *UpdateTeam) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this UpdateTeam form. Returns any validation errors
// for their corresponding form fields.
func (form *UpdateTeam) Validate() (apiErrs []*errors.APIError) {
	if form.Name != nil {
		*form.Name = NormalizeFullName(*form.Name)
		if apiErr := ValidateFullName(*form.Name, "name"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if form.Description != nil {
		if apiErr := ValidateTeamDescription(*form.Description, "description"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	return apiErrs
}

// SetMembership is a form submitted to set or update a user's membership in an
// organizaiton or team.
type SetMembership struct {
	IsAdmin  *bool `json:"isAdmin,omitempty"  description:"Whether the member should be an admin of the organization or team (default false), unchanged if nil or omitted"`
	IsPublic *bool `json:"isPublic,omitempty" description:"Whether the membership is public (default false), unchanged if nil or omitted"`
}

// AuthConfig is a form submitted to set system auth configuration.
type AuthConfig struct {
	Backend string `json:"backend" description:"The name of the auth backend to use" enum:"managed|ldap"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *AuthConfig) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this AuthConfig form. Returns any validation errors
// for their corresponding form fields.
func (form *AuthConfig) Validate() (apiErrs []*errors.APIError) {
	form.Backend = strings.ToLower(strings.TrimSpace(form.Backend))

	if _, ok := config.SupportedAuthBackends[form.Backend]; !ok {
		apiErrs = append(apiErrs, errors.InvalidFormField("backend", "not a valid choice for auth backend"))
	}

	return apiErrs
}

// OpenIDConfig is a form submitted to set system OpenID configuration.
type OpenIDConfig struct {
	IssuerIdentifier string `json:"issuerIdentifier" description:"The Issuer Identifier for the system's OpenID Connect provider"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *OpenIDConfig) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this OpenIDConfig form. Returns any validation errors
// for their corresponding form fields.
func (form *OpenIDConfig) Validate() (apiErrs []*errors.APIError) {
	// An Issuer Identifier is a case sensitive URL using the https scheme
	// that contains scheme, host, and optionally, port number and path
	// components and no query or fragment components.
	if form.IssuerIdentifier == "" {
		return []*errors.APIError{errors.InvalidFormField("issuerIdentifier", "this field can not be left blank")}
	}

	parsedURL, err := url.Parse(form.IssuerIdentifier)
	if err != nil {
		return []*errors.APIError{errors.InvalidFormField("issuerIdentifier", fmt.Sprintf("unable to parse URI: %s", err))}
	}

	if strings.ToLower(parsedURL.Scheme) != "https" {
		apiErrs = append(apiErrs, errors.InvalidFormField("issuerIdentifier", "URL scheme must be https"))
	}

	if parsedURL.Host == "" {
		apiErrs = append(apiErrs, errors.InvalidFormField("issuerIdentifier", "URL must contain a host"))
	}

	if !(parsedURL.Fragment == "" && parsedURL.RawQuery == "") {
		apiErrs = append(apiErrs, errors.InvalidFormField("issuerIdentifier", "URL must not contain query or fragment components"))
	}

	return apiErrs
}

// LDAPSettings is a sub-form used in the AuthConfig form.
type LDAPSettings struct {
	RecoveryAdminUsername string           `json:"recoveryAdminUsername"           description:"The user with this name will be able to authenticate to the system even when the LDAP server is unavailable or if auth is misconfigured"`
	RecoveryAdminPassword *string          `json:"recoveryAdminPassword,omitempty" description:"A secure hash of this password is stored locally so that you can still login with the recovery admin username if the LDAP server is unavailable or if auth is misconfigured. May be nil or omitted to leave the current password unchanged"`
	ServerURL             string           `json:"serverURL"                       description:"The URL of the LDAP server"`
	NoSimplePagination    bool             `json:"noSimplePagination"              description:"The server does not support the Simple Paged Results control extension (RFC 2696)"`
	StartTLS              bool             `json:"startTLS"                        description:"Whether to use StartTLS to secure the connection to the server, ignored if server URL scheme is 'ldaps://'"`
	RootCerts             string           `json:"rootCerts"                       description:"A root certificate PEM bundle to use when establishing a TLS connection to the server"`
	TLSSkipVerify         bool             `json:"tlsSkipVerify"                   description:"Whether to skip verifying of the server's certificate when establishing a TLS connection, not recommended unless testing on a secure network"`
	ReaderDN              string           `json:"readerDN"                        description:"The distinguished name the system will use to bind to the LDAP server when performing searches"`
	ReaderPassword        string           `json:"readerPassword"                  description:"The password that the system will use to bind to the LDAP server when performing searches"`
	UserSearchConfigs     []UserSearchOpts `json:"userSearchConfigs"               description:"One or more settings for syncing users"`
	AdminSyncOpts         MemberSyncOpts   `json:"adminSyncOpts"                   description:"Settings for syncing system admin users"`
	SyncSchedule          string           `json:"syncSchedule"                    description:"The scheduled time for automatic LDAP sync jobs, in CRON format with seconds omitted, default is @hourly if empty or omitted"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *LDAPSettings) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this LDAPSettings form. Returns any validation errors
// for their corresponding form fields. However, this does not validate that
// the recovery admin user exists.
func (form *LDAPSettings) Validate() (apiErrs []*errors.APIError) {
	return form.validate("")
}

// ValidateLDAPServerURL validates that the given rawURL is a valid ldap or
// ldaps URI.
func ValidateLDAPServerURL(fieldName, rawURL string, startTLS bool) *errors.APIError {
	if rawURL == "" {
		return errors.InvalidFormField(fieldName, "this field can not be left blank")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("unable to parse URI: %s", err))
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	if !(scheme == "ldap" || scheme == "ldaps") {
		return errors.InvalidFormField(fieldName, "must begin with 'ldap://' or 'ldaps://'")
	}

	if startTLS && scheme == "ldaps" {
		return errors.InvalidFormField(fieldName, "cannot use StartTLS with 'ldaps://' scheme")
	}

	return nil
}

// validate validates this LDAPSettings form. Returns any validation errors
// for their corresponding form fields.
func (form *LDAPSettings) validate(baseFieldName string) (apiErrs []*errors.APIError) {
	if apiErr := ValidateAccountName(&form.RecoveryAdminUsername, fqFieldName(baseFieldName, "recoveryAdminUsername")); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}
	if form.RecoveryAdminPassword != nil {
		if apiErr := ValidatePassword(*form.RecoveryAdminPassword, fqFieldName(baseFieldName, "recoveryAdminPassword")); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	form.ServerURL = strings.TrimSpace(form.ServerURL)
	if apiErr := ValidateLDAPServerURL(fqFieldName(baseFieldName, "serverURL"), form.ServerURL, form.StartTLS); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}
	if form.RootCerts != "" {
		if !(x509.NewCertPool().AppendCertsFromPEM([]byte(form.RootCerts))) {
			apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "rootCerts"), "unable to successfully parse any PEM certificates"))
		}
	}

	// If either ReaderDN or ReaderPassword is specified, BOTH must be.
	if (form.ReaderDN == "") != (form.ReaderPassword == "") {
		apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "readerDN"), "this field must be specified with readerPassword"))
	}

	if len(form.UserSearchConfigs) == 0 {
		apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "userSearchConfigs"), "there must be at least one config object in this list"))
	}
	for _, subform := range form.UserSearchConfigs {
		if subFormErrs := subform.validate(fqFieldName(baseFieldName, "userSearchConfigs")); subFormErrs != nil {
			apiErrs = append(apiErrs, subFormErrs...)
		}
	}
	if subFormErrs := form.AdminSyncOpts.validate(fqFieldName(baseFieldName, "adminSyncOpts")); subFormErrs != nil {
		apiErrs = append(apiErrs, subFormErrs...)
	}
	if apiErr := ValidateCronSpec(&form.SyncSchedule, fqFieldName(baseFieldName, "syncSchedule")); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	return apiErrs
}

// UserSearchOpts is used in LDAPSettings.
type UserSearchOpts struct {
	BaseDN       string `json:"baseDN"              description:"The distinguished name of the element from which the LDAP server will search for users"`
	ScopeSubtree bool   `json:"scopeSubtree"        description:"Whether to search for users in the entire subtree of the base DN or to only search one level under the base DN (if false)"`
	UsernameAttr string `json:"usernameAttr"        description:"The name of the attribute of the LDAP user element which should be selected as the username"`
	FullNameAttr string `json:"fullNameAttr"        description:"The name of the attribute of the LDAP user element which should be selected as the full name of the user"`
	Filter       string `json:"filter"              description:"The LDAP search filter used to select user elements, may be left blank"`
}

// Validate validates this UserSearchOpts form. Returns any validation errors
// for their corresponding form fields.
func (form *UserSearchOpts) validate(baseFieldName string) (apiErrs []*errors.APIError) {
	if form.BaseDN == "" {
		apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "baseDN"), "this field can not be left blank"))
	}
	if form.UsernameAttr == "" {
		apiErrs = append(apiErrs, errors.InvalidFormField(fqFieldName(baseFieldName, "usernameAttr"), "this field can not be left blank"))
	}

	return apiErrs
}

// TryLdapLogin is a form submitted to try a login using the given LDAP
// settings.
type TryLdapLogin struct {
	Username     string       `json:"username" description:"Value of the specified username attribute for the user entry in the LDAP directory"`
	Password     string       `json:"password" description:"The password that the system will use to bind to the LDAP server when authenticating the user"`
	LDAPSettings LDAPSettings `json:"ldapSettings" description:"Various options for LDAP syncing"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *TryLdapLogin) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this TryLdapLogin form. Returns any validation errors
// for their corresponding form fields. However, this does not validate that
// the ldap recovery admin user exists.
func (form *TryLdapLogin) Validate() (apiErrs []*errors.APIError) {
	if apiErr := ValidateAccountName(&form.Username, "username"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}
	if form.Password == "" { // Don't enforce any rules about LDAP passwords.
		apiErrs = append(apiErrs, errors.InvalidFormField("password", "this field can not be left blank"))
	}
	if subFormErrs := form.LDAPSettings.validate("ldapSettings"); subFormErrs != nil {
		apiErrs = append(apiErrs, subFormErrs...)
	}

	return apiErrs
}

// fqFieldName returns the fully qualified field name for form fields with
// sub-forms.
func fqFieldName(base, name string) string {
	if base == "" {
		return name
	}

	return fmt.Sprintf("%s.%s", base, name)
}

// JobSubmission is a form for scheduling an on-demand job to run.
type JobSubmission struct {
	Action string `json:"action" description:"The action which the job will perform"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *JobSubmission) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this JobSubmission form.
func (form *JobSubmission) Validate() (apiErrs []*errors.APIError) {
	if _, ok := worker.RegisteredActions[form.Action]; !ok {
		return []*errors.APIError{errors.InvalidFormField("action", "not a registered action")}
	}

	return nil
}

// ValidateServiceDescription validates that the given description meets our
// description requirements. If the description is not valid, returns a non-nil
// API error describing the validation error.
func ValidateServiceDescription(description, fieldName string) *errors.APIError {
	if len(description) > maxDescriptonLength {
		return errors.InvalidFormField(fieldName, fmt.Sprintf("must be no longer than %d characters", maxDescriptonLength))
	}

	return nil
}

// ValidateServiceResourceURIs validates the given service resource URIs.
// Simply ensures that the given URLs are parseable and that no more than 5
// resource URIs are specified.
func ValidateServiceResourceURIs(resourceURIs []string, fieldName string) *errors.APIError {
	if len(resourceURIs) == 0 {
		return errors.InvalidFormField(fieldName, "at least one value must be specified")
	}

	if len(resourceURIs) > 7 {
		return errors.InvalidFormField(fieldName, "no more than 7 values may be specified")
	}

	for _, resourceURI := range resourceURIs {
		parsed, err := url.Parse(resourceURI)
		if err != nil {
			return errors.InvalidFormField(fieldName, fmt.Sprintf("unable to parse URL: %s", err))
		}

		if strings.ToLower(parsed.Scheme) != "https" {
			return errors.InvalidFormField(fieldName, "URL scheme must be 'https'")
		}

		if parsed.Host == "" {
			return errors.InvalidFormField(fieldName, "URL does not specify a host address")
		}
	}

	return nil
}

// ValidateProviderIdentities verifies that at least one and no more than 7
// provider Identities are specified and that each one is no longer than 128
// characters.
func ValidateProviderIdentities(providerIdentities []string, fieldName string) *errors.APIError {
	if len(providerIdentities) == 0 {
		return errors.InvalidFormField(fieldName, "at least one value must be specified")
	}

	if len(providerIdentities) > 7 {
		return errors.InvalidFormField(fieldName, "no more than 7 values may be specified")
	}

	for _, providerIdentity := range providerIdentities {
		if len(providerIdentity) > 128 {
			return errors.InvalidFormField(fieldName, "values may be no longer than 128 characters each")
		}
	}

	return nil
}

// CreateService is a form submitted to create a service.
type CreateService struct {
	Name               string   `json:"name"               description:"simple name of this service; unique to the owning account"`
	Description        string   `json:"description"        description:"short description of the service"`
	URL                string   `json:"url"                description:"web address for the service; must use 'https'"`
	Privileged         bool     `json:"privileged"         description:"whether this service has implicit authorization to all accounts; can only be set by a system admin"`
	RedirectURIs       []string `json:"redirectURIs"       description:"a list of URLs which serve as callbacks from Oauth authorization requests; all URLs must use 'https'"`
	JWKsURIs           []string `json:"jwksURIs"           description:"a list of URLs for the service's JSON Web Key set, the set of public keys it used to sign authentication tokens. At least 1 and no more than 7 must be specified. URLs must use 'https'. Provide a caBundle value if they cannot be verified by a public certificate authority"`
	ProviderIdentities []string `json:"providerIdentities" description:"A list of identifiers that this service will use as the 'aud' (Audience) for its authentication tokens (the tokens must have at least one 'aud' value that matches). Identity tokens issued by the provider will also use the first matching 'aud' (Audience) value in the 'iss' (Issuer) field. A single address (hostname[:port]) for the auth provider is recommended, but up to 7 values may be specified to allow services which use multiple addresses (IPs or DNS names). Each may be no longer than 128 characters. At least one must be specified"`
	CABundle           string   `json:"caBundle"           description:"PEM certificate bundle used to verify the JWKs URL endpoints"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *CreateService) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this CreateService form.
func (form *CreateService) Validate() (apiErrs []*errors.APIError) {
	form.Name = NormalizeFullName(form.Name)
	if apiErr := ValidateFullName(form.Name, "name"); apiErr != nil {
		return []*errors.APIError{apiErr}
	}

	if apiErr := ValidateServiceDescription(form.Description, "description"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	if apiErr := ValidateServiceResourceURIs([]string{form.URL}, "url"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	if apiErr := ValidateServiceResourceURIs(form.RedirectURIs, "redirectURIs"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	if apiErr := ValidateServiceResourceURIs(form.JWKsURIs, "jwksURIs"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	if apiErr := ValidateProviderIdentities(form.ProviderIdentities, "providerIdentities"); apiErr != nil {
		apiErrs = append(apiErrs, apiErr)
	}

	if form.CABundle != "" {
		if !(x509.NewCertPool().AppendCertsFromPEM([]byte(form.CABundle))) {
			apiErrs = append(apiErrs, errors.InvalidFormField("caBundle", "unable to successfully parse any PEM certificates"))
		}
	}

	return apiErrs
}

// UpdateService is a form submitted to update a service.
type UpdateService struct {
	Description        *string   `json:"description,omitempty"        description:"short description of the service"`
	URL                *string   `json:"url,omitempty"                description:"web address for the service; must use 'https'"`
	Privileged         *bool     `json:"privileged,omitempty"         description:"whether this service has implicit authorization to all accounts; can only be set by a system admin"`
	RedirectURIs       *[]string `json:"redirectURIs,omitempty"       description:"a list of URLs which serve as callbacks from Oauth authorization requests; all URLs must use 'https'"`
	JWKsURIs           *[]string `json:"jwksURIs,omitempty"           description:"a list of URLs for the service's JSON Web Key set, the set of public keys it used to sign authentication tokens; must use 'https'"`
	ProviderIdentities *[]string `json:"providerIdentities,omitempty" description:"A list of identifiers that this service will use as the 'aud' (Audience) for its authentication tokens (the tokens must have at least one 'aud' value that matches). Identity tokens issued by the provider will also use the first matching 'aud' (Audience) value in the 'iss' (Issuer) field. A single address (hostname[:port]) for the auth provider is recommended, but up to 7 values may be specified to allow services which use multiple addresses (IPs or DNS names). Each may be no longer than 128 characters. At least one must be specified"`
	CABundle           *string   `json:"caBundle,omitempty"           description:"PEM certificate bundle used to authenticate the JWKs URL endpoint"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *UpdateService) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this UpdateService form.
func (form *UpdateService) Validate() (apiErrs []*errors.APIError) {
	if form.Description != nil {
		if apiErr := ValidateServiceDescription(*form.Description, "description"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if form.URL != nil {
		if apiErr := ValidateServiceResourceURIs([]string{*form.URL}, "url"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if form.RedirectURIs != nil {
		if apiErr := ValidateServiceResourceURIs(*form.RedirectURIs, "redirectURIs"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if form.JWKsURIs != nil {
		if apiErr := ValidateServiceResourceURIs(*form.JWKsURIs, "jwksURIs"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if form.ProviderIdentities != nil {
		if apiErr := ValidateProviderIdentities(*form.ProviderIdentities, "providerIdentities"); apiErr != nil {
			apiErrs = append(apiErrs, apiErr)
		}
	}

	if form.CABundle != nil && *form.CABundle != "" {
		if !(x509.NewCertPool().AppendCertsFromPEM([]byte(*form.CABundle))) {
			apiErrs = append(apiErrs, errors.InvalidFormField("caBundle", "unable to successfully parse any PEM certificates"))
		}
	}

	return apiErrs
}

// Login is the form used to login.
type Login struct {
	Username string `json:"username" description:"the username of the account to login as"`
	Password string `json:"password" description:"the password for the user account"`
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *Login) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this Login form.
func (form *Login) Validate() (apiErrs []*errors.APIError) {
	form.Username = NormalizeAccountName(form.Username)
	return nil
}
