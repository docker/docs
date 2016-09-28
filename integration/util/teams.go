package util

import (
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/api/responses"
	"github.com/stretchr/testify/assert"
)

func (u *Util) CreateManagedTeamWithChecks(orgname, teamname string) (*responses.Team, func()) {
	retFunc := func() {}

	team, err := u.API.EnziSession().CreateTeam(orgname, forms.CreateTeam{Name: teamname})
	if err != nil {
		u.T().Fatal(err)
	} else {
		retFunc = func() {
			u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
			err := u.API.EnziSession().DeleteTeam(orgname, teamname)
			assert.Nil(u.T(), err, "%s", err)
		}
		assert.NotNil(u.T(), team)
		assert.NotEmpty(u.T(), team.ID)
		assert.NotEmpty(u.T(), team.OrgID)
		// assert.Equal(u.T(), apiclient.TeamTypeManaged, team.Type)
		assert.Equal(u.T(), teamname, team.Name)
	}

	return team, retFunc
}

func (u *Util) AddTeamMember(orgname, teamname, username string) func() {
	retFunc := func() {}

	member, err := u.API.EnziSession().AddTeamMember(orgname, teamname, username, forms.SetMembership{})
	if err != nil {
		u.T().Fatal(err)
	} else {
		retFunc = func() {
			u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
			err := u.API.EnziSession().DeleteTeamMember(orgname, teamname, username)
			assert.Nil(u.T(), err, "%s", err)
		}
		assert.Equal(u.T(), member.Member.Name, username)
		assert.False(u.T(), member.IsAdmin)
		assert.False(u.T(), member.IsPublic)
	}

	return retFunc
}
