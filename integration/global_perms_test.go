package integration

import (
	// "fmt"
	"testing"

	// "github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	// "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// GlobalPermsTestSuite is an integration test suite for the DTR API.
type GlobalPermsTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	setOldAuth func()
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *GlobalPermsTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	suite.setOldAuth = suite.u.SwitchAuth()
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *GlobalPermsTestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *GlobalPermsTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *GlobalPermsTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

// TODO maybe bring back some of this test?
// func (suite *GlobalPermsTestSuite) TestGlobalPerms() {
// 	if suite.u.IsSuiteRunningInLDAPMode() {
// 		//TODO: the reason we skip this test now is that the default teams become LDAP Type
// 		//when DTR is in LDAP auth mode, which means we'd need infrastructure to add the members
// 		//in LDAP without syncing. While doable, this is not a good investment right now because we'll
// 		//soon be switching to eNZi
// 		return
// 	}

// 	globalOrg := "_global"
// 	org := "nsa"
// 	testRepo := "selfiestick-botnet"

// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

// 	// create a repo under an org that we can test access to
// 	defer suite.u.CreateOrganizationWithChecks(org)()
// 	defer suite.u.CreateRepoWithChecks(org, testRepo, "blah", "blah blah", "private")()

// 	// create a regular user who will be added to the global org
// 	user, deluser := suite.u.CreateActivateRandomUser()
// 	defer deluser()

// 	// get an image we can push
// 	localImg1 := "tianon/true"
// 	err := suite.Docker.PullImage(localImg1, nil)
// 	require.Nil(suite.T(), err)

// 	// get a second image we can push
// 	localImg2 := "vmarmol/false"
// 	err = suite.Docker.PullImage(localImg2, nil)
// 	require.Nil(suite.T(), err)

// 	for _, test := range []util.UserPermissionsOnRepo{
// 		{
// 			AccessLevel:      "",
// 			Push:             false,
// 			Pull:             false,
// 			View:             false,
// 			DeleteTags:       false,
// 			EditDescription:  false,
// 			MakePublic:       false,
// 			ManageTeamAccess: false,
// 			ExplicitlyShared: false,
// 		},
// 		{
// 			AccessLevel:      "read-only",
// 			Push:             false,
// 			Pull:             true,
// 			View:             true,
// 			DeleteTags:       false,
// 			EditDescription:  false,
// 			MakePublic:       false,
// 			ManageTeamAccess: false,
// 			ExplicitlyShared: false,
// 		},
// 		{
// 			AccessLevel:      "read-write",
// 			Push:             true,
// 			Pull:             true,
// 			View:             true,
// 			DeleteTags:       true,
// 			EditDescription:  false,
// 			MakePublic:       false,
// 			ManageTeamAccess: false,
// 			ExplicitlyShared: false,
// 		},
// 		{
// 			AccessLevel:      "admin",
// 			Push:             true,
// 			Pull:             true,
// 			View:             true,
// 			DeleteTags:       true,
// 			EditDescription:  true,
// 			MakePublic:       true,
// 			ManageTeamAccess: true,
// 			ExplicitlyShared: false,
// 		},
// 		// we go back to none again at the end to make sure revocation actually works
// 		{
// 			AccessLevel:      "",
// 			Push:             false,
// 			Pull:             false,
// 			View:             false,
// 			DeleteTags:       false,
// 			EditDescription:  false,
// 			MakePublic:       false,
// 			ManageTeamAccess: false,
// 			ExplicitlyShared: false,
// 		},
// 	} {
// 		testDesc := fmt.Sprintf("Current test: %#v", test)

// 		suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 		if test.AccessLevel != "" {
// 			require.Nil(suite.T(), suite.API.AddTeamMember(globalOrg, test.AccessLevel, user.Name), testDesc)
// 		}
// 		short := "unchanged"
// 		long := "unchanged"
// 		private := "private"
// 		_, err := suite.API.UpdateRepository(org, testRepo, apiclient.RepositoryUpdateForm{
// 			ShortDescription: &short,
// 			LongDescription:  &long,
// 			Visibility:       &private,
// 		})
// 		require.Nil(suite.T(), err, testDesc)

// 		// execute extensive permission checks
// 		suite.u.CheckUserPermissionsOnRepo(testDesc, test, org, testRepo, localImg1, localImg2, "", "", user.Name, user.Password)

// 		if test.AccessLevel != "" {
// 			suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 			require.Nil(suite.T(), suite.API.DeleteTeamMember(globalOrg, test.AccessLevel, user.Name), testDesc)
// 		}
// 	}

// 	util.AppendDockerIgnorableLoggedErrors([]string{"Error streaming logs: unexpected EOF"})
// }

func TestGlobalPermsSuite(t *testing.T) {
	suite.Run(t, new(GlobalPermsTestSuite))
}
