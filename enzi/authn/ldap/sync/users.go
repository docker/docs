package sync

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/authn/ldap"
	"github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/schema"

	goldap "github.com/go-ldap/ldap"
)

// SyncAllUsers syncs all users in the system using the configured LDAP search
// filters to find active users.
func (syncer *LdapSyncer) SyncAllUsers() error {
	logger := context.GetLogger(syncer.ctx)
	logger.Info("beginning sync of all user accounts")

	currentUserSet, err := syncer.makeCurrentUserSet()
	if err != nil {
		return fmt.Errorf("unable to make set of current users: %s", err)
	}

	logger.Infof("number of current users: %d", len(currentUserSet))

	foundUserSet := syncer.findLDAPUsers()
	if len(foundUserSet) == 0 {
		return fmt.Errorf("unable to find LDAP users: search found no users")
	}

	logger.Infof("number of users found via LDAP search: %d", len(foundUserSet))

	/*
		Now we have 2 sets of users:

			A = {all existing database users}
			B = {all found LDAP user objects}

		The set of users which need to be marked as active and ensure
		that their DNs are up-to-date is denoted by the set:

			A ∩ B		(set intersection)

		The set of users which need to be marked as inactive (so that
		an admin may delete if desired) is denoted by the set:

			A - B		(set subtraction)

		The set of users which need to be added is denoted by the set:

			B - A		(set subtraction)

		One could imagine it like a venn diagram as well, but I'm not
		that good at ASCII art.
	*/

	var nCreated, nUpdated, nUnchanged, nDeactivated int

	// We can handle the case of A ∩ B and A - B (updating existing users)
	// in one pass. We will iterate over the items in A and if they are not
	// found in B then they should be marked inactive. If they are found in
	// B but the DNs do not match or the user is not active, they will be
	// updated to match.
	for username, currentData := range currentUserSet {
		newData, found := foundUserSet[username]
		if !found {
			newData = userData{
				ldapDN:   currentData.ldapDN,
				fullName: currentData.fullName,
				isActive: false,
			}
		}

		updated := syncer.ensureUserDataMatches(username, currentData, newData)
		switch {
		case updated && found:
			nUpdated++
		case updated:
			nDeactivated++
		default:
			nUnchanged++
		}

		// Remove from the found user set so that at the end of this
		// loop, it will only contain new users. This reduces the set
		// to B - A.
		delete(foundUserSet, username)
	}

	// Next, we'll handle the case of set B - A: adding new users. We will
	// iterate over the remaining items in the found user set and create
	// the users in the database.
	for username, foundData := range foundUserSet {
		if syncer.createLDAPUser(username, foundData.ldapDN, foundData.fullName) {
			nCreated++
		}
	}

	logger.Infof("users created: %d", nCreated)
	logger.Infof("users updated: %d", nUpdated)
	logger.Infof("users deactivated: %d", nDeactivated)
	logger.Infof("users unchanged: %d", nUnchanged)

	return nil
}

type userData struct {
	ldapDN   string
	fullName string
	isActive bool
}

// makeCurrentUserSet builds a mapping of usernames to userData for every user
// currently in the database.
func (syncer *LdapSyncer) makeCurrentUserSet() (map[string]userData, error) {
	logger := context.GetLogger(syncer.ctx)
	logger.Info("creating set of current users")

	// Use limit=0 to fetch all users in one page.
	users, _, err := syncer.schemaMgr.ListUsers("", 0)
	if err != nil {
		return nil, fmt.Errorf("unable to list users from database: %s", err)
	}

	currentUserSet := make(map[string]userData, len(users))
	for i, user := range users {
		logger.Debugf("current user: %s", user.Name)

		currentUserSet[user.Name] = userData{
			ldapDN:   user.LdapDN,
			fullName: user.FullName,
			isActive: user.IsActive,
		}

		syncer.dnToUserMap.Put(user.LdapDN, &users[i])
	}

	return currentUserSet, nil
}

const userPageSize = 500

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

// findLDAPUsers builds a mapping of usernames to ldapDN/isActive for every
// object in the configured directory service which matches the configured
// user search filters.
func (syncer *LdapSyncer) findLDAPUsers() (foundUserSet map[string]userData) {
	context.GetLogger(syncer.ctx).Info("creating set of users found via LDAP search")

	// Allocate a map to store LDAP DN -> username, initially large enough
	// to fit a full page of results.
	foundUserSet = make(map[string]userData, userPageSize)

	for _, config := range syncer.ldapSettings.UserSearchConfigs {
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
		if config.ScopeSubtree {
			// Search all entries in the subtree of the base DN.
			opts.scope = goldap.ScopeWholeSubtree
		}

		opts.filter = fmt.Sprintf("(%s=*)", config.UsernameAttr)
		if config.Filter != "" {
			// Logical AND with the configured search filter.
			opts.filter = ldap.AndSearchFilter(opts.filter, config.Filter)
		}

		// We only need the username attribute. In the future we might
		// get other fields like full name, email address, etc.
		opts.attributes = []string{config.UsernameAttr}

		if config.FullNameAttr != "" {
			opts.attributes = append(opts.attributes, config.FullNameAttr)
		}

		rootURL := ldap.LDAPSearchURL{
			HostURLString: syncer.ldapSettings.ServerURL,
			BaseDN:        config.BaseDN,
		}

		syncer.searchForUsers(rootURL, opts, func(ctx context.Context, entry *goldap.Entry) {
			addUserSearchResult(ctx, entry, config, foundUserSet)
		})
	}

	// This set now contains all usernames and user DNs found from all
	// LDAP server referrals.
	return foundUserSet
}

// searchForUsers searches for users using the given root LDAP URL and search
// options. The given processEntry function is called on every search result.
func (syncer *LdapSyncer) searchForUsers(rootURL ldap.LDAPSearchURL, opts searchOpts, processEntry func(context.Context, *goldap.Entry)) {
	logger := context.GetLogger(syncer.ctx)

	// Start with 1 referral. Gather paged search results from each server
	// and then move on to any referred servers.
	referrals := list.New()
	referrals.PushBack(rootURL)

	for e := referrals.Front(); e != nil; e = e.Next() {
		referral := e.Value.(ldap.LDAPSearchURL)

		searchResult, err := syncer.searchForUsersWithReferral(referral, opts)
		if err != nil {
			logger.Warnf("unable to perform search on server %q: %s", referral.HostURLString, err)
			continue // Try the next server, if any.
		}

		// Append any referrals to the referalls list.
		for _, referralURL := range searchResult.Referrals {
			logger.Debugf("found referral: %s", referralURL)

			parsedReferral, err := ldap.ParseReferral(referralURL)
			if err != nil {
				logger.Warnf("unable to parse referral URL %q: %s", referralURL, err)
				continue // Skip to the next referral URL.
			}

			referrals.PushBack(parsedReferral)
		}

		// Process entries with the given function.
		for _, entry := range searchResult.Entries {
			logger.Debugf("found user entry: %s", entry.DN)
			processEntry(syncer.ctx, entry)
		}
	}
}

// searchForUsersWithReferral searches for users at the given LDAP referral URL
// using the given options. Returns result with entries and additional
// referrals or an error.
func (syncer *LdapSyncer) searchForUsersWithReferral(referral ldap.LDAPSearchURL, opts searchOpts) (*goldap.SearchResult, error) {
	context.GetLogger(syncer.ctx).Infof("searching for users at %s", referral.HostURLString)

	conn, err := ldap.GetConn(referral.HostURLString, syncer.ldapSettings)
	if err != nil {
		return nil, fmt.Errorf("unable to get LDAP server connection to %s (startTLS=%t): %s", referral.HostURLString, syncer.ldapSettings.StartTLS, err)
	}
	defer conn.Close()

	// Bind if we are given a reader DN and password.
	readerDN := syncer.ldapSettings.ReaderDN
	readerPassword := syncer.ldapSettings.ReaderPassword
	if readerDN != "" && readerPassword != "" {
		if err := conn.Bind(readerDN, readerPassword); err != nil {
			return nil, fmt.Errorf("unable to bind LDAP reader: %s", err)
		}
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

	var searchResult *goldap.SearchResult
	if syncer.ldapSettings.NoSimplePagination {
		// As of 2016-03-28, we've only found one LDAP server which
		// does not support the simple paged results control which has
		// been part of the LDAPv3 spec since 1999.
		// *COUGH*COUGH*IBM BLUEPAGES*COUGH*
		searchResult, err = conn.Search(searchRequest)
	} else {
		searchResult, err = conn.SearchWithPaging(searchRequest, userPageSize)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to perform search with paging (%#v): %s", searchRequest, err)
	}

	return searchResult, nil
}

// addUserSearchResult validates the username from the given user search result
// entry according to the given config and the systems valid username rules. If
// the username is valid it is added to the given foundUserSet set keyed by
// the username.
func addUserSearchResult(ctx context.Context, entry *goldap.Entry, opts config.UserSearchOpts, foundUserSet map[string]userData) {
	logger := context.GetLogger(ctx)

	username := entry.GetAttributeValue(opts.UsernameAttr)

	logger.Debugf("un-normalized username: %q", username)

	// Skip if the account name is not valid. This also normalizes the
	// account name, which includes using the "Composed Normal Form" for
	// Unicode strings as well as make it all lowercase. This should
	// generally be fine, as all LDAP attributes are case insensitive by
	// default, and it is supposedly rare for LDAP attributes to be case
	// sensitive.
	// http://stackoverflow.com/questions/29897684/is-ldap-dn-case-insensitive
	if validationErr := forms.ValidateAccountName(&username, opts.UsernameAttr); validationErr != nil {
		logger.Warnf("invalid username %q: %s", username, validationErr.Detail)
		return
	}

	logger.Debugf("normalized username: %q", username)

	var fullName string
	if opts.FullNameAttr != "" {
		fullName = entry.GetAttributeValue(opts.FullNameAttr)
		if validationErr := forms.ValidateFullName(fullName, opts.FullNameAttr); validationErr != nil {
			logger.Warnf("invalid full name %q: %s", fullName, validationErr.Detail)
			fullName = "" // Ignore it since it's invalid.
		}
	}

	foundUserSet[username] = userData{
		// Normalize the DN to be all lowercase. See comment above for
		// why the DN should be treated as case-insensitive.
		ldapDN:   strings.ToLower(entry.DN),
		fullName: fullName,
		isActive: true,
	}
}

// createLDAPUser attempts to create a new active user with the given username
// and LDAP DN. The username must have already been validated to ensure that it
// matches our valid account name format. If there is an error, it will be
// logged. Returns whether or not the creation was successful.
func (syncer *LdapSyncer) createLDAPUser(username, ldapDN, fullName string) bool {
	user := &schema.Account{
		Name:     username,
		FullName: fullName,
		LdapDN:   ldapDN,
		IsActive: true,
	}

	if err := syncer.schemaMgr.CreateAccount(user); err != nil {
		context.GetLogger(syncer.ctx).Warnf("unable to add new user account: %s", err)
		return false
	}

	syncer.dnToUserMap.Put(ldapDN, user)

	return true
}

// ensureUserDataMatches attempts to update the given user in the database if
// their active status or ldap DN has changed in the new LDAP user data.
// Returns whether or not the user was updated.
func (syncer *LdapSyncer) ensureUserDataMatches(username string, currentData, newData userData) bool {
	logger := context.GetLogger(syncer.ctx)

	// Ensure that the recovery admin user is always marked active.
	if username == syncer.ldapSettings.RecoveryAdminUsername {
		newData.isActive = true
	}

	isActiveUnchanged := currentData.isActive == newData.isActive
	ldapDNUnchanged := currentData.ldapDN == newData.ldapDN
	// Do not consider the new full name if it's empty.
	fullNameUnchanged := newData.fullName == "" || currentData.fullName == newData.fullName

	if isActiveUnchanged && ldapDNUnchanged && fullNameUnchanged {
		return false // Nothing to update.
	}

	user, err := syncer.schemaMgr.GetUserByName(username)
	if err != nil {
		logger.Warnf("unable to get current user %s from database: %s", username, err)
		return false
	}

	user.IsActive = newData.isActive
	user.LdapDN = newData.ldapDN

	updateFields := schema.AccountUpdateFields{
		IsActive: &user.IsActive,
		LdapDN:   &user.LdapDN,
	}

	// Don't update the full name if it's empty.
	if newData.fullName != "" {
		updateFields.FullName = &newData.fullName
	}

	if err := syncer.schemaMgr.UpdateAccount(user.ID, updateFields); err != nil {
		logger.Warnf("unable to update current user %s in database: %s", username, err)
		return false
	}

	// Delete the old ID/DN mappings and add the new ones.
	syncer.dnToUserMap.Delete(currentData.ldapDN)
	syncer.dnToUserMap.Put(user.LdapDN, user)

	return true
}
