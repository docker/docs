package sync

import (
	"fmt"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/authn/ldap"
	"github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/schema"
)

// LdapSyncer is able to sync users and team membership with an LDAP directory
// service.
type LdapSyncer struct {
	ctx          context.Context
	schemaMgr    schema.Manager
	ldapSettings *config.Settings
	dnToUserMap  ldap.DNAccountMap
}

// NewLdapSyncer returns a new LdapSyncer using the given base context, schema
// manager, and LDAP settings.
func NewLdapSyncer(ctx context.Context, schemaMgr schema.Manager, ldapSettings *config.Settings) *LdapSyncer {
	return &LdapSyncer{
		ctx:          ctx,
		schemaMgr:    schemaMgr,
		ldapSettings: ldapSettings,
		dnToUserMap:  ldap.DNAccountMap{},
	}
}

// Run executes an LDAP sync of users, system admins, org admins, team members,
// and org members for the current LDAP configuration.
func (syncer *LdapSyncer) Run() error {
	if err := syncer.SyncAllUsers(); err != nil {
		return fmt.Errorf("unable to sync all users: %s", err)
	}

	if err := syncer.SyncAllMemberships(); err != nil {
		return fmt.Errorf("unable to sync all memberships: %s", err)
	}

	return nil
}
