package sync

import (
	"fmt"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/authn/ldap"
	"github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/schema"

	goldap "github.com/go-ldap/ldap"
)

// SyncAllMemberships syncs all system admins, organization members, and team
// members in the system which have been configured for syncing with LDAP group
// member selection or user search filters.
func (syncer *LdapSyncer) SyncAllMemberships() error {
	context.GetLogger(syncer.ctx).Info("beginning sync of all memberships")

	if err := syncer.syncSystemAdmins(); err != nil {
		return fmt.Errorf("unable to sync system admins: %s", err)
	}

	// Use limit=0 to fetch all orgs in one page.
	orgs, _, err := syncer.schemaMgr.ListOrgs("", 0)
	if err != nil {
		return fmt.Errorf("unable to list orgs from database: %s", err)
	}

	for _, org := range orgs {
		ctx := context.WithValue(syncer.ctx, "org", org.Name)
		ctx = context.WithLogger(ctx, context.GetLogger(ctx, "org"))

		syncer.syncTeamMembershipsInOrg(ctx, org)

		if org.AdminSyncConfig.EnableSync {
			syncer.syncOrgMembership(ctx, org)
		}
	}

	return nil
}

func (syncer *LdapSyncer) syncSystemAdmins() error {
	logger := context.GetLogger(syncer.ctx)

	if !syncer.ldapSettings.AdminSyncOpts.EnableSync {
		logger.Info("not configured to sync system admin users")
		return nil
	}

	logger.Info("beginning sync of system admin users")

	// Use limit=0 to fetch all admins in one page.
	currentAdmins, _, err := syncer.schemaMgr.ListAdmins("", 0)
	if err != nil {
		return fmt.Errorf("unable to list current system admins: %s", err)
	}

	logger.Infof("number of current system admin users: %d", len(currentAdmins))

	foundAdminDNSet := syncer.searchForMembers(syncer.ldapSettings.AdminSyncOpts)

	logger.Infof("number of system admin users found via LDAP search: %d", len(foundAdminDNSet))

	if len(foundAdminDNSet) == 0 {
		return fmt.Errorf("unable to find system admin users: search found no users")
	}

	var nAdded, nRemoved, nUnchanged, nNotSynced int

	// Using our list of current admins and set of DNs of those who should
	// be admins, we can determine those who should no longer be admins and
	// those who need to be promoted admins.

	for _, currentAdmin := range currentAdmins {
		stillAdmin := foundAdminDNSet.Exists(currentAdmin.LdapDN)
		stillAdmin = stillAdmin && currentAdmin.IsActive
		stillAdmin = stillAdmin || currentAdmin.Name == syncer.ldapSettings.RecoveryAdminUsername

		if stillAdmin {
			nUnchanged++
		} else if err := syncer.schemaMgr.UpdateAccount(
			currentAdmin.ID,
			schema.AccountUpdateFields{
				// Note: stillAdmin is false here.
				IsAdmin: &stillAdmin,
			},
		); err != nil {
			logger.Warnf("unable to demote user %s from system admin status: %s", currentAdmin.Name, err)
			nUnchanged++
		} else {
			nRemoved++
		}

		// Remove from the set so that our second pass includes
		// only those who should be promoted.
		foundAdminDNSet.Remove(currentAdmin.LdapDN)
	}

	for adminDN := range foundAdminDNSet {
		user, exists := syncer.dnToUserMap.Get(adminDN)
		if !(exists && user.IsActive) {
			nNotSynced++
			continue
		}

		if err := syncer.schemaMgr.UpdateAccount(
			user.ID,
			schema.AccountUpdateFields{
				// Note: exists is true here.
				IsAdmin: &exists,
			},
		); err != nil {
			logger.Warnf("unable to promote user %s to system admin status: %s", user.Name, err)
			nUnchanged++
		} else {
			nAdded++
		}
	}

	// Ensure that the recovery admin user is marked as a system admin.

	recoveryAdminUsername := syncer.ldapSettings.RecoveryAdminUsername
	recoveryAdmin, err := syncer.schemaMgr.GetUserByName(recoveryAdminUsername)
	if err != nil {
		return fmt.Errorf("unable to get recovery admin user %s from database: %s", recoveryAdminUsername, err)
	}

	if !recoveryAdmin.IsAdmin {
		if err := syncer.schemaMgr.UpdateAccount(
			recoveryAdmin.ID,
			schema.AccountUpdateFields{IsAdmin: &[]bool{true}[0]},
		); err != nil {
			return fmt.Errorf("unable to set recovery admin user %s as system admin: %s", recoveryAdminUsername, err)
		}
	}

	logger.Infof("system admin users added: %d", nAdded)
	logger.Infof("system admin users removed: %d", nRemoved)
	logger.Infof("system admin users unchanged: %d", nUnchanged)
	logger.Infof("system admin users not synced: %d", nNotSynced)

	return nil
}

func (syncer *LdapSyncer) searchForMembers(config config.MemberSyncOpts) (memberDNSet ldap.DNSet) {
	if config.SelectGroupMembers {
		// Perform a simple group member selection.
		return ldap.GetGroupMembers(syncer.ctx, syncer.ldapSettings, config.GroupDN, config.GroupMemberAttr)
	}

	// Otherwise, we will search for members with a user search filter.
	// Prepare some common search options. See details at
	// https://tools.ietf.org/html/rfc4511#section-4.5.1
	opts := searchOpts{
		deferAliases: goldap.DerefAlways,
		sizeLimit:    0,     // No client-requested size limit restrictions.
		timeLimit:    0,     // No client-requested time limit restrictions.
		typesOnly:    false, // We want attribute values, not just descriptions.
	}

	// By default, search only a single level under the base DN.
	opts.scope = goldap.ScopeSingleLevel
	if config.SearchScopeSubtree {
		// Search all entries in the subtree of the base DN.
		opts.scope = goldap.ScopeWholeSubtree
	}

	if config.SearchFilter == "" {
		opts.filter = "(objectClass=*)" // Use the default filter.
	} else {
		opts.filter = ldap.AndSearchFilter(config.SearchFilter)
	}

	// We only need the DN of the user.
	opts.attributes = []string{"DN"}

	rootURL := ldap.LDAPSearchURL{
		HostURLString: syncer.ldapSettings.ServerURL,
		BaseDN:        config.SearchBaseDN,
	}

	memberDNSet = make(ldap.DNSet, userPageSize)

	syncer.searchForUsers(rootURL, opts, func(ctx context.Context, entry *goldap.Entry) {
		memberDNSet.Add(entry.DN)
	})

	return memberDNSet
}

func (syncer *LdapSyncer) syncTeamMembershipsInOrg(ctx context.Context, org schema.Account) {
	logger := context.GetLogger(ctx)
	logger.Info("beginning sync of team members")

	// Use limit=0 to fetch all teams in one page.
	teams, _, err := syncer.schemaMgr.ListLDAPSyncTeamsInOrg(org.ID, "", 0)
	if err != nil {
		logger.Warnf("unable to list LDAP-synced teams: %s", err)
		return
	}

	logger.Infof("found %d teams to sync", len(teams))

	for _, team := range teams {
		teamCtx := context.WithValue(ctx, "team", team.Name)
		teamCtx = context.WithLogger(teamCtx, context.GetLogger(teamCtx, "team"))

		syncer.syncTeam(teamCtx, team)
	}
}

func convertMemberSyncOpts(opts schema.MemberSyncOpts) config.MemberSyncOpts {
	return config.MemberSyncOpts{
		SelectGroupMembers: opts.SelectGroupMembers,
		GroupDN:            opts.GroupDN,
		GroupMemberAttr:    opts.GroupMemberAttr,
		SearchBaseDN:       opts.SearchBaseDN,
		SearchScopeSubtree: opts.SearchScopeSubtree,
		SearchFilter:       opts.SearchFilter,
	}
}

func (syncer *LdapSyncer) syncTeam(ctx context.Context, team schema.Team) {
	logger := context.GetLogger(ctx)
	logger.Info("beginning sync of team members")

	// Use limit=0 to fetch all members in one page.
	currentMembers, _, err := syncer.schemaMgr.ListTeamMembers(team.ID, "", 0)
	if err != nil {
		logger.Warnf("unable to list current team members: %s", err)
		return
	}

	logger.Infof("number of current team members: %d", len(currentMembers))

	foundMemberDNSet := syncer.searchForMembers(convertMemberSyncOpts(team.MemberSyncConfig))

	logger.Infof("number of team members found via LDAP search: %d", len(foundMemberDNSet))

	if len(foundMemberDNSet) == 0 {
		logger.Warnf("unable to find team members: search found no users")
		return
	}

	var nAdded, nRemoved, nUnchanged, nNotSynced int

	// Using our list of current members and set of DNs of those who should
	// be members, we can determine those who should no longer be members
	// and those who need to be added as members.

	for _, currentMember := range currentMembers {
		member := currentMember.Member

		stillMember := foundMemberDNSet.Exists(member.LdapDN)
		if stillMember && member.IsActive {
			nUnchanged++
		} else if err := syncer.schemaMgr.DeleteTeamMembership(team.ID, member.ID); err != nil {
			logger.Warnf("unable to remove team member %s: %s", member.Name, err)
			nUnchanged++
		} else {
			nRemoved++
		}

		// Remove from the set so that our second pass includes
		// only those who should be added as members.
		foundMemberDNSet.Remove(member.LdapDN)
	}

	for memberDN := range foundMemberDNSet {
		user, exists := syncer.dnToUserMap.Get(memberDN)
		if !(exists && user.IsActive) {
			nNotSynced++
			continue
		}

		// Add team membership (adds org membership if not already an
		// org member).
		if err := syncer.schemaMgr.AddTeamMembership(team.OrgID, team.ID, user.ID, nil, nil); err != nil {
			logger.Warnf("unable to add team member %s: %s", user.Name, err)
			nUnchanged++
		} else {
			nAdded++
		}
	}

	logger.Infof("team members added: %d", nAdded)
	logger.Infof("team members removed: %d", nRemoved)
	logger.Infof("team members unchanged: %d", nUnchanged)
	logger.Infof("team members not synced: %d", nNotSynced)
}

func (syncer *LdapSyncer) syncOrgMembership(ctx context.Context, org schema.Account) {
	logger := context.GetLogger(ctx)
	logger.Info("beginning sync of org members")

	// Use limit=0 to fetch all members in one page.
	currentMembers, _, err := syncer.schemaMgr.ListOrgMembers(org.ID, "", 0)
	if err != nil {
		logger.Warnf("unable to list current org members: %s", err)
		return
	}

	logger.Infof("number of current org members: %d", len(currentMembers))

	membersWithTeamsIDs, err := syncer.schemaMgr.ListTeamMembersInOrg(org.ID)
	if err != nil {
		logger.Warnf("unable to list distinct team members in org: %s", err)
		return
	}

	logger.Infof("number of distinct team members in org: %d", len(membersWithTeamsIDs))

	// Convert to a set for quick lookups by user ID.
	membersWithTeamsIDSet := make(map[string]struct{}, len(membersWithTeamsIDs))
	for _, memberID := range membersWithTeamsIDs {
		membersWithTeamsIDSet[memberID] = struct{}{}
	}

	foundAdminDNSet := syncer.searchForMembers(convertMemberSyncOpts(org.AdminSyncConfig))

	logger.Infof("number of org admin members found via LDAP search: %d", len(foundAdminDNSet))

	if len(foundAdminDNSet) == 0 {
		logger.Warn("unable to find org admin members: search found no users")
		return
	}

	var nAdded, nRemoved, nUnchanged, nNotSynced int

	// Using our list of current members, set of IDs of those who are a
	// member of a team, and set of LDAP DNs for those who should be admins
	// of the org, we can determine those who should no longer be members,
	// those who need to be added as members, and those who should be set
	// as org admins.

	// Do a pass on current members to remove those who are neither in a
	// team nor match the admin search.
	for _, membership := range currentMembers {
		member := membership.Member

		_, hasTeam := membersWithTeamsIDSet[member.ID]
		_, isAdmin := foundAdminDNSet[member.LdapDN]

		// The member is no longer a member of the org if they are no
		// longer a member of any team *and* did not match the admin
		// search.
		noLongerMember := !(member.IsActive && (hasTeam || isAdmin))
		// The member is no longer an admin if their current membership
		// info says they were previously but they are not found by
		// this search.
		noLongerAdmin := membership.IsAdmin && !isAdmin

		if noLongerMember {
			logger.Debugf("removing user %s from org because: isActive=%t hasTeam=%t isAdmin=%t", member.Name, member.IsActive, hasTeam, isAdmin)

			// The user should no longer be a member of the org.
			if err := syncer.schemaMgr.DeleteOrgMembership(org.ID, member.ID); err != nil {
				logger.Warnf("unable to remove org member %s: %s", member.Name, err)
				nUnchanged++
			} else {
				nRemoved++
			}
		} else if noLongerAdmin {
			// The member should still be an org member because
			// they are a member of a team, but we need to remove
			// their admin status.
			// Note: isAdmin is false.
			if err := syncer.schemaMgr.AddOrgMembership(org.ID, member.ID, &isAdmin, nil); err != nil {
				logger.Warnf("unable to remove org member %s admin status: %s", member.Name, err)
			}
		} else if !isAdmin {
			// New org admins are synced in the next loop, but this
			// member's status is unchanged.
			nUnchanged++
		}
	}

	// Create a set of current admin user IDs for quick lookups and so that
	// we don't have to make extraneous writes to the DB.
	currentAdminMemberIDSet := make(map[string]struct{}, len(currentMembers))
	for _, member := range currentMembers {
		if member.IsAdmin {
			currentAdminMemberIDSet[member.Member.ID] = struct{}{}
		}
	}

	// Do a pass on found admins to set their org membership admin status.
	for adminDN := range foundAdminDNSet {
		user, exists := syncer.dnToUserMap.Get(adminDN)
		if !(exists && user.IsActive) {
			logger.Debugf("not adding found org admin %s because: exists=%t isActive=%t", exists, exists && user.IsActive)
			nNotSynced++
			continue
		}

		if _, alreadyAdmin := currentAdminMemberIDSet[user.ID]; alreadyAdmin {
			logger.Debugf("user %s is already an org admin", user.Name, org.Name)
			nUnchanged++
			continue // Don't bother with a write to the DB.
		}

		logger.Debugf("adding user %s as org admin", user.Name)

		// Note: exists is true here.
		if err := syncer.schemaMgr.AddOrgMembership(org.ID, user.ID, &exists, nil); err != nil {
			logger.Warnf("unable to add member %s to org: %s", user.Name, err)
			nUnchanged++
		} else {
			nAdded++
		}
	}

	logger.Infof("org admin members added: %d", nAdded)
	logger.Infof("org members removed: %d", nRemoved)
	logger.Infof("org members unchanged: %d", nUnchanged)
	logger.Infof("org admin members not synced: %d", nNotSynced)
}
