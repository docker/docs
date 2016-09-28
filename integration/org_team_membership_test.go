package integration

// TODO do this later...
// import (
// 	"container/list"
// 	"fmt"
// 	"net/http"
// 	"testing"
//
// 	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
// 	"github.com/docker/dhe-deploy/integration/framework"
// 	"github.com/docker/dhe-deploy/integration/util"
//
// 	"github.com/docker/enzi/api/forms"
//
// 	// "github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )
//
// // OrgTeamMembershipAPITestSuite is an integration test suite for the DTR API's
// // team membership endpoints.
// type OrgTeamMembershipAPITestSuite struct {
// 	suite.Suite
// 	*framework.IntegrationFramework
// 	u *util.Util
//
// 	testUsers    []testUser
// 	testOrgs     []testOrg
// 	cleanupFuncs *list.List
// }
//
type testUser struct {
	Name          string
	Password      string
	OrgMembership map[string][]string // maps orgname to a list of team names.
}

type testOrg struct {
	Name           string
	TeamMembership map[string][]string // maps team name to a list of member usernames.
}

//
// func TestOrgTeamMembershipAPITestSuite(t *testing.T) {
// 	suite.Run(t, new(OrgTeamMembershipAPITestSuite))
// }
//
// // SetupSuite handles setting up the API test suite by initializing auth in
// // managed mode with our initial admin user. It also sets up the suite with
// // some other users, orgs, and teams for testing team membership API endpoints.
// func (suite *OrgTeamMembershipAPITestSuite) SetupSuite() {
// 	suite.IntegrationFramework, suite.u = setupFramework(suite)
//
// 	suite.cleanupFuncs = list.New()
// 	suite.cleanupFuncs.PushFront(suite.u.SwitchAuth())
//
// 	// Accounts:
// 	// admin user (already exists)
// 	//
// 	// Organizaitons:
// 	//        org0   org1
// 	//
// 	// Users (org0 teams) (org1 teams):
// 	// user0  owners
// 	// user1  owners
// 	// user2               owners
// 	// user3               owners
// 	// user4  team0
// 	// user5  team0
// 	// user6  team0,1
// 	// user7  team0,1      team0
// 	// user8  team1        team0
// 	// user9  team1        team0,1
// 	// user10              team0,1
// 	// user11              team1
// 	// user12              team1
// 	// user13
// 	// user14
// 	//
// 	// Note: This suite setup does *not* place the users in these teams.
// 	// The teams are only created and the desired structure is created in
// 	// a way that can be quickly iterated on by the test methods.
//
// 	// Create 15 test users.
// 	suite.testUsers = make([]testUser, 15)
// 	for i := range suite.testUsers {
// 		user, cleanupUser := suite.u.CreateActivateRandomUser()
// 		suite.testUsers[i] = testUser{
// 			Name:          user.Name,
// 			Password:      user.Password,
// 			OrgMembership: make(map[string][]string),
// 		}
// 		suite.cleanupFuncs.PushFront(cleanupUser)
// 	}
//
// 	// Create 2 test organizations.
// 	suite.testOrgs = make([]testOrg, 2)
// 	for i := range suite.testOrgs {
// 		org := testOrg{
// 			Name:           "org_" + util.GenerateRandomUsername(8),
// 			TeamMembership: make(map[string][]string),
// 		}
// 		suite.testOrgs[i] = org
//
// 		cleanupOrg := suite.u.CreateOrganizationWithChecks(org.Name)
// 		suite.cleanupFuncs.PushFront(cleanupOrg)
//
// 		// We already have the "owners" team, just add 2 members.
// 		org.TeamMembership["owners"] = make([]string, 2)
// 		for j := range org.TeamMembership["owners"] {
// 			// users 0 and 1 will own the first org while users 2
// 			// and 3 will own the second org.
// 			user := suite.testUsers[2*i+j]
// 			org.TeamMembership["owners"][j] = user.Name
// 			user.OrgMembership[org.Name] = append(user.OrgMembership[org.Name], "owners")
// 		}
//
// 		// Create 2 more teams with 4 users each.
// 		for j := 0; j < 2; j++ {
// 			teamname := fmt.Sprintf("team%d", j)
// 			// No need to call the cleanup func this returns. The
// 			// team will be deleted along with the organization.
// 			suite.u.CreateManagedTeamWithChecks(org.Name, teamname)
//
// 			org.TeamMembership[teamname] = make([]string, 4)
// 			for k := range org.TeamMembership[teamname] {
// 				// Offset by 4 to skip the owners users.
// 				// Offset by 3*i to offset for different orgs.
// 				// Offset by 2*j to offset for different teams.
// 				// This puts all users in a team expect for 2
// 				// which remain solo (not a member of any
// 				// organization).
// 				user := suite.testUsers[4+(3*i)+(2*j)+k]
// 				org.TeamMembership[teamname][k] = user.Name
// 				user.OrgMembership[org.Name] = append(user.OrgMembership[org.Name], teamname)
// 			}
// 		}
// 	}
// }
//
// // TearDownSuite handles tearing down the API test suite by restoring the
// // original auth settings.
// func (suite *OrgTeamMembershipAPITestSuite) TearDownSuite() {
// 	for e := suite.cleanupFuncs.Front(); e != nil; e = e.Next() {
// 		cleanupFunc := e.Value.(func())
// 		cleanupFunc()
// 	}
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) SetupTest() {
// 	util.WipeDTRIgnorableLoggedErrors()
// 	util.WipeDockerIgnorableLoggedErrors()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TearDownTest() {
// 	suite.u.TestLogs()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) setupAllTeamMembership() {
// 	// Login as the system admin to make this easy.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
//
// 	for _, org := range suite.testOrgs {
// 		for teamname, members := range org.TeamMembership {
// 			for _, username := range members {
// 				_, err := suite.API.EnziSession().AddTeamMember(org.Name, teamname, username, forms.SetMembership{})
// 				suite.Assert().Nil(err, "expected AddTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, teamname, username, err)
// 			}
// 		}
// 	}
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) teardownAllTeamMembership() {
// 	// Login as the system admin to make this easy.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
//
// 	for _, org := range suite.testOrgs {
// 		for teamname, members := range org.TeamMembership {
// 			for _, username := range members {
// 				err := suite.API.EnziSession().DeleteTeamMember(org.Name, teamname, username)
// 				suite.Assert().Nil(err, "expected DeleteTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, teamname, username, err)
// 			}
// 		}
// 	}
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestAddTeamMember() {
// 	defer suite.teardownAllTeamMembership()
//
// 	suite.testAddTeamMember()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testAddTeamMember() {
// 	// Logout to ensure that you have to be authenticated to add a user to
// 	// a team.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	user := suite.testUsers[0]
// 	_, err := suite.API.EnziSession().AddTeamMember(org.Name, "owners", user.Name, forms.SetMembership{})
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as an ordinary user (not an org member) to ensure that they
// 	// are not authorized to add a user to a team.
// 	suite.API.Login(user.Name, user.Password)
// 	_, err = suite.API.EnziSession().AddTeamMember(org.Name, "owners", user.Name, forms.SetMembership{})
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as the system admin user to add the first member of the
// 	// "owners" team.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 	_, err = suite.API.EnziSession().AddTeamMember(org.Name, "owners", user.Name, forms.SetMembership{})
// 	suite.Assert().Nil(err, "expected AddTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, "owners", user.Name, err)
//
// 	// Login as the owner user we just added to try to add another user to
// 	// a team.
// 	suite.API.Login(user.Name, user.Password)
//
// 	// Ensure an owner user can only add users to teams within their own
// 	// organization.
// 	_, err = suite.API.EnziSession().AddTeamMember(suite.testOrgs[1].Name, "owners", suite.testUsers[1].Name, forms.SetMembership{})
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// The owner user should be able to add the other users to the teams
// 	// within the org. This loop includes adding the owner themselves which
// 	// tests that adding members is idempotent.
// 	for teamname, members := range org.TeamMembership {
// 		for _, username := range members {
// 			_, err = suite.API.EnziSession().AddTeamMember(org.Name, teamname, username, forms.SetMembership{})
// 			suite.Assert().Nil(err, "expected AddTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, teamname, username, err)
// 		}
// 	}
//
// 	// Test that user[4] (a member of team0) can't add themself to team1
// 	// because they are not in the "owners" team (but they are an org
// 	// member).
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
// 	_, err = suite.API.EnziSession().AddTeamMember(org.Name, "team1", user4.Name, forms.SetMembership{})
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	util.AppendDTRIgnorableLoggedErrors([]string{"duplicate key value violates unique constraint"})
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestDeleteTeamMember() {
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testDeleteTeamMember()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testDeleteTeamMember() {
// 	// Logout to ensure that you have to be authenticated to delete a user
// 	// from a team.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	user := suite.testUsers[0]
// 	err := suite.API.EnziSession().DeleteTeamMember(org.Name, "owners", user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a regular (non-org member) to ensure that they are not
// 	// authorized to delete a team member.
// 	user14 := suite.testUsers[14]
// 	suite.API.Login(user14.Name, user14.Password)
// 	err = suite.API.EnziSession().DeleteTeamMember(org.Name, "owners", user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as an org member (but not an owner) to ensure that they are
// 	// not authorized to delete a team member.
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
// 	err = suite.API.EnziSession().DeleteTeamMember(org.Name, "owners", user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as a system admin and ensure that they can remove a user from
// 	// a team (user4 is a member of team0).
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 	err = suite.API.EnziSession().DeleteTeamMember(org.Name, "team0", user4.Name)
// 	suite.Assert().Nil(err, "expected DeleteTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, "team0", user4.Name, err)
//
// 	// Login as the org owner and delete the remaining users from the teams
// 	// and remove themselves last.
// 	suite.API.Login(user.Name, user.Password)
// 	for teamname, usernames := range org.TeamMembership {
// 		for _, username := range usernames {
// 			if teamname == "owners" && username == user.Name {
// 				// Don't delete yourself from the owners team yet.
// 				continue
// 			}
// 			err = suite.API.EnziSession().DeleteTeamMember(org.Name, teamname, username)
// 			suite.Assert().Nil(err, "expected DeleteTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, teamname, username, err)
// 		}
// 	}
//
// 	// Lastly, delete them from their owners team.
// 	err = suite.API.EnziSession().DeleteTeamMember(org.Name, "owners", user.Name)
// 	suite.Assert().Nil(err, "expected DeleteTeamMember to succeed for org %q, team %q, user %q: %s", org.Name, "owners", user.Name, err)
// }

// func (suite *OrgTeamMembershipAPITestSuite) TestCheckTeamMembership() {
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testCheckTeamMembership()
// }

// func (suite *OrgTeamMembershipAPITestSuite) testCheckTeamMembership() {
// 	// Logout to ensure that you have to be authenticated to check team
// 	// membership.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	user := suite.testUsers[0]
// 	_, err := suite.API.CheckTeamMembership(org.Name, "owners", user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a regular (non-org member) to ensure that they are not
// 	// authorized to check team membership.
// 	user14 := suite.testUsers[14]
// 	suite.API.Login(user14.Name, user14.Password)
// 	_, err = suite.API.CheckTeamMembership(org.Name, "owners", user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as an org member (user4 is a member of team0) and ensure that
// 	// they can check team membership within their org.
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
// 	isMember, err := suite.API.CheckTeamMembership(org.Name, "owners", user.Name)
// 	suite.Assert().Nil(err, "expected CheckTeamMembership to succeed for org %q, team %q, user %q: %s", org.Name, "owners", user.Name, err)
// 	suite.Assert().True(isMember, "expected user %q to be a member of team %q in org %q", user.Name, "owners", org.Name)
//
// 	// They are not in the team themselves though.
// 	isMember, err = suite.API.CheckTeamMembership(org.Name, "owners", user4.Name)
// 	suite.Assert().Nil(err, "expected CheckTeamMembership to succeed for org %q, team %q, user %q: %s", org.Name, "owners", user4.Name, err)
// 	suite.Assert().False(isMember, "expected user %q to not be a member of team %q in org %q", user4.Name, "owners", org.Name)
//
// 	// And they can't check membership in an org they're not in.
// 	otherOrg := suite.testOrgs[1]
// 	_, err = suite.API.CheckTeamMembership(otherOrg.Name, "owners", user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// And, of course, the system admin can check any team membership.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 	isMember, err = suite.API.CheckTeamMembership(org.Name, "team0", user4.Name)
// 	suite.Assert().Nil(err, "expected CheckTeamMembership to succeed for org %q, team %q, user %q: %s", org.Name, "team0", user4.Name, err)
// 	suite.Assert().True(isMember, "expected user %q to be a member of team %q in org %q", user4.Name, "team0", org.Name)
//
// 	isMember, err = suite.API.CheckTeamMembership(org.Name, "team1", user4.Name)
// 	suite.Assert().Nil(err, "expected CheckTeamMembership to succeed for org %q, team %q, user %q: %s", org.Name, "team1", user4.Name, err)
// 	suite.Assert().False(isMember, "expected user %q to not be a member of team %q in org %q", user4.Name, "team1", org.Name)
// }
//
// func assertEqualNames(t *testing.T, expectedNames, actualNames []string) {
// 	expectedNameSet := make(map[string]struct{})
// 	for _, expectedName := range expectedNames {
// 		expectedNameSet[expectedName] = struct{}{}
// 	}
//
// 	actualNameSet := make(map[string]struct{})
// 	for _, actualName := range actualNames {
// 		actualNameSet[actualName] = struct{}{}
// 	}
//
// 	assert.Equal(t, len(expectedNameSet), len(actualNameSet), "number of actual names does not match number of expected names")
// 	for expectedName := range expectedNameSet {
// 		assert.Contains(t, actualNameSet, expectedName)
// 	}
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestListOrganizationMemberTeams() {
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testListOrganizationMemberTeams()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testListOrganizationMemberTeams() {
// 	// Logout to ensure that you have to be authenticated to list
// 	// organization member teams.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	user := suite.testUsers[6] // User 6 should be in 2 teams in the org.
// 	_, err := suite.API.ListOrganizationMemberTeams(org.Name, user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a regular (non-org member) to ensure that they are not
// 	// authorized to list organization member teams.
// 	user14 := suite.testUsers[14]
// 	suite.API.Login(user14.Name, user14.Password)
// 	_, err = suite.API.ListOrganizationMemberTeams(org.Name, user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as an org member (user4 is a member of team0).
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
//
// 	// They can't list teams for members in an org they're not in.
// 	otherOrg := suite.testOrgs[1]
// 	user2 := suite.testUsers[2]
// 	_, err = suite.API.ListOrganizationMemberTeams(otherOrg.Name, user2.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Ensure that they can list teams for other members.
// 	teams, err := suite.API.ListOrganizationMemberTeams(org.Name, user.Name)
// 	suite.Assert().Nil(err, "expected ListOrganizationMemberTeams to succeed for org %q, user %q: %s", org.Name, user.Name, err)
// 	suite.Assert().Len(teams, len(user.OrgMembership[org.Name]), "user is a member of an unexpected number of teams in the organization")
//
// 	actualTeamNames := make([]string, len(teams))
// 	for i, team := range teams {
// 		actualTeamNames[i] = team.Name
// 	}
// 	assertEqualNames(suite.T(), user.OrgMembership[org.Name], actualTeamNames)
//
// 	// And, of course, the system admin can list any org member's teams.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 	for _, user := range suite.testUsers {
// 		for _, org := range suite.testOrgs {
// 			expectedTeamnames := user.OrgMembership[org.Name]
// 			teams, err := suite.API.ListOrganizationMemberTeams(org.Name, user.Name)
// 			suite.Assert().Nil(err, "expected ListOrganizationMemberTeams to succeed for org %q, user %q: %s", org.Name, user.Name, err)
// 			suite.Assert().Len(teams, len(expectedTeamnames), "user is a member of an unexpected number of teams in the organization")
//
// 			actualTeamNames := make([]string, len(teams))
// 			for i, team := range teams {
// 				actualTeamNames[i] = team.Name
// 			}
// 			assertEqualNames(suite.T(), expectedTeamnames, actualTeamNames)
// 		}
// 	}
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestListTeamMembers() {
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testListTeamMembers()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testListTeamMembers() {
// 	// Logout to ensure that you have to be authenticated to list
// 	// team members.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	_, err := suite.API.ListTeamMembers(org.Name, "owners")
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a regular (non-org member) to ensure that they are not
// 	// authorized to list team members.
// 	user14 := suite.testUsers[14]
// 	suite.API.Login(user14.Name, user14.Password)
// 	_, err = suite.API.ListTeamMembers(org.Name, "owners")
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as an org member (user4 is a member of team0).
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
//
// 	// They can't list team members in an org they're not in.
// 	otherOrg := suite.testOrgs[1]
// 	_, err = suite.API.ListTeamMembers(otherOrg.Name, "owners")
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Ensure that they can list members of a team they are in.
// 	members, err := suite.API.ListTeamMembers(org.Name, "team0")
// 	suite.Assert().Nil(err, "expected ListTeamMembers to succeed for org %q, team %q: %s", org.Name, "team0", err)
// 	suite.Assert().Len(members, len(org.TeamMembership["team0"]), "team has an unexpected number of members")
//
// 	actualMembers := make([]string, len(members))
// 	for i, member := range members {
// 		actualMembers[i] = member.Name
// 	}
// 	assertEqualNames(suite.T(), org.TeamMembership["team0"], actualMembers)
//
// 	// Ensure they can list members of any team in their org (even teams
// 	// they are not a member of).
// 	members, err = suite.API.ListTeamMembers(org.Name, "team1")
// 	suite.Assert().Nil(err, "expected ListTeamMembers to succeed for org %q, team %q: %s", org.Name, "team1", err)
// 	suite.Assert().Len(members, len(org.TeamMembership["team1"]), "team has an unexpected number of members")
//
// 	actualMembers = make([]string, len(members))
// 	for i, member := range members {
// 		actualMembers[i] = member.Name
// 	}
// 	assertEqualNames(suite.T(), org.TeamMembership["team1"], actualMembers)
//
// 	// And, of course, the system admin can list any team's members.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 	for _, org := range suite.testOrgs {
// 		for teamname, expectedMembers := range org.TeamMembership {
// 			members, err := suite.API.ListTeamMembers(org.Name, teamname)
// 			suite.Assert().Nil(err, "expected ListTeamMembers to succeed for org %q, team %q: %s", org.Name, teamname, err)
// 			suite.Assert().Len(members, len(expectedMembers), "team has an unexpected number of members")
//
// 			actualMembers = make([]string, len(members))
// 			for i, member := range members {
// 				actualMembers[i] = member.Name
// 			}
// 			assertEqualNames(suite.T(), expectedMembers, actualMembers)
// 		}
// 	}
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestListUserOrganizations() {
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testListUserOrganizations()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testListUserOrganizations() {
// 	// Logout to ensure that you have to be authenticated to list a user's
// 	// organizations.
// 	suite.API.Login("", "")
// 	user := suite.testUsers[0]
// 	_, err := suite.API.ListUserOrganizations(user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a different user to ensure that they can't list another
// 	// user's organizations.
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
// 	_, err = suite.API.ListUserOrganizations(user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// However, a system admin can.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
//
// 	// User7 is in two organizations.
// 	user7 := suite.testUsers[7]
// 	orgs, err := suite.API.ListUserOrganizations(user7.Name)
// 	suite.Assert().Nil(err, "expected ListUserOrganizations to succeed for user %q: %s", user7.Name, err)
// 	suite.Assert().Len(orgs, len(user7.OrgMembership), "user has an unexpected number of organizations")
//
// 	actualOrgnames := make([]string, len(orgs))
// 	for i, org := range orgs {
// 		actualOrgnames[i] = org.Name
// 	}
// 	expectedOrgnames := make([]string, 0, len(user7.OrgMembership))
// 	for orgname := range user7.OrgMembership {
// 		expectedOrgnames = append(expectedOrgnames, orgname)
// 	}
// 	assertEqualNames(suite.T(), expectedOrgnames, actualOrgnames)
//
// 	// And the user can query their own org membership.
// 	suite.API.Login(user7.Name, user7.Password)
// 	orgs, err = suite.API.ListUserOrganizations(user7.Name)
// 	suite.Assert().Nil(err, "expected ListUserOrganizations to succeed for user %q: %s", user7.Name, err)
// 	suite.Assert().Len(orgs, len(user7.OrgMembership), "user has an unexpected number of organizations")
//
// 	actualOrgnames = make([]string, len(orgs))
// 	for i, org := range orgs {
// 		actualOrgnames[i] = org.Name
// 	}
// 	assertEqualNames(suite.T(), expectedOrgnames, actualOrgnames)
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestListOrganizationMembers() {
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testListOrganizationMembers()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testListOrganizationMembers() {
// 	// Logout to ensure that you have to be authenticated to list an
// 	// organization's members.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	_, err := suite.API.ListOrgMembers(org.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a user who is not in the org to ensure that they can't list
// 	// the organization's members.
// 	user14 := suite.testUsers[14]
// 	suite.API.Login(user14.Name, user14.Password)
// 	_, err = suite.API.ListOrgMembers(org.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// However, a system admin can.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
//
// 	expectedMemberNames := []string{
// 		// Org0 should have the following members.
// 		suite.testUsers[0].Name, // owners
// 		suite.testUsers[1].Name, // owners
// 		suite.testUsers[4].Name, // team0
// 		suite.testUsers[5].Name, // team0
// 		suite.testUsers[6].Name, // team0 & team1
// 		suite.testUsers[7].Name, // team0 & team1
// 		suite.testUsers[8].Name, // team1
// 		suite.testUsers[9].Name, // team1
// 	}
//
// 	members, err := suite.API.ListOrgMembers(org.Name)
// 	suite.Assert().Nil(err, "expected ListOrgMembers to succeed for org %q: %s", org.Name, err)
// 	suite.Assert().Len(members, len(expectedMemberNames), "organization has an unexpected number of members")
//
// 	actualMemberNames := make([]string, len(members))
// 	for i, member := range members {
// 		actualMemberNames[i] = member.Name
// 	}
// 	assertEqualNames(suite.T(), expectedMemberNames, actualMemberNames)
//
// 	// And any org member can query the list of users in their own orgs.
// 	user7 := suite.testUsers[7]
// 	suite.API.Login(user7.Name, user7.Password)
// 	members, err = suite.API.ListOrgMembers(org.Name)
// 	suite.Assert().Nil(err, "expected ListOrgMembers to succeed for org %q: %s", org.Name, err)
// 	suite.Assert().Len(members, len(expectedMemberNames), "organization has an unexpected number of members")
//
// 	actualMemberNames = make([]string, len(members))
// 	for i, member := range members {
// 		actualMemberNames[i] = member.Name
// 	}
// 	assertEqualNames(suite.T(), expectedMemberNames, actualMemberNames)
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestCheckOrganizationMembership() {
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testCheckOrganizationMembership()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testCheckOrganizationMembership() {
// 	// Logout to ensure that you have to be authenticated to check org
// 	// membership.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	user := suite.testUsers[0]
// 	_, err := suite.API.CheckOrganizationMembership(org.Name, user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a regular (non-org member) to ensure that they are not
// 	// authorized to check org membership.
// 	user14 := suite.testUsers[14]
// 	suite.API.Login(user14.Name, user14.Password)
// 	_, err = suite.API.CheckOrganizationMembership(org.Name, user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as an org member (user4 is a member of team0) and ensure that
// 	// they can check membership within their org.
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
// 	isMember, err := suite.API.CheckOrganizationMembership(org.Name, user.Name)
// 	suite.Assert().Nil(err, "expected CheckOrganizationMembership to succeed for org %q, user %q: %s", org.Name, user.Name, err)
// 	suite.Assert().True(isMember, "expected user %q to be a member of org %q", user.Name, org.Name)
//
// 	// User2 isn't in their org though.
// 	user2 := suite.testUsers[2]
// 	isMember, err = suite.API.CheckOrganizationMembership(org.Name, user2.Name)
// 	suite.Assert().Nil(err, "expected CheckOrganizationMembership to succeed for org %q, user %q: %s", org.Name, user2.Name, err)
// 	suite.Assert().False(isMember, "expected user %q to not be a member of org %q", user2.Name, org.Name)
//
// 	// And they can't check membership in an org they're not in.
// 	otherOrg := suite.testOrgs[1]
// 	_, err = suite.API.CheckOrganizationMembership(otherOrg.Name, user2.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// And, of course, the system admin can check any org membership.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 	isMember, err = suite.API.CheckOrganizationMembership(org.Name, user.Name)
// 	suite.Assert().Nil(err, "expected CheckOrganizationMembership to succeed for org %q, user %q: %s", org.Name, user.Name, err)
// 	suite.Assert().True(isMember, "expected user %q to be a member of org %q", user.Name, org.Name)
//
// 	isMember, err = suite.API.CheckOrganizationMembership(org.Name, user2.Name)
// 	suite.Assert().Nil(err, "expected CheckOrganizationMembership to succeed for org %q, user %q: %s", org.Name, user2.Name, err)
// 	suite.Assert().False(isMember, "expected user %q to not be a member of org %q", user2.Name, org.Name)
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) TestDeleteOrganizationMember() {
// 	if suite.u.IsSuiteRunningInLDAPMode() {
// 		// This test case doesn't make sense since the API endpoint is disabled in LDAP mode
// 		return
// 	}
// 	defer suite.teardownAllTeamMembership()
// 	suite.setupAllTeamMembership()
//
// 	suite.testDeleteOrganizationMember()
// }
//
// func (suite *OrgTeamMembershipAPITestSuite) testDeleteOrganizationMember() {
// 	// Logout to ensure that you have to be authenticated to delete a user
// 	// from an organization.
// 	suite.API.Login("", "")
// 	org := suite.testOrgs[0]
// 	user := suite.testUsers[0]
// 	err := suite.API.DeleteOrganizationMember(org.Name, user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusUnauthorized, errors.ErrorCodeNotAuthenticated)
//
// 	// Login as a regular (non-org member) to ensure that they are not
// 	// authorized to delete an org member.
// 	user14 := suite.testUsers[14]
// 	suite.API.Login(user14.Name, user14.Password)
// 	err = suite.API.DeleteOrganizationMember(org.Name, user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as an org member (but not an owner) to ensure that they are
// 	// not authorized to delete an org member.
// 	user4 := suite.testUsers[4]
// 	suite.API.Login(user4.Name, user4.Password)
// 	err = suite.API.DeleteOrganizationMember(org.Name, user.Name)
// 	suite.u.AssertErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
//
// 	// Login as a system admin and ensure that they can remove a user from
// 	// an organization.
// 	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
// 	err = suite.API.DeleteOrganizationMember(org.Name, user4.Name)
// 	suite.Assert().Nil(err, "expected DeleteOrganizationMember to succeed for org %q, user %q: %s", org.Name, user4.Name, err)
//
// 	// Login as the org owner and delete the remaining users from the org
// 	// and remove themselves last. This will iterate over some users
// 	// multiple times (those in multiple teams) but should still succeed
// 	// because the operation is idempotent.
// 	suite.API.Login(user.Name, user.Password)
// 	for _, usernames := range org.TeamMembership {
// 		for _, username := range usernames {
// 			if username == user.Name {
// 				// Don't delete yourself from the org yet.
// 				continue
// 			}
// 			err = suite.API.DeleteOrganizationMember(org.Name, username)
// 			suite.Assert().Nil(err, "expected DeleteOrganizationMember to succeed for org %q, user %q: %s", org.Name, username, err)
// 		}
// 	}
//
// 	// Lastly, delete themself from the org.
// 	err = suite.API.DeleteOrganizationMember(org.Name, user.Name)
// 	suite.Assert().Nil(err, "expected DeleteOrganizationMember to succeed for org %q, user %q: %s", org.Name, user.Name, err)
// }
