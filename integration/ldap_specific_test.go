package integration

import (
	// "net/http"
	// "strings"
	"testing"

	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	// enzierrors "github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	goldap "github.com/go-ldap/ldap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LDAPTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	switchConfigToManaged bool
}

func (suite *LDAPTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	if suite.u.IsSuiteRunningInManagedMode() {
		suite.switchConfigToManaged = true
		suite.u.SetSuiteRunInLDAPMode()
		suite.u.RestartLDAPContainer()
	}
}

func (suite *LDAPTestSuite) TearDownSuite() {
	if suite.switchConfigToManaged {
		suite.u.SetSuiteRunInManagedMode()
	}
}

func (suite *LDAPTestSuite) TearDownTest() {
	suite.u.TestLogs()
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

// TestLDAPLogin tests that users can log in with LDAP credentials
func (suite *LDAPTestSuite) TestLDAPLogin() {
	var err error

	suite.u.CreateUserInLDAPServer("adminuser", "password")
	suite.u.CreateUserInLDAPServer("adminUpper", "password")
	suite.u.CreateUserInLDAPServer("user", "password")
	suite.u.CreateUserInLDAPServer("Upper", "password")

	ldapConn := suite.u.GetBoundLDAPConn()
	defer ldapConn.Close()

	// Make adminUpper and adminuser admins according to the ldap server
	addRequest := goldap.NewAddRequest("cn=dtradmins,dc=example,dc=com")
	addRequest.Attribute("objectclass", []string{"groupofnames"})
	addRequest.Attribute("cn", []string{"dtradmins"})
	addRequest.Attribute("member", []string{"cn=adminuser,dc=example,dc=com", "cn=adminUpper,dc=example,dc=com"})
	err = ldapConn.Add(addRequest)
	assert.Nil(suite.T(), err)

	// before committing the config, test it out for each user
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	for _, username := range []string{"adminuser", "adminupper", "adminUpper", "user", "Upper", "upper"} {
		res, err := suite.API.EnziSession().TryLDAPLogin(forms.TryLdapLogin{
			Username:     username,
			Password:     "password",
			LDAPSettings: suite.u.GetLDAPAuthSettings(),
		})
		// if strings.ToLower(username) == username {
		assert.Nil(suite.T(), err, "%s", err)
		assert.NotNil(suite.T(), res)
		// } else {
		// assert.NotNil(suite.T(), err)
		// assert.Nil(suite.T(), res)
		// }
	}

	// now we switch to ldap mode and defer switching back
	setOldAuth := suite.u.SwitchAuth()
	defer setOldAuth()

	// make sure we delete the users at the end
	defer func() {
		suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
		assert.Nil(suite.T(), suite.u.DeleteAccount("user"))
		assert.Nil(suite.T(), suite.u.DeleteAccount("upper"))
		assert.Nil(suite.T(), suite.u.DeleteAccount("adminuser"))
		assert.Nil(suite.T(), suite.u.DeleteAccount("adminupper"))
	}()

	// FIXME: how to do this?? this is the first reason why this test is skipped.
	suite.u.LDAPSync()

	someUsers := []struct {
		username, password string
	}{
		{"user", "password"},
		{"upper", "password"},
		{"Upper", "password"},
		{"Upper", "password"},
		{"adminuser", "password"},
		{"adminUpper", "password"},
		{"adminupper", "password"},
	}
	// try to log in and check if it worked
	for _, test := range someUsers {
		suite.API.Login(test.username, test.password)
		session := suite.API.EnziSession()
		// if strings.ToLower(test.username) == test.username {
		if assert.Nil(suite.T(), err) {
			_, err := session.GetAccount(test.username)
			assert.Nil(suite.T(), err)
		}
		// } else {
		// 	assert.NotNil(suite.T(), err)
		// }
	}

	// Restart the LDAP container in order to wipe it

	// TODO re-enable when either fixed as bug or decided on what the proper behaviour is
	// try to log in and check that it failed
	// for _, test := range someUsers {
	// 	suite.API.Login(test.username, test.password)
	// 	suite.API.EnziSession()
	// 	// TODO
	// 	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, enzierrors.InvalidAuthentication(""))
	// }

	util.AppendDockerIgnorableLoggedErrors([]string{
		"Handler for DELETE /v1.21/containers/ldap_1 return",
		"no such id: ldap_1",        // docker 1.9
		"No such container: ldap_1", // docker 1.10
	})
}

func TestLDAPSuite(t *testing.T) {
	suite.Run(t, new(LDAPTestSuite))
}
