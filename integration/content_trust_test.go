package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"
	"github.com/docker/engine-api/client"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// These tests ensure that content trust is correctly set up with DTR by shelling out to
// docker and pushing and pulling images, because it's important to make sure that Docker's
// particular Notary+Registry+Auth handling works

// All the integration tests are run in a gobase container. So how this test suite works is that
// it uses the engine API to connect to the host docker socket (within the integration test container).
// It then creates a DIND container to run the docker client binary to push/pull from the DTR container.

// ContentTrustTestSuite is an integration test suite for the content trust functionality of DTR
type ContentTrustTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	imageName      string
	adminAccountId string
	setOldAuth     func()

	client     *client.Client
	registryCA []byte
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *ContentTrustTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	suite.setOldAuth = suite.u.SwitchAuth()
	suite.imageName = "tianon/true"

	account, err := suite.API.GetAccount(suite.Config.AdminUsername)
	require.Nil(suite.T(), err, "Getting account returned error: %s", err)
	suite.adminAccountId = account.ID

	err = suite.Docker.PullImage(suite.imageName, nil)
	defer func() {
		if err != nil {
			suite.setOldAuth()
		}
	}()
	require.NoError(suite.T(), err)

	suite.client, err = client.NewClient("unix:///var/run/docker.sock", suite.Config.DockerAPIVersion, nil,
		map[string]string{"User-Agent": "engine-api-cli-1.0"})
	require.NoError(suite.T(), err)
	_, err = suite.client.ServerVersion(context.Background())
	require.NoError(suite.T(), err)

	// get the CA so we ca
	suite.registryCA, err = suite.API.GetCA()
	require.NoError(suite.T(), err)
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *ContentTrustTestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *ContentTrustTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *ContentTrustTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *ContentTrustTestSuite) TestPushPullSuccess() {
	var err error
	testRepo := util.RandStringBytes(20)
	tag := util.RandStringBytes(20)

	// create a private repo under the admin namespace using the admin - private because we want to make
	// sure notary is enforced with auth as well
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	// skip return function since we delete this same repo later
	suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "private")

	// tag the image we want to push, skip return function since we delete this same repo later
	suite.u.TagImageWithChecks(
		suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag, suite.imageName)

	// make sure DTR has no tags
	response, err := suite.API.GetRepositoryTags(suite.Config.AdminUsername, testRepo)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 0, len(response))

	img := fmt.Sprintf("%s/%s/%s:%s", suite.Config.DTRHost, suite.Config.AdminUsername, testRepo, tag)
	// set up docker client and log in
	binDocker, err := util.NewDockerClientWithTrust(suite.client, suite.Config.DTRHost, "docker:dind", suite.registryCA)
	defer binDocker.Cleanup()

	func() {
		// push without content trust, using the Docker client library (not shelling out)
		// defer the deletion until we push to this tag again, otherwise we
		// have dangling manifests
		defer suite.u.PushImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag)()

		// make sure the image appears in the list as unsigned
		func() {
			response, err := suite.API.GetTagTrust(suite.Config.AdminUsername, testRepo, tag)
			require.NoError(suite.T(), err)
			require.Equal(suite.T(), tag, response.Name)
			require.False(suite.T(), response.InNotary)
		}()
		require.NoError(suite.T(), err)
		_, err = binDocker.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
		require.NoError(suite.T(), err)

		// pulling with content trust fails, because even though the image exists, there is no trust data
		_, err = binDocker.Pull(img, true)
		require.Error(suite.T(), err)

		// pulling it without content trust succeeds
		_, err = binDocker.Pull(img, false)
		require.NoError(suite.T(), err)

	}()

	// push to it with content trust
	_, err = binDocker.Push(img, true)
	require.NoError(suite.T(), err)

	// pull with content trust will succeed now
	_, err = binDocker.Pull(img, true)
	require.NoError(suite.T(), err)

	// make sure the image appears in the list as signed and with the same hash
	func() {
		response, err := suite.API.GetTagTrust(suite.Config.AdminUsername, testRepo, tag)
		require.Nil(suite.T(), err)
		require.Equal(suite.T(), tag, response.Name)
		require.Equal(suite.T(), response.InNotary, true)
		// TODO: uncomment this once https://github.com/docker/dhe-deploy/issues/2676 is fixed
		// require.Equal(suite.T(), response.HashMismatch, false)
	}()

	// logout and try to delete the repo using the API, which should fail
	require.NoError(suite.T(), suite.API.Logout())
	require.Error(suite.T(), suite.API.DeleteRepository(suite.Config.AdminUsername, testRepo))

	// login as a random user and make sure we cannot delete the admin repo through the API
	user, cleanup := suite.u.CreateActivateRandomUser()
	defer cleanup()
	require.NoError(suite.T(), suite.API.Login(user.Name, user.Password))
	// deleting as a non-admin should fail
	require.Error(suite.T(), suite.API.DeleteRepository(suite.Config.AdminUsername, testRepo))
	require.NoError(suite.T(), suite.API.Logout())

	// clear all binDocker's trust and login data
	_, err = binDocker.ClearConfigData()
	require.NoError(suite.T(), err)

	// since we haven't logged in on binDocker, pull should fail because the repo is private
	_, err = binDocker.Pull(img, true)
	require.Error(suite.T(), err)

	// after logging in pull should succeed
	_, err = binDocker.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	require.NoError(suite.T(), err)
	_, err = binDocker.Pull(img, true)
	require.NoError(suite.T(), err)

	// log back in as the admin at the API level
	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))

	// delete the repository from the API as an admin, and its associated trust data
	require.NoError(suite.T(), suite.API.DeleteRepository(suite.Config.AdminUsername, testRepo))

	// create and push without trust to the same reponame and check that we cannot pull without trust now since we've deleted past trust data
	// the defer statements below are for cleaning up this image later
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, "new blah same name", "new blah blah blah same name", "private")()
	defer suite.u.TagImageWithChecks(
		suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag, suite.imageName)()
	defer suite.u.PushImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag)()

	// clear all binDocker's trust and login data and then log back in
	_, err = binDocker.ClearConfigData()
	require.NoError(suite.T(), err)
	_, err = binDocker.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	require.NoError(suite.T(), err)

	// Pulling with trust will fail
	_, err = binDocker.Pull(img, true)
	require.Error(suite.T(), err)

	// Ensure that the tag exist but there is no signed data for it
	func() {
		response, err := suite.API.GetTagTrust(suite.Config.AdminUsername, testRepo, tag)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), tag, response.Name)
		require.False(suite.T(), response.InNotary)
	}()
}

func TestContentTrustSuite(t *testing.T) {
	suite.Run(t, new(ContentTrustTestSuite))
}
