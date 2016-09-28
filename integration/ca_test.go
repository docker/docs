package integration

import (
	"bytes"
	"testing"

	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/suite"
)

type CATestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util
}

func (suite *CATestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
}

func (suite *CATestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *CATestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *CATestSuite) TestGetCA() {
	ca, err := suite.API.GetCA()
	if err != nil {
		suite.T().Fatalf("Failed to get ca: %s", err)
	}

	if err != nil {
		suite.T().Fatalf("Couldnt read ca response body: %s", err)
	}

	if !bytes.Contains(ca, []byte("-----BEGIN CERTIFICATE-----")) {
		suite.T().Fatalf("Body should contain 'BEGIN CERTIFICATE' but is: %s", ca)
	}
	if !bytes.Contains(ca, []byte("-----END CERTIFICATE-----")) {
		suite.T().Fatalf("Body should contain 'END CERTIFICATE' but is: %s", ca)
	}
}

func TestCASuite(t *testing.T) {
	suite.Run(t, new(CATestSuite))
}
