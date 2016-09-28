package stress

import (
	"testing"

	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/suite"
)

var (
	ContainerCount = utils.GetStressObjectCount(200)
)

func (s *TestSuite) TestCreateContainers() {
	utils.TestCreateContainers(s.T(), "https://"+s.OrcaClientTestSuite.Client.URL.Host, ContainerCount)
}

type TestSuite struct {
	utils.OrcaClientTestSuite
}

func TestStress(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
