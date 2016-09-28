package integration

import (
	"net/http"
	"path"
	"testing"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"
	"github.com/samalba/dockerclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ReposAPITestSuite is an integration test suite for the DTR API.
type ReposAPITestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	setOldAuth func()
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *ReposAPITestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	suite.setOldAuth = suite.u.SwitchAuth()
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *ReposAPITestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *ReposAPITestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *ReposAPITestSuite) TearDownTest() {
	suite.u.TestLogs()
}

// TestCreateRepo tests creation of repos for users and orgs
func (suite *ReposAPITestSuite) TestCreateRepo() {
	var err error
	testRepo := "testrepo"

	// create a repo under the admin namespace using the admin
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "public")()

	// Creating a repo with the same name should result in an error.
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	_, err = suite.API.CreateRepository(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "public")
	suite.u.RequireErrorCodes(err, http.StatusBadRequest, errors.ErrorCodeRepositoryExists)

	// create a non-admin user
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	user, cleanup := suite.u.CreateActivateRandomUser()
	defer cleanup()
	username := user.Name
	password := user.Password

	// try to create a repo as the non-admin user
	suite.API.Login(username, password)
	defer suite.u.CreateRepoWithChecks(username, testRepo, "blah", "blah blah", "public")()

	// Creating a repo with the same name should result in an error.
	_, err = suite.API.CreateRepository(username, testRepo, "blah", "blah blah", "public")
	suite.u.RequireErrorCodes(err, http.StatusBadRequest, errors.ErrorCodeRepositoryExists)

	// Creating a repo with a regular user without permissions under another user's namespace should not be allowed
	_, err = suite.API.CreateRepository(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "public")
	suite.u.RequireErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)

	// Creating a repo for a non existent account should fail
	_, err = suite.API.CreateRepository("arstareishtairs", testRepo, "blah", "blah blah", "public")
	suite.u.RequireErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchAccount)

	// You must be authenticated to create a repo
	require.Nil(suite.T(), suite.API.Logout())
	_, err = suite.API.CreateRepository(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "public")
	require.NotNil(suite.T(), err, "%s", err)

	// Create an org and make sure the admin can create a repo under it
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	orgname := "testorg"
	defer suite.u.CreateOrganizationWithChecks(orgname)()
	defer suite.u.CreateRepoWithChecks(orgname, testRepo, "blah", "blah blah", "public")()

	util.AppendDTRIgnorableLoggedErrors([]string{"duplicate key value violates unique constraint"})
}

func (suite *ReposAPITestSuite) TestListRepos() {
	testRepo1 := "test123"
	testRepo2 := "test1234"

	// create a non-admin user

	user1, delfunc1 := suite.u.CreateActivateRandomUser()
	defer delfunc1()
	username := user1.Name
	password := user1.Password

	// create another non-admin user
	user2, delfunc2 := suite.u.CreateActivateRandomUser()
	defer delfunc2()
	username2 := user2.Name
	password2 := user2.Password

	// create an org
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	orgname := "testorg"
	defer suite.u.CreateOrganizationWithChecks(orgname)()

	// create some repos for both users and the org
	suite.API.Login(username, password)
	defer suite.u.CreateRepoWithChecks(username, testRepo1, "blah", "blah blah", "public")()
	defer suite.u.CreateRepoWithChecks(username, testRepo2, "blah", "blah blah", "private")()
	suite.API.Login(username2, password2)
	defer suite.u.CreateRepoWithChecks(username2, testRepo1, "blah", "blah blah", "public")()
	defer suite.u.CreateRepoWithChecks(username2, testRepo2, "blah", "blah blah", "private")()
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateRepoWithChecks(orgname, testRepo1, "blah", "blah blah", "public")()
	defer suite.u.CreateRepoWithChecks(orgname, testRepo2, "blah", "blah blah", "private")()

	// the first user can see all of their repos and the public ones
	suite.API.Login(username, password)
	repos, err := suite.API.ListRepositories(username)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 2)
	repos, err = suite.API.ListRepositories(username2)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 1)
	repos, err = suite.API.ListRepositories(orgname)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 1)
	repos, err = suite.API.ListAllRepositories()
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 4)

	// the admin can see all repos
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	repos, err = suite.API.ListRepositories(username)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 2)
	repos, err = suite.API.ListRepositories(username2)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 2)
	repos, err = suite.API.ListRepositories(orgname)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 2)
	repos, err = suite.API.ListAllRepositories()
	assert.Nil(suite.T(), err, "%s", err)
	assert.Len(suite.T(), repos, 6)
}

func (suite *ReposAPITestSuite) TestUpdateAndGetDetailsRepo() {
	testRepo := "bananas"
	shortDesc1 := "blah"
	longDesc1 := "blah blah"
	visibility1 := "private"
	shortDesc2 := "ring"
	longDesc2 := "ring ring ring bananaphone!"
	visibility2 := "public"

	// create user

	user1, delfunc1 := suite.u.CreateActivateRandomUser()
	defer delfunc1()
	username := user1.Name
	password := user1.Password

	// create a repo

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, shortDesc1, longDesc1, visibility1)()

	// get repo details should work

	repo1, err := suite.API.GetRepository(suite.Config.AdminUsername, testRepo)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Equal(suite.T(), suite.Config.AdminUsername, repo1.Namespace)
	assert.Equal(suite.T(), testRepo, repo1.Name)
	assert.Equal(suite.T(), shortDesc1, repo1.ShortDescription)
	assert.Equal(suite.T(), longDesc1, repo1.LongDescription)
	assert.Equal(suite.T(), visibility1, repo1.Visibility)

	// make sure a non-admin can't read the admin's private repo

	suite.API.Login(username, password)
	_, err = suite.API.GetRepository(suite.Config.AdminUsername, testRepo)
	suite.u.RequireErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)

	// make sure a non-admin can't touch the admin's private repo

	suite.API.Login(username, password)
	suite.u.API.Login(username, password)
	_, err = suite.API.UpdateRepository(suite.Config.AdminUsername, testRepo, apiclient.RepositoryUpdateForm{
		ShortDescription: &shortDesc2,
		LongDescription:  &longDesc2,
		Visibility:       &visibility2,
	})
	suite.u.RequireErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)

	// make sure an admin can update their own repo

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	_, err = suite.API.UpdateRepository(suite.Config.AdminUsername, testRepo, apiclient.RepositoryUpdateForm{
		ShortDescription: &shortDesc2,
		LongDescription:  &longDesc2,
		Visibility:       &visibility2,
	})
	assert.Nil(suite.T(), err, "%s", err)

	repo2, err := suite.API.GetRepository(suite.Config.AdminUsername, testRepo)
	assert.Nil(suite.T(), err, "%s", err)
	assert.Equal(suite.T(), repo1.ID, repo2.ID)
	assert.Equal(suite.T(), suite.Config.AdminUsername, repo2.Namespace)
	assert.Equal(suite.T(), testRepo, repo2.Name)
	assert.Equal(suite.T(), shortDesc2, repo2.ShortDescription)
	assert.Equal(suite.T(), longDesc2, repo2.LongDescription)
	assert.Equal(suite.T(), visibility2, repo2.Visibility)

	// make sure you can't set an arbitrary visibility string

	notvisibility := "lulz"
	_, err = suite.API.UpdateRepository(suite.Config.AdminUsername, testRepo, apiclient.RepositoryUpdateForm{
		Visibility: &notvisibility,
	})
	suite.u.RequireErrorCodes(err, http.StatusBadRequest, errors.ErrorCodeInvalidRepositoryVisibility)

	// make sure a non-admin can't touch the admin's repo

	suite.API.Login(username, password)
	_, err = suite.API.UpdateRepository(suite.Config.AdminUsername, testRepo, apiclient.RepositoryUpdateForm{
		ShortDescription: &shortDesc2,
		LongDescription:  &longDesc2,
		Visibility:       &visibility2,
	})
	suite.u.RequireErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)
}

func (suite *ReposAPITestSuite) TestDeleteRepo() {
	// The only thing not covered by the create repo tests is deleting a repo you don't own
	testRepo := "thesilkroad"

	user1, delfunc1 := suite.u.CreateActivateRandomUser()
	defer delfunc1()
	user2, delfunc2 := suite.u.CreateActivateRandomUser()
	defer delfunc2()

	suite.API.Login(user2.Name, user2.Password)
	_, err := suite.API.CreateRepository(user2.Name, testRepo, "blah", "blah blah", "public")
	require.Nil(suite.T(), err)

	// prepare for pushing
	tag := "whatever"
	smallImage := "tianon/true:latest"
	err = suite.Docker.PullImage(smallImage, nil)
	require.Nil(suite.T(), err)

	defer suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+user2.Name, testRepo, tag, smallImage)()

	// make sure we can push before we delete
	authConfig := dockerclient.AuthConfig{
		Username: user2.Name,
		Password: user2.Password,
		Email:    "a@a.a",
	}
	err = suite.Docker.PushImage(path.Join(suite.Config.DTRHost, user2.Name, testRepo), tag, &authConfig)
	require.Nil(suite.T(), err)

	// deleting a repo you don't own should error out
	suite.API.Login(user1.Name, user1.Password)
	err = suite.API.DeleteRepository(user2.Name, testRepo)
	suite.u.RequireErrorCodes(err, http.StatusForbidden, errors.ErrorCodeNotAuthorized)

	// deleting a repo you own should be fine
	suite.API.Login(user2.Name, user2.Password)
	err = suite.API.DeleteRepository(user2.Name, testRepo)
	assert.Nil(suite.T(), err, "%s", err)

	// you shouldn't be able to push after deleting the repo
	err = suite.Docker.PushImage(suite.Config.DTRHost+"/"+user2.Name, tag, &authConfig)
	assert.NotNil(suite.T(), err)
}

func TestReposAPISuite(t *testing.T) {
	suite.Run(t, new(ReposAPITestSuite))
}
