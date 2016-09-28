package integration

import (
	"net/http"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/integration/util"

	enzierrors "github.com/docker/orca/enzi/api/errors"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (suite *SettingsAPITestSuite) TestCreateAccount() {
	passwords := map[string]string{}
	defer suite.u.SwitchAuth()()

	normalUser := "poo"
	passwords[normalUser] = util.GenerateRandomPassword(10)
	normalUser2 := "ooo"
	passwords[normalUser2] = util.GenerateRandomPassword(10)
	defer suite.u.CreateUserWithChecks(normalUser, passwords[normalUser])()
	defer suite.u.CreateUserWithChecks(normalUser2, passwords[normalUser2])()

	require.Nil(suite.T(), suite.API.Logout())
	if err := suite.u.ActivateUser(normalUser); err != nil {
		suite.u.AssertEnziErrorCode(err, http.StatusUnauthorized, enzierrors.AuthenticationRequired())
	} else {
		suite.T().Fatal("unauthenticated user should not be able to activate users")
	}

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	if err := suite.u.ActivateUser(normalUser); err != nil {
		suite.T().Error(err)
	}

	// This only make sense in managed mode
	if suite.u.IsSuiteRunningInManagedMode() {
		err := suite.API.Login(normalUser, passwords[normalUser])
		normalUserSession := suite.u.API.EnziSession()
		require.Nil(suite.u.T(), err)
		if _, err := normalUserSession.UpdateAccount(normalUser, forms.UpdateAccount{IsActive: &[]bool{false}[0]}); err != nil {
			assert.Contains(suite.T(), err.Error(), errors.ErrorCodeNotAuthorized.Code)
		} else {
			suite.T().Error("nonadmin user should not be able to activate users")
		}
	}

	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	shittyPasswordUser := "god"
	passwords[shittyPasswordUser] = ""

	if shittyPasswordAccount, err := suite.u.CreateUser(shittyPasswordUser, passwords[shittyPasswordUser]); err != nil {
		assert.Nil(suite.T(), shittyPasswordAccount)
		assert.Contains(suite.T(), err.Error(), "INVALID_FORM_FIELD")
	} else {
		suite.T().Errorf("should have failed to create %q due to bad password", shittyPasswordUser)
		if err := suite.u.DeleteAccount(shittyPasswordUser); err != nil {
			suite.T().Logf("cleanup error: %v", err)
		}
	}
}
