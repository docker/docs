package config

import (
	"fmt"

	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
)

// UserSearchOpts specifies a search base DN, whether to search one level under
// the base DN or the whole subtree, the username attribute, and a search
// filter to use.
type UserSearchOpts struct {
	BaseDN       string `json:"baseDN"`
	ScopeSubtree bool   `json:"scopeSubtree"` // Default to one level.
	UsernameAttr string `json:"usernameAttr"`
	FullNameAttr string `json:"fullNameAttr"`
	Filter       string `json:"filter"`
}

// MemberSyncOpts specifies options for syncing a subset of local users.
type MemberSyncOpts struct {
	EnableSync bool `json:"enableSync"`
	// Whether to sync using groupDN+memberAttr selection or sync using a
	// user search filter.
	SelectGroupMembers bool `json:"selectGroupMembers"`

	// These fields are used to sync users using direct group membership.
	GroupDN         string `json:"groupDN"`
	GroupMemberAttr string `json:"groupMemberAttr"`

	// These fields are used to sync users using a search.
	SearchBaseDN       string `json:"searchBaseDN"`
	SearchScopeSubtree bool   `json:"searchScopeSubtree"`
	SearchFilter       string `json:"searchFilter"`
}

// Settings holds all of the basic required ldap settings.
type Settings struct {
	Debug bool `json:"-"`

	RecoveryAdminUsername string `json:"recoveryAdminUsername"`

	ServerURL          string `json:"serverURL"`
	NoSimplePagination bool   `json:"noSimplePagination"`
	StartTLS           bool   `json:"startTLS"`
	RootCerts          string `json:"rootCerts"`
	TLSSkipVerify      bool   `json:"tlsSkipVerify"`

	ReaderDN       string `json:"readerDN"`
	ReaderPassword string `json:"readerPassword"`

	UserSearchConfigs []UserSearchOpts `json:"userSearchConfigs"`
	AdminSyncOpts     MemberSyncOpts   `json:"adminSyncOpts"`

	SyncSchedule string `json:"syncSchedule"`
}

func (s Settings) String() string {
	return fmt.Sprintf("%#v", s)
}

// LDAPConfigPropertyKey is the property key which maps to the system's LDAP
// configuration.
const LDAPConfigPropertyKey = config.AuthConfigPropertyKey + ".ldap"

// GetLDAPConfig retrieves the current LDAP configuration using the given
// schema manager.
func GetLDAPConfig(mgr schema.Manager) (ldapConfig *Settings, err error) {
	ldapConfig = new(Settings)

	if err := mgr.GetProperty(LDAPConfigPropertyKey, ldapConfig); err != nil {
		if err == schema.ErrNoSuchProperty {
			// Use an empty Settings struct.
			return &Settings{}, nil
		}

		return nil, fmt.Errorf("unable to lookup LDAP config property: %s", err)
	}

	return ldapConfig, nil
}

// SetLDAPConfig sets the current LDAP configuration to the given LDAP config
// value. If it is nil, the current LDAP config will be deleted.
func SetLDAPConfig(mgr schema.Manager, ldapConfig *Settings) error {
	if ldapConfig == nil {
		return mgr.DeleteProperty(LDAPConfigPropertyKey)
	}

	if err := mgr.SetProperty(LDAPConfigPropertyKey, ldapConfig); err != nil {
		return fmt.Errorf("unable to set LDAP config property: %s", err)
	}

	return nil
}
