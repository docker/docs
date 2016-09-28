package schema

import (
	rethink "gopkg.in/dancannon/gorethink.v2"
)

// Different repository access levels.
const (
	AccessLevelReadOnly  = "read-only"
	AccessLevelReadWrite = "read-write"
	AccessLevelAdmin     = "admin"
)

// AccessLevelRanks provides an ordering of different repository access levels.
// Whenever a user's access on a repository has multiple distinct results, you
// may want to use the one that ranks highest. Notice that access levels that
// do not exist will have a rank of 0.
var AccessLevelRanks = map[string]uint{
	AccessLevelReadOnly:  1,
	AccessLevelReadWrite: 2,
	AccessLevelAdmin:     3,
}

// HighestRankingAccessLevel returns the highest ranking access level among
// the given access levels.
func HighestRankingAccessLevel(accessLevels ...string) (highestRankingAccessLevel string) {
	for _, accessLevel := range accessLevels {
		if AccessLevelRanks[accessLevel] > AccessLevelRanks[highestRankingAccessLevel] {
			highestRankingAccessLevel = accessLevel
		}
	}

	return highestRankingAccessLevel
}

// AccessLevelAtLeast returns whether the first access level is ranked
// greater than or equal to the second access level.
func AccessLevelAtLeast(accessLevel, atLeastAccesLevel string) bool {
	return AccessLevelRanks[accessLevel] >= AccessLevelRanks[atLeastAccesLevel]
}

// RepositoryAccessManager exports CRUDy methods for repository access in the
// database.
type RepositoryAccessManager struct {
	session *rethink.Session
}

// NewRepositoryAccessManager creates a new repository access manager using the
// given rethink session.
func NewRepositoryAccessManager(session *rethink.Session) *RepositoryAccessManager {
	return &RepositoryAccessManager{session}
}
