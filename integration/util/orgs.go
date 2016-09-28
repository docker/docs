package util

import (
	"github.com/docker/orca/enzi/api/forms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (u *Util) CreateOrganizationWithChecks(name string) func() {
	retFunc := func() {}
	if acc, err := u.API.EnziSession().CreateAccount(forms.CreateAccount{
		Name:  name,
		IsOrg: true,
	}); err != nil {
		require.Nil(u.T(), err, "%s", err)
	} else {
		retFunc = func() {
			u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
			err := u.API.EnziSession().DeleteAccount(name)
			assert.Nil(u.T(), err, "%s", err)
		}
		assert.NotNil(u.T(), acc)
		assert.NotEmpty(u.T(), acc.ID)
		assert.True(u.T(), acc.IsOrg)
		assert.Equal(u.T(), name, acc.Name)
		assert.Nil(u.T(), acc.IsActive)
		assert.Nil(u.T(), acc.IsAdmin)
	}
	return retFunc
}
