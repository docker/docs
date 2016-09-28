package util

import (
	"fmt"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/stretchr/testify/require"
)

func (u *Util) IsSuiteRunningInLDAPMode() bool {
	return u.IntegrationFramework.Config.AuthMethodLDAP
}

func (u *Util) SetSuiteRunInLDAPMode() {
	u.IntegrationFramework.Config.AuthMethodLDAP = true
}

func (u *Util) IsSuiteRunningInManagedMode() bool {
	return !u.IntegrationFramework.Config.AuthMethodLDAP
}

func (u *Util) SetSuiteRunInManagedMode() {
	u.IntegrationFramework.Config.AuthMethodLDAP = false
}

func (u *Util) SwitchAuth() func() {
	if u.IsSuiteRunningInLDAPMode() {
		return u.SwitchToLDAPAuth()
	} else if u.IsSuiteRunningInManagedMode() {
		return u.SwitchToManagedAuth()
	}

	u.T().Fatal("Auth method specified was not managed or ldap")
	return func() {}
}

func (u *Util) SwitchToManagedAuth() func() {
	// default should be in managed auth, so this is not necessary???
	u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
	oldAuthSettings, err := u.API.EnziSession().GetAuthConfig()
	require.Nil(u.T(), err, "%s", err)

	// if oldAuthSettings.UseLDAP {
	// 	u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
	// 	_, err := u.API.EnziSession().SetAuthConfig(forms.AuthConfig{
	// 		UseLDAP: false,
	// 	})
	// 	require.Nil(u.T(), err, "%s", err)
	// }
	// require.False(u.T(), oldAuthSettings.UseLDAP, "You need to not be in LDAP auth mode")
	require.NotNil(u.T(), oldAuthSettings)
	return func() {
		u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
		_, err := u.API.EnziSession().SetAuthConfig(forms.AuthConfig{oldAuthSettings.Backend})
		// forms.AuthConfig{
		// UseLDAP: oldAuthSettings.UseLDAP,
		// LdapSettings: oldAuthSettings.LdapSettings, TODO I'm too
		// lazy to copy over all the fields, which is required
		// because they won't typecheck like this. in general these
		// should all be empty anyways
		// })
		require.Nil(u.T(), err, "%s", err)
	}

	// sqlUsers, err := u.API.ListSqlUsers()
	// require.Nil(u.T(), err)

	// sqlUsers = AddFormUsers([]User{{
	// 	Name:     u.Config.AdminUsername,
	// 	Password: u.Config.AdminPassword,
	// 	IsAdmin:  true,
	// }}, sqlUsers)

	// newAuthSettings, err := u.API.SetAuthSettings(&adminserver.AuthSettings{
	// 	Method: schema.AuthMethodManaged,
	// 	Managed: &adminserver.ManagedAuthSettings{
	// 		Users: sqlUsers,
	// 	},
	// })
	// require.Nil(u.T(), err)
	// require.Equal(u.T(), schema.AuthMethodManaged, newAuthSettings.Method)

	// for _, user := range newAuthSettings.Managed.Users {
	// 	if user.Username == u.Config.AdminUsername {
	// 		require.True(u.T(), *user.Account.IsAdmin, "isAdmin should be true: %#v", user)
	// 	}
	// }

	// return func() {
	// 	u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
	// 	_, err := u.API.SetAuthSettings(oldAuthSettings)
	// 	require.Nil(u.T(), err)
	// }
}

func (u *Util) SwitchToLDAPAuth() func() {
	settings := u.GetLDAPAuthSettings()
	return u.SwitchToSpecificLDAPAuth(settings)
}

func (u *Util) SwitchToSpecificLDAPAuth(settings forms.LDAPSettings) func() {
	u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
	oldAuthSettings, err := u.API.EnziSession().GetAuthConfig()
	require.Nil(u.T(), err)

	// require.False(u.T(), oldAuthSettings.UseLDAP, "You need to not be in LDAP auth mode")

	_, err = u.API.EnziSession().SetLDAPSettings(settings)
	require.Nil(u.T(), err)

	newAuthSettings, err := u.API.EnziSession().SetAuthConfig(forms.AuthConfig{
		Backend: "ldap",
	})
	require.Nil(u.T(), err)
	require.Equal(u.T(), newAuthSettings.Backend, "ldap")

	return func() {
		u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
		_, err := u.API.EnziSession().SetAuthConfig(forms.AuthConfig{
			Backend: oldAuthSettings.Backend,
		})
		require.Nil(u.T(), err, fmt.Sprintf("%s", err))
	}
}
