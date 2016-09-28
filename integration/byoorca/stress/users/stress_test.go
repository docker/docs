package stress

import (
	"testing"

	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/suite"
)

var (
	UserCount = utils.GetStressObjectCount(200)
)

// TODO - not idempotent (will break on second run once the auth service integrates)
func (s *TestSuite) TestAddUsers() {
	utils.TestAddUsers(s.T(), "https://"+s.OrcaClientTestSuite.Client.URL.Host, UserCount)
}

type TestSuite struct {
	utils.OrcaClientTestSuite
}

func TestStress(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
