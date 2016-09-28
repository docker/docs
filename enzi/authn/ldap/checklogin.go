package ldap

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/authn"
	"github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/schema"
	goldap "github.com/go-ldap/ldap"
)

// CheckLoginWithSettings attemps to login using the given username, password,
// and settings. It returns a user (on success), along with a buffer of logs
// and an error (nil on success).
func CheckLoginWithSettings(username, password string, settings *config.Settings) (user *authn.Account, logBuf *bytes.Buffer, err error) {
	// Prepare context and logging.
	logBuf = new(bytes.Buffer)
	tryLoginLogger := log.NewEntry(&log.Logger{
		Out:       logBuf,
		Formatter: new(log.TextFormatter),
		Level:     log.DebugLevel,
	})
	ctx := context.WithLogger(
		context.Background(),
		tryLoginLogger,
	)

	var userDN string
	for _, userSearchOpts := range settings.UserSearchConfigs {
		rootURL := LDAPSearchURL{
			HostURLString: settings.ServerURL,
			BaseDN:        userSearchOpts.BaseDN,
		}

		searchScope := goldap.ScopeSingleLevel
		if userSearchOpts.ScopeSubtree {
			searchScope = goldap.ScopeWholeSubtree
		}

		// Make search filter for the given user.
		searchFilter := fmt.Sprintf("(%s=%s)", userSearchOpts.UsernameAttr, username)
		if userSearchOpts.Filter != "" {
			searchFilter = AndSearchFilter(userSearchOpts.Filter, searchFilter)
		}

		searchOpts := SearchOpts{
			RootURL: rootURL,
			Scope:   searchScope,
			Filter:  searchFilter,
		}

		if result := SearchWithReferrals(ctx, settings, searchOpts); result != nil {
			// Found the user, no need to perform other searches.
			userDN = result.Entry.DN
			break
		}
	}

	if userDN == "" {
		return nil, logBuf, errors.New("no LDAP search results found")
	}

	// Attempt to connect to the LDAP server and Bind using the DN/password
	// combination.
	ldapConn, err := GetConn(settings.ServerURL, settings)
	if err != nil {
		return nil, logBuf, err
	}
	defer ldapConn.Close()

	if err := ldapConn.Bind(userDN, password); err != nil {
		settings.ReaderPassword = "****************" // hide the password!
		ctx = context.WithValues(ctx, map[string]interface{}{
			"error":      err.Error(),
			"DN":         userDN,
			"LDAPconfig": settings.String(),
		})
		context.GetLogger(ctx, "error", "DN").Info("failed to bind user")

		if ldapErr, ok := err.(*goldap.Error); ok && ldapErr.ResultCode == goldap.LDAPResultInvalidCredentials {
			return nil, logBuf, authn.ErrInvalidUsernamePassword()
		}

		return nil, logBuf, err
	}

	isAdmin := CheckIfMember(ctx, settings, settings.AdminSyncOpts, userDN)

	user = &authn.Account{
		Account: schema.Account{
			Name:    forms.NormalizeAccountName(username),
			IsAdmin: isAdmin,
		},
	}

	return user, logBuf, nil
}

// CheckIfMember determines if
func CheckIfMember(ctx context.Context, settings *config.Settings, memberSyncOpts config.MemberSyncOpts, memberDN string) (isMember bool) {
	if !memberSyncOpts.EnableSync {
		// Not configured to sync system admins via LDAP.
		return false
	}

	memberDN = strings.ToLower(memberDN) // Normalize the DN.

	// If membership is synced by direct group membership, handle that case
	// first as it's a bit simpler.
	if memberSyncOpts.SelectGroupMembers {
		rootURL := LDAPSearchURL{
			HostURLString: settings.ServerURL,
			BaseDN:        memberSyncOpts.GroupDN,
		}

		// Make search filter for the member attribute to equal the memberDN.
		searchFilter := fmt.Sprintf("(%s=%s)", memberSyncOpts.GroupMemberAttr, memberDN)

		searchOpts := SearchOpts{
			RootURL: rootURL,
			Scope:   goldap.ScopeBaseObject,
			Filter:  searchFilter,
		}

		return SearchWithReferrals(ctx, settings, searchOpts) != nil
	}

	// Otherwise, confirm that the member DN is within the search base DN
	// scope before making sure the member matches the filter.

	searchBaseDN := strings.ToLower(memberSyncOpts.SearchBaseDN)

	memberDNParts := strings.Split(memberDN, ",")
	searchBaseDNParts := strings.Split(searchBaseDN, ",")

	// The member DN must be in the subtree of the search base DN.
	if len(searchBaseDNParts) >= len(memberDNParts) {
		return false
	}
	// If searching only one level, then there may only be one
	// additional component of the member DN, like 'CN=username'.
	if !memberSyncOpts.SearchScopeSubtree && len(memberDNParts) > len(searchBaseDNParts)+1 {
		return false
	}

	// Iterate through components in reverse, ensuring the memberDN
	// contains the components of the search base DN.
	reverseStringSlice(memberDNParts)
	reverseStringSlice(searchBaseDNParts)
	for i, searchBaseDNPart := range searchBaseDNParts {
		if searchBaseDNPart != memberDNParts[i] {
			return false
		}
	}

	// Now we just need to ensure that the member object matches
	// the filter.
	rootURL := LDAPSearchURL{
		HostURLString: settings.ServerURL,
		BaseDN:        memberDN,
	}

	searchFilter := memberSyncOpts.SearchFilter
	if searchFilter == "" {
		searchFilter = "(objectClass=*)"
	}

	searchOpts := SearchOpts{
		RootURL: rootURL,
		Scope:   goldap.ScopeBaseObject,
		Filter:  AndSearchFilter(searchFilter),
	}

	return SearchWithReferrals(ctx, settings, searchOpts) != nil
}

// reverseStringSlice modifies the given slice in-place to reverse the ordering
// of its elements.
func reverseStringSlice(vals []string) {
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}
}
