package integration

import (
	"fmt"
	"testing"

	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// RepoNSTeamAccessTestSuite is an integration test suite for the DTR API.
type RepoNSTeamAccessTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	setOldAuth func()
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *RepoNSTeamAccessTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	suite.setOldAuth = suite.u.SwitchAuth()
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *RepoNSTeamAccessTestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *RepoNSTeamAccessTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *RepoNSTeamAccessTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

// Try actually doing various actions with team grants
func (suite *RepoNSTeamAccessTestSuite) TestRepoNSTeamsAccessTable() {
	org := "nsa"
	teamName := "selfies-patrol"
	team2Name := "partyvans"
	testRepo := "selfiestick-botnet"

	// we set up one org with one team and one repo
	// then we control the access of the team to the repo
	// using namespace access permissions
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
			NSAccessLevel:    "",
			TeamsVisible:     true,
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
			NSAccessLevel:    "read-only",
			TeamsVisible:     true,
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
			NSAccessLevel:    "read-write",
			TeamsVisible:     true,
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
			NSAccessLevel:    "admin",
			TeamsVisible:     true,
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
			NSAccessLevel:    "",
			TeamsVisible:     true,
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
		suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
		if test.NSAccessLevel != "" {
			revoke = suite.u.SetRepoNSTeamAccessWithChecks(org, teamName, test.NSAccessLevel)
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
		suite.u.CheckUserPermissionsOnRepo(testDesc, test, org, testRepo, localImg1, localImg2, teamName, "", user.Name, user.Password)

		revoke()
	}
	removeFromTeam()

	// after being removed from the team we should be back to not being able to do anything
	suite.u.CheckUserPermissionsOnRepo("after being removed from the team", util.UserPermissionsOnRepo{
		AccessLevel:      "",
		NSAccessLevel:    "",
		TeamsVisible:     false,
		Push:             false,
		Pull:             false,
		View:             false,
		DeleteTags:       false,
		EditDescription:  false,
		MakePublic:       false,
		ManageTeamAccess: false,
		ExplicitlyShared: true,
	}, org, testRepo, localImg1, localImg2, teamName, "", user.Name, user.Password)

	util.AppendDockerIgnorableLoggedErrors([]string{"Error streaming logs: unexpected EOF"})
}

func TestRepoNSTeamAccessSuite(t *testing.T) {
	suite.Run(t, new(RepoNSTeamAccessTestSuite))
}
