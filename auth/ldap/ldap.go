package ldap

import (
	"crypto"
	"fmt"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/auth/builtin"
	"github.com/docker/orca/utils"
	"github.com/go-ldap/ldap"
)

var (
	// Membership for LDAP groups is managed through searches
	datastoreVersion          = "orca/v1"
	datastoreAccountTeams     = datastoreVersion + "/ldapaccountteams"
	datastoreTeamMembers      = datastoreVersion + "/ldapteammembers"
	datastoreLastSync         = datastoreVersion + "/ldaplastsynctime"
	datastoreLastSyncMessages = datastoreVersion + "/authlastsyncmessages"
	datastoreSyncLock         = datastoreVersion + "/ldapsynclock"
)

type (
	LDAPAuthenticator struct {
		builtin  *builtin.BuiltinAuthenticator
		settings *auth.LDAPSettings
		store    kvstore.Store
	}

	DtrAuthToken struct {
		SessionSecret string
		CsrfSecret    string
	}
)

var _ auth.Authenticator = (*LDAPAuthenticator)(nil)

func NewAuthenticator(store *kvstore.Store, orcaID string, certPEM, keyPEM []byte, settings *auth.LDAPSettings) *LDAPAuthenticator {
	builtinAuthenticator := builtin.NewAuthenticator(store, orcaID, certPEM, keyPEM)

	return &LDAPAuthenticator{
		builtin:  builtinAuthenticator,
		settings: settings,
		store:    *store,
	}
}

func (a *LDAPAuthenticator) CanChangePassword(ctx *auth.Context) bool {
	return false
}

func (a *LDAPAuthenticator) Name() string {
	return "ldap"
}

func (a *LDAPAuthenticator) Store() kvstore.Store {
	return a.store
}

func (a *LDAPAuthenticator) AuthenticateUsernamePassword(username, password string) (*auth.Context, error) {
	account, err := a.GetUser(nil, username)
	// If we wanted to support creation on first login, we could detect
	// err == auth.ErrAccountDoesNotExist and look it up in ldap on the fly...
	if err != nil {
		return nil, err
	}

	// Reject any attempt to authenticate a non-ldap account
	if account.LdapDN == "" {
		log.Debugf("Attempt to login with non-ldap account rejected: %s", account.Username)
		return nil, auth.ErrUnauthorized
	}

	// check for empty password
	// on AD, an empty password will be sent as NT AUTHORITY\ANONYMOUS LOGON
	// and cause a successful login with an empty password
	if strings.TrimSpace(password) == "" {
		return nil, auth.ErrUnauthorized
	}

	// TODO Might want to cache the connection for faster lookups?
	ldapConn, err := GetConn(a.settings.ServerURL, *a.settings)
	if err != nil {
		log.Warnf("Unable to connect to LDAP: %s", err)
		return nil, err
	}
	defer ldapConn.Close()

	log.Debugf("Attempting bind to %s with %s", account.LdapDN, a.settings.ServerURL)
	if err := ldapConn.Bind(account.LdapDN, password); err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok && ldapErr.ResultCode == ldap.LDAPResultInvalidCredentials {
			return nil, auth.ErrUnauthorized
		}
		// TODO - are there other classes of errors we should detect and return generic Unauthorized?
		return nil, err
	}

	// TODO - consider adding extra claims to include group membership for faster perm validation later...
	tokenStr, err := a.builtin.GenerateToken(username, nil)
	if err != nil {
		return nil, err
	}

	ctx := &auth.Context{
		User:         account,
		SessionToken: tokenStr,
	}

	return ctx, nil
}

func (a *LDAPAuthenticator) AuthenticateSessionToken(tokenStr string) (*auth.Context, error) {
	return a.builtin.AuthenticateSessionToken(tokenStr)
}

func (a *LDAPAuthenticator) AuthenticatePublicKey(pubKey crypto.PublicKey) (*auth.Context, error) {
	return a.builtin.AuthenticatePublicKey(pubKey)
}

func (a *LDAPAuthenticator) AddUserPublicKey(user *auth.Account, label string, publicKey crypto.PublicKey) error {
	return a.builtin.AddUserPublicKey(user, label, publicKey)
}

func (a *LDAPAuthenticator) Logout(ctx *auth.Context) error {
	return a.builtin.Logout(ctx)
}

func (a *LDAPAuthenticator) GetUser(ctx *auth.Context, username string) (*auth.Account, error) {
	account, err := a.builtin.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	teams, err := a.getAccountTeams(account.Username)
	if err != nil {
		return nil, err
	}

	account.DiscoveredTeams = teams

	return account, nil
}

func (a *LDAPAuthenticator) ListUsers(ctx *auth.Context) ([]*auth.Account, error) {
	// We don't add membership when listing all of them for performance reasons
	return a.builtin.ListUsers(ctx)
}

func (a *LDAPAuthenticator) ListTeamMembers(ctx *auth.Context, teamID string) ([]*auth.Account, error) {
	members, err := a.getTeamMembers(teamID) // Discovered Members
	if err != nil {
		return nil, err
	}

	accounts := []*auth.Account{}

	for _, m := range members {
		acct, err := a.GetUser(ctx, m)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acct)
	}
	// Retrive any managed members too
	managedMembers, err := a.builtin.ListTeamMembers(ctx, teamID)
	if err != nil {
		return nil, err
	}
	accounts = append(accounts, managedMembers...)

	return accounts, nil
}

func (a *LDAPAuthenticator) DeleteUser(ctx *auth.Context, account *auth.Account) error {
	// We *could* prevent deletion of LDAP based accounts, but if they've changed
	// the search spec for accounts, they might want to clean up old accounts, so
	// lets allow it.

	// Since membership is built up based on known accounts, we'll let the
	// background sync logic remove the stale membership information (they can't
	// login since the account will be gone. (which might make the account re-appear anyway)

	return a.builtin.DeleteUser(ctx, account)

}

func (a *LDAPAuthenticator) SaveUser(ctx *auth.Context, account *auth.Account) (string, error) {
	account.DiscoveredTeams = []string{}
	// TODO - should we prevent modification of the LDAP DN fields?
	return a.builtin.SaveUser(ctx, account)
}

// Used when we know the account we're adding is new to short-circuit validation during bulk adds
func (a *LDAPAuthenticator) UnvalidatedAddAccount(account *auth.Account) error {
	account.DiscoveredTeams = []string{}
	return a.builtin.UnvalidatedAddAccount(account)
}

func (a *LDAPAuthenticator) ChangePassword(ctx *auth.Context, username, oldPassword, newPassword string) error {
	return auth.ErrUnsupported
}

func (a *LDAPAuthenticator) SaveTeam(ctx *auth.Context, team *auth.Team) (string, error) {
	// Validate input first
	team.DiscoveredMembers = []string{}
	if team.LdapDN != "" {
		if team.LdapMemberAttr == "" {
			team.LdapMemberAttr = "member"
		}
	}
	ret, err := a.builtin.SaveTeam(ctx, team)
	if err != nil {
		return ret, err
	}

	// if it's an LDAP group, do an immediate update, and return an error
	// if it fails outright so the user knows they screwed something up
	if team.LdapDN != "" {
		err = a.SyncTeam(team)
		if err != nil {
			// TODO - consider saving a backup copy of the team in case
			//        this was an update, so we can revert it, else delete
			//        the team if it was a create so we don't have cruft
			//        lying around.  For now, we'll punt that to the user
			//        to clean up.
			return ret, err
		}
	}
	return ret, nil
}
func (a *LDAPAuthenticator) DeleteTeam(ctx *auth.Context, team *auth.Team) error {
	// Figure out which members we've got and update them so they no longer have this team listed
	members, _ := a.getTeamMembers(team.Id)
	for _, memberName := range members {
		k := path.Join(datastoreAccountTeams, memberName, team.Id)
		if err := a.store.Delete(k); err != nil && err != kvstore.ErrKeyNotFound {
			// Should we allow partial failures to fall through here so you can delete a busted team?
			// My concern is that may lead to more ominous corruption going undetected (stale memberships)
			// so we'll fail fast and force the user to take other corrective action to get the deletion done
			return fmt.Errorf("Failed to delelete membership for account %s from team %s: %s", memberName, team.Name, err)
		}
	}
	k := path.Join(datastoreTeamMembers, team.Id)
	if err := a.store.Delete(k); err != nil && err != kvstore.ErrKeyNotFound {
		return fmt.Errorf("Failed to delelete membership for team %s: %s", team.Name, err)
	}

	return a.builtin.DeleteTeam(ctx, team)
}
func (a *LDAPAuthenticator) ListTeams(ctx *auth.Context) ([]*auth.Team, error) {
	// We don't add membership when listing all of them for performance reasons
	return a.builtin.ListTeams(ctx)
}
func (a *LDAPAuthenticator) GetTeam(ctx *auth.Context, id string) (*auth.Team, error) {
	team, err := a.builtin.GetTeam(ctx, id)
	if err != nil {
		return nil, err
	}

	members, err := a.getTeamMembers(team.Id)
	if err != nil {
		return nil, err
	}

	team.DiscoveredMembers = members

	return team, nil

}
func (a *LDAPAuthenticator) AddTeamMember(ctx *auth.Context, teamID, username string) error {
	// These pass straight through since the request is user based (aka, managed membership)
	return a.builtin.AddTeamMember(ctx, teamID, username)
}
func (a *LDAPAuthenticator) DeleteTeamMember(ctx *auth.Context, teamID, username string) error {
	// These pass straight through since the request is user based (aka, managed membership)
	return a.builtin.DeleteTeamMember(ctx, teamID, username)
}

func (a *LDAPAuthenticator) getTeamMembers(teamID string) ([]string, error) {
	members := []string{}

	k := path.Join(datastoreTeamMembers, teamID)
	kvPairs, err := a.store.List(k)
	if err != nil && err == kvstore.ErrKeyNotFound {
		return members, nil
	} else if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		members = append(members, string(kvPair.Value))
	}

	return members, nil
}

func (a *LDAPAuthenticator) getAccountTeams(username string) ([]string, error) {
	teams := []string{}

	k := path.Join(datastoreAccountTeams, username)
	kvPairs, err := a.store.List(k)
	if err != nil && err == kvstore.ErrKeyNotFound {
		return teams, nil
	} else if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		teams = append(teams, string(kvPair.Value))
	}

	return teams, nil
}

func (a *LDAPAuthenticator) ListUserTeams(ctx *auth.Context, username string) ([]*auth.Team, error) {
	teams := []*auth.Team{}

	teamIds, err := a.getAccountTeams(username)
	if err != nil {
		return nil, err
	}

	for _, id := range teamIds {
		t, err := a.GetTeam(ctx, id)
		if err != nil {
			return nil, err
		}

		teams = append(teams, t)
	}

	return teams, nil
}

// Called from the auth subsystem when the periodic task timer fires.
func (a *LDAPAuthenticator) Sync(ctx *auth.Context, force, onlyAdmin bool) error {
	now := time.Now()
	var chDone = make(chan interface{})

	log.Debugf("LDAP sync fired - interval is %d minutes - forced: %t", a.settings.SyncInterval, force)
	lock, err := a.store.NewLock(datastoreSyncLock, &kvstore.LockOptions{TTL: time.Duration(a.settings.SyncInterval) * time.Minute})
	if err != nil {
		log.Errorf("Failed to get locker object from KV store for LDAP sync. %s", err)
		return err
	}

	chLost, err := lock.Lock(nil)
	if err != nil {
		log.Errorf("Failed to take lock from KV store for LDAP sync. %s", err)
		return err
	}

	// Unlock if we don't lose the lock
	go func() {
		select {
		case <-chDone:
			log.Debug("Unlocking sync lock")
			lock.Unlock()
		case <-chLost:
			log.Warn("We lost the sync lock (we probably took too long to sync), we'll let the current sync finish")
			// Note: if we lost the lock, attempting to unlock hangs, so just exit the cleanup routine
		}
	}()
	defer func() {
		// Fire the "we're done"
		chDone <- "done"
		close(chDone)
	}()

	kvPair, err := a.store.Get(datastoreLastSync)
	if err == nil {
		lastSync := time.Time{}
		err = lastSync.UnmarshalText(kvPair.Value)
		if err == nil {
			if lastSync.Add(time.Duration(a.settings.SyncInterval)*time.Minute).After(now) && !force {
				log.Debug("LDAP sync not due yet, not performing sync.")
				return nil
			}
		}
	} // else if for some reason it fails to unmarshal, sync anyway and update to correct

	log.Debug("Updating sync time to now")
	data, err := now.MarshalText()
	if err != nil {
		log.Errorf("Failed to update %s: %s", datastoreLastSync, err)
		return err
	}
	err = a.store.Put(datastoreLastSync, data, nil)
	if err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		log.Errorf("Failed to update %s: %s", datastoreLastSync, err)
		return err
	}

	// If we got this far, lets consider it a success, and log partial failures

	// Reset the sync messages
	_ = a.store.Put(datastoreLastSyncMessages, []byte(fmt.Sprintf("%s: Full Sync started", string(data))), nil)

	log.Infof("Performing LDAP sync to %s", a.settings.ServerURL)
	userErr := a.SyncAllUsers(onlyAdmin)
	if err != nil {
		log.Warnf("User sync failure: %s", err)
	}
	teamErr := a.SyncAllTeams()
	if err != nil {
		log.Warnf("Team sync failure: %s", err)
	}
	a.updateSyncMessages("Sync completed")
	log.Info("LDAP sync completed")
	if userErr != nil {
		return userErr
	}
	return teamErr
}

func (a *LDAPAuthenticator) updateSyncMessages(message string) {
	// ~Ignore failiures
	kvpair, err := a.store.Get(datastoreLastSyncMessages)
	if err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		log.Warnf("Failed to update sync message - %s - %s", err, message)
		return
	}
	timestamp, err := time.Now().MarshalText()
	if err != nil {
		log.Warnf("Failed to update sync message - %s - %s", err, message)
		return
	}
	now := string(timestamp)
	_ = a.store.Put(datastoreLastSyncMessages, []byte(string(kvpair.Value)+"\n"+now+": "+message), nil)
}

func (a *LDAPAuthenticator) LastSyncStatus(ctx *auth.Context) string {
	kvpair, err := a.store.Get(datastoreLastSyncMessages)
	if err != nil {
		return ""
	}
	return string(kvpair.Value)
}
