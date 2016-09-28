package integration

import (
	"fmt"
	"time"

	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/util"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pollDelay = 250 * time.Millisecond
const pollRetries = 40

func (suite *SettingsAPITestSuite) TestAutocomplete() {
	normalUsername := "poopers"
	normalPassword := "asdfaslk"
	adminOrg := "administrators"
	defer suite.u.SwitchAuth()()

	require.Nil(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	// do nothing with this user for now, except have them exist
	defer suite.u.CreateUserWithChecks(normalUsername, normalPassword)()
	require.Nil(suite.T(), suite.u.ActivateUser(normalUsername))
	defer suite.u.CreateOrganizationWithChecks(adminOrg)()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "foobar", "foobar is a foobar", "", "private")()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "foobar2", "foobar is a foobar", "", "public")()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "wat", "foobar is a foobar", "", "private")()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "poohbear", "loves honey", "", "public")()
	defer suite.u.CreateRepoWithChecks(normalUsername, "foobar", "foobar is a foobar", "", "public")()
	defer suite.u.CreateRepoWithChecks(normalUsername, "something", "something else", "", "private")()
	defer suite.u.CreateRepoWithChecks(adminOrg, "something", "something else", "", "private")()

	if err := dtrutil.Poll(pollDelay, pollRetries, func() error {
		results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
			IncludeRepositories: true,
			IncludeAccounts:     true,
			Query:               "fooba",
			Limit:               10,
		}, 0})
		if err != nil {
			return err
		}

		errorChecker := util.NewErrorChecker()
		expectedResults := []string{suite.Config.AdminUsername + "/foobar", suite.Config.AdminUsername + "/foobar2", normalUsername + "/foobar"}
		for _, result := range results.RepositoryResults {
			assert.Contains(errorChecker, expectedResults, result.FullName())
		}
		return errorChecker.Errors()
	}); err != nil {
		suite.T().Error(err)
	}

	if results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
		IncludeRepositories: true,
		IncludeAccounts:     true,
		Namespace:           suite.Config.AdminUsername,
		Query:               "fooba",
		Limit:               10,
	}, 0}); err != nil {
		suite.T().Error(err)
	} else {
		resultNames := make([]string, len(results.RepositoryResults))
		for i, result := range results.RepositoryResults {
			resultNames[i] = result.FullName()
		}
		assert.Equal(suite.T(), []string{suite.Config.AdminUsername + "/foobar", suite.Config.AdminUsername + "/foobar2"}, resultNames)
	}

	// warning: the following tests are not guaranteed to pass without race conditions because we don't wrap then in retries
	// and we don't know for sure that all indices have been updated by now.

	require.Nil(suite.T(), suite.API.Logout()) // after logout, no queries will return any useful data
	require.Nil(suite.T(), suite.API.Login(normalUsername, normalPassword))
	if results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
		IncludeRepositories: true,
		IncludeAccounts:     true,
		Query:               suite.Config.AdminUsername[:3],
		Limit:               10,
	}, 0}); err != nil {
		suite.T().Error(err)
	} else {
		resultNames := make([]string, len(results.RepositoryResults))
		for i, result := range results.RepositoryResults {
			resultNames[i] = result.FullName()
		}
		assert.Contains(suite.T(), resultNames, suite.Config.AdminUsername+"/foobar2")
		if assert.Len(suite.T(), results.AccountResults, 2) {
			assert.Equal(suite.T(), results.AccountResults[0].Name, suite.Config.AdminUsername)
			assert.Equal(suite.T(), results.AccountResults[1].Name, adminOrg)
		}
	}

	if results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
		IncludeRepositories: true,
		IncludeAccounts:     false,
		Query:               suite.Config.AdminUsername[:3],
		Limit:               10,
	}, 0}); err != nil {
		suite.T().Error(err)
	} else {
		resultNames := make([]string, len(results.RepositoryResults))
		for i, result := range results.RepositoryResults {
			resultNames[i] = result.FullName()
		}
		assert.Equal(suite.T(), []string{suite.Config.AdminUsername + "/foobar2", suite.Config.AdminUsername + "/poohbear"}, resultNames)
	}

	if results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
		IncludeRepositories: true,
		IncludeAccounts:     true,
		Query:               suite.Config.AdminUsername + "/",
		Limit:               10,
	}, 0}); err != nil {
		suite.T().Error(err)
	} else {
		resultNames := make([]string, len(results.RepositoryResults))
		for i, result := range results.RepositoryResults {
			resultNames[i] = result.FullName()
		}
		assert.Equal(suite.T(), resultNames, []string{suite.Config.AdminUsername + "/foobar2", suite.Config.AdminUsername + "/poohbear"})
		if assert.Len(suite.T(), results.AccountResults, 2) {
			assert.False(suite.T(), results.AccountResults[0].IsOrg)
			assert.True(suite.T(), results.AccountResults[1].IsOrg)
		}
	}

	if results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
		IncludeRepositories: true,
		IncludeAccounts:     true,
		Query:               suite.Config.AdminUsername + "/foobar2",
		Limit:               10,
	}, 0}); err != nil {
		suite.T().Error(err)
	} else {
		resultNames := make([]string, len(results.RepositoryResults))
		for i, result := range results.RepositoryResults {
			resultNames[i] = result.FullName()
		}
		assert.Equal(suite.T(), []string{suite.Config.AdminUsername + "/foobar2"}, resultNames)
		if assert.Len(suite.T(), results.AccountResults, 2) {
			assert.False(suite.T(), results.AccountResults[0].IsOrg)
			assert.True(suite.T(), results.AccountResults[1].IsOrg)
		}
	}

	if results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
		IncludeRepositories: true,
		IncludeAccounts:     true,
		Query:               "/poo",
		Limit:               10,
	}, 0}); err != nil {
		suite.T().Error(err)
	} else {
		resultNames := make([]string, len(results.RepositoryResults))
		for i, result := range results.RepositoryResults {
			resultNames[i] = result.FullName()
		}
		assert.Equal(suite.T(), []string{suite.Config.AdminUsername + "/poohbear"}, resultNames)
		// expect all of the accounts to be returned here...
		assert.Len(suite.T(), results.AccountResults, 4)
	}

	if results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
		IncludeRepositories: false,
		IncludeAccounts:     true,
		Query:               suite.Config.AdminUsername[:3],
		Limit:               10,
	}, 0}); err != nil {
		suite.T().Error(err)
	} else {
		if assert.Len(suite.T(), results.AccountResults, 2) {
			assert.False(suite.T(), results.AccountResults[0].IsOrg)
			assert.True(suite.T(), results.AccountResults[1].IsOrg)
		}
	}
}

func (suite *SettingsAPITestSuite) TestDockerSearch() {
	// this test's setup is mostly copied from TestSearch
	defer suite.u.SwitchAuth()()

	require.Nil(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	normalUser, deferFunc := suite.u.CreateActivateRandomUser()
	defer deferFunc()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "foobar", "foobar is a foobar", "", "private")()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "foobar2", "foobar is a foobar", "", "public")()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "wat", "foobar is a foobar", "", "private")()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, "something", "something else", "", "private")()
	defer suite.u.CreateRepoWithChecks(normalUser.Name, "something", "something else", "", "private")()
	for i := 0; i < 30; i++ {
		repoName := fmt.Sprintf("asdf%d", i)
		repoDesc := fmt.Sprintf("%s is a %s", repoName, repoName)
		defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, repoName, repoDesc, "", "private")()
	}

	util.AppendDockerIgnorableLoggedErrors([]string{"404"})

	if err := dtrutil.Poll(pollDelay, pollRetries, func() error {
		results, err := suite.API.Autocomplete(apiclient.SearchOptions{forms.SearchOptions{
			IncludeRepositories: true,
			IncludeAccounts:     true,
			Query:               "fooba",
			Limit:               10,
		}, 0})
		if err != nil {
			return err
		}

		errorChecker := util.NewErrorChecker()
		resultNames := make([]string, len(results.RepositoryResults))
		for i, result := range results.RepositoryResults {
			resultNames[i] = result.FullName()
		}
		assert.Equal(errorChecker, []string{suite.Config.AdminUsername + "/foobar", suite.Config.AdminUsername + "/foobar2"}, resultNames)
		return errorChecker.Errors()
	}); err != nil {
		suite.T().Error(err)
	}

	results, err := suite.Docker.SearchImages("foobar", suite.Config.DTRHost, &dockerclient.AuthConfig{
		Username: suite.Config.AdminUsername,
		Password: suite.Config.AdminPassword,
	})
	assert.Nil(suite.T(), err)
	resultRepoNames := make([]string, len(results))
	for i, result := range results {
		resultRepoNames[i] = result.Name
	}
	if assert.Len(suite.T(), resultRepoNames, 2) {
		assert.Contains(suite.T(), resultRepoNames, suite.Config.AdminUsername+"/foobar")
		assert.Contains(suite.T(), resultRepoNames, suite.Config.AdminUsername+"/foobar2")
	}

	results, err = suite.Docker.SearchImages("asdf", suite.Config.DTRHost, &dockerclient.AuthConfig{
		Username: suite.Config.AdminUsername,
		Password: suite.Config.AdminPassword,
	})
	assert.Nil(suite.T(), err)
	assert.Len(suite.T(), results, 25)
}
