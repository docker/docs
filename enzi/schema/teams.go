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
	// ErrNoSuchTeam conveys that a team with the given name or id does
	// not exist.
	ErrNoSuchTeam = errors.New("no such team")
	// ErrTeamExists conveys that a team with the given orgID and name
	// already exists.
	ErrTeamExists = errors.New("team already exists")
)

// Team represents a named subset of members of an organization.
type Team struct {
	OrgID string `gorethink:"orgID"` // ID of the organization.
	Name  string `gorethink:"name"`  // Name of the team.
	PK    string `gorethink:"pk"`    // Hash of organization ID and team name. Primary key.
	ID    string `gorethink:"id"`    // Randomly generated uuid for foreign references.

	Description string `gorethink:"description"` // Description of the team.

	// Options for syncing this team's members with LDAP.
	MemberSyncConfig MemberSyncOpts `gorethink:"memberSyncConfig"`
}

// MemberSyncOpts specifies options for syncing members of an organization or
// team.
type MemberSyncOpts struct {
	EnableSync bool `gorethink:"enableSync"`

	// Whether to sync using groupDN+memberAttr selection or sync using a
	// user search filter.
	SelectGroupMembers bool `gorethink:"selectGroupMembers"`

	// These fields are used to sync users using direct group membership.
	GroupDN         string `gorethink:"groupDN"`
	GroupMemberAttr string `gorethink:"groupMemberAttr"`

	// These fields are used to sync users using a search.
	SearchBaseDN       string `gorethink:"searchBaseDN"`
	SearchScopeSubtree bool   `gorethink:"searchScopeSubtree"`
	SearchFilter       string `gorethink:"searchFilter"`
}

func computeTeamPK(orgID, name string) string {
	hash := sha256.New()

	hash.Write([]byte(orgID))
	hash.Write([]byte(name))

	return hex.EncodeToString(hash.Sum(nil))
}

var teamsTable = table{
	db:         dbName,
	name:       "teams",
	primaryKey: "pk", // Guarantees uniqueness of (org_id, name). Quick lookups.
	secondaryIndexes: map[string][]string{
		"id": nil, // For quick lookups by ID.
		"orgID_memberSyncConfig.enableSync_name": {"orgID", "memberSyncConfig.enableSync", "name"}, // For quickly listing teams in an org which are LDAP synced, ordered by name.
		"orgID_name":                             {"orgID", "name"},                                // For quickly listing teams in an org, ordered by name.
	},
}

// CreateTeam inserts a new team into the *teams* table of the database using
// the values supplied by the given team. The ID field of the team is set to a
// random UUID (if not already set). The PK field of the team is set to a hash
// of the organization ID and team name.
func (m *manager) CreateTeam(team *Team) error {
	if team.ID == "" {
		team.ID = uuid.NewV4().String()
	}
	team.PK = computeTeamPK(team.OrgID, team.Name)

	if resp, err := teamsTable.Term().Insert(team).RunWrite(m.session); err != nil {
		if isDuplicatePrimaryKeyErr(resp) {
			return ErrTeamExists
		}

		return fmt.Errorf("unable to create team: %s", err)
	}

	return nil
}

// RenameTeam changes the name of the given team to the given new name. This
// operation is more advanced than simply updating the team record as changing
// the PK of the record is not possible. A duplicate team with the same ID must
// be made and the original must then be deleted. This operation is not atomic,
// but if the copy succeeds and the delete fails, the delete of the old team by
// name can be retried. On success, the given team object will have its Name
// and PK field changed.
func (m *manager) RenameTeam(team *Team, newName string) error {
	// Set the new name and generate a new PK.
	oldName := team.Name
	team.Name = newName
	team.PK = computeTeamPK(team.OrgID, newName)

	// Create the new team.
	if resp, err := teamsTable.Term().Insert(team).RunWrite(m.session); err != nil {
		if isDuplicatePrimaryKeyErr(resp) {
			return ErrTeamExists
		}

		return fmt.Errorf("unable to create renamed team: %s", err)
	}

	if err := m.DeleteTeam(team.OrgID, oldName); err != nil {
		return fmt.Errorf("unable to delete old-named team: %s", err)
	}

	return nil
}

// getTeamByIndexVal queries the database for a team with the given value using
// the given index. If no such team exists the returned error will be
// ErrNoSuchTeam.
func (m *manager) getTeamByIndexVal(indexName string, val interface{}) (*Team, error) {
	cursor, err := teamsTable.Term().GetAllByIndex(indexName, val).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var team Team
	if err := cursor.One(&team); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchTeam
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &team, nil
}

// GetTeamByName queries the database for a team with the given name in an org
// with the given orgID. If no such team exists the returned error will be
// ErrNoSuchTeam.
func (m *manager) GetTeamByName(orgID, name string) (*Team, error) {
	return m.getTeamByIndexVal("pk", computeTeamPK(orgID, name))
}

// GetTeamByID queries the database for a team with the given id. If no such
// team exists the returned error will be ErrNoSuchTeam.
func (m *manager) GetTeamByID(id string) (*Team, error) {
	return m.getTeamByIndexVal("id", id)
}

// ListTeamMembersInOrg gets a slice of distinct userIDs of members of all teams
// with the given orgID, in no particular order. This method is used as a
// convenience function for syncing org membership with the set of all teams'
// memberships.
func (m *manager) ListTeamMembersInOrg(orgID string) (userIDs []string, err error) {
	query := teamsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "orgID_name"},
	).Between(
		[]interface{}{orgID, rethink.MinVal},
		[]interface{}{orgID, rethink.MaxVal},
	).ConcatMap(func(team rethink.Term) rethink.Term {
		return teamMembershipTable.Term().OrderBy(
			rethink.OrderByOpts{Index: "teamID_userID"},
		).Between(
			[]interface{}{team.Field("id"), rethink.MinVal},
			[]interface{}{team.Field("id"), rethink.MaxVal},
		).Field("userID")
	}).Distinct()

	cursor, err := query.Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	if err := cursor.All(&userIDs); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}

	return userIDs, nil
}

// ListTeamsInOrg queries for a slice of all teams with the given orgID,
// ordered by name, starting from the given startName value. If limit is 0, all
// results are returned, otherwise returns at most limit results. Will return
// the startName of the next page or "" if no teams remain.
func (m *manager) ListTeamsInOrg(orgID, startName string, limit uint) (teams []Team, nextName string, err error) {
	query := teamsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "orgID_name"},
	).Between(
		[]interface{}{orgID, startName},
		[]interface{}{orgID, rethink.MaxVal},
	)

	return paginateTeamsQuery(m.session, query, limit)
}

func paginateTeamsQuery(session *rethink.Session, query rethink.Term, limit uint) (teams []Team, nextName string, err error) {
	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(session)
	if err != nil {
		return nil, "", fmt.Errorf("unable to query db: %s", err)
	}

	teams = []Team{}
	if err := cursor.All(&teams); err != nil {
		return nil, "", fmt.Errorf("unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(teams)) > limit {
		nextName = teams[limit].Name
		teams = teams[:limit]
	}

	return teams, nextName, nil
}

// ListLDAPSyncTeamsInOrg returns a slice of teams with the given orgID which
// are configured to be synced with an LDAP group. Teams are ordered by name,
// starting from the given startName value. If limit is 0, all results are
// returned, otherwise returns at most limit results. Will return the startName
// of the next page or "" if no teams remain.
func (m *manager) ListLDAPSyncTeamsInOrg(orgID, startName string, limit uint) (teams []Team, nextName string, err error) {
	query := teamsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "orgID_memberSyncConfig.enableSync_name"},
	).Between(
		[]interface{}{orgID, true, startName},
		[]interface{}{orgID, true, rethink.MaxVal},
	)

	return paginateTeamsQuery(m.session, query, limit)
}

// TeamUpdateFields holds the fields of a team which are updatable. Any fields
// left nil will not be updated when used with the update method.
type TeamUpdateFields struct {
	Description      *string         `gorethink:"description,omitempty"`
	MemberSyncConfig *MemberSyncOpts `gorethink:"memberSyncConfig,omitempty"`
}

// UpdateTeam updates the values in the database for the team with the given
// id.
func (m *manager) UpdateTeam(id string, updateFields TeamUpdateFields) error {
	if _, err := teamsTable.Term().GetAllByIndex("id", id).Update(
		updateFields,
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to run update query: %s", err)
	}

	return nil
}

// DeleteTeam removes the team with the given id from the database.
func (m *manager) DeleteTeam(orgID, name string) error {
	if _, err := teamsTable.Term().Get(computeTeamPK(orgID, name)).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete team from database: %s", err)
	}

	// TODO: Cleanup other tables which may refer to this team.

	return nil
}

// TeamMembership represents a user's membership in a team.
type TeamMembership struct {
	TeamID   string `gorethink:"teamID"`   // ID of the team.
	OrgID    string `gorethink:"orgID"`    // ID of the team's org.
	UserID   string `gorethink:"userID"`   // ID of the member user.
	PK       string `gorethink:"pk"`       // Hash of team ID and user ID. Primary key.
	IsAdmin  bool   `gorethink:"isAdmin"`  // Whether the user is an admin of the team.
	IsPublic bool   `gorethink:"isPublic"` // Whether the user's membership is visible to non-members.
}

func computeTeamMembershipPK(teamID, userID string) string {
	hash := sha256.New()

	hash.Write([]byte(teamID))
	hash.Write([]byte(userID))

	return hex.EncodeToString(hash.Sum(nil))
}

var teamMembershipTable = table{
	db:         dbName,
	name:       "team_membership",
	primaryKey: "pk", // Guarantees uniqueness of (team_id, user_id). Quick membership checks.
	secondaryIndexes: map[string][]string{
		"teamID_userID":          {"teamID", "userID"},             // For quickly listing members of a team.
		"userID_orgID_teamID":    {"userID", "orgID", "teamID"},    // For quickly listing teams in an org for a user.
		"teamID_isAdmin_userID":  {"teamID", "isAdmin", "userID"},  // For quickly listing admin members of a team.
		"teamID_isPublic_userID": {"teamID", "isPublic", "userID"}, // For quickly listing public members of a team.
	},
}

type teamMembershipAttributes struct {
	TeamID   string `gorethink:"teamID"`             // ID of the team.
	OrgID    string `gorethink:"orgID"`              // ID of the team's org.
	UserID   string `gorethink:"userID"`             // ID of the member user.
	PK       string `gorethink:"pk"`                 // Hash of team ID and user ID. Primary key.
	IsAdmin  *bool  `gorethink:"isAdmin,omitempty"`  // Whether the user is an admin of the team.
	IsPublic *bool  `gorethink:"isPublic,omitempty"` // Whether the user's membership is visible to non-members.
}

// AddTeamMembership updates or creates a team membership record for the given
// userID and teamID. The user is also added as a member of the org if they
// are not already.
func (m *manager) AddTeamMembership(orgID, teamID, userID string, isAdmin, isPublic *bool) error {
	// If the user is not yet a member of the org, add them. This operation
	// will not alter their org membership attributes if they are already
	// a member of the org.
	if err := m.AddOrgMembership(orgID, userID, nil, nil); err != nil {
		return err
	}

	membership := teamMembershipAttributes{
		TeamID:   teamID,
		OrgID:    orgID,
		UserID:   userID,
		PK:       computeTeamMembershipPK(teamID, userID),
		IsAdmin:  isAdmin,
		IsPublic: isPublic,
	}

	if _, err := teamMembershipTable.Term().Insert(
		membership,
		// If the user is already in the org, update the membership
		// attributes (if they are not nil).
		rethink.InsertOpts{Conflict: "update"},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to set org membership in database: %s", err)
	}

	return nil
}

// MemberTeam holds information about a given user's membership in a team. It
// contains all fields of the Team struct along with whether or not the member
// is a team admin and whether the membership is public.
type MemberTeam struct {
	Team     Team `gorethink:"team"`
	IsAdmin  bool `gorethink:"isAdmin"`
	IsPublic bool `gorethink:"isPublic"`
}

// getMemberTeamTermByID uses the id index of the teams table to build a
// memberTeam entry term from a given team membership table term. The given
// membership term should be a term representing a row from the team membership
// table.
func getMemberTeamTermByID(membership rethink.Term) rethink.Term {
	return teamsTable.Term().GetAllByIndex(
		"id", membership.Field("teamID"),
	).Map(
		func(team rethink.Term) interface{} {
			return map[string]interface{}{
				"team":     team,
				"isAdmin":  membership.Field("isAdmin").Default(false),
				"isPublic": membership.Field("isPublic").Default(false),
			}
		},
	)
}

// ListTeamsInOrgForUser queries for a slice of all teams with the given orgID
// that the user with the given userID is a member of. Teams are ordered by id,
// starting from the given startID value. If limit is 0, all results are
// returned, otherwise returns at most limit results. Will return the startID
// of the next page or "" if no teams remain.
func (m *manager) ListTeamsInOrgForUser(orgID, userID, startID string, limit uint) (memberTeams []MemberTeam, nextID string, err error) {
	query := teamMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "userID_orgID_teamID"},
	).Between(
		[]interface{}{userID, orgID, startID},
		[]interface{}{userID, orgID, rethink.MaxVal},
	).ConcatMap(getMemberTeamTermByID)

	return paginateMemberTeamsQuery(m.session, query, limit)
}

func paginateMemberTeamsQuery(session *rethink.Session, query rethink.Term, limit uint) (memberTeams []MemberTeam, nextID string, err error) {
	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(session)
	if err != nil {
		return nil, "", fmt.Errorf("unable to query db: %s", err)
	}

	memberTeams = []MemberTeam{}
	if err := cursor.All(&memberTeams); err != nil {
		return nil, "", fmt.Errorf("unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(memberTeams)) > limit {
		nextID = memberTeams[limit].Team.ID
		memberTeams = memberTeams[:limit]
	}

	return memberTeams, nextID, nil
}

// ListTeamMembers queries the database for a slice of all user accounts which
// are members of the team with the given teamID from the database. Members are
// ordered by id, starting from the given startID value. If limit is 0, all
// results are returned, otherwise returns at most limit results. Will return
// the startID of the next page or "" if no user accounts remain.
func (m *manager) ListTeamMembers(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error) {
	query := teamMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "teamID_userID"},
	).Between(
		[]interface{}{teamID, startID},
		[]interface{}{teamID, rethink.MaxVal},
	).ConcatMap(getMemberTermByUserID)

	return paginateMembersQuery(m.session, query, limit)
}

// ListPublicTeamMembers gets a slice of all user accounts which are public
// members of the team with the given teamID. Members are ordered by id,
// starting from the given startID value. If limit is 0, all results are
// returned, otherwise returns at most limit results. Will return the startID
// of the next page or "" if no public members remain.
func (m *manager) ListPublicTeamMembers(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error) {
	query := teamMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "teamID_isPublic_userID"},
	).Between(
		[]interface{}{teamID, true, startID},
		[]interface{}{teamID, true, rethink.MaxVal},
	).ConcatMap(getMemberTermByUserID)

	return paginateMembersQuery(m.session, query, limit)
}

func (m *manager) listTeamAdmins(teamID, startID string, isAdmin bool, limit uint) (teamMembers []MemberInfo, nextID string, err error) {
	query := teamMembershipTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "teamID_isAdmin_userID"},
	).Between(
		[]interface{}{teamID, isAdmin, startID},
		[]interface{}{teamID, isAdmin, rethink.MaxVal},
	).ConcatMap(getMemberTermByUserID)

	return paginateMembersQuery(m.session, query, limit)
}

// ListTeamAdmins queries for a slice of all accounts marked as admins of the
// team with the given teamID. Admins are ordered by ID, starting from the
// given startID value. If limit is 0, all results are returned, otherwise
// returns at most limit results. Will return the startID of the next page or
// "" if no team admins remain.
func (m *manager) ListTeamAdmins(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error) {
	return m.listTeamAdmins(teamID, startID, true, limit)
}

// ListTeamNonAdmins queries for a slice of all accounts not marked as admins
// of the team with the given teamID. Non-admins are ordered by id, starting
// from the given startID value. If limit is 0, all results are returned,
// otherwise returns at most limit results. Will return the startID of the next
// page or "" if no non-admins remain.
func (m *manager) ListTeamNonAdmins(teamID, startID string, limit uint) (teamMembers []MemberInfo, nextID string, err error) {
	return m.listTeamAdmins(teamID, startID, false, limit)
}

// GetTeamMembership gets the team membership info for the given teamID and
// userID. Returns (nil, nil) if the user is not a member of the team.
func (m *manager) GetTeamMembership(teamID, userID string) (*TeamMembership, error) {
	pk := computeTeamMembershipPK(teamID, userID)

	cursor, err := teamMembershipTable.Term().Get(pk).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var membership TeamMembership
	if err := cursor.One(&membership); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &membership, nil
}

// DeleteTeamMembership removes the team membership for the given teamID and
// userID from the database.
func (m *manager) DeleteTeamMembership(teamID, userID string) error {
	pk := computeTeamMembershipPK(teamID, userID)

	if _, err := teamMembershipTable.Term().Get(pk).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete team membership from database: %s", err)
	}

	return nil
}
