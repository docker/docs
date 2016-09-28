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

// PrivatePublicTestSuite is an integration test suite for the DTR API.
type PrivatePublicTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	setOldAuth func()
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *PrivatePublicTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	suite.setOldAuth = suite.u.SwitchAuth()
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *PrivatePublicTestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *PrivatePublicTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *PrivatePublicTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *PrivatePublicTestSuite) TestOrgRepoPrivatePublic() {
	org := "docker"
	team2Name := "sales"
	testRepo := "docker-trusted-diskfiller"

	// set up just like in BaiscTests
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateOrganizationWithChecks(org)()
	_, deferFunc := suite.u.CreateManagedTeamWithChecks(org, team2Name)
	defer deferFunc()
	defer suite.u.CreateRepoWithChecks(org, testRepo, "blah", "blah blah", "private")()
	// create a regular user
	user, deluser := suite.u.CreateActivateRandomUser()
	defer deluser()

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
			ExplicitlyShared: false,
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
			ExplicitlyShared: false,
		},
	} {
		testDesc := fmt.Sprintf("Current test: %#v", test)

		// keep the repo private the first time, public the second time
		suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
		if test.AccessLevel != "" {
			public := "public"
			_, err := suite.API.UpdateRepository(org, testRepo, apiclient.RepositoryUpdateForm{
				Visibility: &public,
			})
			require.Nil(suite.T(), err)
		}

		// execute extensive permission checks
		suite.u.CheckUserPermissionsOnRepo(testDesc, test, org, testRepo, localImg1, localImg2, "", team2Name, user.Name, user.Password)
	}
}

func TestPrivatePublicSuite(t *testing.T) {
	suite.Run(t, new(PrivatePublicTestSuite))
}
