package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (u *Util) SetRepoTeamAccessWithChecks(org, repo, team, level string) func() {
	teamObj, err := u.API.EnziSession().GetTeam(org, team)
	require.Nil(u.T(), err)
	ra, err := u.API.SetRepositoryTeamAccess(org, repo, team, level)
	require.Nil(u.T(), err)
	assert.Equal(u.T(), level, ra.AccessLevel)
	assert.Equal(u.T(), teamObj.ID, ra.Team.ID)
	assert.Equal(u.T(), repo, ra.Repository.Name)
	assert.Equal(u.T(), org, ra.Repository.Namespace)

	return func() {
		u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
		assert.Nil(u.T(), u.API.RevokeRepositoryTeamAccess(org, repo, team))
	}
}
