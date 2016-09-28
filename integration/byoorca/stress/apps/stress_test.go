package stress

import (
	"testing"

	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/suite"
)

var (
	AppCount = utils.GetStressObjectCount(200)
)

func (s *TestSuite) TestComposeApps() {
	utils.TestComposeApps(s.T(), "https://"+s.OrcaClientTestSuite.Client.URL.Host, AppCount)
}

type TestSuite struct {
	utils.OrcaClientTestSuite
}

func TestStress(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
