package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (u *Util) SetRepoNSTeamAccessWithChecks(org, team, level string) func() {
	teamObj, err := u.API.EnziSession().GetTeam(org, team)
	require.Nil(u.T(), err)

	na, err := u.API.SetRepositoryNamespaceTeamAccess(org, team, level)
	require.Nil(u.T(), err)
	assert.Equal(u.T(), level, na.AccessLevel)
	assert.Equal(u.T(), teamObj.ID, na.Team.ID)
	assert.Equal(u.T(), org, string(na.Namespace))

	return func() {
		u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
		assert.Nil(u.T(), u.API.RevokeRepositoryNamespaceTeamAccess(org, team))
	}
}
