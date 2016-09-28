package priv

import (
	"testing"

	"github.com/docker/orca/integration/utils"
	"github.com/stretchr/testify/suite"
)

func (s *TestSuite) TestEscalationsBlocked() {
	utils.TestEscalationsBlocked(s.T(), "https://"+s.OrcaClientTestSuite.Client.URL.Host)
}

type TestSuite struct {
	utils.OrcaClientTestSuite
}

func TestStress(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
