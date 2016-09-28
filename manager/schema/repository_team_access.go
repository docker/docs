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
	ErrNoSuchRepositoryTeamAccess = errors.New("no such repository team access")
)

// RepositoryTeamAccessDB represents a record of a team's access level on a
// repository owned by the team's organization.
type RepositoryTeamAccessDB struct {
	// PK is a hash of the TeamID and the RepositoryID
	PK string `gorethink:"pk"`
	// RepositoryID is the UUID of the repository the RepositoryTeamAccess
	// is granting access to
	RepositoryID string `gorethink:"repositoryID"`
	// TeamID is the UUID of the team in eNZi to which this
	// RepositoryTeamAccess is granting access
	TeamID string `gorethink:"teamID"`
	// AccessLevel is the access level the team has on this repository, e.g.
	// "read-only", "read-write", "admin"
	AccessLevel string `gorethink:"accessLevel"`
}

// RepositoryTeamAccess represents a team's access to a repository
type RepositoryTeamAccess struct {
	AccessLevel string     `gorethink:"accessLevel"`
	Repository  Repository `gorethink:"repository"`
	TeamID      string     `gorethink:"teamID"`
}

var repositoryTeamAccessTable = table{
	db:         deploy.DTRDBName,
	name:       "repository_team_access",
	primaryKey: "pk", // Guarantees uniqueness of (repositoryID, teamID)
	secondaryIndexes: map[string][]string{
		"teamID":              nil,                        // For quick lookups by team ID
		"repositoryID":        nil,                        // For quick lookups by repository ID
		"teamID_repositoryID": {"teamID", "repositoryID"}, // For quickly listing a team's repository accesses
		"repositoryID_teamID": {"repositoryID", "teamID"}, // For quickly listing a repository's team accessors
	},
}

func (m *RepositoryAccessManager) AddRepositoryTeamAccess(repoID, teamID, accessLevel string) error {
	rta := &RepositoryTeamAccessDB{
		PK:           computeRepositoryTeamAccessPK(repoID, teamID),
		RepositoryID: repoID,
		TeamID:       teamID,
		AccessLevel:  accessLevel,
	}
	if _, err := repositoryTeamAccessTable.Term().Insert(
		rta,
		// If the team/repo combination is already in the db, just
		// update the access level.
		rethink.InsertOpts{Conflict: "update"},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to set team access for repository in database: %s", err)
	}
	return nil
}

func (m *RepositoryAccessManager) GetRepositoryAccessForTeam(repositoryID, teamID string) (*RepositoryTeamAccessDB, error) {
	pk := computeRepositoryTeamAccessPK(repositoryID, teamID)
	var rta RepositoryTeamAccessDB
	if err := repositoryTeamAccessTable.getRowByIndexVal(m.session, "pk", pk, &rta, ErrNoSuchRepositoryTeamAccess); err != nil {
		return nil, err
	}
	return &rta, nil
}

func (m *RepositoryAccessManager) DeleteRepositoryTeamAccess(repositoryID, teamID string) error {
	pk := computeRepositoryTeamAccessPK(repositoryID, teamID)
	if _, err := repositoryTeamAccessTable.Term().Get(pk).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete team access for repository from database: %s", err)
	}
	return nil
}

// ListTeamsWithAccessToRepository returns the repositoryTeamAccesses of the teams that have any access to a
// repository
func (m *RepositoryAccessManager) ListTeamsWithAccessToRepository(repositoryID, startID string, limit uint) (accesses []RepositoryTeamAccess, nextID string, err error) {
	query := repositoryTeamAccessTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "repositoryID_teamID"},
	).Between(
		[]interface{}{repositoryID, startID},
		[]interface{}{repositoryID, rethink.MaxVal},
	).ConcatMap(getRepoTermById)

	getPaginateID := func(rta RepositoryTeamAccess) string { return rta.TeamID }
	return m.paginateRepositoryTeamAccessesQuery(query, limit, getPaginateID)
}

// ListRepositoryAccessForTeam returns the repositoryTeamAccesses of the
// repositories that the team has access to
func (m *RepositoryAccessManager) ListRepositoryAccessForTeam(teamID, startID string, limit uint) (accesses []RepositoryTeamAccess, nextID string, err error) {
	query := repositoryTeamAccessTable.Term().OrderBy(
		rethink.OrderByOpts{Index: "teamID_repositoryID"},
	).Between(
		[]interface{}{teamID, startID},
		[]interface{}{teamID, rethink.MaxVal},
	).ConcatMap(getRepoTermById)

	getPaginateID := func(rta RepositoryTeamAccess) string { return rta.Repository.ID }
	return m.paginateRepositoryTeamAccessesQuery(query, limit, getPaginateID)
}

func (m *RepositoryAccessManager) ListRepositoryAccessForTeamsAndRepositories(teamIDs []string, repositories []*Repository) (accesses []RepositoryTeamAccessDB, err error) {
	if len(teamIDs) == 0 || len(repositories) == 0 {
		return nil, nil
	}

	repoIDInterfaces := make([]interface{}, len(repositories))
	for i, repo := range repositories {
		repoIDInterfaces[i] = repo.ID
	}
	teamIDFilter := rethink.Row.Field("teamID").Eq(teamIDs[0])
	for _, teamID := range teamIDs[1:] {
		teamIDFilter = teamIDFilter.Or(rethink.Row.Field("teamID").Eq(teamID))
	}

	// We can only GetAllByIndex for one index, so we do teams,
	// assuming/guessing that there are less accesses per repo than there
	// are accesses per team. This also is nice because it maintains the
	// repository ordering.
	// An alternative to doing the filter would be to calculate nxm PKs and
	// do a GetAllByIndex on those...
	query := repositoryTeamAccessTable.Term().
		GetAllByIndex("repositoryID", repoIDInterfaces...).
		Filter(teamIDFilter)

	if cursor, err := query.Run(m.session); err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	} else if err := cursor.All(&accesses); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}

	return accesses, nil
}

// ListRepositoryAccessForTeams is not currently used and also doesn't work... TODO remove?
func (m *RepositoryAccessManager) ListRepositoryAccessForTeams(teamIDs []string, startID string, limit uint) (accesses []RepositoryTeamAccess, nextID string, err error) {
	teamIDInterfaces := make([]interface{}, len(teamIDs))
	for i, id := range teamIDs {
		teamIDInterfaces[i] = id
	}

	query := repositoryTeamAccessTable.Term().
		GetAllByIndex("teamID", teamIDInterfaces...).
		Filter(rethink.Row.Field("repositoryID").Ge(startID)).
		OrderBy("repositoryID").
		ConcatMap(getRepoTermById)

	getPaginateID := func(rta RepositoryTeamAccess) string { return rta.Repository.ID }
	return m.paginateRepositoryTeamAccessesQuery(query, limit, getPaginateID)
}

func getRepoTermById(rta rethink.Term) rethink.Term {
	return repositoriesTable.Term().GetAllByIndex(
		"id", rta.Field("repositoryID"),
	).Map(func(repo rethink.Term) interface{} {
		return map[string]interface{}{
			"accessLevel": rta.Field("accessLevel"),
			"repository":  repo,
			"teamID":      rta.Field("teamID"),
		}
	})
}

func (m *RepositoryAccessManager) paginateRepositoryTeamAccessesQuery(query rethink.Term, limit uint, getID func(RepositoryTeamAccess) string) (accesses []RepositoryTeamAccess, nextID string, err error) {
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
		nextID = getID(accesses[limit])
		accesses = accesses[:limit]
	}

	return accesses, nextID, nil
}

func computeRepositoryTeamAccessPK(repositoryID, teamID string) string {
	hash := sha256.New()
	hash.Write([]byte(repositoryID))
	hash.Write([]byte(teamID))
	return hex.EncodeToString(hash.Sum(nil))
}
