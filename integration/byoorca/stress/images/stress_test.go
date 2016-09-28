package stress

import (
	"testing"

	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/suite"
)

var (
	ImageCount = utils.GetStressObjectCount(200)
)

func (s *TestSuite) TestBuildImages() {
	utils.TestBuildImages(s.T(), "https://"+s.OrcaClientTestSuite.Client.URL.Host, ImageCount)
}

type TestSuite struct {
	utils.OrcaClientTestSuite
}

func TestStress(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
