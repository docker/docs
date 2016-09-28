package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/docker/orca/enzi/api/forms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// RepoTeamsAPITestSuite is an integration test suite for the DTR API.
type RepoTeamsAPITestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	setOldAuth func()
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *RepoTeamsAPITestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	suite.setOldAuth = suite.u.SwitchAuth()
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *RepoTeamsAPITestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *RepoTeamsAPITestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *RepoTeamsAPITestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *RepoTeamsAPITestSuite) BasicTests(orgAdmin bool) {
	org := "nsa"
	teamName := "selfies-patrol"
	testRepo := "selfiestick-botnet"
	accessLevel := "read-write"

	// create nsa/selfies-patrol and give them access to nsa/selfiestick-botnet
	require.Nil(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateOrganizationWithChecks(org)()

	// if are are testing org admin permissions, we have to create the admin
	var admin util.User
	if orgAdmin {
		var delAdmin func()
		admin, delAdmin = suite.u.CreateActivateRandomUser()
		defer delAdmin()

		// become an admin of the org by getting added by the global admin
		_, err := suite.API.EnziSession().AddOrganizationMember(org, admin.Name, forms.SetMembership{IsAdmin: &[]bool{true}[0]})
		require.Nil(suite.T(), err)

		// do everything else as the org admin instead of the global admin
		require.Nil(suite.T(), suite.API.Login(admin.Name, admin.Password))
	}
	teamObj, deferFunc := suite.u.CreateManagedTeamWithChecks(org, teamName)
	defer deferFunc()
	defer suite.u.CreateRepoWithChecks(org, testRepo, "blah", "blah blah", "private")()

	// give read-write access and make sure we can see it
	revoke := suite.u.SetRepoTeamAccessWithChecks(org, testRepo, teamName, accessLevel)

	// check that the grant exists in both directions
	repo, teamAccess, err := suite.u.API.ListRepositoryTeamAccess(org, testRepo)
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), testRepo, repo.Name)
	require.Equal(suite.T(), 1, len(teamAccess))
	assert.Equal(suite.T(), accessLevel, teamAccess[0].AccessLevel)
	assert.Equal(suite.T(), teamObj.ID, teamAccess[0].Team.ID)
	team, repoAccess, err := suite.u.API.ListTeamRepositoryAccess(org, teamName)
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), teamObj.ID, team.ID)
	require.Equal(suite.T(), 1, len(repoAccess))
	assert.Equal(suite.T(), accessLevel, repoAccess[0].AccessLevel)
	assert.Equal(suite.T(), testRepo, repoAccess[0].Repository.Name)

	// revoke the grant
	revoke()

	// check that the grant is revoked
	repo, teamAccess, err = suite.u.API.ListRepositoryTeamAccess(org, testRepo)
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), testRepo, repo.Name)
	assert.Equal(suite.T(), 0, len(teamAccess))
	team, repoAccess, err = suite.u.API.ListTeamRepositoryAccess(org, teamName)
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), teamObj.ID, team.ID)
	assert.Equal(suite.T(), 0, len(repoAccess))
}

func (suite *RepoTeamsAPITestSuite) TestRepoTeamsAccessErrors() {
	org := "nsa"
	teamName := "selfies-patrol"
	testRepo := "selfiestick-botnet"
	accessLevel := "read-write"

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

	// test the account 404 cases
	_, err := suite.u.API.SetRepositoryTeamAccess("derp", testRepo, teamName, accessLevel)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchAccount)
	err = suite.u.API.RevokeRepositoryTeamAccess("derp", testRepo, teamName)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchAccount)
	_, _, err = suite.u.API.ListRepositoryTeamAccess("derp", testRepo)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchAccount)
	_, _, err = suite.u.API.ListTeamRepositoryAccess("derp", teamName)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchOrganization)

	defer suite.u.CreateOrganizationWithChecks(org)()

	// test the repo 404 cases
	_, err = suite.u.API.SetRepositoryTeamAccess(org, "derp", teamName, accessLevel)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)
	err = suite.u.API.RevokeRepositoryTeamAccess(org, "derp", teamName)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)
	_, _, err = suite.u.API.ListRepositoryTeamAccess(org, "derp")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)

	defer suite.u.CreateRepoWithChecks(org, testRepo, "blah", "blah blah", "private")()

	// test the team 404 cases
	_, err = suite.u.API.SetRepositoryTeamAccess(org, testRepo, "derp", accessLevel)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchTeam)
	err = suite.u.API.RevokeRepositoryTeamAccess(org, testRepo, "derp")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchTeam)
	_, _, err = suite.u.API.ListTeamRepositoryAccess(org, "derp")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchTeam)

	_, deferFunc := suite.u.CreateManagedTeamWithChecks(org, teamName)
	defer deferFunc()

	// test the repo 404 cases even if the team exists
	_, err = suite.u.API.SetRepositoryTeamAccess(org, "derp", teamName, accessLevel)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)
	err = suite.u.API.RevokeRepositoryTeamAccess(org, "derp", teamName)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)
}

// we try the basic team grant creation, deleting and listing as a global admin
func (suite *RepoTeamsAPITestSuite) TestRepoTeamsAccess() {
	suite.BasicTests(false)
}

// we try the basic team grant creation, deleting and listing as an org admin
func (suite *RepoTeamsAPITestSuite) TestRepoTeamsAccessOrgAdmin() {
	suite.BasicTests(true)
}

// Try actually doing various actions with team grants
func (suite *RepoTeamsAPITestSuite) TestRepoTeamsAccessTable() {
	org := "nsa"
	teamName := "selfies-patrol"
	team2Name := "partyvans"
	testRepo := "selfiestick-botnet"

	// set up just like in BaiscTests
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateOrganizationWithChecks(org)()
	_, deferFunc := suite.u.CreateManagedTeamWithChecks(org, teamName)
	defer deferFunc()
	_, deferFunc2 := suite.u.CreateManagedTeamWithChecks(org, team2Name)
	defer deferFunc2()
	defer suite.u.CreateRepoWithChecks(org, testRepo, "blah", "blah blah", "private")()
	// create a regular user and add them to the team
	user, deluser := suite.u.CreateActivateRandomUser()
	defer deluser()
	// add the user to the team
	removeFromTeam := suite.u.AddTeamMember(org, teamName, user.Name)

	// get an image we can push
	localImg1 := "tianon/true"
	err := suite.Docker.PullImage(localImg1, nil)
	require.Nil(suite.T(), err)

	// get a second image we can push
	localImg2 := "vmarmol/false"
	err = suite.Docker.PullImage(localImg2, nil)
	require.Nil(suite.T(), err)

	// now we try doing things with different access levels
	// TABLES FTW
	for _, test := range []util.UserPermissionsOnRepo{
		{
			AccessLevel:      "",
			Push:             false,
			Pull:             false,
			View:             false,
			DeleteTags:       false,
			EditDescription:  false,
			MakePublic:       false,
			ManageTeamAccess: false,
			ExplicitlyShared: true,
		},
		{
			AccessLevel:      "read-only",
			Push:             false,
			Pull:             true,
			View:             true,
			DeleteTags:       false,
			EditDescription:  false,
			MakePublic:       false,
			ManageTeamAccess: false,
			ExplicitlyShared: true,
		},
		{
			AccessLevel:      "read-write",
			Push:             true,
			Pull:             true,
			View:             true,
			DeleteTags:       true,
			EditDescription:  false,
			MakePublic:       false,
			ManageTeamAccess: false,
			ExplicitlyShared: true,
		},
		{
			AccessLevel:      "admin",
			Push:             true,
			Pull:             true,
			View:             true,
			DeleteTags:       true,
			EditDescription:  true,
			MakePublic:       true,
			ManageTeamAccess: true,
			ExplicitlyShared: true,
		},
		// we go back to none again at the end to make sure revocation actually works
		{
			AccessLevel:      "",
			Push:             false,
			Pull:             false,
			View:             false,
			DeleteTags:       false,
			EditDescription:  false,
			MakePublic:       false,
			ManageTeamAccess: false,
			ExplicitlyShared: true,
		},
	} {
		testDesc := fmt.Sprintf("Current test: %#v", test)

		// reset the state as admin before we start
		revoke := func() {}
		require.Nil(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
		if test.AccessLevel != "" {
			revoke = suite.u.SetRepoTeamAccessWithChecks(org, testRepo, teamName, test.AccessLevel)
		}
		short := "unchanged"
		long := "unchanged"
		private := "private"
		_, err := suite.API.UpdateRepository(org, testRepo, apiclient.RepositoryUpdateForm{
			ShortDescription: &short,
			LongDescription:  &long,
			Visibility:       &private,
		})
		require.Nil(suite.T(), err)

		// execute extensive permission checks
		suite.u.CheckUserPermissionsOnRepo(testDesc, test, org, testRepo, localImg1, localImg2, "", team2Name, user.Name, user.Password)

		revoke()
	}
	removeFromTeam()

	// after being removed from the team we should be back to not being able to do anything
	suite.u.CheckUserPermissionsOnRepo("after being removed from the team", util.UserPermissionsOnRepo{
		AccessLevel:      "",
		Push:             false,
		Pull:             false,
		View:             false,
		DeleteTags:       false,
		EditDescription:  false,
		MakePublic:       false,
		ManageTeamAccess: false,
		ExplicitlyShared: true,
	}, org, testRepo, localImg1, localImg2, "", team2Name, user.Name, user.Password)

	util.AppendDockerIgnorableLoggedErrors([]string{"Error streaming logs: unexpected EOF"})
}

func TestRepoTeamsAPISuite(t *testing.T) {
	suite.Run(t, new(RepoTeamsAPITestSuite))
}
