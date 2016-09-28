package ldap

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/utils"
)

func (a LDAPAuthenticator) SyncAllTeams() error {
	var retErr error

	a.updateSyncMessages("Team sync started")
	teams, err := a.ListTeams(nil)
	if err != nil {
		return err
	}
	for _, team := range teams {
		err := a.SyncTeam(team)
		if err != nil {
			msg := fmt.Sprintf("Failed to sync team %s: %s", team.Name, err)
			a.updateSyncMessages(msg)
			log.Error(msg)
		}
	}

	return retErr
}

func (a LDAPAuthenticator) SyncTeam(team *auth.Team) error {
	if team.LdapDN == "" {
		log.Debugf("Group doesn't require sync - non-ldap %s", team.Name)
		return nil
	}

	if team.LdapMemberAttr == "" {
		return fmt.Errorf("empty LDAP group member attribute for org %q team %q", team.OrgId, team.Name)
	}
	log.Infof("Synchronizing team %s with LDAP DN %s", team.Name, team.LdapDN)

	newMemberDNSet := GetGroupMembers(*a.settings, team.LdapDN, team.LdapMemberAttr)

	// TODO - potential optimization - when processing multiple teams we could cache this
	accounts, err := a.ListUsers(nil)
	if err != nil {
		return err
	}
	accountsByDN := map[string]*auth.Account{}
	for _, account := range accounts {
		if account.LdapDN != "" {
			accountsByDN[account.LdapDN] = account
		}
	}

	// Filter members down to only accounts we know about
	memberAccountsByName := map[string]*auth.Account{}
	for memberDN := range newMemberDNSet {
		if acct, ok := accountsByDN[memberDN]; ok {
			memberAccountsByName[acct.Username] = acct
		} // else skip this member since we don't know them
	}

	existingMemberNameList, _ := a.getTeamMembers(team.Id)
	if existingMemberNameList == nil {
		// Assume error case is empty list and add all we find
		existingMemberNameList = []string{}
	}

	// Now update membership

	// Deletion
	for _, existingMemberName := range existingMemberNameList {
		if _, ok := memberAccountsByName[existingMemberName]; !ok {
			log.Infof("Removing %s from team %s", existingMemberName, team.Name)
			k := path.Join(datastoreTeamMembers, team.Id, existingMemberName)
			if err := a.store.Delete(k); err != nil {
				// TODO - harden for partial failures
				return err
			}
			k = path.Join(datastoreAccountTeams, existingMemberName, team.Id)

			if err := a.store.Delete(k); err != nil {
				// TODO - harden for partial failures
				return err
			}
		}
	}

	// Insertion
	for newMemberName := range memberAccountsByName {
		// TODO Potential optimization - build a hash for existing member names
		found := false
		for _, existingMemberName := range existingMemberNameList {
			if existingMemberName == newMemberName {
				found = true
				break
			}
		}
		if !found {
			log.Infof("Adding %s to team %s", newMemberName, team.Name)
			k := path.Join(datastoreTeamMembers, team.Id, newMemberName)
			if err := a.store.Put(k, []byte(newMemberName), nil); err != nil {
				// TODO - harden for partial failures
				return utils.MaybeWrapEtcdClusterErr(err)
			}

			k = path.Join(datastoreAccountTeams, newMemberName, team.Id)
			if err := a.store.Put(k, []byte(team.Id), nil); err != nil {
				// TODO - harden for partial failures
				return utils.MaybeWrapEtcdClusterErr(err)
			}
		}
	}
	log.Infof("Synchronizing for team %s complete", team.Name)
	return nil
}
