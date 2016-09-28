package integration

import (
	"net/http"
	"testing"

	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
)

// AuthorizationAPITestSuite is an integration test suite for the DTR API.
type AuthorizationAPITestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	setOldAuth func()
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *AuthorizationAPITestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)

	suite.setOldAuth = suite.u.SwitchAuth()
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *AuthorizationAPITestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *AuthorizationAPITestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *AuthorizationAPITestSuite) TearDownTest() {
	suite.u.TestLogs()
}

// TestNotAuthenticated tests that api calls return 401s when attempted while not authenticated
func (suite *AuthorizationAPITestSuite) TestNotAuthenticated() {
	suite.API.Logout()

	_, _, err := suite.API.EnziSession().ListAccounts("", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, err = suite.API.EnziSession().CreateAccount(forms.CreateAccount{
		Name:  "need2authenticate",
		IsOrg: true,
	})
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	err = suite.u.DeleteAccount("need2authenticate")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, err = suite.API.EnziSession().GetAccount("need2authenticate")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	if suite.u.IsSuiteRunningInManagedMode() {
		err = suite.u.ChangePassword("need2authenticate", "old", "new")
		suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

		err = suite.u.ActivateUser("need2authenticate")
		suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

		err = suite.u.DeactivateUser("need2authenticate")
		suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())
	}

	_, _, err = suite.API.EnziSession().ListUserOrganizations("need2authenticate", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, _, err = suite.API.EnziSession().ListTeamMembers("need2authenticate", "need2authenticate", "", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, _, err = suite.API.EnziSession().ListOrganizationMembers("need2authenticate", "", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, err = suite.API.EnziSession().CreateTeam("need2authenticate", forms.CreateTeam{"need2authenticate", ""})
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, _, err = suite.API.EnziSession().ListOrganizationTeams("need2authenticate", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, err = suite.API.EnziSession().GetTeam("need2authenticate", "need2authenticate")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, err = suite.API.EnziSession().UpdateTeam("need2authenticate", "need2authenticate", forms.UpdateTeam{})
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	err = suite.API.EnziSession().DeleteTeam("need2authenticate", "need2authenticate")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, err = suite.API.EnziSession().GetTeamMember("need2authenticate", "need2authenticate", "need2authenticate")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	_, err = suite.API.EnziSession().AddTeamMember("need2authenticate", "need2authenticate", "need2authenticate", forms.SetMembership{})
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())

	err = suite.API.EnziSession().DeleteTeamMember("need2authenticate", "need2authenticate", "need2authenticate")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.AuthenticationRequired())
}

// TestNotAuthorized tests that api calls that require authorization return the correct error message when
// they do not have it
func (suite *AuthorizationAPITestSuite) TestNotAuthorized() {
	user, cleanup := suite.u.CreateActivateRandomUser()
	defer cleanup()

	user2, cleanup2 := suite.u.CreateActivateRandomUser()
	defer cleanup2()

	cleanupOrg := suite.u.CreateOrganizationWithChecks("testorg")
	defer cleanupOrg()

	suite.API.Login(user.Name, user.Password)

	err := suite.u.DeleteAccount(user2.Name)
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	if suite.u.IsSuiteRunningInManagedMode() {
		err = suite.u.ChangePassword(user2.Name, "oldpassword", "newpassword")
		suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

		err = suite.u.ActivateUser(user2.Name)
		suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

		err = suite.u.DeactivateUser(user2.Name)
		suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))
	}

	_, _, err = suite.API.EnziSession().ListUserOrganizations(user2.Name, "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	_, _, err = suite.API.EnziSession().ListOrganizationMembers("testorg", "", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	_, err = suite.API.EnziSession().CreateTeam("testorg", forms.CreateTeam{"notAuthorized", ""})
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	_, _, err = suite.API.EnziSession().ListOrganizationTeams("testorg", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	_, err = suite.API.EnziSession().GetTeam("testorg", "notAuthorized")
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	_, err = suite.API.EnziSession().UpdateTeam("testorg", "notAuthorized", forms.UpdateTeam{})
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	err = suite.API.EnziSession().DeleteTeam("testorg", "notAuthorized")
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	_, err = suite.API.EnziSession().GetTeamMember("testorg", "notAuthorized", "notAuthorized")
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	_, err = suite.API.EnziSession().AddTeamMember("testorg", "notAuthorized", "notAuthorized", forms.SetMembership{})
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	err = suite.API.EnziSession().DeleteTeamMember("testorg", "notAuthorized", "notAuthorized")
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))
}

// TestCreateAccount tests creation of user and organization accounts.
// Currently only in managed mode.
func (suite *AuthorizationAPITestSuite) TestCreateAccount() {
	var err error

	// Try creating a user account.

	// Authentication is not required.
	suite.API.Login("", "")

	// TODO reactivate when the bug for the wrong return code gets fixed.
	// Account names must not be longer than 30 characters.
	// _, err = suite.u.CreateUser("thisusernameisreallyreallyreallyreallylong", "password")
	// suite.u.RequireEnziErrorCode(err, http.StatusBadRequest, errors.CannotCreateUser(""))

	// Account names must match the pattern: ^[a-z0-9]+(?:[._-][a-z0-9]+)*$
	// _, err = suite.u.CreateUser("_.-._xXx_PonyLover_xXx_.-._", "password")
	// suite.u.RequireEnziErrorCode(err, http.StatusBadRequest, errors.CannotCreateUser(""))

	// This test user should be created successfully.
	testUsername := "testuser"
	testUserPassword := util.GenerateRandomPassword(12)
	_, err = suite.u.CreateUser(testUsername, testUserPassword)
	require.Nil(suite.T(), err, "unable to create user: %s", err)

	// Cleanup the user that was just created.
	defer func() {
		suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
		err := suite.u.DeleteAccount(testUsername)
		require.Nil(suite.T(), err, "unable to cleanup created test user: %s", err)
	}()

	// Activate this user to use later. Requires authenticating as a system
	// admin user.
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	err = suite.u.ActivateUser(testUsername)
	require.Nil(suite.T(), err, "unable to activate test user: %s", err)

	// Creating a user with the same name should result in an error.
	util.AppendDTRIgnorableLoggedErrors([]string{"duplicate key value violates unique constraint"})
	_, err = suite.u.CreateUser(testUsername, testUserPassword)
	if suite.u.IsSuiteRunningInManagedMode() {
		suite.u.RequireEnziErrorCode(err, http.StatusBadRequest, errors.AccountExists())
	} else if suite.u.IsSuiteRunningInLDAPMode() {
		// It will be an ldap error
		require.NotNil(suite.T(), err)
	}

	// Try creating an organization account.

	// Authenticating as a non-admin isn't enough.
	suite.API.Login(testUsername, testUserPassword)
	_, err = suite.API.EnziSession().CreateAccount(forms.CreateAccount{Name: "need2beadmin", IsOrg: true})
	suite.u.RequireEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	// Only system admins can create an organization.
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	orgname := "testorg"
	_, err = suite.API.EnziSession().CreateAccount(forms.CreateAccount{Name: orgname, IsOrg: true})
	require.Nil(suite.T(), err, "unable to create organization: %s", err)

	// Cleanup the org that was just created.
	defer func() {
		suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
		err := suite.u.DeleteAccount(orgname)
		require.Nil(suite.T(), err, "unable to cleanup created test org: %s", err)
	}()

	// The org should have an "owners" team created automatically.
	// _, err = suite.API.EnziSession().GetTeam(orgname, "owners")
	// require.Nil(suite.T(), err, "unable to get owners team: %s", err)
}

// TestListGetAccount tests listing and getting details on created accounts.
func (suite *AuthorizationAPITestSuite) TestListGetDeleteAccount() {
	user, cleanup := suite.u.CreateActivateRandomUser()
	// just in case, since we will manually delete this user
	defer cleanup()

	suite.API.Login(user.Name, user.Password)

	// Activated user should be allowed to list accounts
	accounts, _, err := suite.API.EnziSession().ListAccounts("", "", 0)
	require.Nil(suite.T(), err, "unable to list users: %s", err)
	oldLength := len(accounts.Accounts)

	// Activated user should be allowed to get his own account
	account, err := suite.API.EnziSession().GetAccount(user.Name)
	require.Nil(suite.T(), err, "unable to get own user: %s", err)
	require.Equal(suite.T(), account.IsOrg, false, "should be a user")
	require.Equal(suite.T(), account.Name, user.Name, "the name should match")
	require.Equal(suite.T(), *account.IsActive, true, "should be active")

	// Activated user should be allowed to get other accounts
	account, err = suite.API.EnziSession().GetAccount(suite.Config.AdminUsername)
	require.Nil(suite.T(), err, "unable to get other user: %s", err)
	require.Equal(suite.T(), account.IsOrg, false, "should be a user")
	require.Equal(suite.T(), account.Name, suite.Config.AdminUsername, "the name should match")
	require.Equal(suite.T(), *account.IsActive, true, "should be active")

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	err = suite.u.DeleteAccount(user.Name)
	require.Nil(suite.T(), err, "unable to cleanup created test user: %s", err)

	// There's one less account now
	accounts, _, err = suite.API.EnziSession().ListAccounts("", "", 0)
	require.Nil(suite.T(), err, "unable to list users as admin: %s", err)
	newLength := len(accounts.Accounts)
	require.Equal(suite.T(), newLength+1, oldLength, "lengths should differ by one after deleting one user")

	// Should no longer be able to get the user
	_, err = suite.API.EnziSession().GetAccount(user.Name)
	suite.u.RequireEnziErrorCode(err, http.StatusNotFound, errors.NoSuchAccount(""))
}

func (suite *AuthorizationAPITestSuite) TestChangePassword() {
	user, cleanup := suite.u.CreateActivateRandomUser()
	// just in case, since we will manually delete this user
	defer cleanup()

	suite.API.Login(user.Name, user.Password)
	suite.u.ChangePassword(user.Name, user.Password, "newpassword")

	suite.API.Login(user.Name, user.Password)
	err := suite.u.ChangePassword(user.Name, user.Password, "newnewpassword")
	suite.u.AssertEnziErrorCode(err, http.StatusBadRequest, errors.PasswordIncorrect())

	suite.API.Login(user.Name, "newpassword")
	_, err = suite.API.EnziSession().GetAccount(user.Name)
	require.Nil(suite.T(), err, "you are now properly logged in: %s", err)
}

// TestActivateDeactivateAccount tests whether activating and deactivating an account properly allows
// and disallows authenticating a user
func (suite *AuthorizationAPITestSuite) TestActivateDeactivateAccount() {
	// Activate and Deactivate are not supported in LDAP mode
	if suite.u.IsSuiteRunningInLDAPMode() {
		return
	}

	// Create an unactivated user
	cleanup := suite.u.CreateUserWithChecks("user", "password")
	defer cleanup()

	suite.API.Login("user", "password")

	_, err := suite.API.EnziSession().GetAccount("user")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.InactiveAccount())

	_, _, err = suite.API.EnziSession().ListAccounts("user", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.InactiveAccount())

	// User shouldn't be able to activate self
	err = suite.u.ActivateUser("user")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.InactiveAccount())

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

	_, err = suite.API.EnziSession().GetAccount("nonexistantuser")
	suite.u.AssertEnziErrorCode(err, http.StatusNotFound, errors.NoSuchAccount(""))

	account, err := suite.API.EnziSession().GetAccount("user")
	assert.Equal(suite.T(), *account.IsActive, false, "should not be active")

	err = suite.u.ActivateUser("user")
	require.Nil(suite.T(), err, "admin should be able to activate: %s", err)

	account, err = suite.API.EnziSession().GetAccount("user")
	assert.Nil(suite.T(), err, "admin should be able to get user: %s", err)
	assert.Equal(suite.T(), *account.IsActive, true, "should be active")

	suite.API.Login("user", "password")

	_, err = suite.API.EnziSession().GetAccount("user")
	assert.Nil(suite.T(), err, "now that active, user should be able to get self: %s", err)

	// User shouldn't be able to deactivate self
	err = suite.u.DeactivateUser("user")
	suite.u.AssertEnziErrorCode(err, http.StatusForbidden, errors.NotAuthorized(""))

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	err = suite.u.DeactivateUser("user")
	assert.Nil(suite.T(), err, "admin should be able to deactivate: %s", err)

	// Test that user has lost priviledges after deactivation
	suite.API.Login("user", "password")

	_, err = suite.API.EnziSession().GetAccount("user")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.InactiveAccount())

	_, _, err = suite.API.EnziSession().ListAccounts("user", "", 0)
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.InactiveAccount())

	// User shouldn't be able to activate self
	err = suite.u.ActivateUser("user")
	suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, errors.InactiveAccount())
}

// TestActivateDeactivateAccount tests whether activating and deactivating an account properly allows
// and disallows authenticating a user
func (suite *AuthorizationAPITestSuite) TestCreateGetListUpdateDeleteTeam() {
	cleanupOrg := suite.u.CreateOrganizationWithChecks("testorg")
	defer cleanupOrg()

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

	_, err := suite.API.EnziSession().CreateTeam("nonexistantorgname", forms.CreateTeam{"someteamname", "description"})
	suite.u.AssertEnziErrorCode(err, http.StatusNotFound, errors.NoSuchAccount(""))

	// _, err = suite.API.EnziSession().CreateTeam("testorg", forms.CreateTeam{"thisusernameisreallyreallyreallyreallylong", "description"})
	// suite.u.AssertEnziErrorCode(err, http.StatusBadRequest, errors.CannotCreateUser(""))

	// _, err = suite.API.EnziSession().CreateTeam("testorg", forms.CreateTeam{"_.-._xXx_PonyLover_xXx_.-._", "description"})
	// suite.u.AssertEnziErrorCode(err, http.StatusBadRequest, errors.CannotCreateUser(""))

	teams, _, err := suite.API.EnziSession().ListOrganizationTeams("testorg", "", 0)
	assert.Nil(suite.T(), err, "should be able to list teams for the new org: %s", err)
	assert.Equal(suite.T(), 0, len(teams.Teams), "new org should have no teams")

	// _, err = suite.API.EnziSession().GetTeam("testorg", "owners")
	// assert.Nil(suite.T(), err, "new org should have an owners team: %s", err)

	// Team should not yet exist
	_, err = suite.API.EnziSession().GetTeam("testorg", "testteam")
	suite.u.AssertEnziErrorCode(err, http.StatusNotFound, errors.NoSuchTeam(""))

	_, err = suite.API.EnziSession().UpdateTeam("testorg", "testteam", forms.UpdateTeam{})
	suite.u.AssertEnziErrorCode(err, http.StatusNotFound, errors.NoSuchTeam(""))

	// Actually create the team
	_, cleanupTeam := suite.u.CreateManagedTeamWithChecks("testorg", "testteam")
	defer cleanupTeam()

	teams, _, err = suite.API.EnziSession().ListOrganizationTeams("testorg", "", 0)
	assert.Nil(suite.T(), err, "should be able to list teams for the org: %s", err)
	assert.Equal(suite.T(), 1, len(teams.Teams), "new org should have testteam")

	// _, err = suite.API.EnziSession().GetTeam("testorg", "owners")
	// assert.Nil(suite.T(), err, "org should still have an owners team: %s", err)

	team, err := suite.API.EnziSession().GetTeam("testorg", "testteam")
	assert.Nil(suite.T(), err, "org should now have the new team: %s", err)
	assert.NotNil(suite.T(), team)
	assert.NotEmpty(suite.T(), team.ID)
	assert.NotEmpty(suite.T(), team.OrgID)
	assert.Equal(suite.T(), "testteam", team.Name)

	teamName := "newtestteam"
	_, err = suite.API.EnziSession().UpdateTeam("testorg", "testteam", forms.UpdateTeam{Name: &teamName})
	assert.Nil(suite.T(), err, "should be able to update team: %s", err)

	// Team should no longer exist
	_, err = suite.API.EnziSession().GetTeam("testorg", "testteam")
	suite.u.AssertEnziErrorCode(err, http.StatusNotFound, errors.NoSuchTeam(""))

	team, err = suite.API.EnziSession().GetTeam("testorg", "newtestteam")
	assert.Nil(suite.T(), err, "org should now have the new team: %s", err)
	assert.NotNil(suite.T(), team)
	assert.NotEmpty(suite.T(), team.ID)
	assert.NotEmpty(suite.T(), team.OrgID)
	assert.Equal(suite.T(), "newtestteam", team.Name)

	// There should still be two teams
	teams, _, err = suite.API.EnziSession().ListOrganizationTeams("testorg", "", 0)
	assert.Nil(suite.T(), err, "should be able to list teams for the org: %s", err)
	assert.Equal(suite.T(), 1, len(teams.Teams), "new org should have newtestteam")

	// Should fail since there is no longer such a team
	_, err = suite.API.EnziSession().GetTeam("testorg", "testteam")
	suite.u.AssertEnziErrorCode(err, http.StatusNotFound, errors.NoSuchTeam(""))

	err = suite.API.EnziSession().DeleteTeam("testorg", "newtestteam")
	assert.Nil(suite.T(), err, "should be able to delete this team: %s", err)

	// Team should no longer exist
	_, err = suite.API.EnziSession().GetTeam("testorg", "newtestteam")
	suite.u.AssertEnziErrorCode(err, http.StatusNotFound, errors.NoSuchTeam(""))

	teams, _, err = suite.API.EnziSession().ListOrganizationTeams("testorg", "", 0)
	assert.Nil(suite.T(), err, "should be able to list teams for the new org: %s", err)
	assert.Equal(suite.T(), 0, len(teams.Teams), "we should be back to no teams")
}

func TestAuthorizationAPISuite(t *testing.T) {
	suite.Run(t, new(AuthorizationAPITestSuite))
}
