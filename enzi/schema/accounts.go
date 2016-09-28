package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var (
	// ErrNoSuchAccount conveys that an account with the given name or id
	// does not exist.
	ErrNoSuchAccount = errors.New("no such account")
	// ErrAccountExists conveys that an account with the given name
	// already exists.
	ErrAccountExists = errors.New("account already exists")
)

// An Account represents a user or organization.
type Account struct {
	Name     string `gorethink:"name"`     // Name of the account. Primary key.
	ID       string `gorethink:"id"`       // Randomly generated uuid for foreign references.
	FullName string `gorethink:"fullName"` // Full name of the account.
	IsOrg    bool   `gorethink:"isOrg"`    // Whether the account is an organization.

	// Fields for users only.
	IsAdmin      bool   `gorethink:"isAdmin"`      // Whether the user is a system admin. Always false for organizations.
	IsActive     bool   `gorethink:"isActive"`     // Whether the user is active. Always false for organizations.
	PasswordHash string `gorethink:"passwordHash"` // Managed user's salted+hashed password. Always empty for organizations.
	LdapDN       string `gorethink:"ldapDN"`       // LDAP user's synced Distinguished Name. Always empty for organizations.

	// Fields for organizaitons only.
	// Options for syncing the org's admins with LDAP. If enabled, org membership can only be set via teams.
	AdminSyncConfig MemberSyncOpts `gorethink:"adminSyncConfig"`
}

var accountsTable = table{
	db:         dbName,
	name:       "accounts",
	primaryKey: "name", // Guarantees uniqueness. Quick lookups by name.
	secondaryIndexes: map[string][]string{
		"id":           nil,                 // For quick lookups by user ID.
		"ldapDN":       nil,                 // For quick lookups by LDAP DN.
		"isOrg_name":   {"isOrg", "name"},   // For quickly listing all organizations.
		"isAdmin_name": {"isAdmin", "name"}, // For quickly listing all admins.
	},
}

// CreateAccount inserts a new account into the *accounts* table of the
// database using the values supplied by the given Account. The ID field of
// the account is set to a random UUID.
func (m *manager) CreateAccount(acct *Account) error {
	acct.ID = uuid.NewV4().String()

	if resp, err := accountsTable.Term().Insert(acct).RunWrite(m.session); err != nil {
		if isDuplicatePrimaryKeyErr(resp) {
			return ErrAccountExists
		}

		return fmt.Errorf("unable to create account in database: %s", err)
	}

	return nil
}

// getAccountByIndexVal queries the database for an account with the given
// value using the given index. If no such account exists the returned error
// will be ErrNoSuchAccount.
func (m *manager) getAccountByIndexVal(indexName string, val interface{}) (*Account, error) {
	cursor, err := accountsTable.Term().GetAllByIndex(indexName, val).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var acct Account
	if err := cursor.One(&acct); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchAccount
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	if acct.IsOrg {
		// Organization accounts should always be active.
		acct.IsActive = true
	}

	return &acct, nil
}

// GetAccountByName queries the database for an account with the given name. If
// no such account exists the returned error will be ErrNoSuchAccount.
func (m *manager) GetAccountByName(name string) (*Account, error) {
	return m.getAccountByIndexVal("name", name)
}

// GetAccountByID queries the database for an account with the given id. If no
// such account exists the returned error will be ErrNoSuchAccount.
func (m *manager) GetAccountByID(id string) (*Account, error) {
	return m.getAccountByIndexVal("id", id)
}

func (m *manager) ensureAccountType(acct *Account, isOrg bool) (*Account, error) {
	if acct.IsOrg != isOrg {
		return nil, ErrNoSuchAccount
	}

	return acct, nil
}

// GetUserByName queries the database for a user account with the given name.
// If no such user account exists the returned error will be ErrNoSuchAccount.
func (m *manager) GetUserByName(name string) (*Account, error) {
	acct, err := m.GetAccountByName(name)
	if err != nil {
		return nil, err
	}

	return m.ensureAccountType(acct, false)
}

// GetUserByID queries the database for a user account with the given id. If no
// such user account exists the returned error will be ErrNoSuchAccount.
func (m *manager) GetUserByID(id string) (*Account, error) {
	acct, err := m.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	return m.ensureAccountType(acct, false)
}

// GetUserByLdapDN queries the database for a user account with the given
// ldapDN. If no such user account exists the returned error will be
// ErrNoSuchAccount.
func (m *manager) GetUserByLdapDN(ldapDN string) (*Account, error) {
	acct, err := m.getAccountByIndexVal("ldapDN", ldapDN)
	if err != nil {
		return nil, err
	}

	return m.ensureAccountType(acct, false)
}

// GetOrgByName queries the database for an organization account with the given
// name. If no such organization account exists the returned error will be
// ErrNoSuchAccount.
func (m *manager) GetOrgByName(name string) (*Account, error) {
	acct, err := m.GetAccountByName(name)
	if err != nil {
		return nil, err
	}

	return m.ensureAccountType(acct, true)
}

// GetOrgByID queries the database for an organization account with the given
// id. If no such organization account exists the returned error will be
// ErrNoSuchAccount.
func (m *manager) GetOrgByID(id string) (*Account, error) {
	acct, err := m.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	return m.ensureAccountType(acct, true)
}

// ListAccounts queries for a slice of all accounts from the database ordered
// by name, starting from the given startName value. If limit is 0, all results
// are returned, otherwise returns at most limit results. Will return the
// startName of the next page or "" if no accounts remain.
func (m *manager) ListAccounts(startName string, limit uint) (accts []Account, nextName string, err error) {
	query := accountsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "name"},
	).Between(
		startName, rethink.MaxVal,
	)

	return paginateAccountsQuery(m.session, query, limit)
}

func paginateAccountsQuery(session *rethink.Session, query rethink.Term, limit uint) (accts []Account, nextName string, err error) {
	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(session)
	if err != nil {
		return nil, "", fmt.Errorf("unable to query db: %s", err)
	}

	accts = []Account{}
	if err := cursor.All(&accts); err != nil {
		return nil, "", fmt.Errorf("unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(accts)) > limit {
		nextName = accts[limit].Name
		accts = accts[:limit]
	}

	return accts, nextName, nil
}

// ListUsers queries for a slice of all user accounts from the database ordered
// by name, starting from the given startName value. If limit is 0, all results
// are returned, otherwise returns at most limit results. Will return the
// startName of the next page or "" if no user accounts remain.
func (m *manager) ListUsers(startName string, limit uint) (accts []Account, nextName string, err error) {
	query := accountsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "isOrg_name"},
	).Between(
		[]interface{}{false, startName},
		[]interface{}{false, rethink.MaxVal},
	)

	return paginateAccountsQuery(m.session, query, limit)
}

// ListOrgs queries for a slice of all organization accounts from the database
// ordered by name, starting from the given startName value. If limit is 0, all
// results are returned, otherwise returns at most limit results. Will return
// the startName of the next page or "" if no organization accounts remain.
func (m *manager) ListOrgs(startName string, limit uint) (accts []Account, nextName string, err error) {
	query := accountsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "isOrg_name"},
	).Between(
		[]interface{}{true, startName},
		[]interface{}{true, rethink.MaxVal},
	)

	return paginateAccountsQuery(m.session, query, limit)
}

// ListAdmins queries for a slice of all users marked as admins ordered
// by name, starting from the given startName value. If limit is 0, all results
// are returned, otherwise returns at most limit results. Will return the
// startName of the next page or "" if no admin users remain.
func (m *manager) ListAdmins(startName string, limit uint) (accts []Account, nextName string, err error) {
	query := accountsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "isAdmin_name"},
	).Between(
		[]interface{}{true, startName},
		[]interface{}{true, rethink.MaxVal},
	).Filter(map[string]interface{}{
		"isOrg": false,
	})

	return paginateAccountsQuery(m.session, query, limit)
}

// ListNonAdmins queries for a slice of all users not marked as admins
// ordered by name, starting from the given startName value. If limit is 0, all
// results are returned, otherwise returns at most limit results. Will return
// the startName of the next page or "" if no non-admin users remain.
func (m *manager) ListNonAdmins(startName string, limit uint) (accts []Account, nextName string, err error) {
	query := accountsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "isAdmin_name"},
	).Between(
		[]interface{}{false, startName},
		[]interface{}{false, rethink.MaxVal},
	).Filter(map[string]interface{}{
		"isOrg": false,
	})

	return paginateAccountsQuery(m.session, query, limit)
}

// AccountUpdateFields holds the fields of an account which are updatable. Any
// fields left nil will not be updated when used with the update method.
type AccountUpdateFields struct {
	FullName        *string         `gorethink:"fullName,omitempty"`
	IsAdmin         *bool           `gorethink:"isAdmin,omitempty"`
	IsActive        *bool           `gorethink:"isActive,omitempty"`
	PasswordHash    *string         `gorethink:"passwordHash,omitempty"`
	LdapDN          *string         `gorethink:"ldapDN,omitempty"`
	AdminSyncConfig *MemberSyncOpts `gorethink:"adminSyncConfig,omitempty"`
}

// UpdateAccount updates the values in the database for the account with the
// given id.
func (m *manager) UpdateAccount(id string, updateFields AccountUpdateFields) error {
	if _, err := accountsTable.Term().GetAllByIndex("id", id).Update(
		updateFields,
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to run update query: %s", err)
	}

	return nil
}

// DeleteAccount removes the account with the given id from the database.
func (m *manager) DeleteAccount(id string) error {
	if _, err := accountsTable.Term().GetAllByIndex("id", id).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete account from database: %s", err)
	}

	// TODO: Cleanup other tables which may refer to this account.

	return nil
}

// OrgMembership represents a user's membership in an organization.
type OrgMembership struct {
	OrgID    string `gorethink:"orgID"`    // ID of the organization.
	UserID   string `gorethink:"userID"`   // ID of the member user.
	PK       string `gorethink:"pk"`       // Hash of organization ID and user ID. Primary key.
	IsAdmin  bool   `gorethink:"isAdmin"`  // Whether the user is an admin of the organization.
	IsPublic bool   `gorethink:"isPublic"` // Whether the user's membership is visible to non-members.
}

func computeOrgMembershipPK(orgID, userID string) string {
	hash := sha256.New()

	hash.Write([]byte(orgID))
	hash.Write([]byte(userID))

	return hex.EncodeToString(hash.Sum(nil))
}

var orgMembershipTable = table{
	db:         dbName,
	name:       "org_membership",
	primaryKey: "pk", // Guarantees uniqueness of (org_id, user_id). Quick membership checks.
	secondaryIndexes: map[string][]string{
		"orgID_userID":          {"orgID", "userID"},             // For quickly listing members of an organization.
		"userID_orgID":          {"userID", "orgID"},             // For quickly listing organizations for a user.
		"orgID_isAdmin_userID":  {"orgID", "isAdmin", "userID"},  // For quicly listing organization admins.
		"orgID_isPublic_userID": {"orgID", "isPublic", "userID"}, // For quickly listing public members of an organization.
	},
}

type orgMembershipAttributes struct {
	OrgID    string `gorethink:"orgID"`              // ID of the organization.
	UserID   string `gorethink:"userID"`             // ID of the member user.
	PK       string `gorethink:"pk"`                 // Hash of organization ID and user ID. Primary key.
	IsAdmin  *bool  `gorethink:"isAdmin,omitempty"`  // Whether the user is an admin of the organization.
	IsPublic *bool  `gorethink:"isPublic,omitempty"` // Whether the user's membership is visible to non-members.
}

// AddOrgMembership creates or updates a given membership record for the given
// userID and orgID.
func (m *manager) AddOrgMembership(orgID, userID string, isAdmin, isPublic *bool) error {
	membership := orgMembershipAttributes{
		OrgID:    orgID,
		UserID:   userID,
		PK:       computeOrgMembershipPK(orgID, userID),
		IsAdmin:  isAdmin,
		IsPublic: isPublic,
	}

	if _, err := orgMembershipTable.Term().Insert(
		membership,
		// If the user is already in the org, update the membership
		// attributes (if they are not nil).
		rethink.InsertOpts{Conflict: "update"},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to set org membership in database: %s", err)
	}

	return nil
}

// MemberOrg holds information about a given user's membership in an org. It
// contains all fields of the Account struct along with whether or not the
// member is an org admin and whether the membership is public.
type MemberOrg struct {
	Org      Account `gorethink:"org"`
	IsAdmin  bool    `gorethink:"isAdmin"`
	IsPublic bool    `gorethink:"isPublic"`
}

// getMemberOrgTermByID uses the id index of the accounts table to build a
// memberOrg entry term from a given org membership table term. The given
// membership term should be a term representing a row from the org membership
// table.
func getMemberOrgTermByID(membership rethink.Term) rethink.Term {
	return accountsTable.Term().GetAllByIndex(
		"id", membership.Field("orgID"),
	).Map(
		func(org rethink.Term) interface{} {
			return map[string]interface{}{
				"org":      org,
				"isAdmin":  membership.Field("isAdmin").Default(false),
				"isPublic": membership.Field("isPublic").Default(false),
			}
		},
	)
}

// ListOrgsForUser queries for a slice of all orgs that the user with the given
// userID is a member of, ordered by id, starting from the given startID value.
// If limit is 0, all results are returned, otherwise returns at most limit
// results. Will return the startID of the next page or "" if no orgs remain.
func (m *manager) ListOrgsForUser(userID, startID string, limit uint) (memberOrgs []MemberOrg, nextID string, err error) {
	query := orgMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "userID_orgID"},
	).Between(
		[]interface{}{userID, startID},
		[]interface{}{userID, rethink.MaxVal},
	).ConcatMap(getMemberOrgTermByID)

	return paginateMemberOrgsQuery(m.session, query, limit)
}

func paginateMemberOrgsQuery(session *rethink.Session, query rethink.Term, limit uint) (memberOrgs []MemberOrg, nextID string, err error) {
	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(session)
	if err != nil {
		return nil, "", fmt.Errorf("unable to query db: %s", err)
	}

	memberOrgs = []MemberOrg{}
	if err := cursor.All(&memberOrgs); err != nil {
		return nil, "", fmt.Errorf("unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(memberOrgs)) > limit {
		nextID = memberOrgs[limit].Org.ID
		memberOrgs = memberOrgs[:limit]
	}

	return memberOrgs, nextID, nil
}

// MemberInfo holds information about a user's membership in an org or team. It
// contains all fields of the Account struct along with whether the member is
// an admin of the org or team and whether the membership is public.
type MemberInfo struct {
	Member   Account `gorethink:"member"`
	IsAdmin  bool    `gorethink:"isAdmin"`
	IsPublic bool    `gorethink:"isPublic"`
}

// getMemberTermByUserID uses the id index of the accounts table to builda member
// entry term from a given membership table term. The given membership term
// should be a term representing a row from the team or org membership tables.
func getMemberTermByUserID(membership rethink.Term) rethink.Term {
	return accountsTable.Term().GetAllByIndex(
		"id", membership.Field("userID"),
	).Map(
		func(member rethink.Term) interface{} {
			return map[string]interface{}{
				"member":   member,
				"isAdmin":  membership.Field("isAdmin").Default(false),
				"isPublic": membership.Field("isPublic").Default(false),
			}
		},
	)
}

// ListOrgMembers queries the database for a slice of all user accounts which
// are members of the org with the given orgID from the database. Members are
// ordered by id, starting from the given startID value. If limit is 0, all
// results are returned, otherwise returns at most limit results. Will return
// the startID of the next page or "" if no user accounts remain.
func (m *manager) ListOrgMembers(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error) {
	query := orgMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "orgID_userID"},
	).Between(
		[]interface{}{orgID, startID},
		[]interface{}{orgID, rethink.MaxVal},
	).ConcatMap(getMemberTermByUserID)

	return paginateMembersQuery(m.session, query, limit)
}

func paginateMembersQuery(session *rethink.Session, query rethink.Term, limit uint) (members []MemberInfo, nextID string, err error) {
	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(session)
	if err != nil {
		return nil, "", fmt.Errorf("unable to query db: %s", err)
	}

	members = []MemberInfo{}
	if err := cursor.All(&members); err != nil {
		return nil, "", fmt.Errorf("unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(members)) > limit {
		nextID = members[limit].Member.ID
		members = members[:limit]
	}

	return members, nextID, nil
}

// ListPublicOrgMembers gets a slice of all user accounts which are public
// members of the org with the given orgID from the database. Members are
// ordered by id, starting from the given startID value. If limit is 0, all
// results are returned, otherwise returns at most limit results. Will return
// the startID of the next page or "" if no public members remain.
func (m *manager) ListPublicOrgMembers(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error) {
	query := orgMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "orgID_isPublic_userID"},
	).Between(
		[]interface{}{orgID, true, startID},
		[]interface{}{orgID, true, rethink.MaxVal},
	).ConcatMap(getMemberTermByUserID)

	return paginateMembersQuery(m.session, query, limit)
}

func (m *manager) listOrgAdmins(orgID, startID string, isAdmin bool, limit uint) (orgMembers []MemberInfo, nextID string, err error) {
	query := orgMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "orgID_isAdmin_userID"},
	).Between(
		[]interface{}{orgID, isAdmin, startID},
		[]interface{}{orgID, isAdmin, rethink.MaxVal},
	).ConcatMap(getMemberTermByUserID)

	return paginateMembersQuery(m.session, query, limit)
}

// ListOrgAdmins queries for a slice of all accounts marked as admins of the
// org with the given orgID. Admins are ordered by ID, starting from the given
// startID value. If limit is 0, all results are returned, otherwise returns at
// most limit results. Will return the startID of the next page or "" if no org
// admins remain.
func (m *manager) ListOrgAdmins(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error) {
	return m.listOrgAdmins(orgID, startID, true, limit)
}

// ListOrgNonAdmins queries for a slice of all accounts not marked as admins of
// the org with the given orgID. Non-admins are ordered by id, starting from
// the given startID value. If limit is 0, all results are returned, otherwise
// returns at most limit results. Will return the startID of the next page or
// "" if no non-admins remain.
func (m *manager) ListOrgNonAdmins(orgID, startID string, limit uint) (orgMembers []MemberInfo, nextID string, err error) {
	return m.listOrgAdmins(orgID, startID, false, limit)
}

// GetOrgMembership gets the org membership info for the given orgID and
// userID. Returns (nil, nil) if the user is not a member of the org.
func (m *manager) GetOrgMembership(orgID, userID string) (*OrgMembership, error) {
	pk := computeOrgMembershipPK(orgID, userID)

	cursor, err := orgMembershipTable.Term().Get(pk).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var membership OrgMembership
	if err := cursor.One(&membership); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &membership, nil
}

// DeleteOrgMembership removes the org membership for the given orgID and
// userID from the database.
func (m *manager) DeleteOrgMembership(orgID, userID string) error {
	pk := computeOrgMembershipPK(orgID, userID)

	// First, remove the user from all teams in the org.
	if _, err := teamMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "userID_orgID_teamID"},
	).Between(
		[]interface{}{userID, orgID, rethink.MinVal},
		[]interface{}{userID, orgID, rethink.MaxVal},
	).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete org team membership from database: %s", err)
	}

	if _, err := orgMembershipTable.Term().Get(pk).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete org membership from database: %s", err)
	}

	return nil
}
