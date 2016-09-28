package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/docker/dhe-deploy"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var (
	ErrNoSuchNamespaceTeamAccess = errors.New("no such namespace team access")
)

// NamespaceTeamAccess represents a record of a team's access level on a
// repository namespace owned by the team's organization.
type NamespaceTeamAccess struct {
	// PK is a hash of the TeamID and the NamespaceID
	PK string `gorethink:"pk"`
	// NamespaceID is the UUID of the organization in eNZi whose namespace
	// the NamespaceTeamAccess is granting access to
	NamespaceID string `gorethink:"namespaceID"`
	// TeamID is the UUID of the team in eNZi to which this
	// NamespaceTeamAccess is granting access
	TeamID string `gorethink:"teamID"`
	// AccessLevel is the access level the team has on this repository, e.g.
	// "read-only", "read-write", "admin"
	AccessLevel string `gorethink:"accessLevel"`
}

var namespaceTeamAccessTable = table{
	db:         deploy.DTRDBName,
	name:       "namespace_team_access",
	primaryKey: "pk", // Guarantees uniqueness of (namespaceID, teamID)
	secondaryIndexes: map[string][]string{
		"teamID_namespaceID": {"teamID", "namespaceID"}, // For quickly listing a team's namespace accesses
		"namespaceID_teamID": {"namespaceID", "teamID"}, // For quickly listing a namespace's team accessors
	},
}

func (m *RepositoryAccessManager) SetNamespaceTeamAccess(nta *NamespaceTeamAccess) error {
	nta.PK = computeNamespaceTeamAccessPK(nta.NamespaceID, nta.TeamID)
	if _, err := namespaceTeamAccessTable.Term().Insert(
		nta,
		// If the team/namespace combination is already in the db, just
		// update the access level.
		rethink.InsertOpts{Conflict: "update"},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to set team access for namespace in database: %s", err)
	}
	return nil
}

func (m *RepositoryAccessManager) GetNamespaceAccessForTeam(namespaceID, teamID string) (*NamespaceTeamAccess, error) {
	pk := computeNamespaceTeamAccessPK(namespaceID, teamID)
	var nta NamespaceTeamAccess
	if err := namespaceTeamAccessTable.getRowByIndexVal(m.session, "pk", pk, &nta, ErrNoSuchNamespaceTeamAccess); err != nil {
		return nil, err
	}
	return &nta, nil
}

func (m *RepositoryAccessManager) DeleteNamespaceTeamAccess(namespaceID, teamID string) error {
	pk := computeNamespaceTeamAccessPK(namespaceID, teamID)
	if _, err := namespaceTeamAccessTable.Term().Get(pk).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete team access for namespace from database: %s", err)
	}
	return nil
}

// ListTeamsWithAccessToNamespace returns the teamIDs of the teams that have any access to a
// namespace
func (m *RepositoryAccessManager) ListTeamsWithAccessToNamespace(namespaceID, startID string, limit uint) (accesses []NamespaceTeamAccess, nextID string, err error) {
	query := namespaceTeamAccessTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "namespaceID_teamID"},
	).Between(
		[]interface{}{namespaceID, startID},
		[]interface{}{namespaceID, rethink.MaxVal},
	)
	return m.paginateNamespaceTeamAccessesQuery(query, limit)
}

// ListNamespaceAccessForTeam returns the teamIDs of the teams that have any access to a
// namespace
func (m *RepositoryAccessManager) ListNamespaceAccessForTeam(teamID, startID string, limit uint) (accesses []NamespaceTeamAccess, nextID string, err error) {
	query := namespaceTeamAccessTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "teamID_namespaceID"},
	).Between(
		[]interface{}{teamID, startID},
		[]interface{}{teamID, rethink.MaxVal},
	)
	return m.paginateNamespaceTeamAccessesQuery(query, limit)
}

func (m *RepositoryAccessManager) paginateNamespaceTeamAccessesQuery(query rethink.Term, limit uint) (accesses []NamespaceTeamAccess, nextID string, err error) {
	if limit > 0 {
		query = query.Limit(limit + 1)
	}

	cursor, err := query.Run(m.session)
	if err != nil {
		return nil, "", fmt.Errorf("unable to query db: %s", err)
	}

	if err := cursor.All(&accesses); err != nil {
		return nil, "", fmt.Errorf("unable to scan query results: %s", err)
	}

	if limit != 0 && uint(len(accesses)) > limit {
		nextID = accesses[limit].NamespaceID
		accesses = accesses[:limit]
	}

	return accesses, nextID, nil
}

func computeNamespaceTeamAccessPK(namespaceID, teamID string) string {
	hash := sha256.New()
	hash.Write([]byte(namespaceID))
	hash.Write([]byte(teamID))
	return hex.EncodeToString(hash.Sum(nil))
}
