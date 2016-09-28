package integration

import (
	"container/list"
	"fmt"
	"path"
	"testing"

	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/docker/orca/enzi/api/forms"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LastTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util
}

func (suite *LastTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
}

func (suite *LastTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *LastTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *LastTestSuite) TestDumpData() {
	// Don't run unless we want a data dump
	if suite.Config.CleanDataDump {
		return
	}

	cleanupFuncs := list.New()
	cleanupFuncs.PushFront(suite.u.SwitchAuth())

	testUsers := []testUser{}
	testOrgs := []testOrg{}

	// The structure of users/teams/orgs after the dump

	// Accounts:
	// admin user (already exists)
	//
	// Organizaitons:
	//        org0   org1
	//
	// Users (org0 teams)      (org1 teams):
	// user0  org0team0
	// user1  org0team0
	// user2  org0team0,1
	// user3  org0team0,1      org1team0
	// user4  org0team1        org1team0
	// user5  org0team1        org1team0,1
	// user6                   org1team0,1
	// user7                   org1team1
	// user8                   org1team1
	// user9
	// user10

	// Create 11 test users.
	testUsers = make([]testUser, 11)
	for i := range testUsers {
		username := fmt.Sprintf("user%d", i)
		cleanupUser := suite.u.CreateUserWithChecks(username, "password")
		suite.API.EnziSession().UpdateAccount(username, forms.UpdateAccount{IsActive: &[]bool{true}[0]})

		testUsers[i] = testUser{
			Name:          username,
			Password:      "password",
			OrgMembership: make(map[string][]string),
		}
		cleanupFuncs.PushFront(cleanupUser)
	}

	// Create 2 test organizations.
	testOrgs = make([]testOrg, 2)
	for i := range testOrgs {
		org := testOrg{
			Name:           fmt.Sprintf("org%d", i),
			TeamMembership: make(map[string][]string),
		}
		testOrgs[i] = org

		cleanupOrg := suite.u.CreateOrganizationWithChecks(org.Name)
		cleanupFuncs.PushFront(cleanupOrg)

		// Create 2 more teams with 4 users each.
		for j := 0; j < 2; j++ {
			teamname := fmt.Sprintf("org%dteam%d", i, j)
			// No need to call the cleanup func this returns. The
			// team will be deleted along with the organization.
			suite.u.CreateManagedTeamWithChecks(org.Name, teamname)

			org.TeamMembership[teamname] = make([]string, 4)
			for k := range org.TeamMembership[teamname] {
				// Offset by 3*i to offset for different orgs.
				// Offset by 2*j to offset for different teams.
				// This puts all users in a team expect for 2
				// which remain solo (not a member of any
				// organization).
				user := testUsers[(3*i)+(2*j)+k]
				org.TeamMembership[teamname][k] = user.Name
				user.OrgMembership[org.Name] = append(user.OrgMembership[org.Name], teamname)
			}
		}
	}

	// Populate the teams
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

	for _, org := range testOrgs {
		for teamname, members := range org.TeamMembership {
			for _, username := range members {
				_, err := suite.API.EnziSession().AddTeamMember(org.Name, teamname, username, forms.SetMembership{})
				suite.Assert().Nil(err, "expected AddTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, teamname, username, err)
			}
		}
	}

	// The structure of repo ownership/permissions

	// Repo   owner   priv/pub  rw         r          x
	// repo0  org0    priv      org0team0  org0team1
	// repo1  org0    priv                            org0team0
	// repo2  org0    pub       org0team0             org0team1
	//
	// repo3  user0   priv      user1      user2
	// repo4  user0   priv                            user1
	// repo5  user0   pub       user1                 user2

	// Pull some small images to push
	localImg := "tianon/true"
	err := suite.Docker.PullImage(localImg, nil)
	require.Nil(suite.T(), err)

	suite.createRepoTagPushImage(cleanupFuncs, "org0", "private", "repo0", localImg)
	cleanupFuncs.PushFront(suite.u.SetRepoTeamAccessWithChecks("org0", "repo0", "org0team0", "read-write"))
	cleanupFuncs.PushFront(suite.u.SetRepoTeamAccessWithChecks("org0", "repo0", "org0team1", "read-only"))

	suite.createRepoTagPushImage(cleanupFuncs, "org0", "private", "repo1", localImg)
	cleanupFuncs.PushFront(suite.u.SetRepoTeamAccessWithChecks("org0", "repo1", "org0team1", "admin"))

	suite.createRepoTagPushImage(cleanupFuncs, "org0", "public", "repo2", localImg)
	cleanupFuncs.PushFront(suite.u.SetRepoTeamAccessWithChecks("org0", "repo2", "org0team0", "read-write"))
	cleanupFuncs.PushFront(suite.u.SetRepoTeamAccessWithChecks("org0", "repo2", "org0team1", "admin"))

	suite.createRepoTagPushImage(cleanupFuncs, "user0", "private", "repo3", localImg)
	suite.API.SetRepositoryUserAccess("user0", "repo3", "user1", "read-write")
	suite.API.SetRepositoryUserAccess("user0", "repo3", "user2", "read-only")

	suite.createRepoTagPushImage(cleanupFuncs, "user0", "private", "repo4", localImg)
	suite.API.SetRepositoryUserAccess("user0", "repo4", "user1", "admin")

	suite.createRepoTagPushImage(cleanupFuncs, "user0", "public", "repo5", localImg)
	suite.API.SetRepositoryUserAccess("user0", "repo5", "user1", "read-write")
	suite.API.SetRepositoryUserAccess("user0", "repo5", "user2", "admin")

	for e := cleanupFuncs.Front(); e != nil; e = e.Next() {
		cleanupFunc := e.Value.(func())
		cleanupFunc()
	}
}

func (suite *LastTestSuite) createRepoTagPushImage(cleanupFuncs *list.List, namespace, privacy, repoName, image string) {
	cleanupFuncs.PushFront(suite.u.CreateRepoWithChecks(namespace, repoName, "blah", "blah blah", privacy))
	imageTag := "latest"
	cleanupFuncs.PushFront(suite.u.TagImageWithChecks(path.Join(suite.Config.DTRHost, namespace), repoName, imageTag, image))
	cleanupFuncs.PushFront(suite.u.PushImageWithChecks(path.Join(suite.Config.DTRHost, namespace), repoName, imageTag))

	return
}

func TestDataDumpSuite(t *testing.T) {
	suite.Run(t, new(LastTestSuite))
}
