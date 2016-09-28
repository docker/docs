package ldap

import (
	"container/list"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"

	goldap "github.com/go-ldap/ldap"
)

const userPageSize = 500

type userData struct {
	ldapDN   string
	disabled bool
}

// makeUserDataSet builds a mapping of usernames to ldapDN/disabled for every
// user currently in the database.
func (a LDAPAuthenticator) makeUserDataSet() (map[string]userData, error) {
	log.Info("Creating set of current users")

	allUsers, err := a.ListUsers(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to list users from database: %s", err)
	}

	currentUserSet := make(map[string]userData, len(allUsers))

	for _, user := range allUsers {
		log.Infof("Have current user: %s", user.Username)

		currentUserSet[user.Username] = userData{
			ldapDN:   user.LdapDN,
			disabled: user.Disabled,
		}
	}

	return currentUserSet, nil
}

// Common search optios (a subset of ldap.SearchRequest).
type searchOpts struct {
	scope        int
	deferAliases int
	sizeLimit    int
	timeLimit    int
	typesOnly    bool
	filter       string
	attributes   []string
}

func (a LDAPAuthenticator) searchForUsers(referral LDAPSearchURL, opts searchOpts) (*goldap.SearchResult, error) {
	log.Infof("Searching for users at %s", referral.HostURLString)

	conn, err := GetConn(referral.HostURLString, *a.settings)
	if err != nil {
		return nil, fmt.Errorf("unable to get LDAP server connection to %s (startTLS=%t): %s", referral.HostURLString, a.settings.StartTLS, err)
	}
	defer conn.Close()

	if err := conn.Bind(a.settings.ReaderDN, a.settings.ReaderPassword); err != nil {
		return nil, fmt.Errorf("unable to bind LDAP reader: (%s) %s", a.settings.ReaderDN, err)
	}

	searchRequest := &goldap.SearchRequest{
		BaseDN:       referral.BaseDN,
		Scope:        opts.scope,
		DerefAliases: opts.deferAliases,
		SizeLimit:    opts.sizeLimit,
		TimeLimit:    opts.timeLimit,
		TypesOnly:    opts.typesOnly,
		Filter:       opts.filter,
		Attributes:   opts.attributes,
	}

	searchResult, err := conn.SearchWithPaging(searchRequest, userPageSize)
	if err != nil {
		return nil, fmt.Errorf("unable to perform search with paging (%#v): %s", searchRequest, err)
	}

	return searchResult, nil
}

// findLDAPUsers builds a mapping of usernames to ldapDN/disabled for every
// object in the configured directory service which matches the configured
// search filter for users.
func (a LDAPAuthenticator) findLDAPUsers() map[string]userData {
	log.Info("Creating set of users found on LDAP server")

	// Prepare some common search arguments.
	// See https://tools.ietf.org/html/rfc4511#section-4.5.1 for full
	// nitty-gritty details.
	opts := searchOpts{
		scope:        goldap.ScopeWholeSubtree,
		deferAliases: goldap.DerefAlways,
		sizeLimit:    0,     // No client-requested size limit restrictions.
		timeLimit:    0,     // No client-requested time limit restrictions.
		typesOnly:    false, // We want attribute values, not just descriptions.
		attributes:   []string{"DN", a.settings.UserLoginAttrName},
	}

	opts.filter = fmt.Sprintf("(%s=*)", a.settings.UserLoginAttrName)
	if a.settings.UserSearchFilter != "" {
		// AND with the configured search filter.
		opts.filter = AndSearchFilter(a.settings.UserSearchFilter, opts.filter)
	}

	// Alocate a map to store LDAP usernames -> DN, initially large enough
	// to fit a full page of results.
	foundUserSet := make(map[string]userData, userPageSize)

	// Start with 1 referral. Gather paged search results from each
	// server and then move on to any referred servers.
	referrals := list.New()
	referrals.PushBack(LDAPSearchURL{
		HostURLString: a.settings.ServerURL,
		BaseDN:        a.settings.UserBaseDN,
	})

	for e := referrals.Front(); e != nil; e = e.Next() {
		referral := e.Value.(LDAPSearchURL)

		searchResult, err := a.searchForUsers(referral, opts)
		if err != nil {
			msg := fmt.Sprintf("unable to perform search on server %q: %s", referral.HostURLString, err)
			log.Warnf(msg)
			a.updateSyncMessages(msg)
			continue // Try the next server, if any.
		}

		// Append any referrals to the referalls list.
		for _, referralURL := range searchResult.Referrals {
			parsedReferral, err := ParseReferral(referralURL)
			if err != nil {
				msg := fmt.Sprintf("unable to parse referral URL %q: %s", referralURL, err)
				log.Warnf(msg)
				a.updateSyncMessages(msg)
				continue // Skip to the next referral URL.
			}

			referrals.PushBack(parsedReferral)
		}

		// Add result entries to the found user set.
		for _, entry := range searchResult.Entries {
			// Make all usernames lowercase. This should generally
			// be fine, as all LDAP attributes are case insensitive
			// by default, and it is supposedly rare for LDAP
			// attributes to be case sensitive.
			// http://stackoverflow.com/questions/29897684/is-ldap-dn-case-insensitive
			username := strings.ToLower(entry.GetAttributeValue(a.settings.UserLoginAttrName))
			foundUserSet[username] = userData{
				ldapDN:   entry.DN,
				disabled: false, // Active because they were found on the LDAP server.
			}
		}
	}
	// This set now contains all usernames and user DNs found from all
	// LDAP server referrals.
	return foundUserSet

}

// createLDAPUser attempts to create a new active user with the given
// username and LDAP DN. The username is validated to ensure that it matches
// our valid account name format. Any errors are logged and the user is not
// created.
func (a LDAPAuthenticator) createLDAPUser(username, ldapDN string) error {
	// TODO - we don't actually validate the account name format...
	acct := &auth.Account{
		Username: username,
		LdapDN:   ldapDN,
		Admin:    username == a.settings.AdminUsername,
		Role:     a.settings.UserDefaultRole, // No need to special case for admin, since that trumps
	}

	err := a.UnvalidatedAddAccount(acct)
	if err != nil {
		log.Errorf("error saving account: %s", err)
		return err
	}
	return nil
}

// ensureUserDataMatches attempts to update the given user in the database if
// their active status or ldap DN has changed in the new LDAP user data.
func (a LDAPAuthenticator) ensureUserDataMatches(username string, currentData, newData userData) error {
	activeMismatch := currentData.disabled != newData.disabled
	ldapDNMismatch := currentData.ldapDN != newData.ldapDN

	if !(activeMismatch || ldapDNMismatch) {
		return nil // Nothing to update.
	}

	userAcct, err := a.GetUser(nil, username)
	if err != nil {
		return fmt.Errorf("unable to get current user %q from database: %s", err)
	}

	// Never disable an admin account with the sync algorithm
	if !userAcct.Admin {
		userAcct.Disabled = newData.disabled
	}

	userAcct.LdapDN = newData.ldapDN
	userAcct.Password = "" // Don't change any stored builtin passwords
	if _, err := a.SaveUser(nil, userAcct); err != nil {
		return fmt.Errorf("unable to update current user %q in database: %s", err)
	}

	return nil
}

// SyncAllUsers syncs all users in the system using the configured LDAP search
// filter to find active users.
func (a LDAPAuthenticator) SyncAllUsers(onlyAdmin bool) error {
	if onlyAdmin {
		log.Info("Beginning sync of admin account: %s", a.settings.AdminUsername)
	} else {
		log.Info("Beginning sync of all user accounts")
	}
	a.updateSyncMessages("User sync started")

	currentUserSet, err := a.makeUserDataSet()
	if err != nil {
		return fmt.Errorf("unable to make set of current users: %s", err)
	}

	foundUserSet := a.findLDAPUsers()
	if len(foundUserSet) == 0 {
		msg := "unable to find LDAP users: search found no users"
		a.updateSyncMessages(msg)
		return fmt.Errorf(msg)
	}
	a.updateSyncMessages(fmt.Sprintf("Evaluating %d discovered users", len(foundUserSet)))

	/*
		Now we have 2 sets of users:

			A = {all existing database users}
			B = {all found LDAP user objects}

		The set of users which need to be added is denoted by the set:

			B - A		(set subtraction)

		The set of users which need to be marked as active and ensure
		that their DNs are up-to-date is denoted by the set:

			A ∩ B		(set intersection)

		The set of users which need to be marked as inactive (so that
		an admin may delete if desired) is denoted by the set:

			A - B		(set subtraction)

		One could imagine it like a venn diagram as well, but I'm not
		that good at ASCII art.
	*/

	// First, we'll handle the case of set B - A: adding new users. We will
	// iterate over the items in B and if they are *not* found in A, then
	// we will create that user in the database.
	addedUsernames := make([]string, 0, len(foundUserSet))
	for username, foundData := range foundUserSet {
		if _, isCurrent := currentUserSet[username]; isCurrent {
			continue // The user already exists.
		}

		// If we're only looking at the admin right now, skip all others
		if onlyAdmin && username != a.settings.AdminUsername {
			continue
		}

		// The user will either be created or skipped due to an error
		// (which will be logged).
		if err := a.createLDAPUser(username, foundData.ldapDN); err != nil {
			log.Errorf("unable to create LDAP user: %s", err)
		} else {
			addedUsernames = append(addedUsernames, username)
		}

		// If we have lots to add, give progress updates every 10%
		if len(addedUsernames) > 1000 && len(addedUsernames)%(len(foundUserSet)/10) == 0 {
			a.updateSyncMessages(fmt.Sprintf("Processed %d users", len(addedUsernames)))
		}
	}

	// We can handle the case of A ∩ B and A - B (updating existing users)
	// in one pass. We will iterate over the items in A and if they are not
	// found in B then they should be marked inactive. If they are found in
	// B but the DNs do not match or the user is not active, they will be
	// updated to match.
	for username, currentData := range currentUserSet {

		// If we're only looking at the admin right now, skip all others
		if onlyAdmin && username != a.settings.AdminUsername {
			continue
		}

		newData, found := foundUserSet[username]
		if !found {
			newData = userData{
				ldapDN:   currentData.ldapDN,
				disabled: true,
			}
		}

		if err := a.ensureUserDataMatches(username, currentData, newData); err != nil {
			log.Errorf("unable to update user data: %s", err)
		}
	}

	return nil
}
